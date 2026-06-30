package handlers

import (
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"encoding/json"
	"net/http"
)

type AuthService interface {
	Login(email string, password string) (string, error)
	GetMe(id int) (dto.GetMeResponse, error)
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
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	tokenString, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(dto.LoginResponse{Token: tokenString})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ClaimsKey).(*domain.MyClaims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userId := claims.Id
	user, err := h.authService.GetMe(userId)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(user)
}
