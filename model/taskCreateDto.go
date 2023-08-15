package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskCreateDto struct {
	TaskName                    string     `json:"task_name"`
	Notes                       *string    `json:"notes,omitempty"`
	StartDate                   time.Time  `json:"start_date"`
	EndDate                     *time.Time `json:"end_date,omitempty"`
	RequiredCompletionsNeeded   *int       `json:"required_completions_needed,omitempty"`
	IntervalBetweenWindowsCount int        `json:"interval_between_windows_count"`
	IntervalBetweenWindowsUnit  string     `json:"interval_between_windows_unit"`
	WindowDurationCount         int        `json:"window_duration_count"`
	WindowDurationUnit          string     `json:"window_duration_unit"`
	TeamId                      uuid.UUID  `json:"team_id"`
	CreatorId                   uuid.UUID  `json:"creator_id"`
	Status                      string     `json:"status"`
}
