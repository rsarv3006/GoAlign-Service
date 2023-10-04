package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskUpdateDto struct {
	TaskId                      uuid.UUID        `json:"task_id"`
	TaskName                    *string          `json:"task_name,omitempty"`
	Notes                       *string          `json:"notes,omitempty"`
	StartDate                   *time.Time       `json:"start_date,omitempty"`
	EndDate                     *time.Time       `json:"end_date,omitempty"`
	RequiredCompletionsNeeded   *int             `json:"required_completions_needed,omitempty"`
	IntervalBetweenWindowsCount *int             `json:"interval_between_windows_count,omitempty"`
	IntervalBetweenWindowsUnit  *IntervalVariant `json:"interval_between_windows_unit,omitempty"`
	WindowDurationCount         *int             `json:"window_duration_count,omitempty"`
	WindowDurationUnit          *IntervalVariant `json:"window_duration_unit,omitempty"`
	AssignedUserId              *uuid.UUID       `json:"assigned_user_id,omitempty"`
}
