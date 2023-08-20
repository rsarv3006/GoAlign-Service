package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamInvite struct {
	TeamInviteId    uuid.UUID `json:"team_invite_id"`
	TeamId          uuid.UUID `json:"team_id"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Status          string    `json:"status"`
	InviteCreatorId uuid.UUID `json:"invite_creator_id"`
}
