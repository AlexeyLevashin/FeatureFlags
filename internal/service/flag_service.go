package service

import (
	"FeatureFlags/internal/domain"
)

type FlagRepository interface {
	GetAll(filter domain.FlagFilter) ([]domain.FeatureFlag, error)
	GetById(id int) (domain.FeatureFlag, error)
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

func (f FlagService) GetById(id int) (domain.FeatureFlag, error) {
	flag, err := f.repo.GetById(id)
	if err != nil {
		return domain.FeatureFlag{}, err
	}
	return flag, nil
}

/*type FlagService interface {
	CreateFlag(ctx context.Context, flag *domain.FeatureFlag) error
}
*/
