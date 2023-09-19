package helper

import "time"

func GetToday() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	midnight := today.Add(time.Second + time.Nanosecond)
	return midnight
}
