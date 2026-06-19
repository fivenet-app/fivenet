package calendar

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/stretchr/testify/assert"
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
