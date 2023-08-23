package model

import "github.com/google/uuid"

type LogModel struct {
	AppLogId   uuid.UUID  `json:"app_log_id"`
	LogMessage string     `json:"log_message"`
	LogLevel   string     `json:"log_level"`
	LogDate    string     `json:"log_date"`
	LogData    *string    `json:"log_data,omitempty"`
	UserId     *uuid.UUID `json:"user_id,omitempty"`
}
