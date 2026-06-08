package calendar

import (
	"testing"
	"time"
)

func TestBirthdayForYear(t *testing.T) {
	got := birthdayForYear(2025, time.February, 29)

	if got.Year() != 2025 || got.Month() != time.February || got.Day() != 28 {
		t.Fatalf("expected leap day to clamp to february 28, got %s", got)
	}
	if got.Hour() != 12 || got.Minute() != 0 {
		t.Fatalf("expected birthday occurrences to use noon UTC, got %s", got)
	}
}
