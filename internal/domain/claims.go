package domain

import "github.com/golang-jwt/jwt/v5"

type MyClaims struct {
	jwt.RegisteredClaims
	Id     int    `json:"id"`
	Email  string `json:"email"`
	TeamId int    `json:"team_id"`
}
