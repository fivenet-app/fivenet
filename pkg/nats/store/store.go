package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/nats/locks"
	"github.com/fivenet-app/fivenet/pkg/server/admin"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var ErrNotFound = errors.New("store: not found")

var metricDataMapCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "nats_store",
	Name:      "datamap_count",
	Help:      "Count of data map entries.",
}, []string{"bucket"})

type protoMessage[T any] interface {
	*T
	proto.Message

	Merge(in *T) *T
}

type Store[T any, U protoMessage[T]] struct {
	logger *zap.Logger
	bucket string
	kv     jetstream.KeyValue
	l      *locks.Locks

	mu   *xsync.MapOf[string, *sync.Mutex]
	data *xsync.MapOf[string, U]

	OnUpdate   OnUpdateFn[T, U]
	OnDelete   OnDeleteFn[T, U]
	OnNotFound OnNotFoundFn[T, U]
}

type Option[T any, U protoMessage[T]] func(s *Store[T, U]) error

type (
	OnUpdateFn[T any, U protoMessage[T]]   func(U) (U, error)
	OnDeleteFn[T any, U protoMessage[T]]   func(jetstream.KeyValueEntry, U) error
	OnNotFoundFn[T any, U protoMessage[T]] func(ctx context.Context, key string) (U, error)
)

func mutexCompute() *sync.Mutex {
	return &sync.Mutex{}
}

func NewWithLocks[T any, U protoMessage[T]](ctx context.Context, logger *zap.Logger, js *events.JSWrapper, bucket string, l *locks.Locks, opts ...Option[T, U]) (*Store[T, U], error) {
	storeKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      bucket,
		Description: fmt.Sprintf("%s Store", bucket),
		History:     1,
		Storage:     jetstream.MemoryStorage,
	})
	if err != nil {
		return nil, err
	}

	s := &Store[T, U]{
		logger: logger.Named("store").With(zap.String("bucket", bucket)),
		bucket: bucket,
		kv:     storeKV,
		l:      l,

		mu:   xsync.NewMapOf[string, *sync.Mutex](),
		data: xsync.NewMapOf[string, U](),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func New[T any, U protoMessage[T]](ctx context.Context, logger *zap.Logger, js *events.JSWrapper, bucket string, opts ...Option[T, U]) (*Store[T, U], error) {
	lockBucket := fmt.Sprintf("%s_locks", bucket)

	locksKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      lockBucket,
		Description: fmt.Sprintf("%s Store Locks", bucket),
		History:     2,
		Storage:     jetstream.MemoryStorage,
		TTL:         5 * locks.LockTimeout,
	})
	if err != nil {
		return nil, err
	}

	l, err := locks.New(logger, locksKV, lockBucket, 4*locks.LockTimeout)
	if err != nil {
		return nil, err
	}

	return NewWithLocks(ctx, logger, js, bucket, l, opts...)
}

// Put upload the message to kv and local
func (s *Store[T, U]) Put(ctx context.Context, key string, msg U) error {
	key = events.SanitizeKey(key)
	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, locks.LockTimeout)
	defer cancel()

	if s.l != nil {
		if err := s.l.Lock(ctx, key); err != nil {
			return err
		}
		defer s.l.Unlock(ctx, key)
	}

	if err := s.put(ctx, key, msg); err != nil {
		return err
	}

	return nil
}

func (s *Store[T, U]) put(ctx context.Context, key string, msg U) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := s.kv.Put(ctx, key, data); err != nil {
		return err
	}

	s.updateFromType(key, msg)

	return nil
}

func (s *Store[T, U]) ComputeUpdate(ctx context.Context, key string, load bool, fn func(key string, existing U) (U, bool, error)) error {
	key = events.SanitizeKey(key)

	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, locks.LockTimeout)
	defer cancel()

	if s.l != nil {
		if err := s.l.Lock(ctx, key); err != nil {
			return err
		}
		defer s.l.Unlock(ctx, key)
	}

	var existing U
	if load {
		var err error
		existing, err = s.load(ctx, key)
		if err != nil {
			// Return if it isn't a not found value
			if !errors.Is(err, jetstream.ErrKeyNotFound) {
				return err
			} else if s.OnNotFound != nil {
				// Try to load value using not found method
				existing, err = s.OnNotFound(ctx, key)
				if err != nil && !errors.Is(err, ErrNotFound) {
					return err
				}
			}
		}
	} else {
		var ok bool
		existing, ok = s.get(key)
		if !ok {
			return fmt.Errorf("no item for key %s found in local store", key)
		}
	}

	computed, changed, err := fn(key, existing)
	if err != nil {
		return fmt.Errorf("store compute update function call returned error. %w", err)
	}

	// Skip computed nil updates for now
	if computed == nil {
		s.logger.Error("compute update returned nil, skipping store put", zap.String("key", key))
		return nil
	}

	// Only update key if no existing key it was changed (indicated by the compute update function call)
	if changed {
		if err := s.put(ctx, key, computed); err != nil {
			return err
		}
	} else {
		s.logger.Debug("store compute update has not changed state", zap.String("key", key))
	}

	return nil
}

// Get copy of data from local data
func (s *Store[T, U]) Get(key string) (U, bool) {
	key = events.SanitizeKey(key)

	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	return s.get(key)
}

func (s *Store[T, U]) get(key string) (U, bool) {
	i, ok := s.data.Load(key)
	if !ok {
		return nil, ok
	}

	return proto.Clone(i).(U), true
}

// Load data from kv store (this will add/update any existing local data entry)
// If no key is found, the original nats error is returned.
func (s *Store[T, U]) Load(ctx context.Context, key string) (U, error) {
	key = events.SanitizeKey(key)

	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	return s.load(ctx, key)
}

func (s *Store[T, U]) load(ctx context.Context, key string) (U, error) {
	entry, err := s.kv.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	return s.update(entry)
}

func (s *Store[T, U]) GetOrLoad(ctx context.Context, key string) (U, error) {
	key = events.SanitizeKey(key)

	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	item, ok := s.get(key)
	if !ok || item == nil {
		var err error
		item, err = s.load(ctx, key)
		if err != nil {
			return nil, err
		}

		return s.updateFromType(key, item), nil
	}

	return item, nil
}

func (s *Store[T, U]) update(entry jetstream.KeyValueEntry) (U, error) {
	data := U(new(T))
	if err := proto.Unmarshal(entry.Value(), data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal store watcher update: %w", err)
	}

	var err error
	item := s.updateFromType(entry.Key(), data)
	if s.OnUpdate != nil {
		item, err = s.OnUpdate(item)
	}

	return proto.Clone(item).(U), err
}

func (s *Store[T, U]) updateFromType(key string, updated U) U {
	current, loaded := s.data.LoadOrStore(key, updated)
	if loaded && current != nil {
		// Compare using protobuf magic and merge if not equal
		if !proto.Equal(current, updated) {
			current.Merge(updated)
		}

		return current
	}

	return updated
}

func (s *Store[T, U]) Delete(ctx context.Context, key string) error {
	key = events.SanitizeKey(key)

	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, locks.LockTimeout)
	defer cancel()

	if s.l != nil {
		if err := s.l.Lock(ctx, key); err != nil {
			return err
		}
		defer s.l.Unlock(ctx, key)
	}

	if err := s.kv.Purge(ctx, key); err != nil {
		return err
	}

	s.data.Delete(key)
	s.mu.Delete(key)

	return nil
}

func (s *Store[T, U]) Keys(ctx context.Context, prefix string) ([]string, error) {
	hasPrefix := prefix != ""
	if hasPrefix {
		prefix = events.SanitizeKey(prefix)
	}

	keys := []string{}
	s.data.Range(func(key string, _ U) bool {
		if hasPrefix {
			if strings.HasPrefix(key, prefix+".") {
				keys = append(keys, key)
			}
		} else {
			keys = append(keys, key)
		}

		return true
	})

	return keys, nil
}

func (s *Store[T, U]) List() ([]U, error) {
	list := []U{}

	s.data.Range(func(key string, value U) bool {
		if value == nil {
			return true
		}

		item, ok := s.Get(key)
		if !ok {
			return true
		}

		list = append(list, item)
		return true
	})

	return list, nil
}

func (s *Store[T, U]) Start(ctx context.Context) error {
	watcher, err := s.kv.WatchAll(ctx)
	if err != nil {
		return err
	}

	go func() {
		updateCh := watcher.Updates()
		for {
			select {
			case <-ctx.Done():
				if err := watcher.Stop(); err != nil {
					if !errors.Is(err, jetstream.ErrConsumerNotFound) {
						s.logger.Error("error while stopping watcher", zap.Error(err))
					}
				} else {
					s.logger.Debug("store watcher done")
				}
				return

			case entry := <-updateCh:
				// After all initial keys have been received, a nil entry is returned
				if entry == nil {
					continue
				}

				s.logger.Debug("key update received via watcher", zap.String("key", entry.Key()), zap.Uint64("delta", entry.Delta()))

				switch entry.Operation() {
				case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
					func() {
						mu, _ := s.mu.LoadOrCompute(entry.Key(), mutexCompute)
						mu.Lock()
						defer mu.Unlock()

						if s.OnDelete != nil {
							item, _ := s.data.LoadAndDelete(entry.Key())
							if err := s.OnDelete(entry, item); err != nil {
								s.logger.Error("failed to run on delete logic in store watcher", zap.Error(err))
							}
						}

						s.mu.Delete(entry.Key())
					}()

				case jetstream.KeyValuePut:
					func() {
						mu, _ := s.mu.LoadOrCompute(entry.Key(), mutexCompute)
						mu.Lock()
						defer mu.Unlock()

						if _, err := s.update(entry); err != nil {
							s.logger.Error("failed to run on update logic in store watcher", zap.Error(err))
						}
					}()

				default:
					s.logger.Error("unknown key operation received", zap.String("key", entry.Key()), zap.Uint8("op", uint8(entry.Operation())))
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-time.After(15 * time.Second):
				metricDataMapCount.WithLabelValues(s.bucket).Set(float64(s.data.Size()))
			}
		}
	}()

	return nil
}
