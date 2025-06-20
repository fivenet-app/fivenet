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

// LeaderElector provides a generic leader election using JetStream KV with per-key TTL.
// It invokes callbacks when this instance becomes leader or loses leadership.
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

	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
}

// New creates a LeaderElector.
// logger: zap.logger for logging events
// js: events.JetStream context for KV operations
// bucket: KV bucket name for election
// key: key under which leadership is claimed
// ttl: how long before a leader record expires
// heartbeat: interval to refresh TTL
// onStarted: called once when this instance becomes leader, receives a context that cancels on leadership loss
// onStopped: called when this instance loses leadership or stops
func New(logger *zap.Logger, js *events.JSWrapper, bucket, key string, ttl, heartbeat time.Duration, onStarted func(ctx context.Context), onStopped func()) (*LeaderElector, error) {
	ctx, cancel := context.WithCancel(context.Background())
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: bucket,
		TTL:    ttl,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create %s key value store. %w", bucket, err)
	}

	return &LeaderElector{
		logger:            logger.Named("leaderelection").With(zap.String("bucket", bucket), zap.String("key", key)),
		kv:                kv,
		js:                js,
		bucket:            bucket,
		key:               key,
		ttl:               ttl,
		heartbeatInterval: heartbeat,
		id:                uuid.NewString(),
		onStarted:         onStarted,
		onStopped:         onStopped,
		ctx:               ctx,
		cancel:            cancel,
	}, nil
}

// Start begins the leader election process. Non-blocking.
func (le *LeaderElector) Start() error {
	le.mu.Lock()
	if le.running {
		le.mu.Unlock()
		return nil
	}
	le.running = true
	le.mu.Unlock()

	// Watch for leadership changes
	if err := le.watch(); err != nil {
		return err
	}

	// Try initial election
	if ok, _ := le.tryAcquire(); ok {
		le.startLeader()
	}

	return nil
}

// Stop halts the election, revoking leadership and stopping heartbeats.
func (le *LeaderElector) Stop() {
	le.mu.Lock()
	defer le.mu.Unlock()
	if !le.running {
		return
	}
	le.running = false
	le.cancel()
}

func (le *LeaderElector) tryAcquire() (bool, error) {
	_, err := le.kv.Create(le.ctx, le.key, []byte(le.id), jetstream.KeyTTL(le.ttl))
	if err == jetstream.ErrKeyExists {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (le *LeaderElector) watch() error {
	watcher, err := le.kv.WatchAll(le.ctx)
	if err != nil {
		return err
	}
	for entry := range watcher.Updates() {
		if entry.Key() == le.key && entry.Operation() == jetstream.KeyValueDelete {
			// leadership expired
			if ok, _ := le.tryAcquire(); ok {
				le.startLeader()
			}
		}
	}

	return nil
}

func (le *LeaderElector) startLeader() {
	// Create a cancellable context for leadership
	ctx, cancel := context.WithCancel(context.Background())
	le.mu.Lock()
	le.cancel = cancel
	le.mu.Unlock()

	// Invoke onStarted callback
	if le.onStarted != nil {
		go le.onStarted(ctx)
	}

	// Heartbeat loop
	ticker := time.NewTicker(le.heartbeatInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				// Refresh TTL
				if _, err := le.kv.Put(le.ctx, le.key, []byte(le.id)); err != nil {
					le.logger.Error("failed to refresh leadership TTL", zap.String("key", le.key), zap.Error(err))
				}

			case <-ctx.Done():
				ticker.Stop()
				// Callback on stopped
				if le.onStopped != nil {
					le.onStopped()
				}
				return
			}
		}
	}()
}
