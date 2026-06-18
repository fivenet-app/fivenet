package calendarstore

import (
	"database/sql"
	"testing"
	"time"

	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarEntriesQueryOmitLimitWhenNil(t *testing.T) {
	t.Parallel()

	stmt := calendarEntriesQuery(
		&userinfo.UserInfo{UserId: 1},
		mysql.Bool(true),
		mysql.Bool(true),
		false,
		nil,
	)

	sql, _ := stmt.Sql()
	assert.NotContains(t, sql, "LIMIT", "expected no limit in query, got %s", sql)
}

func TestCalendarEntriesQueryUsesExplicitLimit(t *testing.T) {
	t.Parallel()

	stmt := calendarEntriesQuery(
		&userinfo.UserInfo{UserId: 1},
		mysql.Bool(true),
		mysql.Bool(true),
		false,
		new(maxCalendarEntriesLimit),
	)

	sql, args := stmt.Sql()
	require.Contains(t, sql, "LIMIT ?", "expected explicit limit placeholder in query, got %s", sql)
	require.NotEmpty(t, args, "expected limit arguments")
	assert.Equal(
		t,
		maxCalendarEntriesLimit,
		args[len(args)-1],
		"expected limit argument %d, got %#v",
		maxCalendarEntriesLimit,
		args,
	)
}

func TestNextRecurringOccurrence(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		interval int32
		every    calendarentries.CalendarEntryRecurringEvery
		want     time.Time
	}{
		{
			name:     "day",
			interval: 2,
			every:    calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_DAY,
			want:     time.Date(2026, time.January, 17, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "week",
			interval: 1,
			every:    calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_WEEK,
			want:     time.Date(2026, time.January, 22, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "month",
			interval: 1,
			every:    calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_MONTH,
			want:     time.Date(2026, time.February, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "year",
			interval: 1,
			every:    calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_YEAR,
			want:     time.Date(2027, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := nextRecurringOccurrence(start, tc.interval, tc.every)
			assert.True(t, got.Equal(tc.want), "expected %s, got %s", tc.want, got)
		})
	}
}

func TestEntryOverlapsRange(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	end := timestamp.New(time.Date(2026, time.January, 15, 11, 0, 0, 0, time.UTC))

	require.True(t, entryOverlapsRange(
		start,
		end,
		time.Date(2026, time.January, 15, 9, 0, 0, 0, time.UTC),
		time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
	), "expected range to overlap")

	require.False(t, entryOverlapsRange(
		start,
		end,
		time.Date(2026, time.January, 15, 12, 1, 0, 0, time.UTC),
		time.Date(2026, time.January, 15, 13, 0, 0, 0, time.UTC),
	), "expected range not to overlap")
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
		{
			Id:        5,
			Title:     "not visible",
			StartTime: timestamp.New(time.Date(2026, time.January, 19, 10, 0, 0, 0, time.UTC)),
		},
	}

	store := New(new(sql.DB))

	filtered := store.FilterUpcomingCalendarEntries(entries, userInfo)
	require.Len(t, filtered, 3)
	assert.Equal(t, "birthday", filtered[0].GetTitle())
	assert.Equal(t, "override yes", filtered[1].GetTitle())
	assert.Equal(t, "own entry", filtered[2].GetTitle())
}
