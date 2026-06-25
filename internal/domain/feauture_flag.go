package domain

import "time"

type FlagStatus string
type EnvironmentType string

const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"

	EnvDev  = "dev"
	EnvStag = "staging"
	EnvProd = "production"
)

type FeatureFlag struct {
	Id          int             `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	Status      FlagStatus      `json:"status" db:"status"`
	Environment EnvironmentType `json:"environment" db:"environment"`
	OwnerUserId int             `json:"ownerUserId,omitempty" db:"owner_user_id"`
	OwnerTeamId int             `json:"ownerTeamId,omitempty" db:"owner_team_id"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`
}
