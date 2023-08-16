package helper

import "time"

func GetToday() time.Time {
	return time.Now().UTC()
}
