package units

import (
	"context"
	"testing"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPutStatusIfNewerDoesNotRegressDatabaseSnapshot(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	statusStore, err := store.New[
		centrumunits.UnitStatus,
		*centrumunits.UnitStatus,
	](ctx, zap.NewNop(), js, "unit_status_snapshot_guard")
	require.NoError(t, err)

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, statusStore.Start(storeCtx, true))
	require.NoError(t, statusStore.Put(ctx, "7", &centrumunits.UnitStatus{
		Id:     2,
		UnitId: 7,
		Status: centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE,
	}))

	db := &UnitDB{statusStore: statusStore}
	require.NoError(t, db.putStatusIfNewer(ctx, "7", &centrumunits.UnitStatus{
		Id:     1,
		UnitId: 7,
		Status: centrumunits.StatusUnit_STATUS_UNIT_BUSY,
	}))

	actual, err := statusStore.Get("7")
	require.NoError(t, err)
	require.Equal(t, int64(2), actual.GetId())
	require.Equal(t, centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE, actual.GetStatus())
}
