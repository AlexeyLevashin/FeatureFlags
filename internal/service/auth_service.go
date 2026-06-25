package service

import (
	"FeatureFlags/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRepository interface {
	FindByEmail(email string) (domain.User, error)
}

type MyClaims struct {
	jwt.RegisteredClaims
	Id    string `json:"user_id"`
	Email string `json:"email"`
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
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
	token, err := generateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateToken(user domain.User) (string, error) {
	myClaims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		Id:    user.Id,
		Email: user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func checkPassword(user domain.User, password string) error {
	// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	if user.Password != password {
		return errors.New("Неверный пароль")
	}
	return nil
}
