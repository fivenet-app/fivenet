package vehiclesstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

func TestLockVehicleParent(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectBegin()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_vehicle WHERE .*plate = \? FOR UPDATE;`).
		WithArgs("ABC DEF1").
		WillReturnRows(sqlmock.NewRows([]string{"plate"}).AddRow("ABC DEF1"))
	mock.ExpectCommit()

	require.NoError(t, lockVehicleParent(t.Context(), tx, "ABC DEF1"))
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestLockVehicleParentMissing(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectBegin()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_vehicle WHERE .*plate = \? FOR UPDATE;`).
		WithArgs("ABC DEF1").
		WillReturnRows(sqlmock.NewRows([]string{"plate"}))
	mock.ExpectRollback()

	err = lockVehicleParent(t.Context(), tx, "ABC DEF1")
	require.ErrorIs(t, err, qrm.ErrNoRows)
	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}
