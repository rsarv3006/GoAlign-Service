package model

import "github.com/google/uuid"

type User struct {
	UserId          uuid.UUID `json:"user_id"`
	UserName        string    `json:"user_name"`
	Email           string    `json:"email"`
	IsActive        bool      `json:"is_active"`
	IsEmailVerified bool      `json:"is_email_verified"`
}

type UserList struct {
	Users []User `json:"users"`
}
