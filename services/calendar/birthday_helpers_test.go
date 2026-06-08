package calendar

import (
	"testing"

	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
)

func TestBirthdayCalendarDisplayName(t *testing.T) {
	got := birthdayCalendarDisplayName(func(key string, vars map[string]any) string {
		return key + ":" + vars["job"].(string)
	}, "police", &jobs.Job{Label: "Police Department"})

	if got != "components.calendar.birthday_calendar_name:Police Department" {
		t.Fatalf("unexpected birthday calendar title: %s", got)
	}
}

func TestBirthdayCalendarAccessEntries(t *testing.T) {
	job := &jobs.Job{
		Grades: []*jobs.JobGrade{
			{Grade: 1},
			{Grade: 4},
			{Grade: 7},
		},
	}

	entries := birthdayCalendarAccessEntries(42, "police", job)
	if len(entries) != 2 {
		t.Fatalf("expected 2 access entries, got %d", len(entries))
	}

	if entries[0].GetAccess() != calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW ||
		entries[0].GetMinimumGrade() != 1 {
		t.Fatalf("unexpected view access entry: %+v", entries[0])
	}

	if entries[1].GetAccess() != calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT ||
		entries[1].GetMinimumGrade() != 7 {
		t.Fatalf("unexpected edit access entry: %+v", entries[1])
	}
}
