package handlers

import (
	"FeatureFlags/internal/domain"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func GenerateToken(user domain.User) (string, error) {
	myClaims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		UserId: user.UserId,
		Email:  user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := domain.User{
		UserId: "1",
		Email:  "test@test.com",
		Name:   "Test User",
	}
	tokenString, err := GenerateToken(user)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
