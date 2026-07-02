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

//go:embed queries/feature_flag/get_all_flags_base.sql
var getAllFlagsQuery string

func (repo *FlagRepo) GetAll(ctx context.Context, filter domain.FlagFilter) ([]domain.FeatureFlagDetails, error) {
	flags := []domain.FeatureFlagDetails{}

	query := getAllFlagsQuery
	args := []interface{}{}
	i := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND f.name ILIKE $%d", i)
		args = append(args, "%"+filter.Search+"%")
		i++
	}

	if filter.Environment != "" {
		query += fmt.Sprintf(" AND f.environment = $%d", i)
		args = append(args, filter.Environment)
		i++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND f.status = $%d", i)
		args = append(args, filter.Status)
		i++
	}

	err := repo.db.SelectContext(ctx, &flags, query, args...)

	if err != nil {
		return []domain.FeatureFlagDetails{}, err
	}

	return flags, nil
}

//go:embed queries/feature_flag/get_feature_flag_details_by_id.sql
var getFlagDetailsByIdQuery string

func (repo *FlagRepo) GetFlagDetailsById(ctx context.Context, id int) (domain.FeatureFlagDetails, error) {
	flag := domain.FeatureFlagDetails{}
	err := repo.db.GetContext(ctx, &flag, getFlagDetailsByIdQuery, id)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.FeatureFlagDetails{}, apperror.NotFound("флаг не найден")
	}

	if err != nil {
		return domain.FeatureFlagDetails{}, err
	}

	return flag, nil
}

//go:embed queries/feature_flag/get_flag_by_id.sql
var getFlagByIdQuery string

func (repo *FlagRepo) GetById(ctx context.Context, id int) (domain.FeatureFlag, error) {
	flag := domain.FeatureFlag{}
	err := repo.db.GetContext(ctx, &flag, getFlagByIdQuery, id)

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

//go:embed queries/feature_flag_updates/insert_flag_update.sql
var insertFlagUpdateLogQuery string

func (repo *FlagRepo) UpdateFlagById(ctx context.Context, flagId int, userId int, featureFlag *domain.FeatureFlag) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateFeatureFlagQuery, featureFlag.Name,
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

	_, err = tx.ExecContext(ctx,
		insertFlagUpdateLogQuery,
		userId, flagId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

//go:embed queries/feature_flag/update_feature_flag_status.sql
var updateFeatureFlagStatusQuery string

func (repo *FlagRepo) UpdateFlagStatusById(ctx context.Context, flagId int, userId int, featureFlagStatus domain.FlagStatus) error {

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateFeatureFlagStatusQuery, featureFlagStatus, flagId)
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

	_, err = tx.ExecContext(ctx,
		insertFlagUpdateLogQuery,
		userId, flagId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
