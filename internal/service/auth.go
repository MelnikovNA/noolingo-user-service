package service

import (
	"context"
	"errors"

	"github.com/noolingo/proto/codegen/go/apierrors"
	enc "github.com/noolingo/sha256password"
	"github.com/noolingo/user-service/internal/domain"
	"github.com/noolingo/user-service/internal/pkg/tokens"
	"github.com/noolingo/user-service/internal/repository"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	logger       *logrus.Logger
	config       *domain.Config
	repository   repository.Repository
	accessToken  tokens.JWTToken
	refreshToken tokens.JWTToken
}

func NewAuthService(p *Params) *AuthService {
	if p.Config.Auth.RefreshSecretKey == "" || p.Config.Auth.AccessSecretKey == "" {
		panic("RefreshSecretKey or AccessSecretKey is not set")
	}

	return &AuthService{
		logger:       p.Logger,
		config:       p.Config,
		repository:   *p.Repository,
		accessToken:  *tokens.New(&p.Config.Auth),
		refreshToken: *tokens.New(&p.Config.Auth),
	}
}

func (a *AuthService) DeleteUser(ctx context.Context, id string) error {
	err := a.repository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) SignIn(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error) {
	user, err := a.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", apierrors.ErrNotFound
	}
	if !enc.CompWithEncrypted(password, user.Password) {
		return "", "", apierrors.ErrForbidden
	}
	return a.makeToken(user)
}

func (a *AuthService) SignOut(ctx context.Context, accessToken, refreshToken string) error {
	panic("not implemented") //TODO
}

func (a *AuthService) SignUp(ctx context.Context, user *domain.User) (string, error) {
	user.Password, _ = enc.EncryptPassword(user.Password)
	return a.repository.CreateUser(ctx, user)
}

func (a *AuthService) makeToken(user *domain.User) (accessToken string, refreshToken string, err error) {
	accessToken, err = a.accessToken.NewToken(user.ID, a.config.Auth.AccessKeyTtl)
	if err != nil {
		a.logger.WithError(err).Errorf("error generating access token")
		return "", "", apierrors.ErrInternalServerError
	}
	refreshToken, err = a.refreshToken.NewToken(user.ID, a.config.Auth.RefreshKeyTtl)
	if err != nil {
		a.logger.WithError(err).Errorf("error generating refresh token")
	}
	return accessToken, refreshToken, err
}

func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	ParsedUserID, err := a.refreshToken.ParseToken(refreshToken)
	if err != nil {
		a.logger.WithError(err).Warn("token parse error")
		return "", "", apierrors.ErrTokenExpired
	}
	user, err := a.repository.GetUserByID(ctx, ParsedUserID)
	if err != nil {
		a.logger.WithError(err).Errorf("error in db")
		return "", "", errors.New("error in DB")
	}
	newAccessToken, newRefreshToken, err = a.makeToken(user)

	return
}
