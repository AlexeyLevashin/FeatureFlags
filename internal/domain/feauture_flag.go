package domain

import (
	"database/sql"
	"time"
)

type FlagStatus string
type EnvironmentType string

const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"

	EnvDev  = "dev"
	EnvStag = "staging"
	EnvProd = "production"
)

func (s FlagStatus) IsValid() bool {
	switch s {
	case StatusEnabled, StatusDisabled:
		return true
	}
	return false
}

func (e EnvironmentType) IsValid() bool {
	switch e {
	case EnvDev, EnvStag, EnvProd:
		return true
	}
	return false
}

type FeatureFlag struct {
	Id          int             `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	Status      FlagStatus      `json:"status" db:"status"`
	Environment EnvironmentType `json:"environment" db:"environment"`
	OwnerUserId int             `json:"ownerUserId" db:"owner_user_id"`
	OwnerTeamId int             `json:"ownerTeamId" db:"owner_team_id"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
}

type FeatureFlagDetails struct {
	Id             int             `db:"id"`
	Name           string          `db:"name"`
	Description    string          `db:"description"`
	Environment    EnvironmentType `db:"environment"`
	Status         FlagStatus      `db:"status"`
	OwnerTeam      string          `db:"owner_team"`
	CreatorName    string          `db:"creator_name"`
	CreatorSurname string          `db:"creator_surname"`
	CreatedAt      time.Time       `db:"created_at"`

	UpdaterName    sql.NullString `db:"updater_name"`
	UpdaterSurname sql.NullString `db:"updater_surname"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
}
