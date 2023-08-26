package model

type IntervalObj struct {
	IntervalCount int             `json:"interval_count"`
	IntervalUnit  IntervalVariant `json:"interval_unit"`
}
