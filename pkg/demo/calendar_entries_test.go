package demo

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	calendarstore "github.com/fivenet-app/fivenet/v2026/stores/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newCalendarSeedDemo(seed uint64, db *sql.DB) *Demo {
	cfg := &config.Config{}
	cfg.Demo.Seed = seed
	cfg.Demo.TargetJob = PoliceJob
	cfg.Demo.Features.CalendarEntries = true

	d := &Demo{
		cfg:       cfg,
		db:        db,
		logger:    zap.NewNop(),
		calendars: nil,
	}
	if db != nil {
		d.accessResolver = access.NewSubjectResolver(db)
		d.calendarAccess = access.NewCalendarSubjectObjectAccess(db)
		d.calendars = calendarstore.New(calendarstore.Params{
			DB:     db,
			Access: access.NewCalendarSubjectObjectAccess(db),
		})
	}
	d.initRandomizers()

	return d
}

type demoCalendarEntrySnapshot struct {
	Title    string
	Start    time.Time
	End      *time.Time
	AllDay   bool
	Content  string
	Calendar int64
}

func snapshotDemoCalendarEntries(
	entries []*calendarentries.CalendarEntry,
) []demoCalendarEntrySnapshot {
	out := make([]demoCalendarEntrySnapshot, 0, len(entries))
	for _, entry := range entries {
		if entry == nil || entry.GetStartTime() == nil {
			out = append(out, demoCalendarEntrySnapshot{})
			continue
		}

		snap := demoCalendarEntrySnapshot{
			Title:    entry.GetTitle(),
			Start:    entry.GetStartTime().AsTime(),
			AllDay:   entry.GetAllDay(),
			Calendar: entry.GetCalendarId(),
		}
		if entry.GetEndTime() != nil {
			end := entry.GetEndTime().AsTime()
			snap.End = &end
		}
		if entry.GetContent() != nil {
			snap.Content = entry.GetContent().GetRawHtml()
		}

		out = append(out, snap)
	}

	return out
}

func TestBuildDemoCalendarEntriesDeterministic(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 21, 12, 0, 0, 0, time.UTC)

	d1 := newCalendarSeedDemo(42, nil)
	d2 := newCalendarSeedDemo(42, nil)

	entries1 := d1.buildDemoCalendarEntriesAt(99, now)
	entries2 := d2.buildDemoCalendarEntriesAt(99, now)

	require.Len(t, entries1, len(entries2))
	assert.GreaterOrEqual(t, len(entries1), 32)
	assert.LessOrEqual(t, len(entries1), 52)
	assert.Equal(t, snapshotDemoCalendarEntries(entries1), snapshotDemoCalendarEntries(entries2))

	for _, entry := range entries1 {
		require.NotNil(t, entry)
		assert.False(t, strings.HasPrefix(entry.GetTitle(), demoCalendarEntryPrefix))
		require.NotNil(t, entry.GetContent())
		assert.Contains(t, entry.GetContent().GetRawHtml(), demoCalendarEntryPrefix)
	}
}

func TestBuildDemoCalendarEntriesStayWithinExpectedWindows(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 21, 12, 0, 0, 0, time.UTC)
	d := newCalendarSeedDemo(1337, nil)
	entries := d.buildDemoCalendarEntriesAt(42, now)

	require.NotEmpty(t, entries)

	pastLimit := windowAnchor(now.AddDate(0, -1, 0))
	futureLimit := now.AddDate(0, 2, 0)
	weekStart := currentWeekStart(now)
	currentWeekEnd := weekStart.AddDate(0, 0, 7)

	shortCount := 0
	allDayCount := 0
	multiDayCount := 0
	currentWeekCount := 0
	recurringCount := 0
	monthlyRecurringCount := 0
	weeklyTuesdayRecurringCount := 0

	for _, entry := range entries {
		require.NotNil(t, entry)
		require.NotNil(t, entry.GetStartTime())

		start := entry.GetStartTime().AsTime()
		if !start.Before(weekStart) && start.Before(currentWeekEnd) {
			currentWeekCount++
		}
		if start.Before(now) {
			assert.False(
				t,
				start.Before(pastLimit),
				"past entry start %s exceeds one month back",
				start,
			)
		} else {
			assert.False(
				t,
				start.After(futureLimit),
				"future entry start %s exceeds two months forward",
				start,
			)
		}

		if recurring := entry.GetRecurring(); recurring != nil {
			recurringCount++
			require.NotNil(t, recurring.GetUntil())
			switch recurring.GetEvery() {
			case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_MONTH:
				monthlyRecurringCount++
			case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_WEEK:
				weeklyTuesdayRecurringCount++
				assert.Equal(t, time.Tuesday, start.Weekday())
			}
		}

		if !entry.GetAllDay() {
			shortCount++
			require.NotNil(t, entry.GetEndTime())
			duration := entry.GetEndTime().AsTime().Sub(start)
			assert.GreaterOrEqual(t, duration, 30*time.Minute)
			assert.LessOrEqual(t, duration, 90*time.Minute)
			assert.NotEqual(t, 0, start.Hour()%24, "short meeting should not start at midnight")
			continue
		}

		allDayCount++
		require.NotNil(t, entry.GetEndTime())
		assert.Equal(t, 0, start.Hour())
		assert.Equal(t, 0, start.Minute())
		duration := entry.GetEndTime().AsTime().Sub(start)
		if duration <= 26*time.Hour {
			assert.Equal(t, 24*time.Hour, duration, "all-day entry should span a single day")
		} else {
			multiDayCount++
			assert.GreaterOrEqual(t, duration, 48*time.Hour)
			assert.LessOrEqual(t, duration, 96*time.Hour)
		}
	}

	assert.Positive(t, shortCount)
	assert.Positive(t, allDayCount)
	assert.Positive(t, multiDayCount)
	assert.GreaterOrEqual(t, currentWeekCount, 6)
	assert.Equal(t, 2, recurringCount)
	assert.Equal(t, 1, monthlyRecurringCount)
	assert.Equal(t, 1, weeklyTuesdayRecurringCount)
}

func TestLookupDemoCalendarAndCreateFallback(t *testing.T) {
	t.Parallel()

	job := PoliceJob
	fallbackName := fmt.Sprintf(demoCalendarFallbackName, titleizeJob(job))

	t.Run("existing", func(t *testing.T) {
		t.Parallel()

		d := newCalendarSeedDemo(9, nil)
		stmt := d.clearDemoCalendarEntriesStmt(17)
		sql, _ := stmt.Sql()
		assert.Contains(t, sql, "DELETE FROM fivenet_calendar_entries")
		assert.Contains(t, sql, "fivenet_calendar_entries.content")
	})

	t.Run("fallback", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		d := newCalendarSeedDemo(10, db)
		d.calendarAccess = nil
		d.accessResolver = nil
		seedUser := &userinfo.UserInfo{
			UserId:   7,
			Job:      job,
			JobGrade: d.highestJobGrade(job),
		}
		mock.ExpectBegin()
		mock.ExpectExec(`(?s).*INSERT INTO fivenet_calendar.*`).
			WillReturnResult(sqlmock.NewResult(99, 1))
		mock.ExpectCommit()

		cal, err := d.createFallbackDemoCalendar(t.Context(), job, seedUser)
		require.NoError(t, err)
		require.NotNil(t, cal)
		assert.Equal(t, int64(99), cal.GetId())
		assert.Equal(t, fallbackName, cal.GetName())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSeedDemoCalendarEntriesForCalendarIsIdempotent(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 21, 12, 0, 0, 0, time.UTC)
	calendarID := int64(17)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	seedUser := &userinfo.UserInfo{
		UserId:   4,
		Job:      PoliceJob,
		JobGrade: 19,
	}

	dExpect := newCalendarSeedDemo(77, nil)
	expectedCount := len(dExpect.buildDemoCalendarEntriesAt(calendarID, now))

	expectRun := func(startID int64) {
		mock.ExpectBegin()
		mock.ExpectExec(`(?s).*DELETE FROM fivenet_calendar_entries.*`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		for i := range expectedCount {
			mock.ExpectExec(`(?s).*INSERT INTO fivenet_calendar_entries.*`).
				WillReturnResult(sqlmock.NewResult(startID+int64(i), 1))
		}
		mock.ExpectCommit()
	}

	expectRun(1)
	expectRun(100)

	d1 := newCalendarSeedDemo(77, db)
	require.NoError(
		t,
		d1.seedDemoCalendarEntriesForCalendar(t.Context(), calendarID, seedUser, now),
	)

	d2 := newCalendarSeedDemo(77, db)
	require.NoError(
		t,
		d2.seedDemoCalendarEntriesForCalendar(t.Context(), calendarID, seedUser, now),
	)

	require.NoError(t, mock.ExpectationsWereMet())
}
