package dto

import "FeatureFlags/internal/domain"

type FlagResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Environment string  `json:"environment"`
	Owner       string  `json:"owner"`
	CreatedBy   string  `json:"createdBy"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedBy   *string `json:"updatedBy"`
	UpdatedAt   *string `json:"updatedAt"`
}

type SaveFlagRequest struct {
	Name        string                 `json:"name" validate:"required,min=3,max=100"`
	Description string                 `json:"description" validate:"required,max=500"`
	Status      domain.FlagStatus      `json:"status" validate:"required"`
	Environment domain.EnvironmentType `json:"environment" validate:"required"`
}

type UpdateFlagStatusRequest struct {
	Status domain.FlagStatus `json:"status" validate:"required"`
}
