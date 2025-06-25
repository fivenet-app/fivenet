/*
 * Modified version of https://github.com/HeavyHorst/certmagic-nats/blob/b27fd6c010166e396b6f9e1c651ba7b02ce6c01f/nats.go#L114
 * which is licensed under [MIT License](https://github.com/HeavyHorst/certmagic-nats/blob/b27fd6c010166e396b6f9e1c651ba7b02ce6c01f/LICENSE)
 *
 * The code has been modified to use the per-key TTL feature of NATS JetStream KeyValue store that was added in NATS 2.11.0.
 */
package locks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

const (
	// LockTimeout is the default timeout for acquiring a lock
	LockTimeout = 750 * time.Millisecond

	// KeyPrefix is the prefix used for all lock keys in the KV store
	KeyPrefix = "LOCK."
)

var (
	// Prometheus metric: total number of lock acquisition attempts
	lockAcquireTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "nats_locks",
		Name:      "acquire_total",
		Help:      "Total number of lock acquisition attempts.",
	}, []string{"bucket", "result"})

	// Prometheus metric: histogram of lock held durations
	lockHeldDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "nats_locks",
		Name:      "held_duration_seconds",
		Help:      "Histogram of lock held durations.",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
	}, []string{"bucket"})
)

// Locks provides distributed locking using NATS JetStream KeyValue store
// with lock expiration and metrics.
type Locks struct {
	// logger instance
	logger *zap.Logger
	// NATS JetStream KeyValue store
	kv jetstream.KeyValue
	// bucket name for locks
	bucket string
	// instanceName is the name of the instance using the locks
	instanceName string

	// maximum age for a lock before considered stale
	maxLockAge time.Duration

	// map of lock key to last revision
	revMap map[string]uint64
	// mutex for revMap
	maplock sync.RWMutex
}

// New creates a new Locks instance for a given bucket and max lock age (`_locks` is automatically added to the bucket).
func New(ctx context.Context, logger *zap.Logger, js *events.JSWrapper, bucket string, maxLockAge time.Duration) (*Locks, error) {
	lBucket := fmt.Sprintf("%s_locks", bucket)
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:         lBucket,
		Description:    fmt.Sprintf("%s Locks", bucket),
		History:        1,
		MaxBytes:       -1,
		Storage:        jetstream.MemoryStorage,
		LimitMarkerTTL: maxLockAge, // Set a limit marker TTL to avoid stale locks
	})
	if err != nil {
		return nil, err
	}

	return NewWithKV(logger, kv, bucket, maxLockAge), nil
}

// NewWithKV creates a new Locks instance using an existing KeyValue store.
// This is useful for testing or when you already have a KeyValue store set up.
func NewWithKV(logger *zap.Logger, kv jetstream.KeyValue, bucket string, maxLockAge time.Duration) *Locks {
	lBucket := fmt.Sprintf("%s_locks", bucket)
	if kv.Bucket() != lBucket {
		logger.Warn("using Locks with a KeyValue store that does not match the expected bucket",
			zap.String("expected_bucket", lBucket), zap.String("actual_bucket", kv.Bucket()))
	}

	return &Locks{
		logger:       logger.Named("locks").With(zap.String("bucket", lBucket)),
		kv:           kv,
		bucket:       lBucket,
		instanceName: instance.ID(),

		maxLockAge: maxLockAge,

		revMap:  map[string]uint64{},
		maplock: sync.RWMutex{},
	}
}

type lockInfo struct {
	Owner string `json:"owner"` // Instance name
	TS    int64  `json:"ts"`    // Unix-nanos
}

func (l *Locks) debugPayload() []byte {
	b, _ := json.Marshal(lockInfo{
		Owner: l.instanceName,
		TS:    time.Now().UnixNano(),
	})
	return b
}

// Lock acquires the lock for key, blocking until the lock
// can be obtained or an error is returned. Note that, even
// after acquiring a lock, an idempotent operation may have
// already been performed by another process that acquired
// the lock before - so always check to make sure idempotent
// operations still need to be performed after acquiring the
// lock.
//
// The actual implementation of obtaining of a lock must be
// an atomic operation so that multiple Lock calls at the
// same time always results in only one caller receiving the
// lock at any given time.
//
// To prevent deadlocks, all implementations (where this concern
// is relevant) should put a reasonable expiration on the lock in
// case Unlock is unable to be called due to some sort of network
// failure or system crash.
func (l *Locks) Lock(ctx context.Context, key string) error {
	start := time.Now()
	lockKey := makeLockKey(key)
	l.logger.Debug("lock", zap.String("key", lockKey))

loop:
	for {
		revision, err := l.kv.Get(ctx, lockKey)
		if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			return err
		}
		if revision == nil {
			break // Nobody holds the lock
		}

		// Wait a little and retry
		select {
		case <-time.After(jitterDelay()):
		case <-ctx.Done():
			lockAcquireTotal.WithLabelValues(l.bucket, "fail").Inc()
			return ctx.Err()
		}
	}

	// Lock doesn't exist, create it
	nrev, err := l.kv.Create(ctx, lockKey, l.debugPayload(), jetstream.KeyTTL(l.maxLockAge))
	if err != nil && isWrongSequence(err) {
		// another process created the lock in the meantime
		// try again
		goto loop
	}

	if err != nil {
		lockAcquireTotal.WithLabelValues(l.bucket, "fail").Inc()
		return err
	}

	l.setRev(lockKey, nrev)

	lockAcquireTotal.WithLabelValues(l.bucket, "success").Inc()
	lockHeldDuration.WithLabelValues(l.bucket).Observe(time.Since(start).Seconds())
	return nil
}

// TryLock attempts to acquire the lock for key without blocking.
// Returns true if the lock was acquired, false otherwise.
func (l *Locks) TryLock(ctx context.Context, key string) (bool, error) {
	lockKey := makeLockKey(key)

	revision, err := l.kv.Get(ctx, lockKey)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return false, err
	}
	if revision != nil {
		return false, nil // someone else holds the lock
	}

	// Create the lock with TTL + debug payload in one shot
	nrev, err := l.kv.Create(ctx, lockKey, l.debugPayload(), jetstream.KeyTTL(l.maxLockAge))
	if err != nil {
		if isWrongSequence(err) {
			return false, nil
		}
		return false, err
	}

	l.setRev(lockKey, nrev)
	return true, nil
}

// IsLocked checks if the lock for key currently exists.
func (l *Locks) IsLocked(ctx context.Context, key string) (bool, error) {
	// Check for existing lock
	revision, err := l.kv.Get(ctx, makeLockKey(key))
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return false, err
	}

	return revision != nil, nil
}

// Unlock releases the lock for key. This method must ONLY be
// called after a successful call to Lock, and only after the
// critical section is finished, even if it errored or timed
// out. Unlock cleans up any resources allocated during Lock.
func (l *Locks) Unlock(ctx context.Context, key string) error {
	lockKey := makeLockKey(key)
	l.logger.Debug("unlock", zap.String("key", lockKey))
	return l.kv.Delete(ctx, lockKey, jetstream.LastRevision(l.getRev(lockKey)))
}

// setRev stores the last revision for a lock key.
func (l *Locks) setRev(key string, value uint64) {
	l.maplock.Lock()
	defer l.maplock.Unlock()
	l.revMap[key] = value
}

// getRev retrieves the last revision for a lock key.
func (l *Locks) getRev(key string) uint64 {
	l.maplock.RLock()
	defer l.maplock.RUnlock()
	return l.revMap[key]
}

// isWrongSequence checks if the error is due to a wrong last sequence in JetStream.
func isWrongSequence(err error) bool {
	if err, ok := err.(*jetstream.APIError); ok {
		if err.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence {
			return true
		}
	}

	// Fallback to checking the error message contents
	return strings.Contains(err.Error(), "wrong last sequence")
}

// makeLockKey returns the full key for a lock.
func makeLockKey(key string) string {
	return KeyPrefix + key
}

// jitterDelay returns a random delay for lock retry backoff.
func jitterDelay() time.Duration {
	return time.Duration(50+rand.Float64()*150) * time.Millisecond
}
