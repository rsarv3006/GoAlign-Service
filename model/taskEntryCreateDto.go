package model

import "github.com/google/uuid"

type TaskEntryCreateDto struct {
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	Notes          string    `json:"notes"`
	AssignedUserId uuid.UUID `json:"assigned_user_id"`
	Status         string    `json:"status"`
	TaskId         uuid.UUID `json:"task_id"`
}
