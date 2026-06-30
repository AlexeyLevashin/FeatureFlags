package repository

import (
	"FeatureFlags/internal/domain"
	"context"
	_ "embed"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (domain.User, error) {
	user := domain.User{}
	err := r.db.Get(&user,
		"SELECT id, email, password_hash, name, surname, team_id FROM users WHERE email = $1",
		email)
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

func (r *UserRepository) FindById(userId int) (domain.User, error) {
	user := domain.User{}
	err := r.db.Get(&user,
		"SELECT email, name, surname FROM users WHERE id = $1",
		userId)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
