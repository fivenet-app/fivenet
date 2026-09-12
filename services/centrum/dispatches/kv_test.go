package dispatches

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
)

func TestScheduleProjectionCleanup(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	idleKV, err := js.CreateOrUpdateKeyValue(t.Context(), jetstream.KeyValueConfig{
		Bucket:         "test_dispatch_projection_cleanup",
		Storage:        jetstream.MemoryStorage,
		History:        1,
		LimitMarkerTTL: InactiveLimitMarkerTTL,
	})
	require.NoError(t, err)

	dispatches := &DispatchDB{idleKV: idleKV}
	createdAt := timestamp.New(time.Now())
	require.NoError(t, dispatches.ScheduleProjectionCleanup(t.Context(), 42, createdAt))
	require.NoError(t, dispatches.ScheduleProjectionCleanup(t.Context(), 42, createdAt))

	entry, err := idleKV.Get(t.Context(), projectionCleanupKey(42))
	require.NoError(t, err)
	require.NotNil(t, entry)
}
