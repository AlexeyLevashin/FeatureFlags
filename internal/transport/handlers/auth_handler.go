package handlers

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"context"
	"encoding/json"
	"net/http"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	GetMe(ctx context.Context, id int) (dto.GetMeResponse, error)
}
type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

// Login авторизует пользователя
// @Summary Вход в систему
// @Description Возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("Неверный формат JSON"))
		return
	}

	tokenString, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.LoginResponse{Token: tokenString})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ClaimsKey).(*domain.MyClaims)
	if !ok {
		apperror.HandleError(w, apperror.Unauthorized("unauthorized"))
		return
	}
	userId := claims.Id
	user, err := h.authService.GetMe(r.Context(), userId)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
