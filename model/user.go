package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId          uuid.UUID `json:"user_id"`
	UserName        string    `json:"username"`
	Email           string    `json:"email"`
	IsActive        bool      `json:"is_active"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserList struct {
	Users []User `json:"users"`
}
