package repository

import "FeatureFlags/internal/domain"

type UserRepository struct{}

// func NewUserRepository() *UserRepository {}

func (r *UserRepository) FindByEmail(email string) (domain.User, error) {
	return domain.User{}, nil
}
