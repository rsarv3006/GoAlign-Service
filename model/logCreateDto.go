package model

type LogCreateDto struct {
	LogMessage string  `json:"log_message"`
	LogLevel   string  `json:"log_level"`
	LogData    *string `json:"log_data,omitempty"`
}
