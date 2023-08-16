package model

import "strings"

type IntervalVariant string

const (
	Minute IntervalVariant = "minute(s)"
	Hour   IntervalVariant = "hour(s)"
	Day    IntervalVariant = "day(s)"
	Week   IntervalVariant = "week(s)"
	Month  IntervalVariant = "month(s)"
	Year   IntervalVariant = "year(s)"
)

func IsValidVariant(s string) bool {
	s = strings.ToLower(s)

	switch s {
	case strings.ToLower(string(Minute)):
		fallthrough
	case strings.ToLower(string(Hour)):
		fallthrough
	case strings.ToLower(string(Day)):
		fallthrough
	case strings.ToLower(string(Week)):
		fallthrough
	case strings.ToLower(string(Month)):
		fallthrough
	case strings.ToLower(string(Year)):
		return true
	}

	return false
}
