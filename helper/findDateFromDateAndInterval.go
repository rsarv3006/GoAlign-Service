package helper

import (
	"errors"
	"time"

	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func FindDateFromDateAndInterval(date time.Time, interval model.IntervalObj) (*time.Time, error) {
	if interval.IntervalCount == 0 {
		return nil, errors.New("Interval count cannot be 0")
	}

	if interval.IntervalUnit == model.Minute {
		date = date.Add(time.Minute * time.Duration(interval.IntervalCount))
	} else if interval.IntervalUnit == model.Hour {
		date = date.Add(time.Hour * time.Duration(interval.IntervalCount))
	} else if interval.IntervalUnit == model.Day {
		date = date.AddDate(0, 0, interval.IntervalCount)
	} else if interval.IntervalUnit == model.Week {
		date = date.AddDate(0, 0, interval.IntervalCount*7)
	} else if interval.IntervalUnit == model.Month {
		date = date.AddDate(0, interval.IntervalCount, 0)
	} else if interval.IntervalUnit == model.Year {
		date = date.AddDate(interval.IntervalCount, 0, 0)
	} else {
		return nil, errors.New("Interval unit not found")
	}

	return &date, nil
}
