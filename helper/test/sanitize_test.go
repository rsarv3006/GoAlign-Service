package helper_test

import (
	"testing"

	"gitlab.com/donutsahoy/yourturn-fiber/helper"
)

func TestSanitizeInput(t *testing.T) {

	testcases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello World!",
			expected: "HelloWorld",
		},
		{
			input:    "123abc@#$*",
			expected: "123abc",
		},
		{
			input:    "abcdefghijklmnopqrstuvwxyz1234567890",
			expected: "abcdefghijklmnopqrstuvwxyz1234567890",
		},
	}

	for _, tc := range testcases {
		output := helper.SanitizeInput(tc.input)
		if output != tc.expected {
			t.Errorf("RemoveNonAlphaNum(%s) = %s, expected %s", tc.input, output, tc.expected)
		}
	}

}
