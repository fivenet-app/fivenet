package timeutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartOfDay(t *testing.T) {
	timestamp := time.Date(2023, 10, 1, 15, 30, 45, 123, time.UTC)
	expected := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)

	result := StartOfDay(timestamp)
	assert.Equal(t, expected, result)
}

func TestEndOfDay(t *testing.T) {
	timestamp := time.Date(2023, 10, 1, 15, 30, 45, 123, time.UTC)
	expected := time.Date(2023, 10, 1, 23, 59, 59, 999, time.UTC)

	result := EndOfDay(timestamp)
	assert.Equal(t, expected, result)
}

func TestInTimeSpan(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		check    time.Time
		expected bool
	}{
		{
			name:     "Within range",
			start:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Before range",
			start:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "After range",
			start:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 21, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "At start of range",
			start:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "At end of range",
			start:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Start equals end, match",
			start:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Start equals end, no match",
			start:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 11, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Start after end, within range",
			start:    time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 11, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Start after end, outside range",
			start:    time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			check:    time.Date(2023, 10, 1, 22, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InTimeSpan(tt.start, tt.end, tt.check)
			assert.Equal(t, tt.expected, result)
		})
	}
}
