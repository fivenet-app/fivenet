package jobsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestStoreCountConductEntries(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(deletedFilterQueryRegex()).
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

func TestStoreCountConductEntriesFiltersDeletedByDefault(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(deletedFilterQueryRegex()).
		WithArgs("police").
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))

	total, err := store.CountConductEntries(
		t.Context(),
		store.db,
		ConductQuery{Job: "police", AllAccess: true, ShowDrafts: true, IncludeDeleted: false},
	)
	require.NoError(t, err)
	require.Equal(t, int64(0), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetConductEntryHonorsDeletedVisibility(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(deletedFilterQueryRegex()).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	entry, err := store.GetConductEntry(t.Context(), store.db, 42, false)
	require.NoError(t, err)
	require.Nil(t, entry)

	require.NoError(t, mock.ExpectationsWereMet())
}

func deletedFilterQueryRegex() string {
	return "(?s).*" + regexp.QuoteMeta("FROM fivenet_job_conduct AS conduct_entry") +
		".*" + regexp.QuoteMeta("conduct_entry.deleted_at IS NULL") + ".*"
}
