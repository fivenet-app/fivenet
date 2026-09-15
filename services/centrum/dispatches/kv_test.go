package dispatches

import (
	"context"
	"errors"
	"testing"
	"time"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type failingAssignmentTimerKV struct {
	jetstream.KeyValue

	err error
}

func (s failingAssignmentTimerKV) Create(
	context.Context,
	string,
	[]byte,
	...jetstream.KVCreateOpt,
) (uint64, error) {
	return 0, s.err
}

func (s failingAssignmentTimerKV) Delete(context.Context, string, ...jetstream.KVDeleteOpt) error {
	return s.err
}

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

	entry, err := idleKV.Get(t.Context(), projectionCleanupKey(42))
	require.NoError(t, err)
	require.NotNil(t, entry)
	firstRevision := entry.Revision()

	// A later projection update must not extend the original cleanup deadline.
	require.NoError(t, dispatches.ScheduleProjectionCleanup(
		t.Context(),
		42,
		timestamp.New(time.Now().Add(time.Hour)),
	))
	entry, err = idleKV.Get(t.Context(), projectionCleanupKey(42))
	require.NoError(t, err)
	assert.Equal(t, firstRevision, entry.Revision())
}

func TestCancelScheduledCleanupRemovesProjectionEntry(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	idleKV, err := js.CreateOrUpdateKeyValue(t.Context(), jetstream.KeyValueConfig{
		Bucket:         "test_dispatch_projection_cancel",
		Storage:        jetstream.MemoryStorage,
		History:        1,
		LimitMarkerTTL: InactiveLimitMarkerTTL,
	})
	require.NoError(t, err)

	dispatches := &DispatchDB{idleKV: idleKV}
	require.NoError(t, dispatches.ScheduleCleanup(t.Context(), 42))
	require.NoError(t, dispatches.CancelScheduledCleanup(t.Context(), 42))

	_, err = idleKV.Get(t.Context(), cleanupKey(42))
	assert.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}

func TestAssignmentExpirationTimerCanBeScheduledAndCancelled(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	idleKV, err := js.CreateOrUpdateKeyValue(t.Context(), jetstream.KeyValueConfig{
		Bucket:         "test_dispatch_assignment_timer",
		Storage:        jetstream.MemoryStorage,
		LimitMarkerTTL: InactiveLimitMarkerTTL,
	})
	require.NoError(t, err)

	dispatches := &DispatchDB{idleKV: idleKV}
	require.NoError(
		t,
		dispatches.ScheduleAssignmentExpiration(t.Context(), 42, 7, time.Now().Add(time.Minute)),
	)

	_, err = idleKV.Get(t.Context(), assignmentExpirationKey(42, 7))
	require.NoError(t, err)

	require.NoError(t, dispatches.CancelAssignmentExpiration(t.Context(), 42, 7))
	_, err = idleKV.Get(t.Context(), assignmentExpirationKey(42, 7))
	assert.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}

func TestAssignmentTimerFailureIsBestEffort(t *testing.T) {
	t.Parallel()

	timerErr := errors.New("KV unavailable")
	dispatches := &DispatchDB{
		idleKV: failingAssignmentTimerKV{err: timerErr},
		logger: zap.NewNop(),
	}

	// The DB transaction and projection update run before this best-effort
	// operation. A timer outage must therefore only be logged.
	dispatches.updateAssignmentExpirationTimers(
		t.Context(),
		42,
		nil,
		map[int64]*centrumunits.Unit{7: {Id: 7}},
		time.Now().Add(time.Minute),
	)
}
