package model

import "github.com/google/uuid"

type TeamSettings struct {
	TeamSettingsId        uuid.UUID `json:"team_settings_id"`
	TeamId                uuid.UUID `json:"team_id"`
	CanAllMembersAddTasks bool      `json:"can_all_members_add_tasks"`
}

type TeamSettingsCreateDto struct {
	TeamId                uuid.UUID `json:"team_id"`
	CanAllMembersAddTasks bool      `json:"can_all_members_add_tasks"`
}
