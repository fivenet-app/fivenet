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
	start := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		check    time.Time
		expected bool
	}{
		{
			name:     "Within range",
			check:    time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Before range",
			check:    time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "After range",
			check:    time.Date(2023, 10, 1, 21, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "At start of range",
			check:    start,
			expected: true,
		},
		{
			name:     "At end of range",
			check:    end,
			expected: true,
		},
		{
			name:     "Start equals end, match",
			check:    start,
			expected: true,
		},
		{
			name:     "Start equals end, no match",
			check:    time.Date(2023, 10, 1, 11, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Start after end, within range",
			check:    time.Date(2023, 10, 1, 5, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Start after end, outside range",
			check:    time.Date(2023, 10, 1, 22, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InTimeSpan(start, end, tt.check)
			assert.Equal(t, tt.expected, result)
		})
	}
}
