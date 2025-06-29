package leaderelection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

const BucketName = "leader_election"

// LeaderElector provides a generic leader election using JetStream KV with per-key TTL.
// It invokes callbacks when this instance becomes leader or loses leadership, and supports
// cancellation via both external context cancellation and Stop(), and can be restarted.
type LeaderElector struct {
	logger            *zap.Logger
	kv                jetstream.KeyValue
	js                *events.JSWrapper
	bucket            string
	key               string
	ttl               time.Duration
	heartbeatInterval time.Duration
	id                string

	onStarted func(ctx context.Context)
	onStopped func()

	// contexts for watcher and current leadership
	parentCtx    context.Context
	parentCancel context.CancelFunc
	leaderCtx    context.Context
	leaderCancel context.CancelFunc

	// watcher cancel to stop watch routine
	watcherCancel context.CancelFunc

	mu      sync.Mutex
	running bool
}

// New creates a LeaderElector. The provided parentCtx will be used to derive
// a watcher context. Stop() will end watcher and leadership, but parentCtx is preserved
// so Start() can be called multiple times.
func New(
	parentCtx context.Context,
	logger *zap.Logger,
	js *events.JSWrapper,
	bucket, key string,
	ttl, heartbeat time.Duration,
	onStarted func(ctx context.Context),
	onStopped func(),
) (*LeaderElector, error) {
	if heartbeat >= ttl {
		panic(fmt.Sprintf("invalid heartbeat (%s): must be less than ttl (%s)", heartbeat, ttl))
	}

	// derive a watcher context
	wCtx, wCancel := context.WithCancel(parentCtx)

	// ensure KV bucket exists
	kv, err := js.CreateOrUpdateKeyValue(wCtx, jetstream.KeyValueConfig{
		Bucket:         bucket,
		TTL:            ttl,
		Description:    "Leader election bucket",
		LimitMarkerTTL: 2 * ttl, // ensure we have enough time to handle re-election
	})
	if err != nil {
		wCancel()
		return nil, fmt.Errorf("failed to create/update KV bucket %s: %w", bucket, err)
	}

	return &LeaderElector{
		logger:            logger.Named("leader_election").With(zap.String("bucket", bucket), zap.String("key", key)),
		kv:                kv,
		js:                js,
		bucket:            bucket,
		key:               key,
		ttl:               ttl,
		heartbeatInterval: heartbeat,
		id:                uuid.NewString(),

		onStarted: onStarted,
		onStopped: onStopped,

		parentCtx:    wCtx,
		parentCancel: wCancel,
	}, nil
}

// Start begins the leader election process. Non-blocking. Can be called multiple times.
func (le *LeaderElector) Start() {
	le.logger.Info("starting leader election", zap.String("key", le.key), zap.String("id", le.id))
	le.mu.Lock()
	if le.running {
		le.mu.Unlock()
		return
	}
	le.running = true
	// reset contexts if previously stopped
	if le.watcherCancel != nil {
		// previous watcherCtx cancelled; create new
		wCtx, wCancel := context.WithCancel(context.Background())
		le.parentCtx = wCtx
		le.parentCancel = wCancel
	}
	le.mu.Unlock()

	// start watching for key deletions (loss of leadership)
	watchCtx, watchCancel := context.WithCancel(le.parentCtx)
	le.mu.Lock()
	le.watcherCancel = watchCancel
	le.mu.Unlock()
	go le.watchParent(watchCtx)

	// attempt initial acquisition
	if ok, _ := le.tryAcquire(); ok {
		le.becomeLeader()
	}
}

// Stop halts election and leadership. Parent context remains for future restarts.
func (le *LeaderElector) Stop() {
	le.mu.Lock()
	if !le.running {
		le.mu.Unlock()
		return
	}
	le.running = false
	// stop watcher
	if le.watcherCancel != nil {
		le.watcherCancel()
		le.watcherCancel = nil
	}
	// cancel current leadership if any
	if le.leaderCancel != nil {
		le.leaderCancel()
		le.leaderCancel = nil
	}
	le.mu.Unlock()
}

// tryAcquire attempts to create the leader key with TTL.
func (le *LeaderElector) tryAcquire() (bool, error) {
	_, err := le.kv.Create(le.parentCtx, le.key, []byte(le.id), jetstream.KeyTTL(le.ttl))
	if err == jetstream.ErrKeyExists {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// watchParent listens for deletion events on the leader key and triggers re-election.
func (le *LeaderElector) watchParent(ctx context.Context) {
	watcher, err := le.kv.WatchAll(ctx, jetstream.UpdatesOnly())
	if err != nil {
		le.logger.Error("watch error", zap.Error(err))
		return
	}
	for entry := range watcher.Updates() {
		if entry == nil {
			continue
		}

		if entry.Key() == le.key && entry.Operation() == jetstream.KeyValueDelete {
			le.logger.Info("leadership lost by expiration, retrying election")
			// cleanup old leadership context
			le.mu.Lock()
			if le.leaderCancel != nil {
				le.leaderCancel()
				le.leaderCancel = nil
			}
			le.mu.Unlock()
			// try to acquire again
			if ok, _ := le.tryAcquire(); ok {
				le.becomeLeader()
			}
		}
	}
}

// becomeLeader sets up leadership context, invokes onStarted, and starts heartbeat.
func (le *LeaderElector) becomeLeader() {
	le.logger.Info("became leader")
	// create a cancelable leader context
	leaderCtx, leaderCancel := context.WithCancel(le.parentCtx)
	le.mu.Lock()
	le.leaderCtx = leaderCtx
	le.leaderCancel = leaderCancel
	le.mu.Unlock()

	// invoke startup callback
	if le.onStarted != nil {
		go le.onStarted(leaderCtx)
	}

	// start heartbeat loop
	ticker := time.NewTicker(le.heartbeatInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if _, err := le.kv.Put(le.parentCtx, le.key, []byte(le.id)); err != nil {
					le.logger.Error("failed to refresh leadership TTL", zap.Error(err))
				}

			case <-leaderCtx.Done():
				ticker.Stop()
				le.logger.Info("leadership context cancelled")
				if le.onStopped != nil {
					le.onStopped()
				}
				return
			}
		}
	}()
}
