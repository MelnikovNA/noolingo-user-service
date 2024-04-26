package service

import (
	"context"

	"github.com/noolingo/user-service/internal/domain"
	"github.com/noolingo/user-service/internal/repository"
	"github.com/sirupsen/logrus"
)

type User interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type Auth interface {
	SignUp(ctx context.Context, user *domain.User) (string, error)
	SignIn(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error)
	SignOut(ctx context.Context, accessToken, refreshToken string) error
	DeleteUser(ctx context.Context, id string) error
}

type Params struct {
	Logger     *logrus.Logger
	Config     *domain.Config
	Repository *repository.Repository
}

type Services struct {
	User User
	Auth Auth
}

func New(p *Params) *Services {
	return &Services{
		Auth: NewAuthService(p),
		User: NewUserService(p),
	}
}
