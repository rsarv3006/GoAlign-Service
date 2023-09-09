package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskCreateDto struct {
	TaskName                  string      `json:"task_name"`
	Notes                     *string     `json:"notes,omitempty"`
	StartDate                 time.Time   `json:"start_date"`
	EndDate                   *time.Time  `json:"end_date,omitempty"`
	RequiredCompletionsNeeded *int        `json:"required_completions_needed,omitempty"`
	IntervalBetweenWindows    IntervalObj `json:"interval_between_windows"`
	WindowDuration            IntervalObj `json:"window_duration"`
	TeamId                    uuid.UUID   `json:"team_id"`
	Status                    string      `json:"status"`
	AssignedUserId            uuid.UUID   `json:"assigned_user_id"`
}
