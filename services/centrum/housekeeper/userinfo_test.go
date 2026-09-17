package housekeeper

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	testnats "github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	pkguserinfo "github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	centrummetrics "github.com/fivenet-app/fivenet/v2026/services/centrum/metrics"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type userInfoReconcileTestMsg struct {
	data       []byte
	metadata   *jetstream.MsgMetadata
	acked      bool
	terminated bool
	nakDelay   time.Duration
}

func (m *userInfoReconcileTestMsg) Metadata() (*jetstream.MsgMetadata, error) {
	return m.metadata, nil
}
func (m *userInfoReconcileTestMsg) Data() []byte                    { return m.data }
func (m *userInfoReconcileTestMsg) Headers() nats.Header            { return nil }
func (m *userInfoReconcileTestMsg) Subject() string                 { return "userinfo.1.changes" }
func (m *userInfoReconcileTestMsg) Reply() string                   { return "" }
func (m *userInfoReconcileTestMsg) Ack() error                      { m.acked = true; return nil }
func (m *userInfoReconcileTestMsg) DoubleAck(context.Context) error { return nil }
func (m *userInfoReconcileTestMsg) Nak() error                      { return nil }
func (m *userInfoReconcileTestMsg) NakWithDelay(delay time.Duration) error {
	m.nakDelay = delay
	return nil
}
func (m *userInfoReconcileTestMsg) InProgress() error           { return nil }
func (m *userInfoReconcileTestMsg) Term() error                 { m.terminated = true; return nil }
func (m *userInfoReconcileTestMsg) TermWithReason(string) error { m.terminated = true; return nil }

var _ jetstream.Msg = (*userInfoReconcileTestMsg)(nil)

func userInfoReconcileMessage(t *testing.T, userID int32, job string) *userInfoReconcileTestMsg {
	t.Helper()
	event := &pbuserinfo.UserInfoChanged{UserId: userID}
	event.SetJob(job)
	data, err := protojson.Marshal(event)
	require.NoError(t, err)
	return &userInfoReconcileTestMsg{data: data}
}

func userInfoRetrieverForTest(userID int32, job string) pkguserinfo.UserInfoRetriever {
	return pkguserinfo.NewMockUserInfoRetriever(map[int32]*pbuserinfo.UserInfo{
		userID: {UserId: userID, Job: job},
	})
}

func TestReconcileUserJobChangeRemovesOnlyCrossJobState(t *testing.T) {
	t.Parallel()

	unitChanges := make(chan struct {
		userID int32
		job    string
	}, 1)
	dispatcherRemovals := make(chan string, 1)
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		tracer:   noop.NewTracerProvider().Tracer("test"),
		unitUserState: fakeUnitUserState{
			reconcile: func(_ context.Context, userID int32, job string) (bool, error) {
				unitChanges <- struct {
					userID int32
					job    string
				}{userID: userID, job: job}
				return true, nil
			},
		},
		dispatcherUserState: fakeDispatcherUserState{
			rangeFn: func(fn func(string, *centrumdispatchers.Dispatchers) bool) {
				fn("ambulance", &centrumdispatchers.Dispatchers{
					Job:         "ambulance",
					Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
				})
				fn("police", &centrumdispatchers.Dispatchers{
					Job:         "police",
					Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
				})
			},
			set: func(_ context.Context, job string, userID int32, active bool) error {
				assert.Equal(t, int32(42), userID)
				assert.False(t, active)
				dispatcherRemovals <- job
				return nil
			},
		},
	}

	require.NoError(t, h.reconcileUserJobChange(t.Context(), 42, "police"))
	require.Equal(t, struct {
		userID int32
		job    string
	}{userID: 42, job: "police"}, <-unitChanges)
	assert.Equal(t, "ambulance", <-dispatcherRemovals)
}

func TestHandleUserInfoReconcileMessageAcknowledgesOnlyAfterSuccess(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{
			reconcile: func(context.Context, int32, string) (bool, error) {
				return false, nil
			},
		},
		dispatcherUserState: fakeDispatcherUserState{},
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.True(t, msg.acked)
	assert.Zero(t, msg.nakDelay)
	assert.False(t, msg.terminated)
}

func TestUserInfoReconcileConsumerRetainsEventAcrossLeadershipHandoff(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	js := testnats.NewServer(t, testnats.ServerOptions{InProcess: true}).GetJS()

	var reconciled atomic.Int32
	newHousekeeper := func() *Housekeeper {
		return &Housekeeper{
			logger:   zap.NewNop(),
			metrics:  centrummetrics.Get(),
			js:       js,
			userinfo: userInfoRetrieverForTest(42, "police"),
			unitUserState: fakeUnitUserState{
				reconcile: func(_ context.Context, userID int32, job string) (bool, error) {
					if userID == 42 && job == "police" {
						reconciled.Add(1)
					}
					return false, nil
				},
			},
			dispatcherUserState: fakeDispatcherUserState{},
		}
	}

	first := newHousekeeper()
	require.NoError(t, first.ensureUserInfoReconcileConsumer(ctx))
	firstLeaderCtx, stopFirstLeader := context.WithCancel(ctx)
	firstConsume, err := first.startUserInfoReconcileConsumer(firstLeaderCtx)
	require.NoError(t, err)

	// Simulate the elected leader stepping down before a canonical event is
	// published. The durable consumer must retain that event for its successor.
	stopFirstLeader()
	firstConsume.ctx.Stop()

	event := &pbuserinfo.UserInfoChanged{AccountId: 1, UserId: 42}
	event.SetJob("ambulance") // Deliberately stale: the handler reloads "police".
	data, err := protojson.Marshal(event)
	require.NoError(t, err)
	_, err = js.Publish(ctx, "userinfo.1.changes", data)
	require.NoError(t, err)

	second := newHousekeeper()
	require.NoError(t, second.ensureUserInfoReconcileConsumer(ctx))
	secondLeaderCtx, stopSecondLeader := context.WithCancel(ctx)
	t.Cleanup(stopSecondLeader)
	secondConsume, err := second.startUserInfoReconcileConsumer(secondLeaderCtx)
	require.NoError(t, err)
	t.Cleanup(secondConsume.ctx.Stop)

	require.Eventually(
		t,
		func() bool { return reconciled.Load() == 1 },
		time.Second,
		10*time.Millisecond,
	)
	time.Sleep(100 * time.Millisecond)
	require.Equal(
		t,
		int32(1),
		reconciled.Load(),
		"the durable event must be reconciled exactly once",
	)
}

func TestHandleUserInfoReconcileMessageRetriesFailure(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{
			reconcile: func(context.Context, int32, string) (bool, error) {
				return false, errors.New("database unavailable")
			},
		},
		dispatcherUserState: fakeDispatcherUserState{},
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.False(t, msg.acked)
	assert.Equal(t, userInfoReconcileRetryDelay, msg.nakDelay)
	assert.False(t, msg.terminated)
}

func TestHandleUserInfoReconcileMessageSlowsRepeatedFailure(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	msg.metadata = &jetstream.MsgMetadata{NumDelivered: userInfoReconcileSlowAfter}
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{
			reconcile: func(context.Context, int32, string) (bool, error) {
				return false, errors.New("database unavailable")
			},
		},
		dispatcherUserState: fakeDispatcherUserState{},
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.False(t, msg.acked)
	assert.Equal(t, userInfoReconcileSlowRetry, msg.nakDelay)
	assert.False(t, msg.terminated)
}

func TestHandleUserInfoReconcileMessageSlowsRepeatedLookupFailure(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	msg.metadata = &jetstream.MsgMetadata{NumDelivered: userInfoReconcileSlowAfter}
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(7, "police"),
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.False(t, msg.acked)
	assert.Equal(t, userInfoReconcileSlowRetry, msg.nakDelay)
	assert.False(t, msg.terminated)
}

func TestReconcileUserJobChangeUsesKnownDispatcherJobs(t *testing.T) {
	t.Parallel()

	dispatcherRemovals := make(chan string, 1)
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{
			reconcile: func(context.Context, int32, string) (bool, error) {
				return false, nil
			},
		},
		dispatcherUserState: fakeDispatcherUserState{
			rangeFn: func(func(string, *centrumdispatchers.Dispatchers) bool) {
				require.Fail(t, "dispatcher projections must not be scanned when jobs are known")
			},
			set: func(_ context.Context, job string, userID int32, active bool) error {
				assert.Equal(t, int32(42), userID)
				assert.False(t, active)
				dispatcherRemovals <- job
				return nil
			},
		},
	}

	require.NoError(t, h.reconcileUserJobChangeWithDispatcherJobs(
		t.Context(),
		42,
		"police",
		map[string]struct{}{"ambulance": {}, "police": {}},
	))
	assert.Equal(t, "ambulance", <-dispatcherRemovals)
}

func TestReconcileUserJobChangeRemovesAllOldDispatcherJobs(t *testing.T) {
	t.Parallel()

	removed := map[string]bool{}
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{
			reconcile: func(context.Context, int32, string) (bool, error) {
				return false, nil
			},
		},
		dispatcherUserState: fakeDispatcherUserState{
			set: func(_ context.Context, job string, userID int32, active bool) error {
				assert.Equal(t, int32(42), userID)
				assert.False(t, active)
				removed[job] = true
				return nil
			},
		},
	}

	require.NoError(t, h.reconcileUserJobChangeWithDispatcherJobs(
		t.Context(),
		42,
		"police",
		map[string]struct{}{
			"ambulance": {},
			"fire":      {},
			"police":    {},
		},
	))

	assert.Equal(t, map[string]bool{"ambulance": true, "fire": true}, removed)
}

func TestHandleUserInfoReconcileMessageUsesCurrentAuthoritativeJob(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	reconciledJob := ""
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "fire"),
		unitUserState: fakeUnitUserState{
			reconcile: func(_ context.Context, _ int32, job string) (bool, error) {
				reconciledJob = job
				return false, nil
			},
		},
		dispatcherUserState: fakeDispatcherUserState{},
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.Equal(t, "fire", reconciledJob)
	assert.True(t, msg.acked)
}

func TestHandleUserInfoReconcileMessageRetriesDispatcherCleanupFailure(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{
			reconcile: func(context.Context, int32, string) (bool, error) {
				return true, nil
			},
		},
		dispatcherUserState: fakeDispatcherUserState{
			rangeFn: func(fn func(string, *centrumdispatchers.Dispatchers) bool) {
				fn("ambulance", &centrumdispatchers.Dispatchers{
					Job:         "ambulance",
					Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
				})
			},
			set: func(context.Context, string, int32, bool) error {
				return errors.New("dispatcher store unavailable")
			},
		},
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.False(t, msg.acked)
	assert.Equal(t, userInfoReconcileRetryDelay, msg.nakDelay)
	assert.False(t, msg.terminated)
}

func TestHandleUserInfoReconcileMessageTerminatesWhenAuthoritativeJobIsEmpty(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, ""),
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.False(t, msg.acked)
	assert.Zero(t, msg.nakDelay)
	assert.True(t, msg.terminated)
}

func TestHandleUserInfoReconcileMessageUsesAuthoritativeJobAcrossStaleEvents(t *testing.T) {
	t.Parallel()

	jobs := make([]string, 0, 2)
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "fire"),
		unitUserState: fakeUnitUserState{
			reconcile: func(_ context.Context, _ int32, job string) (bool, error) {
				jobs = append(jobs, job)
				return false, nil
			},
		},
		dispatcherUserState: fakeDispatcherUserState{},
	}

	first := userInfoReconcileMessage(t, 42, "police")
	second := userInfoReconcileMessage(t, 42, "ambulance")
	h.handleUserInfoReconcileMessage(t.Context(), first)
	h.handleUserInfoReconcileMessage(t.Context(), second)

	assert.Equal(t, []string{"fire", "fire"}, jobs)
	assert.True(t, first.acked)
	assert.True(t, second.acked)
}

func TestHandleUserInfoReconcileMessageDoesNotAckAfterContextCancellation(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
	}
	h.handleUserInfoReconcileMessage(ctx, msg)

	assert.False(t, msg.acked)
	assert.Zero(t, msg.nakDelay)
	assert.False(t, msg.terminated)
}

func TestHandleUserInfoReconcileMessageTerminatesMalformedEvent(t *testing.T) {
	t.Parallel()

	msg := &userInfoReconcileTestMsg{data: []byte("not-json")}
	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)

	assert.False(t, msg.acked)
	assert.Zero(t, msg.nakDelay)
	assert.True(t, msg.terminated)
}

func TestHandleUserInfoReconcileMessageTerminatesInvalidEvent(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 0, "police")
	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)

	assert.False(t, msg.acked)
	assert.Zero(t, msg.nakDelay)
	assert.True(t, msg.terminated)
}

func TestReconcileUserInfoStateDiscoversAndDeduplicatesProjectionUsers(t *testing.T) {
	t.Parallel()

	var reconciled []int32
	var removedDispatchers []string
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{rangeFn: func(fn func(string, *centrumunits.Unit) bool) {
			fn("unit.1", &centrumunits.Unit{
				Users: []*centrumunits.UnitAssignment{{UserId: 42}},
			})
			fn("unit.2", &centrumunits.Unit{
				Users: []*centrumunits.UnitAssignment{{UserId: 42}},
			})
		}, reconcile: func(_ context.Context, userID int32, job string) (bool, error) {
			reconciled = append(reconciled, userID)
			assert.Equal(t, "police", job)
			return false, nil
		}},
		dispatcherUserState: fakeDispatcherUserState{
			rangeFn: func(fn func(string, *centrumdispatchers.Dispatchers) bool) {
				fn("ambulance", &centrumdispatchers.Dispatchers{
					Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
				})
				fn("fire", &centrumdispatchers.Dispatchers{
					Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
				})
				fn("police", &centrumdispatchers.Dispatchers{
					Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
				})
			},
			set: func(_ context.Context, job string, userID int32, active bool) error {
				assert.Equal(t, int32(42), userID)
				assert.False(t, active)
				removedDispatchers = append(removedDispatchers, job)
				return nil
			},
		},
	}

	processed, discovered, err := h.reconcileUserInfoState(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, processed)
	assert.Equal(t, 1, discovered)
	assert.Equal(t, []int32{42}, reconciled)
	assert.ElementsMatch(t, []string{"ambulance", "fire"}, removedDispatchers)
}

func TestReconcileUserInfoStateContinuesAfterLookupFailure(t *testing.T) {
	t.Parallel()

	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{rangeFn: func(fn func(string, *centrumunits.Unit) bool) {
			fn("unit.42", &centrumunits.Unit{Users: []*centrumunits.UnitAssignment{{UserId: 42}}})
			fn("unit.7", &centrumunits.Unit{Users: []*centrumunits.UnitAssignment{{UserId: 7}}})
		}, reconcile: func(context.Context, int32, string) (bool, error) {
			return false, nil
		}},
		dispatcherUserState: fakeDispatcherUserState{},
	}

	processed, discovered, err := h.reconcileUserInfoState(t.Context())
	require.Error(t, err)
	assert.Equal(t, 1, processed)
	assert.Equal(t, 2, discovered)
}

func TestReconcileUserInfoStateStopsOnCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		unitUserState: fakeUnitUserState{rangeFn: func(fn func(string, *centrumunits.Unit) bool) {
			fn("unit.42", &centrumunits.Unit{Users: []*centrumunits.UnitAssignment{{UserId: 42}}})
		}},
		dispatcherUserState: fakeDispatcherUserState{},
	}

	processed, discovered, err := h.reconcileUserInfoState(ctx)
	require.ErrorIs(t, err, context.Canceled)
	assert.Zero(t, processed)
	assert.Equal(t, 1, discovered)
}
