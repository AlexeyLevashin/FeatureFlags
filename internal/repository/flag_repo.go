package repository

import (
	"FeatureFlags/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FlagRepo struct {
	db *sqlx.DB
}

func NewFlagRepository(db *sqlx.DB) *FlagRepo {
	return &FlagRepo{db: db}
}

func (repo FlagRepo) GetAll(filter domain.FlagFilter) ([]domain.FeatureFlag, error) {
	flags := []domain.FeatureFlag{}
	query := "SELECT id, name, description, status, environment, owner_user_id, owner_team_id, updated_at FROM feature_flags WHERE 1=1"
	args := []interface{}{}
	i := 1
	if filter.Search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", i)
		args = append(args, "%"+filter.Search+"%")
		i++
	}
	if filter.Environment != "" {
		query += fmt.Sprintf(" AND environment = $%d", i)
		args = append(args, filter.Environment)
		i++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, filter.Status)
		i++
	}
	err := repo.db.Select(&flags, query, args...)
	if err != nil {
		return []domain.FeatureFlag{}, err
	}
	return flags, nil
}

func (repo FlagRepo) GetById(id int) (domain.FeatureFlag, error) {
	flag := domain.FeatureFlag{}
	err := repo.db.Get(&flag,
		"SELECT id, name, description, status, environment, owner_user_id, owner_team_id, updated_at FROM feature_flags WHERE id = $1",
		id)
	if err != nil {
		return domain.FeatureFlag{}, err
	}
	return flag, nil
}
//func Create
