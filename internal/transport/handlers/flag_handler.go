package handlers

import (
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type FlagService interface {
	GetAll(filter domain.FlagFilter) ([]dto.FlagResponse, error)
	GetById(id int) (dto.FlagResponse, error)
	Create(ctx context.Context, ownerUserId int, ownerTeamId int, request dto.SaveFlagRequest) (int, error)
	UpdateFlagById(ctx context.Context, flagId int, ownerUserId int, ownerTeamId int, request dto.SaveFlagRequest) error
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
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(ClaimsKey).(*domain.MyClaims)
	if !ok {
		http.Error(w, "Ошибка авторизации: нет данных пользователя", http.StatusUnauthorized)
		return
	}

	flagId, err := h.service.Create(r.Context(), claims.Id, claims.TeamId, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
// @Router /flags [get]
func (h *FlagHandler) GetAllFlags(w http.ResponseWriter, r *http.Request) {
	var flags []dto.FlagResponse
	filter := domain.FlagFilter{
		Search:      r.URL.Query().Get("search"),
		Environment: r.URL.Query().Get("environment"),
		Status:      r.URL.Query().Get("status"),
	}
	flags, err := h.service.GetAll(filter)
	if err != nil {
		http.Error(w, "Ошибка при получении списка флагов", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(flags)
}

func (h *FlagHandler) GetFlagById(w http.ResponseWriter, r *http.Request) {
	var flag dto.FlagResponse
	idStr := r.PathValue("id")
	id, er := strconv.Atoi(idStr)
	if er != nil {
		http.Error(w,
			"Ошибка при преобразовании id строки в int",
			http.StatusInternalServerError)
		return
	}
	flag, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, "Ошибка при получении флага по id", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(flag)
}

// UpdateFlagById Update обновляет поля фича-флага
// @Summary Редактировать фича-флаг
// @Description Редактирует уже существующий фича-флаг
// @Tags flags
// @Accept json
// @Produce json
// @Param request body dto.SaveFlagRequest true "Данные нового флага"
// @Security ApiKeyAuth
// @Router /flags [post]
func (h *FlagHandler) UpdateFlagById(w http.ResponseWriter, r *http.Request) {
	var request dto.SaveFlagRequest
	idStr := r.PathValue("id")
	id, er := strconv.Atoi(idStr)
	if er != nil {
		http.Error(w, "Ошибка при преобразовании id строки в int", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(ClaimsKey).(*domain.MyClaims)
	if !ok {
		http.Error(w, "Ошибка авторизации: нет данных пользователя", http.StatusUnauthorized)
		return
	}

	err := h.service.UpdateFlagById(r.Context(), id, claims.Id, claims.TeamId, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
