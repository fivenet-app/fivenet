package syncstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/stretchr/testify/require"
)

func TestEndActiveJobTimeclocks(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(`(?s)UPDATE fivenet_job_timeclock SET .*spent_time = \(COALESCE\(` + "`spent_time`" + `, 0\) \+ CAST\(\(TIMESTAMPDIFF\(SECOND, ` + "`start_time`" + `, CURRENT_TIMESTAMP\) / 3600\) AS DECIMAL\(10,2\)\)\).*end_time = CURRENT_TIMESTAMP.*WHERE .*start_time IS NOT NULL.*end_time IS NULL.*;`).
		WillReturnResult(sqlmock.NewResult(0, 3))

	resp, err := store.EndActiveJobTimeclocks(t.Context(), &pbsync.EndActiveJobTimeclocksRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(3), resp.GetRowsAffected())
	require.NoError(t, mock.ExpectationsWereMet())
}
