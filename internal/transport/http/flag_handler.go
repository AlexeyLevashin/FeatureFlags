package handlers

import (
	"FeatureFlags/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"
)

type FlagService interface {
	GetAll(filter domain.FlagFilter) ([]domain.FeatureFlag, error)
	GetById(id int) (domain.FeatureFlag, error)
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

func (h *FlagHandler) GetFlagById(w http.ResponseWriter, r *http.Request) {
	var flag domain.FeatureFlag
	idStr := r.PathValue("id")
	id, er := strconv.Atoi(idStr)
	if er != nil {
		http.Error(w,
			"Ошибка при преобразовании id строки в int",
			http.StatusInternalServerError)
	}
	flag, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, "Ошибка при получении флага по id", http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(flag)
}
