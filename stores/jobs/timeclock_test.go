package jobsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestStoreCleanupTimeclock(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_timeclock SET start_time = NULL WHERE ( (fivenet_job_timeclock.date BETWEEN (CURRENT_DATE - INTERVAL 14 DAY) AND (CURRENT_DATE - INTERVAL 2 DAY)) AND (fivenet_job_timeclock.start_time IS NOT NULL) AND (fivenet_job_timeclock.end_time IS NULL) ) LIMIT ?;`)).
		WithArgs(int64(1000)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.CleanupTimeclock(t.Context(), store.db))
	require.NoError(t, mock.ExpectationsWereMet())
}
