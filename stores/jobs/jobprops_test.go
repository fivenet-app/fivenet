package jobs

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestStoreSetMOTD(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_props`) + `(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.SetMOTD(t.Context(), store.db, "police", "hello"))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetMOTD(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_props`)).
		WithArgs("police", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"get_motd_response.motd"}).AddRow("hello"))

	motd, err := store.GetMOTD(t.Context(), store.db, "police")
	require.NoError(t, err)
	require.Equal(t, "hello", motd)
	require.NoError(t, mock.ExpectationsWereMet())
}
