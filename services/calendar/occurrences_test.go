package calendar

import (
	"testing"
	"time"

	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
)

func TestNextRecurringOccurrence(t *testing.T) {
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
			got := nextRecurringOccurrence(start, tc.interval, tc.every)
			if !got.Equal(tc.want) {
				t.Fatalf("expected %s, got %s", tc.want, got)
			}
		})
	}
}

func TestEntryOverlapsRange(t *testing.T) {
	start := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	end := timestamp.New(time.Date(2026, time.January, 15, 11, 0, 0, 0, time.UTC))

	if !entryOverlapsRange(
		start,
		end,
		time.Date(2026, time.January, 15, 9, 0, 0, 0, time.UTC),
		time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
	) {
		t.Fatal("expected range to overlap")
	}

	if entryOverlapsRange(
		start,
		end,
		time.Date(2026, time.January, 15, 12, 1, 0, 0, time.UTC),
		time.Date(2026, time.January, 15, 13, 0, 0, 0, time.UTC),
	) {
		t.Fatal("expected range not to overlap")
	}
}
