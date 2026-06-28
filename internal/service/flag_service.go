package service

import (
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"context"
	"errors"
	"fmt"
	"time"
)

type FlagRepository interface {
	GetAll(filter domain.FlagFilter) ([]domain.FeatureFlag, error)
	GetById(id int) (domain.FeatureFlag, error)
	Create(ctx context.Context, featureFlag *domain.FeatureFlag) (int, error)
}

type UserRepository interface {
	CheckExists(ctx context.Context, userId int) (bool, error)
}

type TeamRepository interface {
	CheckExists(ctx context.Context, teamId int) (bool, error)
}

type FlagService struct {
	flagRepo FlagRepository
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewFlagService(f FlagRepository, u UserRepository, t TeamRepository) *FlagService {
	return &FlagService{
		flagRepo: f,
		userRepo: u,
		teamRepo: t,
	}
}

func (f FlagService) GetAll(filter domain.FlagFilter) ([]dto.FlagResponse, error) {
	flags, err := f.flagRepo.GetAll(filter)
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
	flag, err := f.flagRepo.GetById(id)
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

func (f *FlagService) Create(ctx context.Context, request dto.CreateFlagRequest, ownerUserId int, ownerTeamId int) (int, error) {
	if request.Status != domain.StatusEnabled && request.Status != domain.StatusDisabled {
		return 0, errors.New("недопустимый статус флага")
	}

	if request.Environment != domain.EnvDev && request.Environment != domain.EnvStag && request.Environment != domain.EnvProd {
		return 0, errors.New("недопустимое окружение флага")
	}

	checkUserExists, err := f.userRepo.CheckExists(ctx, ownerUserId)
	if err != nil || !checkUserExists {
		return 0, errors.New("пользователь не найден")
	}

	checkTeamExists, err := f.teamRepo.CheckExists(ctx, ownerTeamId)
	if err != nil || !checkTeamExists {
		return 0, errors.New("команда не найдена")
	}

	flagDb := createFlagRequestToDomain(request)
	flagDb.OwnerUserId = ownerUserId
	flagDb.OwnerTeamId = ownerTeamId
	flagId, err := f.flagRepo.Create(ctx, flagDb)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания флага: %w", err)
	}

	return flagId, nil
}

func createFlagRequestToDomain(request dto.CreateFlagRequest) *domain.FeatureFlag {
	return &domain.FeatureFlag{
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
		Environment: request.Environment,
	}
}
