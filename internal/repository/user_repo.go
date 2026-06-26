package repository

import (
	"FeatureFlags/internal/domain"

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
		"SELECT * FROM users WHERE email = $1",
		email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
