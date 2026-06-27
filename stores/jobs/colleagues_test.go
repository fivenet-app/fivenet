package jobsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/require"
)

func TestStoreCreateColleagueActivity(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_colleague_activity`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(
		t,
		store.CreateColleagueActivity(
			t.Context(),
			store.db,
			&colleaguesactivity.ColleagueActivity{Job: "police"},
		),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCountColleagueActivityJoinsTargetUserForAccessFilters(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	tTargetUser := table.FivenetUser.AS("target_user")

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT COUNT(DISTINCT colleague_activity.id) AS "data_count.total" FROM fivenet_job_colleague_activity AS colleague_activity INNER JOIN fivenet_user AS target_user ON`,
		)+
			`(?s).*`+
			regexp.QuoteMeta(
				`target_user.id = colleague_activity.target_user_id`,
			)+
			`(?s).*`+
			regexp.QuoteMeta(
				`WHERE (colleague_activity.job = ?) AND (target_user.id = ?);`,
			),
	).
		WithArgs("police", int32(113031)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	count, err := store.CountColleagueActivity(
		t.Context(),
		store.db,
		ListQuery{
			Job:   "police",
			Where: tTargetUser.ID.EQ(mysql.Int32(113031)),
		},
	)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
	require.NoError(t, mock.ExpectationsWereMet())
}
