// Package timeutils provides utility functions for time manipulation and comparing.
package timeutils

import (
	"time"
)

// StartOfDay Based upon https://stackoverflow.com/a/36988882 from "VinGarcia".
// It returns the start of the day, which is 00:00:00.000.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay Based upon https://stackoverflow.com/a/36988882 from "VinGarcia".
// It returns the end of the day, which is 23:59:59.999.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, t.Location())
}

// InTimeSpan Based upon https://stackoverflow.com/a/55093788 from "Kamil Dziedzic".
// It checks if the current time is within the start and end times, inclusive.
// If the start time is after the end time, it checks if the current time is not before the end time and not after the start time.
// If the start time is before the end time, it checks if the current time is not before the start time and not after the end time.
// If both times are equal, it checks if the current time is equal to the start time.
func InTimeSpan(start time.Time, end time.Time, current time.Time) bool {
	if start.After(end) {
		return !current.Before(end) && !current.After(start)
	}

	if start.Before(end) {
		return !current.Before(start) && !current.After(end)
	}

	if start.Equal(end) {
		return current.Equal(start)
	}

	return !start.After(current) || !end.Before(current)
}
