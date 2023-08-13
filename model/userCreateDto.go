package model

type UserCreateDto struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
