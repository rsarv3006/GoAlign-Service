package model

import (
	"github.com/google/uuid"
	"time"
)

type TaskEntryCreateDto struct {
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Notes          string    `json:"notes"`
	AssignedUserId uuid.UUID `json:"assigned_user_id"`
	TaskId         uuid.UUID `json:"task_id"`
}
