package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBirthdayForYear(t *testing.T) {
	t.Parallel()

	got := birthdayForYear(2025, time.February, 29)

	assert.Equal(t, 2025, got.Year(), "expected leap day to clamp to february 28, got %s", got)
	assert.Equal(
		t,
		time.February,
		got.Month(),
		"expected leap day to clamp to february 28, got %s",
		got,
	)
	assert.Equal(t, 28, got.Day(), "expected leap day to clamp to february 28, got %s", got)
	assert.Equal(t, 12, got.Hour(), "expected birthday occurrences to use noon UTC, got %s", got)
	assert.Equal(t, 0, got.Minute(), "expected birthday occurrences to use noon UTC, got %s", got)
}
