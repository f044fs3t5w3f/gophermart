package handler

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalOrderResponse(t *testing.T) {
	type testCase struct {
		name     string
		order    orderResponse
		expected string
	}
	tests := []testCase{
		{
			name: "new",
			order: orderResponse{
				Number: "123456",
				Status: "NEW",
			},

			expected: `{"number":"123456","status":"NEW"}`,
		},
		{
			name: "processed with 0 accrual",
			order: orderResponse{
				Number: "654321",
				Status: "PROCESSED",
			},
			expected: `{"number":"654321","status":"PROCESSED","accrual":0}`,
		},
		{
			name: "processed with not 0 accrual",
			order: orderResponse{
				Number:  "654321",
				Status:  "PROCESSED",
				Accrual: 10,
			},
			expected: `{"number":"654321","status":"PROCESSED","accrual":10}`,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {

		})
		marshaled, err := json.Marshal(test.order)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, string(marshaled))
	}
}
