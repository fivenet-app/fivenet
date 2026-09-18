package housekeeper

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adhocore/gronx"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type assignmentExpirationWriterStub struct {
	dispatchID  int64
	unitIDs     []int64
	deleteCalls int
	updateErr   error
}

func (s *assignmentExpirationWriterStub) UpdateAssignments(
	_ context.Context,
	_ *string,
	_ *int32,
	dispatchID int64,
	_ []int64,
	unitIDs []int64,
	_ time.Time,
) error {
	s.dispatchID = dispatchID
	s.unitIDs = append([]int64(nil), unitIDs...)
	return s.updateErr
}

func (s *assignmentExpirationWriterStub) DeleteExpiredAssignments(
	_ context.Context,
	_ int64,
	_ []int64,
) (int64, error) {
	s.deleteCalls++
	return 0, nil
}

func TestHandleDispatchAssignmentExpirationFallsBackToSQLWithoutKVTimer(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	rows := sqlmock.NewRows([]string{"dispatch_id", "unit_id", "job"})
	rows.AddRow(42, 7, "ambulance")
	rows.AddRow(42, 8, "ambulance")
	mock.ExpectQuery("SELECT .*dispatch_id AS.*unit_id AS.*job AS.*FROM.*expires_at.*LIMIT").
		WillReturnRows(rows)

	writer := &assignmentExpirationWriterStub{}
	h := &Housekeeper{
		db:                         db,
		logger:                     zap.NewNop(),
		assignmentExpirationWriter: writer,
	}

	expired, dispatches, jobs, units, backlog, err := h.handleDispatchAssignmentExpiration(
		t.Context(),
	)
	require.NoError(t, err)
	require.Equal(t, 2, expired)
	require.Equal(t, 1, dispatches)
	require.Equal(t, 1, jobs)
	require.Equal(t, 2, units)
	require.False(t, backlog)
	require.Equal(t, int64(42), writer.dispatchID)
	require.ElementsMatch(t, []int64{7, 8}, writer.unitIDs)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleDispatchAssignmentExpirationCleansArchivedDispatch(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	rows := sqlmock.NewRows([]string{"dispatch_id", "unit_id", "job"})
	rows.AddRow(509171, 7, "ambulance")
	mock.ExpectQuery("SELECT .*dispatch_id AS.*unit_id AS.*job AS.*FROM.*expires_at.*LIMIT").
		WillReturnRows(rows)

	writer := &assignmentExpirationWriterStub{updateErr: jetstream.ErrKeyNotFound}
	h := &Housekeeper{
		db:                         db,
		logger:                     zap.NewNop(),
		assignmentExpirationWriter: writer,
	}

	expired, dispatches, jobs, units, backlog, err := h.handleDispatchAssignmentExpiration(
		t.Context(),
	)
	require.NoError(t, err)
	assert.Equal(t, 1, expired)
	assert.Equal(t, 1, dispatches)
	assert.Equal(t, 1, jobs)
	assert.Equal(t, 1, units)
	assert.False(t, backlog)
	assert.Equal(t, 1, writer.deleteCalls)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDispatchHousekeeperSchedules(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, schedule string
		want           time.Time
	}{
		{"targeted cancellation recovery", cancelOldDispatchesSchedule, time.Date(2026, time.January, 1, 0, 5, 0, 0, time.UTC)},
		{"empty unit dispatch recovery", auditEmptyUnitDispatchesSchedule, time.Date(2026, time.January, 1, 0, 5, 0, 0, time.UTC)},
		{"kv recovery audit", deleteOldDispatchesKVSchedule, time.Date(2026, time.January, 1, 0, 30, 0, 0, time.UTC)},
		{"authoritative user info recovery audit", reconcileUserInfoStateSchedule, time.Date(2026, time.January, 1, 4, 0, 0, 0, time.UTC)},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			next, err := gronx.NextTickAfter(
				test.schedule,
				time.Date(2026, time.January, 1, 0, 0, 1, 0, time.UTC),
				false,
			)
			require.NoError(t, err)
			assert.Equal(t, test.want, next)
		})
	}
}

func TestDeleteOldDispatchesSelectsOldestFirst(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectQuery("(?s)ORDER BY.*created_at.*ASC.*id.*ASC.*LIMIT \\?").
		WithArgs(75).
		WillReturnRows(sqlmock.NewRows([]string{"dispatch_id"}))
	h := &Housekeeper{db: db}
	deleted, err := h.deleteOldDispatches(t.Context())
	require.NoError(t, err)
	assert.Zero(t, deleted)
	require.NoError(t, mock.ExpectationsWereMet())
}
