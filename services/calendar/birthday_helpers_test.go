package calendar

import (
	"testing"

	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBirthdayCalendarDisplayName(t *testing.T) {
	t.Parallel()

	got := birthdayCalendarDisplayName(func(key string, vars map[string]any) string {
		return key + ":" + vars["job"].(string)
	}, "police", &jobs.Job{Label: "Police Department"})

	assert.Equal(
		t,
		"components.calendar.birthday_calendar_name:Police Department",
		got,
		"unexpected birthday calendar title: %s",
		got,
	)
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
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
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
		calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT,
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
