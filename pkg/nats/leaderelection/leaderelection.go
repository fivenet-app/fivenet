// leader_election_example.go
//
// Self-healing leader-election helper for NATS **JetStream** Key-Value.
//
// Public surface mirrors your original helper but now:
//   - All JetStream calls that require a **context.Context** receive one
//     (bucket creation, Create/Put, WatchAll…).
//   - `onStarted(ctx)` gets a leadership-lifetime context that is cancelled the
//     moment this instance loses leadership or Stop() is called.
//
// Usage:
// ```go
// js     := events.NewJSWrapper(nc)
// logger, _ := zap.NewProduction()
// defer logger.Sync()
//
// le, _ := leaderelection.NewLeaderElector(
//
//	ctx,
//	js,
//	"leader_bucket", // bucket name
//	"leader",        // key
//	8*time.Second,    // per-key TTL
//	3*time.Second,    // heartbeat
//	logger,
//	func(ctx context.Context) {
//	    go runScheduler(ctx) // workload tied to leadership ctx
//	},
//	nil, // onStopped optional
//
// )
//
// le.Start()
// <-ctx.Done()
// le.Stop()
// ```
package leaderelection

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type LeaderElector struct {
	logger    *zap.Logger
	mu        sync.Mutex
	kv        jetstream.KeyValue
	key       string
	ttl       time.Duration
	heartbeat time.Duration

	ctx    context.Context // parent-lifetime
	cancel context.CancelFunc

	leadershipCtx    context.Context // cancelled on step-down
	leadershipCancel context.CancelFunc

	isLeader bool

	onStarted func(context.Context)
	onStopped func()
}

// New builds the helper and prepares it for Start().
// Returns an error if the KV bucket cannot be created or accessed.
func New(
	parent context.Context,
	l *zap.Logger,
	js *events.JSWrapper,
	bucket, key string,
	ttl, heartbeat time.Duration,
	onStarted func(context.Context),
	onStopped func(),
) (*LeaderElector, error) {
	kv, err := js.CreateOrUpdateKeyValue(parent, jetstream.KeyValueConfig{
		Bucket:         bucket,
		TTL:            ttl,
		Description:    "Leader election bucket",
		LimitMarkerTTL: 2 * ttl, // ensure we have enough time to handle re-election
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(parent)

	return &LeaderElector{
		logger:    l,
		kv:        kv,
		key:       key,
		ttl:       ttl,
		heartbeat: heartbeat,
		ctx:       ctx,
		cancel:    cancel,
		onStarted: onStarted,
		onStopped: onStopped,
	}, nil
}

// Lifecycle

// Start launches background loops - call once.
func (le *LeaderElector) Start() {
	le.tryAcquire()     // quick attempt
	go le.watchParent() // resilient watcher
	go le.retryLoop()   // periodic retries
}

// Stop terminates everything and cancels any leadership context.
func (le *LeaderElector) Stop() {
	le.demote("Stop called")
	le.cancel()
}

// Internal helpers

// demote revokes leadership if held, fires hooks & cancels ctx.
func (le *LeaderElector) demote(reason string) {
	le.mu.Lock()
	defer le.mu.Unlock()

	if !le.isLeader {
		return
	}

	le.logger.Info("stepping down", zap.String("reason", reason))
	le.isLeader = false
	if le.leadershipCancel != nil {
		le.leadershipCancel()
		le.leadershipCancel = nil
		le.leadershipCtx = nil
	}
	if le.onStopped != nil {
		le.onStopped()
	}
}

// retryLoop periodically re-tries acquisition to cover missed events.
func (le *LeaderElector) retryLoop() {
	t := time.NewTicker(le.heartbeat)
	defer t.Stop()
	for {
		select {
		case <-le.ctx.Done():
			return

		case <-t.C:
			if !le.isLeader {
				le.tryAcquire()
			}
		}
	}
}

// watchParent opens a watcher and restarts on closure.
func (le *LeaderElector) watchParent() {
	for {
		if le.ctx.Err() != nil {
			return
		}

		watcher, err := le.kv.WatchAll(le.ctx, jetstream.UpdatesOnly())
		if err != nil {
			le.logger.Warn("watch error (retrying)", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		for entry := range watcher.Updates() {
			if entry == nil {
				continue // reset marker
			}

			switch entry.Operation() {
			case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
				if entry.Key() == le.key {
					le.demote("KV key vanished")
					le.tryAcquire() // immediate re-take attempt
				}
			}
		}

		le.logger.Info("watcher closed; restarting")
		time.Sleep(time.Second)
	}
}

// tryAcquire attempts to write the leader key with TTL.
func (le *LeaderElector) tryAcquire() {
	_, err := le.kv.Create(le.ctx, le.key, nil, jetstream.KeyTTL(le.ttl))
	if err == nil {
		le.promote()
		return
	}
	if !errors.Is(err, jetstream.ErrKeyExists) {
		le.logger.Warn("tryAcquire", zap.Error(err))
	}
}

// promote marks this instance leader, fires callback, and starts heartbeat.
func (le *LeaderElector) promote() {
	le.mu.Lock()
	defer le.mu.Unlock()

	if le.isLeader {
		return
	}

	le.isLeader = true
	le.leadershipCtx, le.leadershipCancel = context.WithCancel(le.ctx)

	le.logger.Info("became leader")
	if le.onStarted != nil {
		go le.onStarted(le.leadershipCtx)
	}

	go le.heartbeatLoop()
}

// heartbeatLoop refreshes the key while leader.
func (le *LeaderElector) heartbeatLoop() {
	t := time.NewTicker(le.heartbeat)
	defer t.Stop()

	for {
		select {
		case <-le.ctx.Done():
			return
		case <-t.C:
			if _, err := le.kv.Put(le.ctx, le.key, nil); err != nil {
				le.logger.Warn("failed to refresh key", zap.Error(err))
				le.demote("refresh failed")
				return
			}
		}
	}
}
