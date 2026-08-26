package syncers

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccountsSyncStampsAttemptAndSuccess(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
	})

	mock.ExpectQuery("SELECT license FROM accounts").
		WillReturnRows(sqlmock.NewRows([]string{"license"}))

	state := &dbsyncconfig.TableSyncState{}
	syncer := NewAccountsSync(&Syncer{
		logger: zap.NewNop(),
		db:     db,
		cfg: &dbsyncconfig.DBSyncConfig{
			Limits: dbsyncconfig.SyncLimits{
				Accounts: 10,
			},
			Tables: dbsyncconfig.DBSyncSourceTables{
				Accounts: dbsyncconfig.AccountsTable{
					DBSyncTable: dbsyncconfig.DBSyncTable{
						Query: new("SELECT license FROM accounts $whereCondition LIMIT $limit"),
					},
				},
			},
		},
	}, state)

	fetched, err := syncer.Sync(t.Context())
	require.NoError(t, err)
	require.Zero(t, fetched)
	require.NotNil(t, state.GetLastAttemptAt())
	require.NotNil(t, state.GetLastSyncedAt())
	require.Nil(t, state.GetLastCheck())
	require.Nil(t, state.GetLastID())
	require.Nil(t, state.GetLastError())
}

func TestAccountsSyncErrorRecordsLastError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
	})

	mock.ExpectQuery("SELECT license FROM accounts").
		WillReturnError(errors.New("boom"))

	state := &dbsyncconfig.TableSyncState{}
	syncer := NewAccountsSync(&Syncer{
		logger: zap.NewNop(),
		db:     db,
		cfg: &dbsyncconfig.DBSyncConfig{
			Limits: dbsyncconfig.SyncLimits{
				Accounts: 10,
			},
			Tables: dbsyncconfig.DBSyncSourceTables{
				Accounts: dbsyncconfig.AccountsTable{
					DBSyncTable: dbsyncconfig.DBSyncTable{
						Query: new("SELECT license FROM accounts $whereCondition LIMIT $limit"),
					},
				},
			},
		},
	}, state)

	_, err = syncer.Sync(t.Context())
	require.Error(t, err)
	require.NotNil(t, state.GetLastAttemptAt())
	require.Nil(t, state.GetLastSyncedAt())
	require.NotNil(t, state.GetLastError())
	require.Equal(t, "failed to query accounts. jet: boom", *state.GetLastError())
}
