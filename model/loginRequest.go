package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	LoginRequestId             uuid.UUID `json:"login_request_id"`
	UserId                     uuid.UUID `json:"user_id"`
	LoginRequestDate           time.Time `json:"login_request_date"`
	LoginRequestExpirationDate time.Time `json:"login_request_expiration_date"`
	LoginRequestToken          string    `json:"login_request_token"`
	LoginRequestStatus         string    `json:"login_request_status"`
}
