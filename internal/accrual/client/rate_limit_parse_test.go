package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractRateLimit(t *testing.T) {
	type testCase struct {
		name      string
		text      string
		wantError bool
		wantRate  int
	}

	cases := []testCase{
		{
			name:      "Correct",
			text:      "No more than 30 requests per minute allowed",
			wantError: false,
			wantRate:  30,
		},
		{
			name:      "Wrong pattern",
			text:      "No more than 30 requests per minute",
			wantError: true,
		},
		{
			name:      "wrong number",
			text:      "No more than a requests per minute allowed",
			wantError: true,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			rate, err := extractRateLimit(testCase.text)
			if testCase.wantError {
				assert.Error(t, err, "Expected error, got nil")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, testCase.wantRate, rate)
			}
		})
	}
}
