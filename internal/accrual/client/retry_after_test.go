package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractRetryAfter(t *testing.T) {
	type testCase struct {
		name      string
		header    http.Header
		wantError bool
		wantRetry int
	}

	cases := []testCase{
		{
			name: "Correct",
			header: http.Header{
				"Retry-After": []string{"10"},
			},
			wantError: false,
			wantRetry: 10,
		},
		{
			name:      "Missing header",
			header:    http.Header{},
			wantError: true,
		},
		{
			name: "Invalid header",
			header: http.Header{
				"Retry-After": []string{"abc"},
			},
			wantError: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			retry, err := extractRetryAfter(testCase.header)

			if testCase.wantError {
				assert.Error(t, err, "Expected error, got nil")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, testCase.wantRetry, retry)
			}
		})
	}
}
