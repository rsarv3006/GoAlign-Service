package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskEntry struct {
	TaskEntryId    uuid.UUID  `json:"task_entry_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	Notes          string     `json:"notes"`
	AssignedUserId uuid.UUID  `json:"assigned_user_id"`
	Status         string     `json:"status"`
	CompletedDate  *time.Time `json:"completed_date,omitempty"`
	TaskId         uuid.UUID  `json:"task_id"`
}
