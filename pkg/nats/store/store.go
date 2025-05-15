package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/locks"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/multierr"
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

// (Mostly) Read Only store version - Technically the `Get` action can result in data being written to the store's internal cache
type StoreRO[T any, U protoutils.ProtoMessageWithMerge[T]] interface {
	Get(key string) (U, bool)
	Keys(ctx context.Context, prefix string) []string
	List() []U
	Range(ctx context.Context, fn func(key string, value U) bool)
}

type Store[T any, U protoutils.ProtoMessageWithMerge[T]] struct {
	logger *zap.Logger
	bucket string
	kv     jetstream.KeyValue
	l      *locks.Locks
	cl     bool

	prefix string

	mu   *xsync.Map[string, *sync.Mutex]
	data *xsync.Map[string, U]

	onUpdate   OnUpdateFn[T, U]
	onDelete   OnDeleteFn[T, U]
	onNotFound OnNotFoundFn[T, U]
}

type Option[T any, U protoutils.ProtoMessageWithMerge[T]] func(s *Store[T, U])

type (
	OnUpdateFn[T any, U protoutils.ProtoMessageWithMerge[T]]   func(s *Store[T, U], value U) (U, error)
	OnDeleteFn[T any, U protoutils.ProtoMessageWithMerge[T]]   func(s *Store[T, U], entry jetstream.KeyValueEntry, value U) error
	OnNotFoundFn[T any, U protoutils.ProtoMessageWithMerge[T]] func(s *Store[T, U], ctx context.Context, key string) (U, error)
)

func mutexCompute() (*sync.Mutex, bool) {
	return &sync.Mutex{}, false
}

func New[T any, U protoutils.ProtoMessageWithMerge[T]](ctx context.Context, logger *zap.Logger, js *events.JSWrapper, bucket string, opts ...Option[T, U]) (*Store[T, U], error) {
	s := &Store[T, U]{
		logger: logger.Named("store").With(zap.String("bucket", bucket)),
		bucket: bucket,

		cl: true,

		mu:   xsync.NewMap[string, *sync.Mutex](),
		data: xsync.NewMap[string, U](),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.kv == nil {
		storeKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:      bucket,
			Description: fmt.Sprintf("%s Store", bucket),
			History:     1,
			Storage:     jetstream.MemoryStorage,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create kv (bucket %s) for store. %w", bucket, err)
		}

		s.kv = storeKV
	}

	// Create locks only if not overriden by option
	if s.cl && s.l == nil {
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
		s.l = l
	}

	if s.prefix != "" {
		s.prefix = s.prefix + "."
	}

	return s, nil
}

// Put upload the message to kv and local
func (s *Store[T, U]) Put(ctx context.Context, key string, msg U) error {
	key = s.prefix + events.SanitizeKey(key)
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
	key = s.prefix + events.SanitizeKey(key)

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
			} else if s.onNotFound != nil {
				// Try to load value using not found method
				existing, err = s.onNotFound(s, ctx, key)
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
	key = s.prefix + events.SanitizeKey(key)

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
	key = s.prefix + events.SanitizeKey(key)

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
	key = s.prefix + events.SanitizeKey(key)

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
	if s.onUpdate != nil {
		item, err = s.onUpdate(s, item)
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
	key = s.prefix + events.SanitizeKey(key)

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

func (s *Store[T, U]) Keys(ctx context.Context, prefix string) []string {
	hasPrefix := (s.prefix + prefix) != ""
	if hasPrefix {
		if prefix != "" {
			prefix += "."
		}
		prefix = s.prefix + events.SanitizeKey(prefix)
	}

	keys := []string{}
	s.data.Range(func(key string, _ U) bool {
		if hasPrefix {
			if strings.HasPrefix(key, prefix) {
				if s.prefix == "" {
					keys = append(keys, key)
				} else {
					after, _ := strings.CutPrefix(key, s.prefix)
					keys = append(keys, after)
				}
			}
		} else {
			keys = append(keys, key)
		}

		return true
	})

	return keys
}

func (s *Store[T, U]) List() []U {
	list := []U{}

	s.data.Range(func(key string, value U) bool {
		if value == nil {
			return true
		}

		if s.prefix != "" {
			key, _ = strings.CutPrefix(key, s.prefix)
		}
		item, ok := s.Get(key)
		if !ok {
			return true
		}

		list = append(list, item)
		return true
	})

	return list
}

func (s *Store[T, U]) Start(ctx context.Context, wait bool) error {
	watcher, err := s.kv.WatchAll(ctx)
	if err != nil {
		return err
	}

	var ready atomic.Bool
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		updateCh := watcher.Updates()
		for {
			select {
			case <-ctx.Done():
				if err := watcher.Stop(); err != nil {
					if !errors.Is(err, nats.ErrConsumerNotFound) {
						s.logger.Error("error while stopping watcher", zap.Error(err))
					}
				} else {
					s.logger.Debug("store watcher done")
				}
				return

			case entry := <-updateCh:
				// After all initial keys have been received, a nil entry is returned
				if entry == nil {
					if !ready.Swap(true) {
						wg.Done()
					}
					continue
				}

				// Ignore keys that don't have the specified prefix (if any is set)
				if s.prefix != "" && !strings.HasPrefix(entry.Key(), s.prefix) {
					continue
				}

				s.logger.Debug("key update received via watcher", zap.String("key", entry.Key()), zap.Uint64("delta", entry.Delta()), zap.Uint8("op", uint8(entry.Operation())))

				switch entry.Operation() {
				case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
					// Handle delete and purge operations
					func() {
						mu, _ := s.mu.LoadOrCompute(entry.Key(), mutexCompute)
						mu.Lock()
						defer mu.Unlock()

						if s.onDelete != nil {
							item, _ := s.data.LoadAndDelete(entry.Key())
							if err := s.onDelete(s, entry, item); err != nil {
								s.logger.Error("failed to run on delete logic in store watcher", zap.String("key", entry.Key()), zap.Error(err))
							}
						}

						s.mu.Delete(entry.Key())
					}()

				case jetstream.KeyValuePut:
					// Handle put operations
					func() {
						mu, _ := s.mu.LoadOrCompute(entry.Key(), mutexCompute)
						mu.Lock()
						defer mu.Unlock()

						if _, err := s.update(entry); err != nil {
							s.logger.Error("failed to run on update logic in store watcher", zap.String("key", entry.Key()), zap.Error(err))
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

	if wait {
		wg.Wait()
	}

	return nil
}

func (s *Store[T, U]) Range(ctx context.Context, fn func(key string, value U) bool) {
	keys := s.Keys(ctx, "")

	for _, key := range keys {
		v, ok := s.Get(key)
		if !ok {
			continue
		}

		if !fn(key, v) {
			break
		}
	}
}

func (s *Store[T, U]) Clear(ctx context.Context) error {
	keys := s.Keys(ctx, "")

	errs := multierr.Combine()
	for _, key := range keys {
		if err := s.Delete(ctx, key); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (s *Store[T, U]) ReadOnly() StoreRO[T, U] {
	return s
}
