package sync

import (
	"testing"
	"time"

	dbsyncstate "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/dbsync"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/stretchr/testify/require"
)

func TestAuthFuncOverrideMarksAPITokenAuth(t *testing.T) {
	t.Parallel()

	server := &Server{
		tokens: []string{"sync-token"},
	}

	ctx := auth.SetTokenInGRPCContext(t.Context(), "sync-token")
	out, err := server.AuthFuncOverride(ctx, "/services.sync.SyncService/SendData")
	require.NoError(t, err)

	kind, ok := auth.GetAuthKindFromContext(out)
	require.True(t, ok)
	require.Equal(t, auth.AuthKindAPIToken, kind)
}

func TestGetDBSyncStatusMirrorsNestedSyncState(t *testing.T) {
	t.Parallel()

	lastCheck := time.Date(2026, 8, 26, 12, 34, 56, 0, time.UTC)
	lastSyncedAt := time.Date(2026, 8, 26, 12, 40, 0, 0, time.UTC)
	lastAttemptAt := time.Date(2026, 8, 26, 12, 41, 0, 0, time.UTC)
	lastError := "failed to fetch vehicles"
	lastID := "42"

	server := &Server{
		cfg: &config.Config{
			Sync: config.Sync{Enabled: true},
		},
	}
	server.lastDBSyncState.Store(&pbsync.ClientSyncState{
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
				LastError:     nil,
				Checkpoint:    nil,
			},
		},
	})

	status := server.GetDBSyncStatus()
	require.NotNil(t, status)
	require.NotNil(t, status.GetSyncState())
	require.Len(t, status.GetSyncState().GetTables(), 2)
	require.True(t, status.GetEnabled())

	vehicles := status.GetSyncState().GetTables()[0]
	require.Equal(t, "vehicles", vehicles.GetTable())
	require.NotNil(t, vehicles.GetCheckpoint())
	require.True(t, vehicles.GetCheckpoint().GetLastCheck().AsTime().Equal(lastCheck))
	require.Equal(t, "42", vehicles.GetCheckpoint().GetLastId())
	require.True(t, vehicles.GetLastSyncedAt().AsTime().Equal(lastSyncedAt))
	require.True(t, vehicles.GetLastAttemptAt().AsTime().Equal(lastAttemptAt))
	require.Equal(t, lastError, vehicles.GetLastError())
	require.False(t, status.GetStreamConnected())

	accounts := status.GetSyncState().GetTables()[1]
	require.Equal(t, "accounts", accounts.GetTable())
	require.Nil(t, accounts.GetCheckpoint())
	require.True(t, accounts.GetLastSyncedAt().AsTime().Equal(lastSyncedAt))
	require.True(t, accounts.GetLastAttemptAt().AsTime().Equal(lastAttemptAt))
	require.Empty(t, accounts.GetLastError())
}

func TestGetDBSyncStatusWithoutStateKeepsSyncStateNil(t *testing.T) {
	t.Parallel()

	server := &Server{
		cfg: &config.Config{Sync: config.Sync{Enabled: true}},
	}

	status := server.GetDBSyncStatus()
	require.NotNil(t, status)
	require.Nil(t, status.GetSyncState())
	require.False(t, status.GetStreamConnected())
}

func TestGetDBSyncStatusReflectsStreamLifecycle(t *testing.T) {
	t.Parallel()

	server := &Server{
		cfg: &config.Config{Sync: config.Sync{Enabled: true}},
	}

	require.False(t, server.GetDBSyncStatus().GetStreamConnected())

	server.markDBSyncStreamConnected()
	require.True(t, server.GetDBSyncStatus().GetStreamConnected())

	server.markDBSyncStreamDisconnected()
	require.False(t, server.GetDBSyncStatus().GetStreamConnected())
}
