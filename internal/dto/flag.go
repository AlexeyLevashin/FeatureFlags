package dto

type FlagResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Environment string `json:"environment"`
	OwnerUserId int    `json:"ownerUserId"`
	OwnerTeamId int    `json:"ownerTeamId"`
	UpdatedAt   string `json:"updatedAt"`
}
