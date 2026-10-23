package dispatches

import (
	"context"
	"testing"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPutStatusIfActiveDoesNotRegressTerminalStatus(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	statusStore, err := store.New[
		centrumdispatches.DispatchStatus,
		*centrumdispatches.DispatchStatus,
	](ctx, zap.NewNop(), js, "status_guard")
	require.NoError(t, err)

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, statusStore.Start(storeCtx, true))
	require.NoError(t, statusStore.Put(ctx, "7", &centrumdispatches.DispatchStatus{
		DispatchId: 7,
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	}))

	db := &DispatchDB{statusStore: statusStore}
	require.NoError(t, db.putStatusIfActive(ctx, "7", &centrumdispatches.DispatchStatus{
		DispatchId: 7,
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED,
	}))

	actual, err := statusStore.Get("7")
	require.NoError(t, err)
	require.True(t, centrumutils.IsStatusDispatchComplete(actual.GetStatus()))
	require.Equal(t, centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED, actual.GetStatus())
}

func TestPutStatusIfNewerDoesNotRegressDatabaseSnapshot(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	statusStore, err := store.New[
		centrumdispatches.DispatchStatus,
		*centrumdispatches.DispatchStatus,
	](ctx, zap.NewNop(), js, "status_snapshot_guard")
	require.NoError(t, err)

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, statusStore.Start(storeCtx, true))
	require.NoError(t, statusStore.Put(ctx, "7", &centrumdispatches.DispatchStatus{
		Id:         2,
		DispatchId: 7,
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	}))

	db := &DispatchDB{statusStore: statusStore}
	require.NoError(t, db.putStatusIfNewer(ctx, "7", &centrumdispatches.DispatchStatus{
		Id:         1,
		DispatchId: 7,
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED,
	}))

	actual, err := statusStore.Get("7")
	require.NoError(t, err)
	require.Equal(t, int64(2), actual.GetId())
	require.Equal(t, centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED, actual.GetStatus())
}

func TestPutStatusIfNewerPreservesExistingWhenSnapshotHasNoStatus(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	statusStore, err := store.New[
		centrumdispatches.DispatchStatus,
		*centrumdispatches.DispatchStatus,
	](ctx, zap.NewNop(), js, "status_snapshot_nil_guard")
	require.NoError(t, err)

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, statusStore.Start(storeCtx, true))
	require.NoError(t, statusStore.Put(ctx, "7", &centrumdispatches.DispatchStatus{
		Id:         2,
		DispatchId: 7,
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	}))

	db := &DispatchDB{statusStore: statusStore}
	require.NoError(t, db.putStatusIfNewer(ctx, "7", nil))

	actual, err := statusStore.Get("7")
	require.NoError(t, err)
	require.Equal(t, int64(2), actual.GetId())
}
