package service

import (
	"context"
	"errors"

	"github.com/noolingo/user-service/internal/domain"
	"github.com/noolingo/user-service/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	logger     *logrus.Logger
	config     *domain.Config
	repository repository.Repository
}

var (
	ErrNoUserFound = errors.New("no such user found")
)

func NewUserService(p *Params) *UserService {
	return &UserService{
		logger:     p.Logger,
		config:     p.Config,
		repository: *p.Repository,
	}
}

func (u *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUserFound
	}

	return user, err
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUserFound
	}
	return user, err
}

func (u *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	user2, err := u.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}
	user.Password = user2.Password
	return u.repository.UpdateUser(ctx, user)
}
