package repository

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"context"
	"database/sql"
	"errors"

	"fmt"

	_ "embed"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type FlagRepo struct {
	db *sqlx.DB
}

func NewFlagRepository(db *sqlx.DB) *FlagRepo {
	return &FlagRepo{db: db}
}

func (repo FlagRepo) GetAll(ctx context.Context, filter domain.FlagFilter) ([]domain.FeatureFlag, error) {
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

	err := repo.db.SelectContext(ctx, &flags, query, args...)

	if err != nil {
		return []domain.FeatureFlag{}, err
	}

	return flags, nil
}

func (repo FlagRepo) GetById(ctx context.Context, id int) (domain.FeatureFlag, error) {
	flag := domain.FeatureFlag{}
	err := repo.db.GetContext(ctx, &flag,
		"SELECT id, name, description, status, environment, owner_user_id, owner_team_id, updated_at FROM feature_flags WHERE id = $1",
		id)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.FeatureFlag{}, apperror.NotFound("флаг не найден")
	}

	if err != nil {
		return domain.FeatureFlag{}, err
	}

	return flag, nil
}

//go:embed queries/feature_flag/create_feature_flag.sql
var createFeatureFlagQuery string

func (repo *FlagRepo) Create(ctx context.Context, featureFlag *domain.FeatureFlag) (int, error) {
	var featureFlagId int
	err := repo.db.QueryRowContext(ctx, createFeatureFlagQuery, featureFlag.Name,
		featureFlag.Description, featureFlag.Status, featureFlag.Environment,
		featureFlag.OwnerUserId, featureFlag.OwnerTeamId).Scan(&featureFlagId)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, apperror.Conflict("флаг с таким именем уже существует")
		}
		return 0, err
	}

	return featureFlagId, nil
}

//go:embed queries/feature_flag/update_feature_flag.sql
var updateFeatureFlagQuery string

func (repo *FlagRepo) UpdateFlagById(ctx context.Context, flagId int, featureFlag *domain.FeatureFlag) error {
	result, err := repo.db.ExecContext(ctx, updateFeatureFlagQuery, featureFlag.Name,
		featureFlag.Description, featureFlag.Status, featureFlag.Environment, flagId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return apperror.Conflict("флаг с таким именем уже существует")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperror.NotFound("фич-флаг не найден")
	}

	return nil
}

//go:embed queries/feature_flag/update_feature_flag_status.sql
var updateFeatureFlagStatusQuery string

func (repo *FlagRepo) UpdateFlagStatusById(ctx context.Context, flagId int, featureFlag domain.FlagStatus) error {
	result, err := repo.db.ExecContext(ctx, updateFeatureFlagStatusQuery, featureFlag, flagId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperror.NotFound("фич-флаг не найден")
	}

	return nil
}
