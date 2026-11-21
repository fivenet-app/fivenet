package citizens

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

func TestLockCitizenParent(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectBegin()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user WHERE .*id = \? FOR UPDATE;`).
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int32(42)))
	mock.ExpectCommit()

	require.NoError(t, lockCitizenParent(t.Context(), tx, 42))
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestLockCitizenParentMissing(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectBegin()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user WHERE .*id = \? FOR UPDATE;`).
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectRollback()

	err = lockCitizenParent(t.Context(), tx, 42)
	require.ErrorIs(t, err, qrm.ErrNoRows)
	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}
