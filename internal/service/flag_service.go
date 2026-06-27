package service

import (
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"time"
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

func (f FlagService) GetAll(filter domain.FlagFilter) ([]dto.FlagResponse, error) {
	flags, err := f.repo.GetAll(filter)
	if err != nil {
		return []dto.FlagResponse{}, err
	}
	result := make([]dto.FlagResponse, len(flags))
	for i, flag := range flags {
		result[i] = toFlagResponse(flag)
	}
	return result, nil
}

func (f FlagService) GetById(id int) (dto.FlagResponse, error) {
	flag, err := f.repo.GetById(id)
	if err != nil {
		return dto.FlagResponse{}, err
	} 
	return toFlagResponse(flag), nil
}

func toFlagResponse(flag domain.FeatureFlag) dto.FlagResponse {
	return dto.FlagResponse{
		Id:          flag.Id,
		Name:        flag.Name,
		Description: flag.Description,
		Status:      string(flag.Status),
		Environment: string(flag.Environment),
		OwnerUserId: flag.OwnerUserId,
		OwnerTeamId: flag.OwnerTeamId,
		UpdatedAt:   flag.UpdatedAt.Format(time.RFC3339),
	}
}
/*type FlagService interface {
	CreateFlag(ctx context.Context, flag *domain.FeatureFlag) error
}
*/
