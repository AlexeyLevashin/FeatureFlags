package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GetMeResponse struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	TeamId  int    `json:"teamId"`
}
