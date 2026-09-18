package dispatches

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDeleteExpiredAssignmentsKeepsExpiryPredicate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec("(?s)DELETE FROM.*fivenet_centrum_dispatches_asgmts.*expires_at.*CURRENT_TIMESTAMP").WithArgs(int64(42), int64(7), int64(8), int64(2)).WillReturnResult(sqlmock.NewResult(0, 1))

	dispatches := &DispatchDB{db: db}
	deleted, err := dispatches.DeleteExpiredAssignments(t.Context(), 42, []int64{7, 8})
	require.NoError(t, err)
	require.Equal(t, int64(1), deleted)
	require.NoError(t, mock.ExpectationsWereMet())
}
