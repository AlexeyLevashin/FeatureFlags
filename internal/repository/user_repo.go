package repository

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user := domain.User{}
	err := r.db.GetContext(ctx, &user,
		"SELECT id, email, password_hash, name, surname, team_id FROM users WHERE email = $1",
		email)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, apperror.NotFound("Неверный email или пароль")
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

//go:embed queries/user/check_user_exists.sql
var CheckUserExists string

func (r *UserRepository) CheckExists(ctx context.Context, userId int) (bool, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, CheckUserExists, userId)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) FindById(ctx context.Context, userId int) (domain.User, error) {
	user := domain.User{}
	err := r.db.GetContext(ctx, &user,
		"SELECT email, name, surname FROM users WHERE id = $1",
		userId)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, apperror.NotFound("пользователь не найден")
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
