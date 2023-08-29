package model

type UserCreateDto struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}
