package leaderelection

import (
	"context"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// waitForCtx waits for the next context sent on ch or fails after timeout.
func waitForCtx(t *testing.T, ch <-chan context.Context, timeout time.Duration) context.Context {
	t.Helper()
	select {
	case ctx := <-ch:
		return ctx
	case <-time.After(timeout):
		require.Fail(t, "timed out waiting for leader start")
		return nil
	}
}

func TestLeaderElector_StopsHeartbeatAfterStepDown(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	natsServer := nats.NewServer(t, nats.ServerOptions{InProcess: true})
	js := natsServer.GetJS()

	started := make(chan context.Context, 2)

	le, err := New(
		ctx,
		zap.NewNop(),
		js,
		"le_test_bucket",
		"leader",
		1500*time.Millisecond,
		150*time.Millisecond,
		func(c context.Context) { started <- c },
		nil,
	)
	require.NoError(t, err)

	le.Start()
	firstCtx := waitForCtx(t, started, 2*time.Second)

	kv, err := js.KeyValue(ctx, "le_test_bucket")
	require.NoError(t, err)

	// Force step-down and ensure leadership context cancels.
	le.demote("test demote")
	select {
	case <-firstCtx.Done():
	case <-time.After(2 * time.Second):
		require.Fail(t, "leadership context was not canceled after demote")
	}

	// Competitor keeps the key alive with a distinct payload.
	stopCompetitor := make(chan struct{})
	go func() {
		ticker := time.NewTicker(60 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopCompetitor:
				return
			case <-ticker.C:
				_, _ = kv.Put(ctx, "leader", []byte("B"))
			}
		}
	}()

	// Give enough time for potential stray heartbeats to run.
	time.Sleep(400 * time.Millisecond)

	entry, err := kv.Get(ctx, "leader")
	require.NoError(t, err)
	require.Equal(t, []byte("B"), entry.Value(), "leader heartbeat should stop after demote")

	close(stopCompetitor)
	le.Stop()
}

func TestLeaderElector_ReacquiresAfterCompetitorLeaves(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	natsServer := nats.NewServer(t, nats.ServerOptions{InProcess: true})
	js := natsServer.GetJS()

	started := make(chan context.Context, 3)

	le, err := New(
		ctx,
		zap.NewNop(),
		js,
		"le_test_bucket2",
		"leader",
		1*time.Second,
		150*time.Millisecond,
		func(c context.Context) { started <- c },
		nil,
	)
	require.NoError(t, err)

	le.Start()
	firstCtx := waitForCtx(t, started, 2*time.Second)

	kv, err := js.KeyValue(ctx, "le_test_bucket2")
	require.NoError(t, err)

	// Step down current leader and simulate another leader holding the key for a while.
	le.demote("hand over to competitor")
	select {
	case <-firstCtx.Done():
	case <-time.After(2 * time.Second):
		require.Fail(t, "leadership context was not canceled after demote")
	}

	stopCompetitor := make(chan struct{})
	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopCompetitor:
				return
			case <-ticker.C:
				_, _ = kv.Put(ctx, "leader", []byte("B"))
			}
		}
	}()

	// Keep competitor active briefly, then release the key.
	time.Sleep(350 * time.Millisecond)
	close(stopCompetitor)

	// Wait for the key to expire naturally (TTL is 1s) and allow retries to kick in.
	time.Sleep(1200 * time.Millisecond)

	secondCtx := waitForCtx(t, started, 2*time.Second)
	require.NotEqual(
		t,
		firstCtx,
		secondCtx,
		"elector should reacquire leadership after key is free",
	)

	le.Stop()
}

func TestLeaderElectorsFailOverWithoutOverlappingLeadership(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	natsServer := nats.NewServer(t, nats.ServerOptions{InProcess: true})
	js := natsServer.GetJS()

	type leadershipStart struct {
		id  string
		ctx context.Context
	}
	starts := make(chan leadershipStart, 4)
	stops := make(chan string, 4)
	newElector := func(id string) *LeaderElector {
		elector, err := New(
			ctx,
			zap.NewNop(),
			js,
			"le_two_electors",
			"leader",
			time.Second,
			100*time.Millisecond,
			func(leadershipCtx context.Context) { starts <- leadershipStart{id: id, ctx: leadershipCtx} },
			func() { stops <- id },
		)
		require.NoError(t, err)
		return elector
	}

	first := newElector("first")
	second := newElector("second")
	first.Start()
	second.Start()

	initial := func() leadershipStart {
		select {
		case start := <-starts:
			return start
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for initial leader")
			return leadershipStart{}
		}
	}()

	select {
	case duplicate := <-starts:
		t.Fatalf("both electors became leader before failover: %s and %s", initial.id, duplicate.id)
	case <-time.After(300 * time.Millisecond):
	}

	if initial.id == "first" {
		first.Stop()
	} else {
		second.Stop()
	}

	select {
	case stopped := <-stops:
		require.Equal(t, initial.id, stopped)
	case <-time.After(time.Second):
		t.Fatal("stopped callback was not called for the initial leader")
	}
	select {
	case <-initial.ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("initial leadership context was not canceled after Stop")
	}

	var failover leadershipStart
	select {
	case failover = <-starts:
	case <-time.After(3 * time.Second):
		t.Fatal("remaining elector did not take over after the leader key expired")
	}
	require.NotEqual(t, initial.id, failover.id)

	if failover.id == "first" {
		first.Stop()
	} else {
		second.Stop()
	}
	select {
	case stopped := <-stops:
		require.Equal(t, failover.id, stopped)
	case <-time.After(time.Second):
		t.Fatal("stopped callback was not called for the failover leader")
	}
}
