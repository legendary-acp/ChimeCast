package models

type RegisterRequest struct {
	Name     string `json:"name"`
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}
