package domain

type User struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
