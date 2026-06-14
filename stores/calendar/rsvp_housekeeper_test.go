package calendarstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCleanupCalendarRSVPOccurrencesRunsBothDeletesInTransaction(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectBegin()
	mock.ExpectExec(`(?s).*fivenet_calendar_rsvp_occurrence.*INNER JOIN.*fivenet_calendar_entries.*created_at.*recurrence_version.*`).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`(?s).*fivenet_calendar_rsvp_occurrence.*LEFT JOIN.*fivenet_calendar_entries.*created_at.*IS NULL.*`).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	rows, err := store.CleanupCalendarRSVPOccurrences(t.Context())
	require.NoError(t, err)
	require.Equal(t, int64(5), rows)
	require.NoError(t, mock.ExpectationsWereMet())
}
