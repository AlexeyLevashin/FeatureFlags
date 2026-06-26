package domain

type User struct {
	Id           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	Name         string `json:"name" db:"name"`
	Surname      string `json:"surname" db:"surname"`
	PasswordHash string `json:"-" db:"password_hash"`
	TeamId       int    `json:"team_id" db:"team_id"`
}
