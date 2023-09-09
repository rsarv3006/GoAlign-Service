package helper_test

import (
	"testing"

	"gitlab.com/donutsahoy/yourturn-fiber/helper"
)

func TestIsValidEmailAddress(t *testing.T) {

	testcases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "Hello World!",
			expected: false,
		},
		{
			input:    "123abc@#$*",
			expected: false,
		},
		{
			input:    "test@test.com",
			expected: true,
		},
	}

	for _, tc := range testcases {
		output := helper.IsValidEmailAddress(tc.input)
		if output != tc.expected {
			t.Errorf("Expected %v, got %v", tc.expected, output)
		}
	}

}
