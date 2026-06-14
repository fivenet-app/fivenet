package jobs

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
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
