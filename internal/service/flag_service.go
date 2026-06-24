package service

import (
	"FeatureFlags/internal/domain"
	"context"
)

type FlagService interface {
	CreateFlag(ctx context.Context, flag *domain.FeatureFlag) error
}
