package dto

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id uint `json:"id"`
}
