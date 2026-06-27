package service

import (
	"FeatureFlags/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	FindByEmail(email string) (domain.User, error)
}

type AuthService struct {
	repo      AuthRepository
	jwtSecret string
}

func NewAuthService(repo AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(email string, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "Неверный email или пароль!", err
	}
	err = checkPassword(user, password)
	if err != nil {
		return "Неверный email или пароль!", err
	}
	token, err := generateToken(user, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateToken(user domain.User, secret string) (string, error) {
	myClaims := domain.MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		Id:     user.Id,
		Email:  user.Email,
		TeamId: user.TeamId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func checkPassword(user domain.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return errors.New("Неверный пароль")
	}
	return nil
}
