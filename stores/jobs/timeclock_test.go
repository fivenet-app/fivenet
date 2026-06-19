package jobsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
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

func TestStoreListTimeclockDefaultOrderByUsesAggregatedColumns(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(`(?s)SELECT COUNT\(DISTINCT .*timeclock_entry\.user_id.* AS .*data_count\.total.*FROM .*fivenet_job_timeclock.*INNER JOIN .*fivenet_user.*ON .*colleague\.id = timeclock_entry\.user_id.*WHERE .*timeclock_entry\.job = \?.*;`).
		WithArgs("police").
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	mock.ExpectQuery(`(?s)SELECT .*ORDER BY .*agg\.date DESC, agg\.spent_time DESC.*`).
		WillReturnRows(sqlmock.NewRows(nil))

	_, err := store.ListTimeclock(t.Context(), store.db, TimeclockQuery{
		Job:      "police",
		UserMode: jobstimeclock.TimeclockViewMode_TIMECLOCK_VIEW_MODE_ALL,
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
