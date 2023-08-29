package model

import "github.com/google/uuid"

type LoginCodeRequestDto struct {
	LoginCodeRequestId uuid.UUID `json:"login_code_request_id"`
	UserId             uuid.UUID `json:"user_id"`
	LoginRequestToken  string    `json:"login_request_token"`
}
