package model

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	TeamId        uuid.UUID `json:"team_id"`
	TeamName      string    `json:"team_name"`
	CreatorUserId uuid.UUID `json:"creator_user_id"`
	Status        string    `json:"status"`
	TeamManagerId uuid.UUID `json:"team_manager_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Teams struct {
	Teams []Team `json:"teams"`
}
