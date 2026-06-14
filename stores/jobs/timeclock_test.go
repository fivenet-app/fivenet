package jobs

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestStoreCleanupTimeclock(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_timeclock AS timeclock_entry SET start_time = NULL WHERE`)).
		WithArgs(int64(1000)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.CleanupTimeclock(t.Context(), store.db))
	require.NoError(t, mock.ExpectationsWereMet())
}
