package handlers

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/domain"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if token == "" {
				apperror.HandleError(w, apperror.Unauthorized("Токен не передан"))
				return
			}

			parsedToken, err := jwt.ParseWithClaims(token, &domain.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !parsedToken.Valid {
				apperror.HandleError(w, apperror.Unauthorized("Невалидный токен"))
				return
			}

			claims, ok := parsedToken.Claims.(*domain.MyClaims)
			if !ok {
				apperror.HandleError(w, apperror.Unauthorized("Невалидный токен"))
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
