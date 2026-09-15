package housekeeper

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type assignmentExpirationWriterStub struct {
	dispatchID int64
	unitIDs    []int64
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
	return nil
}

func TestHandleDispatchAssignmentExpirationFallsBackToSQLWithoutKVTimer(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	rows := sqlmock.NewRows([]string{"dispatch_id", "unit_id", "job"})
	rows.AddRow(42, 7, "ambulance")
	rows.AddRow(42, 8, "ambulance")
	mock.ExpectQuery("SELECT .*dispatch_id AS.*unit_id AS.*job AS.*FROM.*expires_at.*LIMIT").WillReturnRows(rows)

	writer := &assignmentExpirationWriterStub{}
	h := &Housekeeper{
		db:                         db,
		logger:                     zap.NewNop(),
		assignmentExpirationWriter: writer,
	}

	expired, dispatches, jobs, units, backlog, err := h.handleDispatchAssignmentExpiration(t.Context())
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
