package handlers

import (
	"FeatureFlags/internal/domain"
	"encoding/json"
	"net/http"
)

type FlagService interface {
	GetAll(filter domain.FlagFilter) ([]domain.FeatureFlag, error)
}

type FlagHandler struct {
	service FlagService
}

func NewFlagHandler(service FlagService) *FlagHandler {
	return &FlagHandler{service: service}
}

func (h *FlagHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *FlagHandler) GetAllFlags(w http.ResponseWriter, r *http.Request) {
	var flags []domain.FeatureFlag
	filter := domain.FlagFilter{
		Search:      r.URL.Query().Get("search"),
		Environment: r.URL.Query().Get("environment"),
		Status:      r.URL.Query().Get("status"),
	}
	flags, err := h.service.GetAll(filter)
	if err != nil {
		http.Error(w, "Ошибка при получении списка флагов", http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(flags)
}
