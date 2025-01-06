package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const DateTemplate = "01/02/2006" // Example format for dates

func TestGetTodayAndYesterday(t *testing.T) {

	tests := []struct {
		name              string
		expectedToday     string
		expectedYesterday string
	}{
		{
			name:              "Regular Day",
			expectedToday:     time.Now().Format(DateTemplate),
			expectedYesterday: time.Now().Add(-24 * time.Hour).Format(DateTemplate),
		},
	}

	// Iterate through test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			actualYesterday, actualToday := getTodayAndYesterday()

			// Validate the results with assertions
			assert.Equal(t, tc.expectedToday, actualToday, "Mismatch for today in test: %s", tc.name)
			assert.Equal(t, tc.expectedYesterday, actualYesterday, "Mismatch for yesterday in test: %s", tc.name)
		})
	}
}
