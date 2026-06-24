package domain

import "time"

const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"

	EnvDev  = "dev"
	EnvStag = "staging"
	EnvProd = "production"
)

type FeatureFlag struct {
	Id          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	Environment string    `json:"environment" db:"environment"`
	OwnerUserID *string   `json:"ownerUserId,omitempty" db:"owner_user_id"`
	OwnerTeamID *string   `json:"ownerTeamId,omitempty" db:"owner_team_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
