package domain

import "github.com/golang-jwt/jwt/v5"

type MyClaims struct {
	jwt.RegisteredClaims
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
