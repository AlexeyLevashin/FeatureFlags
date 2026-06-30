package domain

import (
	"context"
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
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`
}

type FlagRepository interface {
	Create(ctx context.Context, flag *FeatureFlag) error
	GetByID(ctx context.Context, id string) (*FeatureFlag, error)
}
