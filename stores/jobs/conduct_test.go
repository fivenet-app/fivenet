package jobs

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestStoreCountConductEntries(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_conduct AS conduct_entry`)).
		WithArgs("police").
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))

	total, err := store.CountConductEntries(
		t.Context(),
		store.db,
		ConductQuery{Job: "police", AllAccess: true, ShowDrafts: true},
	)
	require.NoError(t, err)
	require.Equal(t, int64(0), total)
	require.NoError(t, mock.ExpectationsWereMet())
}
