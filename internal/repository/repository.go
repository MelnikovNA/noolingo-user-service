package repository

import (
	"context"
	"database/sql"

	"github.com/noolingo/user-service/internal/domain"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (userID string, err error)
	UpdateUser(ctx context.Context, user *domain.User) (err error)
	DeleteUser(ctx context.Context, id string) (err error)
}

func New(db *sql.DB) Repository {
	return &user{db: db}
}
