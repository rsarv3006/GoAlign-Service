package model

import "github.com/google/uuid"

type TeamInviteCreateDto struct {
	TeamId uuid.UUID `json:"team_id"`
	Email  string    `json:"email"`
}
