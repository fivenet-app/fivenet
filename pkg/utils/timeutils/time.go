package timeutils

import (
	"time"
)

// Based upon https://stackoverflow.com/a/36988882 from `VinGarcia`
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, t.Location())
}

// Based upon https://stackoverflow.com/a/55093788 from `Kamil Dziedzic`
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
