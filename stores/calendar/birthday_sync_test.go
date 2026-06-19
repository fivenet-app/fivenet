package calendarstore

import (
	"testing"
	"time"

	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestBirthdayCalendarAccessEntries(t *testing.T) {
	t.Parallel()

	job := &jobs.Job{
		Grades: []*jobs.JobGrade{
			{Grade: 1},
			{Grade: 4},
			{Grade: 7},
		},
	}

	entries := birthdayCalendarAccessEntries(42, "police", job)
	require.Len(t, entries, 2, "expected 2 access entries")

	assert.Equal(
		t,
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW),
		entries[0].GetAccess(),
		"unexpected view access entry: %+v",
		entries[0],
	)
	assert.Equal(
		t,
		int32(1),
		entries[0].GetMinimumGrade(),
		"unexpected view access entry: %+v",
		entries[0],
	)

	assert.Equal(
		t,
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT),
		entries[1].GetAccess(),
		"unexpected edit access entry: %+v",
		entries[1],
	)
	assert.Equal(
		t,
		int32(7),
		entries[1].GetMinimumGrade(),
		"unexpected edit access entry: %+v",
		entries[1],
	)
}
