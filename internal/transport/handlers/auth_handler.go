package handlers

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AuthService interface {
	Login(email string, password string) (string, error)
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
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	tokenString, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
