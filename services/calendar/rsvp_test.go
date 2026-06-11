package calendar

import (
	"errors"
	"fmt"
	"testing"
	"time"

	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRecurringOccurrenceKey(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)

	entry := &calendarentries.CalendarEntry{
		Id:                123,
		StartTime:         timestamp.New(start),
		RecurrenceVersion: 2,
		Recurring: &calendarentries.CalendarEntryRecurring{
			Every: calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_WEEK,
			Count: 1,
			Until: timestamp.New(start.AddDate(0, 1, 0)),
		},
	}

	tests := []struct {
		name string
		key  string
		want error
	}{
		{
			name: "valid first occurrence",
			key:  "recurring:123:2:1780308000",
			want: nil,
		},
		{
			name: "valid later occurrence",
			key:  fmt.Sprintf("recurring:123:2:%d", start.AddDate(0, 0, 14).Unix()),
			want: nil,
		},
		{
			name: "wrong entry id",
			key:  fmt.Sprintf("recurring:999:2:%d", start.Unix()),
			want: errorscalendar.ErrNoPerms,
		},
		{
			name: "wrong recurrence version",
			key:  fmt.Sprintf("recurring:123:1:%d", start.Unix()),
			want: errorscalendar.ErrNoPerms,
		},
		{
			name: "not an occurrence",
			key:  fmt.Sprintf("recurring:123:2:%d", start.AddDate(0, 0, 1).Unix()),
			want: errorscalendar.ErrNoPerms,
		},
		{
			name: "after until",
			key:  fmt.Sprintf("recurring:123:2:%d", start.AddDate(0, 2, 0).Unix()),
			want: errorscalendar.ErrNoPerms,
		},
		{
			name: "old key shape rejected",
			key:  fmt.Sprintf("recurring:123:%d", start.Unix()),
			want: errorscalendar.ErrNoPerms,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateRecurringOccurrenceKey(entry, tt.key)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestFilterUpcomingCalendarEntries(t *testing.T) {
	t.Parallel()

	userInfo := &userinfo.UserInfo{UserId: 7}

	entries := []*calendarentries.CalendarEntry{
		{
			Id:        1,
			Title:     "birthday",
			StartTime: timestamp.New(time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)),
			Occurrence: &calendarentries.CalendarEntryOccurrence{
				Kind: calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_BIRTHDAY,
			},
		},
		{
			Id:        2,
			Title:     "not visible",
			StartTime: timestamp.New(time.Date(2026, time.January, 16, 10, 0, 0, 0, time.UTC)),
			Occurrence: &calendarentries.CalendarEntryOccurrence{
				Kind: calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING,
			},
		},
		{
			Id:        3,
			Title:     "override yes",
			StartTime: timestamp.New(time.Date(2026, time.January, 17, 10, 0, 0, 0, time.UTC)),
			Occurrence: &calendarentries.CalendarEntryOccurrence{
				Kind: calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING,
			},
			Rsvp: &calendarentries.CalendarEntryRSVP{
				Response: calendarentries.RsvpResponses_RSVP_RESPONSES_YES,
			},
		},
		{
			Id:        4,
			Title:     "own entry",
			StartTime: timestamp.New(time.Date(2026, time.January, 18, 10, 0, 0, 0, time.UTC)),
			CreatorId: func() *int32 { v := int32(7); return &v }(),
		},
	}

	filtered := filterUpcomingCalendarEntries(entries, userInfo)
	require.Len(t, filtered, 3)
	assert.Equal(t, "birthday", filtered[0].GetTitle())
	assert.Equal(t, "override yes", filtered[1].GetTitle())
	assert.Equal(t, "own entry", filtered[2].GetTitle())
}

func TestCalendarEntryVisibilityIncludesRecurringOverrides(t *testing.T) {
	t.Parallel()

	stmt := tCalendarEntry.
		SELECT(mysql.Int(1)).
		WHERE(calendarEntryVisibility(
			&userinfo.UserInfo{UserId: 7, Job: "test", JobGrade: 0},
			calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
			calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN,
		))

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "fivenet_calendar_rsvp_occurrence")
}
