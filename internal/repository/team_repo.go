package repository

import (
	"context"
	_ "embed"

	"github.com/jmoiron/sqlx"
)

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository { return &TeamRepository{db: db} }

//go:embed queries/team/check_team_exists.sql
var CheckTeamExists string

func (repo *TeamRepository) CheckExists(ctx context.Context, teamId int) (bool, error) {
	var exists bool

	err := repo.db.GetContext(ctx, &exists, CheckTeamExists, teamId)
	if err != nil {
		return false, err
	}

	return exists, nil
}
