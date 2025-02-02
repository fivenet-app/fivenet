package timeutils

import (
	"time"
)

// Based upon https://stackoverflow.com/a/36988882 from `VinGarcia`
func TruncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func TruncateToNight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, t.Location())
}

// Based upon https://stackoverflow.com/a/55093788 from `Kamil Dziedzic`
func InTimeSpan(start time.Time, end time.Time, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}

	if start.Equal(end) {
		return check.Equal(start)
	}

	return !start.After(check) || !end.Before(check)
}
