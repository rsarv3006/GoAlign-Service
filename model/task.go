package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	TaskId                    uuid.UUID   `json:"task_id"`
	TaskName                  string      `json:"task_name"`
	Notes                     string      `json:"notes"`
	StartDate                 time.Time   `json:"start_date"`
	EndDate                   *time.Time  `json:"end_date,omitempty"`
	RequiredCompletionsNeeded int         `json:"required_completions_needed"`
	CompletionCount           int         `json:"completion_count"`
	IntervalBetweenWindows    IntervalObj `json:"interval_between_windows"`
	WindowDuration            IntervalObj `json:"window_duration"`
	TeamId                    uuid.UUID   `json:"team_id"`
	CreatorId                 uuid.UUID   `json:"creator_id"`
	CreatedAt                 time.Time   `json:"created_at"`
	UpdatedAt                 time.Time   `json:"updated_at"`
	Status                    string      `json:"status"`
}
