package statusmapper

import (
	"testing"
	"time"

	dbsyncstate "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/dbsync"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/stretchr/testify/require"
)

func TestFromRuntimeStateBuildsClientSyncState(t *testing.T) {
	t.Parallel()

	lastCheck := time.Date(2026, 8, 26, 12, 34, 56, 0, time.UTC)
	lastSyncedAt := time.Date(2026, 8, 26, 12, 40, 0, 0, time.UTC)
	lastAttemptAt := time.Date(2026, 8, 26, 12, 41, 0, 0, time.UTC)
	lastError := "failed to fetch accounts"
	lastID := "42"

	snapshot := FromRuntimeState(&dbsyncconfig.State{
		Accounts: &dbsyncconfig.TableSyncState{
			LastSyncedAt:  &lastSyncedAt,
			LastAttemptAt: &lastAttemptAt,
			LastError:     &lastError,
		},
		Vehicles: &dbsyncconfig.TableSyncState{
			LastCheck:     &lastCheck,
			LastID:        &lastID,
			LastSyncedAt:  &lastSyncedAt,
			LastAttemptAt: &lastAttemptAt,
		},
	})
	require.NotNil(t, snapshot)

	out := snapshot.ToClientSyncState()
	require.NotNil(t, out)
	require.Len(t, out.GetTables(), 2)

	accounts := out.GetTables()[0]
	require.Equal(t, "accounts", accounts.GetTable())
	require.Nil(t, accounts.GetCheckpoint())
	require.True(t, accounts.GetLastSyncedAt().AsTime().Equal(lastSyncedAt))
	require.True(t, accounts.GetLastAttemptAt().AsTime().Equal(lastAttemptAt))
	require.Equal(t, lastError, accounts.GetLastError())

	vehicles := out.GetTables()[1]
	require.Equal(t, "vehicles", vehicles.GetTable())
	require.NotNil(t, vehicles.GetCheckpoint())
	require.True(t, vehicles.GetCheckpoint().GetLastCheck().AsTime().Equal(lastCheck))
	require.Equal(t, "42", vehicles.GetCheckpoint().GetLastId())
	require.True(t, vehicles.GetLastSyncedAt().AsTime().Equal(lastSyncedAt))
	require.True(t, vehicles.GetLastAttemptAt().AsTime().Equal(lastAttemptAt))
}

func TestFromStreamRequestSharesTablesWithSettingsStatus(t *testing.T) {
	t.Parallel()

	lastCheck := time.Date(2026, 8, 26, 12, 34, 56, 0, time.UTC)
	lastSyncedAt := time.Date(2026, 8, 26, 12, 40, 0, 0, time.UTC)
	lastAttemptAt := time.Date(2026, 8, 26, 12, 41, 0, 0, time.UTC)
	lastError := "failed to fetch vehicles"
	lastID := "99"

	snapshot := FromClientSyncState(&pbsync.ClientSyncState{
		Tables: []*dbsyncstate.DBSyncTableSyncState{
			{
				Table: "vehicles",
				Checkpoint: &dbsyncstate.DBSyncCheckpoint{
					LastCheck: timestamp.New(lastCheck),
					LastId:    &lastID,
				},
				LastSyncedAt:  timestamp.New(lastSyncedAt),
				LastAttemptAt: timestamp.New(lastAttemptAt),
				LastError:     &lastError,
			},
			{
				Table:         "accounts",
				LastSyncedAt:  timestamp.New(lastSyncedAt),
				LastAttemptAt: timestamp.New(lastAttemptAt),
			},
		},
	})
	require.NotNil(t, snapshot)

	settingsState := snapshot.ToSettingsSyncState()
	require.NotNil(t, settingsState)
	status := &settings.DBSyncStatus{SyncState: settingsState}
	require.Len(t, status.GetSyncState().GetTables(), 2)

	vehicles := status.GetSyncState().GetTables()[0]
	require.Equal(t, "vehicles", vehicles.GetTable())
	require.NotNil(t, vehicles.GetCheckpoint())
	require.True(t, vehicles.GetCheckpoint().GetLastCheck().AsTime().Equal(lastCheck))
	require.Equal(t, "99", vehicles.GetCheckpoint().GetLastId())
	require.True(t, vehicles.GetLastSyncedAt().AsTime().Equal(lastSyncedAt))
	require.True(t, vehicles.GetLastAttemptAt().AsTime().Equal(lastAttemptAt))
	require.True(t, vehicles.HasLastError())
	require.Equal(t, lastError, vehicles.GetLastError())

	accounts := status.GetSyncState().GetTables()[1]
	require.Equal(t, "accounts", accounts.GetTable())
	require.Nil(t, accounts.GetCheckpoint())
	require.True(t, accounts.GetLastSyncedAt().AsTime().Equal(lastSyncedAt))
	require.True(t, accounts.GetLastAttemptAt().AsTime().Equal(lastAttemptAt))
	require.False(t, accounts.HasLastError())
	require.Empty(t, accounts.GetLastError())
}

func TestNilInputsReturnNilSnapshot(t *testing.T) {
	t.Parallel()

	require.Nil(t, FromRuntimeState(nil))
	require.Nil(t, FromClientSyncState(nil))
	require.Nil(t, (*Snapshot)(nil).ToClientSyncState())
	require.Nil(t, (*Snapshot)(nil).ToSettingsSyncState())
}
