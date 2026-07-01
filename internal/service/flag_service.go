package service

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"context"
	"time"
)

type FlagRepository interface {
	GetAll(ctx context.Context, filter domain.FlagFilter) ([]domain.FeatureFlagDetails, error)
	GetById(ctx context.Context, id int) (domain.FeatureFlag, error)
	GetFlagDetailsById(ctx context.Context, id int) (domain.FeatureFlagDetails, error)
	Create(ctx context.Context, featureFlag *domain.FeatureFlag) (int, error)
	UpdateFlagById(ctx context.Context, flagId int, userId int, featureFlag *domain.FeatureFlag) error
	UpdateFlagStatusById(ctx context.Context, flagId int, userId int, featureFlag domain.FlagStatus) error
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

func (f FlagService) GetAll(ctx context.Context, filter domain.FlagFilter) ([]dto.FlagResponse, error) {
	flags, err := f.flagRepo.GetAll(ctx, filter)
	if err != nil {
		return []dto.FlagResponse{}, err
	}

	result := make([]dto.FlagResponse, len(flags))
	for i, flag := range flags {
		result[i] = toFlagResponse(flag)
	}

	return result, nil
}

func (f FlagService) GetFlagDetailsById(ctx context.Context, id int) (dto.FlagResponse, error) {
	flag, err := f.flagRepo.GetFlagDetailsById(ctx, id)
	if err != nil {
		return dto.FlagResponse{}, err
	}
	return toFlagResponse(flag), nil
}

func toFlagResponse(flag domain.FeatureFlagDetails) dto.FlagResponse {
	var updatedBy *string
	var updatedAt *string

	if flag.UpdaterName.Valid {
		name := flag.UpdaterName.String + " " + flag.UpdaterSurname.String
		timeStr := flag.UpdatedAt.Time.Format(time.RFC3339)

		updatedBy = &name
		updatedAt = &timeStr
	}

	return dto.FlagResponse{
		Id:          flag.Id,
		Name:        flag.Name,
		Description: flag.Description,
		Status:      string(flag.Status),
		Environment: string(flag.Environment),

		Owner:     flag.OwnerTeam,
		CreatedBy: flag.CreatorName + " " + flag.CreatorSurname,
		CreatedAt: flag.CreatedAt.Format(time.RFC3339),

		UpdatedBy: updatedBy,
		UpdatedAt: updatedAt,
	}
}

func (f *FlagService) Create(ctx context.Context, ownerUserId int,
	ownerTeamId int, request dto.SaveFlagRequest) (int, error) {
	if !request.Status.IsValid() {
		return 0, apperror.BadRequest("недопустимый статус флага")
	}

	if !request.Environment.IsValid() {
		return 0, apperror.BadRequest("недопустимое окружение флага")
	}

	checkUserExists, err := f.userRepo.CheckExists(ctx, ownerUserId)
	if err != nil {
		return 0, err
	}
	if !checkUserExists {
		return 0, apperror.NotFound("пользователь не найден")
	}

	checkTeamExists, err := f.teamRepo.CheckExists(ctx, ownerTeamId)
	if err != nil {
		return 0, err
	}
	if !checkTeamExists {
		return 0, apperror.NotFound("команда не найдена")
	}

	flagDb := saveFlagRequestToDomain(request)
	flagDb.OwnerUserId = ownerUserId
	flagDb.OwnerTeamId = ownerTeamId
	flagId, err := f.flagRepo.Create(ctx, flagDb)
	if err != nil {
		return 0, err
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

func (f *FlagService) UpdateFlagById(ctx context.Context, flagId int, userId,
	ownerTeamId int, request dto.SaveFlagRequest) error {
	if !request.Status.IsValid() {
		return apperror.BadRequest("недопустимый статус флага")
	}

	if !request.Environment.IsValid() {
		return apperror.BadRequest("недопустимое окружение флага")
	}

	existingFlag, err := f.flagRepo.GetById(ctx, flagId)
	if err != nil {
		return err
	}

	if existingFlag.OwnerTeamId != ownerTeamId {
		return apperror.Forbidden("редактирование флагов других команд запрещено")
	}

	flagDb := saveFlagRequestToDomain(request)

	err = f.flagRepo.UpdateFlagById(ctx, flagId, userId, flagDb)
	if err != nil {
		return err
	}

	return nil
}

func (f *FlagService) UpdateFlagStatusById(ctx context.Context, flagId int, userId int,
	ownerTeamId int, request dto.UpdateFlagStatusRequest) error {
	if !request.Status.IsValid() {
		return apperror.BadRequest("недопустимый статус фич флага")
	}

	existingFlag, err := f.flagRepo.GetById(ctx, flagId)
	if err != nil {
		return err
	}

	if existingFlag.OwnerTeamId != ownerTeamId {
		return apperror.Forbidden("редактирование статуса фич флага других команд запрещено")
	}

	err = f.flagRepo.UpdateFlagStatusById(ctx, flagId, userId, request.Status)
	if err != nil {
		return err
	}

	return nil
}
