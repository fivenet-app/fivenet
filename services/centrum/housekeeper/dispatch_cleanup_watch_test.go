package housekeeper

import (
	"context"
	"testing"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	centrummetrics "github.com/fivenet-app/fivenet/v2026/services/centrum/metrics"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDispatchCleanupWatcherOnlyActsOnPurge(t *testing.T) {
	t.Parallel()

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

	kv, err := js.CreateOrUpdateKeyValue(t.Context(), jetstream.KeyValueConfig{
		Bucket:         "test_dispatch_cleanup_watch",
		Storage:        jetstream.MemoryStorage,
		LimitMarkerTTL: time.Second,
	})
	require.NoError(t, err)

	lookups := make(chan struct{}, 2)
	deletions := make(chan int64, 2)
	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
		dispatchLifecycle: fakeDispatchLifecycle{
			get: func(context.Context, int64) (*centrumdispatches.Dispatch, error) {
				lookups <- struct{}{}
				return &centrumdispatches.Dispatch{
					Id: 42,
					Status: &centrumdispatches.DispatchStatus{
						Status: centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
					},
				}, nil
			},
			delete: func(_ context.Context, id int64, _ bool) error {
				deletions <- id
				return nil
			},
		},
	}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- h.watchDispatchCleanup(ctx, kv) }()
	t.Cleanup(func() {
		cancel()
		require.ErrorIs(t, <-done, context.Canceled)
	})
	// Watch establishes an ephemeral consumer asynchronously. Wait until it is
	// attached so the test cannot publish the timer operations before the
	// watcher is able to observe them.
	require.Eventually(t, func() bool {
		stream, err := js.Stream(t.Context(), "KV_test_dispatch_cleanup_watch")
		if err != nil {
			return false
		}
		for range stream.ConsumerNames(t.Context()).Name() {
			return true
		}
		return false
	}, time.Second, 10*time.Millisecond)

	// A normal cancellation writes DELETE and must not look up or delete the
	// dispatch projection.
	_, err = kv.Create(t.Context(), "cleanup.42", nil, jetstream.KeyTTL(time.Minute))
	require.NoError(t, err)
	require.NoError(t, kv.Delete(t.Context(), "cleanup.42"))
	assert.Never(t, func() bool {
		select {
		case <-lookups:
			return true
		default:
			return false
		}
	}, 250*time.Millisecond, 10*time.Millisecond)

	// A purge marker is the operation emitted by KeyTTL expiry and performs the
	// completed-projection cleanup.
	_, err = kv.Create(t.Context(), "cleanup.42", nil, jetstream.KeyTTL(time.Minute))
	require.NoError(t, err)
	require.NoError(t, kv.Purge(t.Context(), "cleanup.42"))

	select {
	case id := <-deletions:
		assert.Equal(t, int64(42), id)
	case <-time.After(time.Second):
		t.Fatal("cleanup purge was not handled")
	}
}
