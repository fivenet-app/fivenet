package settingsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	"github.com/stretchr/testify/require"
)

func TestStoreSetJobPropsUpsertsProps(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_props`) + `(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.SetJobProps(t.Context(), &jobsprops.JobProps{Job: "police"}))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteJobPropsSoftDeletesRow(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_props SET deleted_at = CURRENT_TIMESTAMP WHERE fivenet_job_props.job = ? LIMIT ?;`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.DeleteJobProps(t.Context(), "police"))
	require.NoError(t, mock.ExpectationsWereMet())
}
