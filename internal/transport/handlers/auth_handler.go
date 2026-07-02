package handlers

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"FeatureFlags/internal/transport/httputil"
	"FeatureFlags/internal/transport/middleware"
	"context"
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
// @Param request body dto.LoginRequest true "Данные для входа"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := httputil.ReadAndValidate[dto.LoginRequest](r)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	tokenString, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.LoginResponse{Token: tokenString})
}

// GetMe
// @Summary Получить данные текущего пользователя
// @Description Возвращает информацию о пользователе по токену из заголовка Authorization
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Router /me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*domain.MyClaims)
	if !ok {
		httputil.HandleError(w, apperror.Unauthorized("unauthorized"))
		return
	}

	userId := claims.Id
	user, err := h.authService.GetMe(r.Context(), userId)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, user)
}
