package calendarstore

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRSVPCalendarEntryUsesOccurrenceTable(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	mock.ExpectQuery(regexp.QuoteMeta("fivenet_calendar_rsvp_occurrence")).
		WillReturnRows(sqlmock.NewRows([]string{"entry_id"}))
	mock.ExpectQuery(regexp.QuoteMeta("fivenet_calendar_rsvp AS calendar_entry_rsvp")).
		WillReturnRows(sqlmock.NewRows([]string{"entry_id"}))

	entry, err := store.GetRSVPCalendarEntry(t.Context(), 1, 2, "recurring:1:1:10")
	require.NoError(t, err)
	assert.Nil(t, entry)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestValidateRecuringOccurrenceKey(t *testing.T) {
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

	store := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := store.ValidateRecurringOccurrenceKey(entry, tt.key)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}
