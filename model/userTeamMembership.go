package model

import (
	"time"

	"github.com/google/uuid"
)

type UserTeamMembership struct {
	UserTeamMembershipId uuid.UUID `json:"user_team_membership_id"`
	UserId               uuid.UUID `json:"user_id"`
	TeamId               uuid.UUID `json:"team_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	Status               string    `json:"status"`
}

type UserTeamMemberships struct {
	UserTeamMemberships []UserTeamMembership `json:"user_team_memberships"`
}
