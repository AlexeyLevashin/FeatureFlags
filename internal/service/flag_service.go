package service

import (
	"FeatureFlags/internal/domain"
)

type FlagRepository interface {
	GetAll(filter domain.FlagFilter) ([]domain.FeatureFlag, error)
}

type FlagService struct {
	repo FlagRepository
}

func NewFlagService(repo FlagRepository) *FlagService {
	return &FlagService{repo: repo}
}

func (f FlagService) GetAll(filter domain.FlagFilter) ([]domain.FeatureFlag, error) {
	flags, err := f.repo.GetAll(filter)
	if err != nil {
		return []domain.FeatureFlag{}, err
	}
	return flags, nil
}

/*type FlagService interface {
	CreateFlag(ctx context.Context, flag *domain.FeatureFlag) error
}
*/
