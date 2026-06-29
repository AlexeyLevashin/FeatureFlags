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
	UpdateFlagById(ctx context.Context, flagId int, featureFlag *domain.FeatureFlag) error
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

func (f *FlagService) Create(ctx context.Context, ownerUserId int, ownerTeamId int, request dto.SaveFlagRequest) (int, error) {
	if !request.Status.IsValid() {
		return 0, errors.New("недопустимый статус флага")
	}

	if !request.Environment.IsValid() {
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

	flagDb := saveFlagRequestToDomain(request)
	flagDb.OwnerUserId = ownerUserId
	flagDb.OwnerTeamId = ownerTeamId
	flagId, err := f.flagRepo.Create(ctx, flagDb)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания флага: %w", err)
	}

	return flagId, nil
}

func saveFlagRequestToDomain(request dto.SaveFlagRequest) *domain.FeatureFlag {
	return &domain.FeatureFlag{
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
		Environment: request.Environment,
	}
}

func (f *FlagService) UpdateFlagById(ctx context.Context, flagId int, ownerUserId int, ownerTeamId int, request dto.SaveFlagRequest) error {
	if !request.Status.IsValid() {
		return errors.New("недопустимый статус флага")
	}

	if !request.Environment.IsValid() {
		return errors.New("недопустимое окружение флага")
	}

	existingFlag, err := f.flagRepo.GetById(flagId)
	if err != nil {
		return errors.New("флаг не найден")
	}

	if existingFlag.OwnerTeamId != ownerTeamId {
		return errors.New("редактирование флагов других команд запрещено")
	}

	checkUserExists, err := f.userRepo.CheckExists(ctx, ownerUserId)
	if err != nil || !checkUserExists {
		return errors.New("пользователь не найден")
	}

	checkTeamExists, err := f.teamRepo.CheckExists(ctx, ownerTeamId)
	if err != nil || !checkTeamExists {
		return errors.New("команда не найдена")
	}

	flagDb := saveFlagRequestToDomain(request)

	err = f.flagRepo.UpdateFlagById(ctx, flagId, flagDb)
	if err != nil {
		return fmt.Errorf("ошибка редактирования флага: %w", err)
	}

	return nil
}
