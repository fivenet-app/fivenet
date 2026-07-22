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

func TestStoreCountInactiveEmployeesIgnoresUserPropsJoin(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(`(?s)SELECT COUNT\(fivenet_user_jobs\.user_id\).*FROM fivenet_user_jobs.*INNER JOIN fivenet_user AS colleague.*LEFT JOIN fivenet_job_colleague_props AS colleague_props.*WHERE .*fivenet_user_jobs\.job = \?.*EXISTS .*fivenet_job_timeclock.*NOT \(EXISTS .*fivenet_job_timeclock.*;`).
		WithArgs("police", "police", "police", "police").
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(3)))

	got, err := store.CountInactiveEmployees(t.Context(), store.db, InactiveEmployeesQuery{
		Job:  "police",
		Days: 14,
	})
	require.NoError(t, err)
	require.Equal(t, int64(3), got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListInactiveEmployeesUsesUserJobsBase(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user_jobs.*INNER JOIN fivenet_user AS colleague.*LEFT JOIN fivenet_user_props.*LEFT JOIN fivenet_job_colleague_props AS colleague_props.*LEFT JOIN fivenet_files AS profile_picture.*WHERE .*fivenet_user_jobs\.job = \?.*EXISTS .*fivenet_job_timeclock.*NOT \(EXISTS .*fivenet_job_timeclock.*ORDER BY .*colleague\.job_grade ASC.*LIMIT \?.*;`).
		WithArgs("police", "police", "police", "police", int64(20), int64(0)).
		WillReturnRows(sqlmock.NewRows(nil))

	_, err := store.ListInactiveEmployees(t.Context(), store.db, InactiveEmployeesQuery{
		Job:  "police",
		Days: 14,
		Limit: 20,
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
