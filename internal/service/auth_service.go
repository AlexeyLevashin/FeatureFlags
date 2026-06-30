package service

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"FeatureFlags/internal/dto"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int) (domain.User, error)
}

type AuthService struct {
	repo      AuthRepository
	jwtSecret string
}

func NewAuthService(repo AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	err = checkPassword(user, password)
	if err != nil {
		return "", err
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
		return apperror.Unauthorized("Неверный email или пароль")
	}
	return nil
}

func (s *AuthService) GetMe(ctx context.Context, id int) (dto.GetMeResponse, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return dto.GetMeResponse{}, err
	}
	userResponse := toUserResponse(user)
	return userResponse, nil
}

func toUserResponse(user domain.User) dto.GetMeResponse {
	return dto.GetMeResponse{
		Email:   user.Email,
		Name:    user.Name,
		Surname: user.Surname,
		TeamId:  user.TeamId,
	}
}
