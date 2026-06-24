package domain

type User struct {
	Id    string `json:"user_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
