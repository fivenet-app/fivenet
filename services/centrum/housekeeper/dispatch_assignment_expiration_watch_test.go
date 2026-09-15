package housekeeper

import (
	"context"
	"sync"
	"testing"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	centrummetrics "github.com/fivenet-app/fivenet/v2026/services/centrum/metrics"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type assignmentExpirationSourceStub struct {
	kv jetstream.KeyValue

	mu       sync.Mutex
	dispatch *centrumdispatches.Dispatch
	getCalls chan struct{}
	removals chan []int64
}

func (s *assignmentExpirationSourceStub) IdleStore() jetstream.KeyValue {
	return s.kv
}

func (s *assignmentExpirationSourceStub) Get(_ context.Context, _ int64) (*centrumdispatches.Dispatch, error) {
	select {
	case s.getCalls <- struct{}{}:
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dispatch, nil
}

func (s *assignmentExpirationSourceStub) UpdateAssignments(
	_ context.Context,
	_ *string,
	_ *int32,
	_ int64,
	_ []int64,
	toRemove []int64,
	_ time.Time,
) error {
	s.mu.Lock()
	s.dispatch.Units = removeDispatchAssignments(s.dispatch.GetUnits(), toRemove)
	s.mu.Unlock()

	select {
	case s.removals <- append([]int64(nil), toRemove...):
	default:
	}
	return nil
}

func removeDispatchAssignments(
	assignments []*centrumdispatches.DispatchAssignment,
	toRemove []int64,
) []*centrumdispatches.DispatchAssignment {
	remaining := make([]*centrumdispatches.DispatchAssignment, 0, len(assignments))
	for _, assignment := range assignments {
		remove := false
		for _, unitID := range toRemove {
			if assignment.GetUnitId() == unitID {
				remove = true
				break
			}
		}
		if !remove {
			remaining = append(remaining, assignment)
		}
	}

	return remaining
}

func newAssignmentExpirationWatcherTest(
	t *testing.T,
	dispatch *centrumdispatches.Dispatch,
) (*assignmentExpirationSourceStub, context.CancelFunc, <-chan error) {
	t.Helper()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	kv, err := js.CreateOrUpdateKeyValue(t.Context(), jetstream.KeyValueConfig{
		Bucket:         "test_assignment_expiration",
		Storage:        jetstream.MemoryStorage,
		LimitMarkerTTL: time.Second,
	})
	require.NoError(t, err)

	source := &assignmentExpirationSourceStub{
		kv:       kv,
		dispatch: dispatch,
		getCalls: make(chan struct{}, 4),
		removals: make(chan []int64, 4),
	}
	h := &Housekeeper{
		logger:                     zap.NewNop(),
		metrics:                    centrummetrics.Get(),
		assignmentExpirationSource: source,
	}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- h.dispatchAssignmentExpirationWatcher(ctx) }()

	t.Cleanup(func() {
		cancel()
		require.ErrorIs(t, <-done, context.Canceled)
	})

	return source, cancel, done
}

func TestDispatchAssignmentExpirationWatcherRemovesOnlyExpiredAssignment(t *testing.T) {
	expiredAt := timestamp.New(time.Now().Add(-time.Second))
	futureAt := timestamp.New(time.Now().Add(time.Minute))
	source, _, _ := newAssignmentExpirationWatcherTest(t, &centrumdispatches.Dispatch{
		Id: 42,
		Units: []*centrumdispatches.DispatchAssignment{
			{DispatchId: 42, UnitId: 7, ExpiresAt: expiredAt},
			{DispatchId: 42, UnitId: 8, ExpiresAt: futureAt},
		},
	})

	_, err := source.kv.Create(t.Context(), "assignment.42.7", nil, jetstream.KeyTTL(time.Second))
	require.NoError(t, err)

	select {
	case removed := <-source.removals:
		assert.Equal(t, []int64{7}, removed)
	case <-time.After(5 * time.Second):
		t.Fatal("assignment TTL expiry was not handled")
	}

	source.mu.Lock()
	defer source.mu.Unlock()
	require.Len(t, source.dispatch.GetUnits(), 1)
	assert.Equal(t, int64(8), source.dispatch.GetUnits()[0].GetUnitId())
}

func TestDispatchAssignmentExpirationWatcherIgnoresCancelledAndStaleTimers(t *testing.T) {
	for _, test := range []struct {
		name        string
		expiresAt   *timestamp.Timestamp
		cancelTimer bool
	}{
		{
			name:        "cancelled timer",
			expiresAt:   timestamp.New(time.Now().Add(-time.Second)),
			cancelTimer: true,
		},
		{
			name:      "accepted assignment",
			expiresAt: nil,
		},
		{
			name:      "reassigned assignment",
			expiresAt: timestamp.New(time.Now().Add(time.Minute)),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			source, _, _ := newAssignmentExpirationWatcherTest(t, &centrumdispatches.Dispatch{
				Id:    42,
				Units: []*centrumdispatches.DispatchAssignment{{DispatchId: 42, UnitId: 7, ExpiresAt: test.expiresAt}},
			})

			_, err := source.kv.Create(t.Context(), "assignment.42.7", nil, jetstream.KeyTTL(time.Second))
			require.NoError(t, err)
			if test.cancelTimer {
				require.NoError(t, source.kv.Delete(t.Context(), "assignment.42.7"))
			}

			if !test.cancelTimer {
				select {
				case <-source.getCalls:
				case <-time.After(5 * time.Second):
					t.Fatal("assignment TTL expiry was not observed")
				}
			} else {
				time.Sleep(100 * time.Millisecond)
			}

			select {
			case removed := <-source.removals:
				t.Fatalf("unexpected assignment removal: %v", removed)
			case <-time.After(100 * time.Millisecond):
			}
		})
	}
}
