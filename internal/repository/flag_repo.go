package repository

import (
	"FeatureFlags/internal/domain"
	"context"
)

type FlagRepository interface {
	Create(ctx context.Context, flag *domain.FeatureFlag) error
	GetByID(ctx context.Context, id string) (*domain.FeatureFlag, error)
}
