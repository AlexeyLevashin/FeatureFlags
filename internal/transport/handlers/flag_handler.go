package handlers

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type FlagService interface {
	GetAll(ctx context.Context, filter domain.FlagFilter) ([]dto.FlagResponse, error)
	GetById(ctx context.Context, id int) (dto.FlagResponse, error)
	Create(ctx context.Context, ownerUserId int, ownerTeamId int, request dto.SaveFlagRequest) (int, error)
	UpdateFlagById(ctx context.Context, flagId int, ownerTeamId int, request dto.SaveFlagRequest) error
	UpdateFlagStatusById(ctx context.Context, flagId int, ownerTeamId int, request dto.UpdateFlagStatusRequest) error
}

type FlagHandler struct {
	service FlagService
}

func NewFlagHandler(service FlagService) *FlagHandler {
	return &FlagHandler{service: service}
}

// CreateFlag Create создает новый фича-флаг
// @Summary Создать новый флаг
// @Description Создает новый фича-флаг и привязывает его к создателю и команде
// @Tags flags
// @Accept json
// @Produce json
// @Param request body dto.SaveFlagRequest true "Данные нового флага"
// @Security ApiKeyAuth
// @Router /flags [post]
func (h *FlagHandler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	var request dto.SaveFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.HandleError(w, apperror.BadRequest("Неверный формат JSON"))
		return
	}

	claims, ok := r.Context().Value(ClaimsKey).(*domain.MyClaims)
	if !ok {
		apperror.HandleError(w, apperror.Unauthorized("Ошибка авторизации: нет данных пользователя"))
		return
	}

	flagId, err := h.service.Create(r.Context(), claims.Id, claims.TeamId, request)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": flagId})
}

// GetAllFlags
// @Summary Получить список всех флагов
// @Description Возвращает массив фича-флагов с возможностью фильтрации
// @Tags flags
// @Produce json
// @Param search query string false "Поиск по имени"
// @Param environment query string false "Фильтр по окружению (dev, staging, production)"
// @Param status query string false "Фильтр по статусу (enabled, disabled)"
// @Security ApiKeyAuth
// @Router /flags [get]
func (h *FlagHandler) GetAllFlags(w http.ResponseWriter, r *http.Request) {
	var flags []dto.FlagResponse
	filter := domain.FlagFilter{
		Search:      r.URL.Query().Get("search"),
		Environment: r.URL.Query().Get("environment"),
		Status:      r.URL.Query().Get("status"),
	}

	flags, err := h.service.GetAll(r.Context(), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(flags)
}

// GetFlagById
// @Summary Получить фич-флаг по id
// @Description Возвращает фич-флаг по его идентификатору
// @Tags flags
// @Produce json
// @Param id path int true "ID флага"
// @Security ApiKeyAuth
// @Router /flags/{id} [get]
func (h *FlagHandler) GetFlagById(w http.ResponseWriter, r *http.Request) {
	var flag dto.FlagResponse
	idStr := r.PathValue("id")
	id, er := strconv.Atoi(idStr)
	if er != nil {
		apperror.HandleError(w, apperror.BadRequest("Ошибка при преобразовании id строки в int"))
		return
	}
	flag, err := h.service.GetById(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(flag)
}

// UpdateFlagById Update обновляет поля фича-флага
// @Summary Редактировать фича-флаг
// @Description Редактирует уже существующий фича-флаг
// @Tags flags
// @Accept json
// @Param id path int true "ID флага"
// @Param request body dto.SaveFlagRequest true "Данные нового флага"
// @Security ApiKeyAuth
// @Router /flags/{id} [put]
func (h *FlagHandler) UpdateFlagById(w http.ResponseWriter, r *http.Request) {
	var request dto.SaveFlagRequest
	idStr := r.PathValue("id")
	id, er := strconv.Atoi(idStr)
	if er != nil {
		apperror.HandleError(w, apperror.BadRequest("Ошибка при преобразовании id строки в int"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.HandleError(w, apperror.BadRequest("Неверный формат JSON"))
		return
	}

	claims, ok := r.Context().Value(ClaimsKey).(*domain.MyClaims)
	if !ok {
		apperror.HandleError(w, apperror.Unauthorized("Ошибка авторизации: нет данных пользователя"))
		return
	}

	err := h.service.UpdateFlagById(r.Context(), id, claims.TeamId, request)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateFlagStatusById UpdateFlagStatus обновляет статус фича-флага
// @Summary Изменить статус фича-флаг
// @Description Редактирует уже существующий фича-флаг
// @Tags flags
// @Accept json
// @Param id path int true "ID флага"
// @Param request body dto.UpdateFlagStatusRequest true "Новый статус"
// @Security ApiKeyAuth
// @Router /flags/{id}/status [patch]
func (h *FlagHandler) UpdateFlagStatusById(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateFlagStatusRequest
	idStr := r.PathValue("id")
	id, er := strconv.Atoi(idStr)
	if er != nil {
		apperror.HandleError(w, apperror.BadRequest("Ошибка при преобразовании id строки в int"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.HandleError(w, apperror.BadRequest("Неверный формат JSON"))
		return
	}

	claims, ok := r.Context().Value(ClaimsKey).(*domain.MyClaims)
	if !ok {
		apperror.HandleError(w, apperror.Unauthorized("Ошибка авторизации: нет данных пользователя"))
		return
	}

	err := h.service.UpdateFlagStatusById(r.Context(), id, claims.TeamId, request)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
