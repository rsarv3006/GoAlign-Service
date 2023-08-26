package helper_test

import (
	"testing"
	"time"

	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func TestFindDateFromDateAndInterval(t *testing.T) {

	testCases := []struct {
		date        time.Time
		interval    model.IntervalObj
		expected    time.Time
		description string
	}{
		{
			description: "should return a date time advanced by 1 minute",
			date:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval: model.IntervalObj{
				IntervalCount: 1,
				IntervalUnit:  model.Minute,
			},
			expected: time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC),
		},
		{
			description: "should return a date time advanced by 1 hour",
			date:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval: model.IntervalObj{
				IntervalCount: 1,
				IntervalUnit:  model.Hour,
			},
			expected: time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		},
		{
			description: "should return a date time advanced by 1 day",
			date:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval: model.IntervalObj{
				IntervalCount: 1,
				IntervalUnit:  model.Day,
			},
			expected: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			description: "should return a date time advanced by 1 week",
			date:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval: model.IntervalObj{
				IntervalCount: 1,
				IntervalUnit:  model.Week,
			},
			expected: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			description: "should return a date time advanced by 1 month",
			date:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval: model.IntervalObj{
				IntervalCount: 1,
				IntervalUnit:  model.Month,
			},
			expected: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			description: "should return a date time advanced by 1 year",
			date:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval: model.IntervalObj{
				IntervalCount: 1,
				IntervalUnit:  model.Year,
			},
			expected: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		output, err := helper.FindDateFromDateAndInterval(tc.date, tc.interval)
		expectedDateString := tc.expected.Format(time.RFC3339)
		outputDateString := output.Format(time.RFC3339)
		if err != nil || outputDateString != expectedDateString {
			t.Errorf("FindDateFromDateAndInterval(%v, %v) = %v, want %v", tc.date, tc.interval, output, tc.expected)
		}
	}
}
