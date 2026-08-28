package syncstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/stretchr/testify/require"
)

func testSyncUser() *syncdata.DataUser {
	return &syncdata.DataUser{
		UserId:     42,
		Identifier: "license:external-42",
		Firstname:  "External",
		Lastname:   new("User"),
		Job:        "police",
		JobGrade:   2,
	}
}

func TestHandleUsersDataSkipsInvalidIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		user *syncdata.DataUser
	}{
		{
			name: "nil user",
		},
		{
			name: "zero id",
			user: &syncdata.DataUser{Identifier: "license:zero"},
		},
		{
			name: "negative id",
			user: &syncdata.DataUser{UserId: -1, Identifier: "license:negative"},
		},
		{
			name: "empty identifier",
			user: &syncdata.DataUser{UserId: 42},
		},
		{
			name: "whitespace identifier",
			user: &syncdata.DataUser{UserId: 42, Identifier: " \t"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, mock := newTestStore(t)

			rows, err := store.handleUsersData(t.Context(), []*syncdata.DataUser{tt.user})
			require.NoError(t, err)
			require.Zero(t, rows)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReconcileUserRowInsertsExternalIdentity(t *testing.T) {
	store, mock := newTestStore(t)
	user := testSyncUser()

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}))
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}))
	mock.ExpectExec(`(?s)INSERT INTO fivenet_user`).
		WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectCommit()

	tx, err := store.db.Begin()
	require.NoError(t, err)
	rows, err := store.reconcileUserRow(t.Context(), tx, user)
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestReconcileUserRowUpdatesExternalIDMatch(t *testing.T) {
	store, mock := newTestStore(t)
	user := testSyncUser()

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}).AddRow(42, "license:external-42"))
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}).AddRow(42, "license:external-42"))
	mock.ExpectExec(`(?s)UPDATE fivenet_user`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := store.db.Begin()
	require.NoError(t, err)
	rows, err := store.reconcileUserRow(t.Context(), tx, user)
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestReconcileUserRowRekeysIdentifierMatch(t *testing.T) {
	store, mock := newTestStore(t)
	user := testSyncUser()

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}))
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}).AddRow(7, "license:external-42"))
	mock.ExpectExec(`(?s)UPDATE fivenet_user`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := store.db.Begin()
	require.NoError(t, err)
	rows, err := store.reconcileUserRow(t.Context(), tx, user)
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestReconcileUserRowExternalIDWinsIdentityConflict(t *testing.T) {
	store, mock := newTestStore(t)
	user := testSyncUser()

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}).AddRow(42, "license:old-42"))
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user.*FOR UPDATE`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}).
			AddRow(7, "license:external-42"))
	mock.ExpectExec(`(?s)DELETE FROM fivenet_user`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)UPDATE fivenet_user`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := store.db.Begin()
	require.NoError(t, err)
	rows, err := store.reconcileUserRow(t.Context(), tx, user)
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}
