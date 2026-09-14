package housekeeper

import (
	"context"
	"errors"
	"testing"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
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
		reconcileUserUnitJobChange: func(_ context.Context, userID int32, job string) (bool, error) {
			unitChanges <- struct {
				userID int32
				job    string
			}{userID: userID, job: job}
			return true, nil
		},
		rangeDispatchers: func(fn func(string, *centrumdispatchers.Dispatchers) bool) {
			fn("ambulance", &centrumdispatchers.Dispatchers{
				Job:         "ambulance",
				Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
			})
			fn("police", &centrumdispatchers.Dispatchers{
				Job:         "police",
				Dispatchers: []*jobscolleagues.Colleague{{UserId: 42}},
			})
		},
		setDispatcherState: func(_ context.Context, job string, userID int32, active bool) error {
			assert.Equal(t, int32(42), userID)
			assert.False(t, active)
			dispatcherRemovals <- job
			return nil
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
		reconcileUserUnitJobChange: func(context.Context, int32, string) (bool, error) {
			return false, nil
		},
		rangeDispatchers: func(func(string, *centrumdispatchers.Dispatchers) bool) {},
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.True(t, msg.acked)
	assert.Zero(t, msg.nakDelay)
	assert.False(t, msg.terminated)
}

func TestHandleUserInfoReconcileMessageRetriesFailure(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "police"),
		reconcileUserUnitJobChange: func(context.Context, int32, string) (bool, error) {
			return false, errors.New("database unavailable")
		},
		rangeDispatchers: func(func(string, *centrumdispatchers.Dispatchers) bool) {},
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
		reconcileUserUnitJobChange: func(context.Context, int32, string) (bool, error) {
			return false, errors.New("database unavailable")
		},
		rangeDispatchers: func(func(string, *centrumdispatchers.Dispatchers) bool) {},
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
		reconcileUserUnitJobChange: func(context.Context, int32, string) (bool, error) {
			return false, nil
		},
		rangeDispatchers: func(func(string, *centrumdispatchers.Dispatchers) bool) {
			require.Fail(t, "dispatcher projections must not be scanned when jobs are known")
		},
		setDispatcherState: func(_ context.Context, job string, userID int32, active bool) error {
			assert.Equal(t, int32(42), userID)
			assert.False(t, active)
			dispatcherRemovals <- job
			return nil
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

func TestHandleUserInfoReconcileMessageUsesCurrentAuthoritativeJob(t *testing.T) {
	t.Parallel()

	msg := userInfoReconcileMessage(t, 42, "police")
	reconciledJob := ""
	h := &Housekeeper{
		logger:   zap.NewNop(),
		metrics:  centrummetrics.Get(),
		userinfo: userInfoRetrieverForTest(42, "fire"),
		reconcileUserUnitJobChange: func(_ context.Context, _ int32, job string) (bool, error) {
			reconciledJob = job
			return false, nil
		},
		rangeDispatchers: func(func(string, *centrumdispatchers.Dispatchers) bool) {},
	}

	h.handleUserInfoReconcileMessage(t.Context(), msg)
	assert.Equal(t, "fire", reconciledJob)
	assert.True(t, msg.acked)
}
