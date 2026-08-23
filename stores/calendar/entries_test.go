package calendarstore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
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
		nil,
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
		nil,
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

func TestCalendarEntriesQueryUsesAliasedCalendarEntryColumnForBirthdayVisibility(t *testing.T) {
	t.Parallel()

	store := New(testParams(new(sql.DB))).(*Store)
	stmt := calendarEntriesQuery(
		&userinfo.UserInfo{UserId: 1, Job: "police", Superuser: true},
		mysql.Bool(true),
		store.birthdayCalendarVisible(
			tCalendarEntry.CalendarID,
			calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
			&userinfo.UserInfo{UserId: 1, Job: "police", Superuser: true},
		),
		nil,
		false,
		nil,
	)

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "calendar_entry.calendar_id")
	assert.NotContains(t, sql, "fivenet_calendar_entries.calendar_id")
}

func TestCalendarEntriesQuerySelectsIconColumns(t *testing.T) {
	t.Parallel()

	stmt := calendarEntriesQuery(
		&userinfo.UserInfo{UserId: 1, Superuser: true},
		mysql.Bool(true),
		mysql.Bool(true),
		nil,
		false,
		nil,
	)

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "calendar.icon")
	assert.Contains(t, sql, "calendar_entry.icon")
}

func TestUpsertCalendarEntryInsertIncludesIconColumn(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(testParams(db)).(*Store)
	icon := "CalendarStarIcon"

	mock.ExpectExec(`(?s)INSERT INTO fivenet_calendar_entries .*\bicon\b.*`).
		WillReturnResult(sqlmock.NewResult(101, 1))

	id, err := store.UpsertCalendarEntry(
		t.Context(),
		db,
		&calendarentries.CalendarEntry{
			CalendarId: 1,
			Title:      "Test entry",
			Icon:       &icon,
			StartTime:  timestamp.New(time.Date(2026, time.August, 23, 10, 0, 0, 0, time.UTC)),
			Closed:     false,
		},
		nil,
		&userinfo.UserInfo{UserId: 7, Job: "police"},
	)
	require.NoError(t, err)
	assert.Equal(t, int64(101), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertCalendarEntryUpdateIncludesIconColumn(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(testParams(db)).(*Store)
	icon := "CalendarStarIcon"

	mock.ExpectExec(`(?s)UPDATE fivenet_calendar_entries .*\bicon\b.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	id, err := store.UpsertCalendarEntry(
		t.Context(),
		db,
		&calendarentries.CalendarEntry{
			Id:         77,
			CalendarId: 1,
			Title:      "Updated entry",
			Icon:       &icon,
			StartTime:  timestamp.New(time.Date(2026, time.August, 23, 10, 0, 0, 0, time.UTC)),
			Closed:     false,
		},
		&calendarentries.CalendarEntry{Id: 77},
		&userinfo.UserInfo{UserId: 7, Job: "police"},
	)
	require.NoError(t, err)
	assert.Equal(t, int64(77), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCalendarEntryVisibilityAllowsCreatorOwnedPrivateCalendars(t *testing.T) {
	t.Parallel()

	store := New(testParams(new(sql.DB))).(*Store)
	stmt := mysql.
		SELECT(mysql.Int(1)).
		FROM(tCalendar).
		WHERE(calendarEntryVisibility(
			store.access,
			&userinfo.UserInfo{UserId: 1},
			calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
			calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN,
		))

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "job IS NULL")
	assert.Contains(t, sql, "creator_id = ?")
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
			CreatorId: new(int32(7)),
		},
		{
			Id:        5,
			Title:     "not visible",
			StartTime: timestamp.New(time.Date(2026, time.January, 19, 10, 0, 0, 0, time.UTC)),
		},
	}

	store := New(testParams(new(sql.DB)))

	filtered := store.FilterUpcomingCalendarEntries(entries, userInfo)
	require.Len(t, filtered, 3)
	assert.Equal(t, "birthday", filtered[0].GetTitle())
	assert.Equal(t, "override yes", filtered[1].GetTitle())
	assert.Equal(t, "own entry", filtered[2].GetTitle())
}
