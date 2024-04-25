package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/noolingo/user-service/internal/domain"
)

var ErrNotFound = errors.New("user not found")

type user struct {
	db *sql.DB
}

func (u *user) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	err := u.db.QueryRowContext(ctx, "select * from user where id=?", id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *user) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := u.db.QueryRowContext(ctx, "select * from user where email=?", email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *user) CreateUser(ctx context.Context, user2 *domain.User) (userID string, err error) {
	ins, err := u.db.PrepareContext(ctx,
		"insert into user(name,email,password) values (?,?,?)")
	if err != nil {
		return "", err
	}
	res, err := ins.ExecContext(ctx, user2.Name, user2.Email, user2.Password)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(res.LastInsertId()), nil
}

func (u *user) UpdateUser(ctx context.Context, user2 *domain.User) (err error) {
	ins, err := u.db.PrepareContext(ctx,
		"update user set name=?,email=?,password=? where id=?")
	if err != nil {
		return err
	}
	_, err = ins.ExecContext(ctx, user2.Name, user2.Email, user2.Password, user2.ID)
	return err
}

func (u *user) DeleteUser(ctx context.Context, id string) (err error) {
	ins, err := u.db.PrepareContext(ctx, "delete from user where id=?")
	if err != nil {
		return err
	}
	_, err = ins.ExecContext(ctx, id)
	return err
}
