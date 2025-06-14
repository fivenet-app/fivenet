package store

import (
	"context"
	"errors"
	"fmt"
	"slices"
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

// ErrPrefixAmbiguous is returned when more than one key matches the prefix.
var ErrPrefixAmbiguous = fmt.Errorf("multiple keys found with given prefix")

var metricDataMapCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "nats_store",
	Name:      "datamap_count",
	Help:      "Count of data map entries.",
}, []string{"bucket"})

// (Mostly) Read Only store version - Technically the `Get` action can result in data being written to the store's internal cache
type StoreRO[T any, U protoutils.ProtoMessageWithMerge[T]] interface {
	Get(key string) (U, error)
	GetBySegmentOne(prefix string) (U, error)
	Keys(prefix string) []string
	KeysFiltered(prefix string, filter func(string) bool) []string
	List() []U
	ListFiltered(ctx context.Context, prefix string, filter func(string) bool) []U
	Range(fn func(key string, value U) bool)
	Count() int
	WatchAll(ctx context.Context) (chan *KeyValueEntry[T, U], error)
}

type Store[T any, U protoutils.ProtoMessageWithMerge[T]] struct {
	logger *zap.Logger
	bucket string
	kv     jetstream.KeyValue
	l      *locks.Locks
	cl     bool

	prefix      string
	ignoredKeys []string

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
		s.prefix = strings.TrimSuffix(s.prefix, ".") + "."
	}

	return s, nil
}

func (s *Store[T, U]) Start(ctx context.Context, wait bool) error {
	watcher, err := s.kv.Watch(ctx, s.prefix+">")
	if err != nil {
		return fmt.Errorf("failed to start store kv. %w", err)
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

				key := entry.Key()
				if s.ignoredKeys != nil && slices.Contains(s.ignoredKeys, key) {
					continue
				}

				switch entry.Operation() {
				case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
					// Handle delete and purge operations
					func() {
						mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
						mu.Lock()
						defer mu.Unlock()

						if s.onDelete != nil {
							item, _ := s.data.LoadAndDelete(key)
							if err := s.onDelete(s, entry, item); err != nil {
								s.logger.Error("failed to run on delete logic in store watcher", zap.String("key", key), zap.Error(err))
							}
						}

						s.mu.Delete(key)
					}()

				case jetstream.KeyValuePut:
					// Parse and set value "locally"
					func() {
						mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
						mu.Lock()
						defer mu.Unlock()

						if _, err := s.update(entry); err != nil {
							s.logger.Error("failed to run on update logic in store watcher", zap.String("key", key), zap.Error(err))
						}
					}()

				default:
					s.logger.Warn("unknown key operation received", zap.String("key", key), zap.Uint8("op", uint8(entry.Operation())))
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
	var err error
	if load {
		existing, err = s.load(ctx, key)
		if err != nil {
			// Return if it isn't a not found value
			if !errors.Is(err, jetstream.ErrKeyNotFound) {
				return err
			} else if s.onNotFound != nil {
				// Try to load value using not found method
				existing, err = s.onNotFound(s, ctx, key)
				if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
					return err
				}
			}
		}
	} else {
		existing, err = s.get(key)
		if err != nil {
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
func (s *Store[T, U]) Get(key string) (U, error) {
	key = s.prefix + events.SanitizeKey(key)

	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	return s.get(key)
}

func (s *Store[T, U]) get(key string) (U, error) {
	i, ok := s.data.Load(key)
	if !ok {
		return nil, jetstream.ErrKeyNotFound
	}

	return proto.Clone(i).(U), nil
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

// GetBySegmentOne finds exactly one entry whose key segments
// start with the given prefix.  E.g. prefix="JOB.GRADE" will match
// only "JOB.GRADE.USER_ID", not "JOB.GRADE2.X" or "JOB.GRADE_USER_ID".
func (s *Store[T, U]) GetBySegmentOne(prefix string) (U, error) {
	// ensure we only match whole segments:
	//   "JOB.GRADE" -> "JOB.GRADE."
	//   ""           -> ""  (no prefix => match everything)
	seg := prefix
	if seg != "" && !strings.HasSuffix(seg, ".") {
		seg = seg + "."
	}

	// 1) find matching keys in-memory
	var candidates []string
	s.data.Range(func(internalKey string, _ U) bool {
		// strip your store.prefix to get the user-key
		userKey := internalKey
		if s.prefix != "" {
			userKey = strings.TrimPrefix(internalKey, s.prefix)
		}
		if seg == "" || strings.HasPrefix(userKey, seg) {
			candidates = append(candidates, userKey)
		}
		return true
	})

	switch len(candidates) {
	case 0:
		var zero U
		return zero, jetstream.ErrKeyNotFound

	case 1:
		// delegate your normal Get (cache → KV)
		return s.Get(candidates[0])

	default:
		return *new(U), fmt.Errorf("%w: %v", ErrPrefixAmbiguous, candidates)
	}
}

func (s *Store[T, U]) GetOrLoad(ctx context.Context, key string) (U, error) {
	key = s.prefix + events.SanitizeKey(key)

	mu, _ := s.mu.LoadOrCompute(key, mutexCompute)
	mu.Lock()
	defer mu.Unlock()

	item, err := s.get(key)
	if err != nil || item == nil {
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
		return nil, fmt.Errorf("failed to unmarshal store watcher update. %w", err)
	}

	var err error
	item := s.updateFromType(entry.Key(), data)
	if s.onUpdate != nil {
		item, err = s.onUpdate(s, item)
	}

	return proto.Clone(item).(U), err
}

func (s *Store[T, U]) updateFromType(key string, incoming U) U {
	current, loaded := s.data.LoadOrStore(key, incoming)
	if loaded && current != nil {
		// Compare using protobuf magic and merge if not equal
		if !proto.Equal(current, incoming) {
			current.Merge(incoming)
		}

		return current
	}

	return incoming
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

func (s *Store[T, U]) Keys(prefix string) []string {
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
				if s.prefix != "" {
					after, _ := strings.CutPrefix(key, s.prefix)
					key = after
				}

				if s.ignoredKeys != nil && slices.Contains(s.ignoredKeys, key) {
					return true
				}

				keys = append(keys, key)
			}
		} else {
			keys = append(keys, key)
		}

		return true
	})

	return keys
}

// KeysFiltered streams directly over the cache and applies both
// the prefix and the arbitrary filter in one pass.
func (s *Store[T, U]) KeysFiltered(prefix string, filter func(string) bool) []string {
	// Prepare the full internal prefix
	full := s.prefix
	if prefix != "" {
		full += events.SanitizeKey(prefix) + "."
	}

	var out []string
	s.data.Range(func(internalKey string, _ U) bool {
		// must start with the bucket-prefix + user prefix
		if full != "" && !strings.HasPrefix(internalKey, full) {
			return true
		}

		// strip off store.prefix so we return the “user” key
		userKey := internalKey
		if s.prefix != "" {
			userKey = strings.TrimPrefix(internalKey, s.prefix)
		}

		if s.ignoredKeys != nil && slices.Contains(s.ignoredKeys, userKey) {
			return true
		}

		if filter(userKey) {
			out = append(out, userKey)
		}
		return true
	})

	return out
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

		if s.ignoredKeys != nil && slices.Contains(s.ignoredKeys, key) {
			return true
		}

		item, err := s.Get(key)
		if err != nil {
			return true
		}

		list = append(list, item)
		return true
	})

	return list
}

func (s *Store[T, U]) ListFiltered(ctx context.Context, prefix string, filter func(string) bool) []U {
	// Prepare the full internal prefix
	full := s.prefix
	if prefix != "" {
		full += events.SanitizeKey(prefix) + "."
	}

	var list []U
	s.data.Range(func(internalKey string, val U) bool {
		// must start with the bucket-prefix + user prefix
		if full != "" && !strings.HasPrefix(internalKey, full) {
			return true
		}

		// strip off store.prefix so we return the “user” key
		userKey := internalKey
		if s.prefix != "" {
			userKey = strings.TrimPrefix(internalKey, s.prefix)
		}

		if s.ignoredKeys != nil && slices.Contains(s.ignoredKeys, userKey) {
			return true
		}

		if filter(userKey) {
			list = append(list, val)
		}
		return true
	})

	return list
}

func (s *Store[T, U]) Count() int {
	return s.data.Size()
}

// Range calls fn on each cached value that passes the ignore/prefix check.
// It touches every entry exactly once and never clones unless the caller asks.
func (s *Store[T, U]) Range(fn func(key string, value U) bool) {
	if s.data == nil {
		return
	}

	var skip map[string]struct{}
	if len(s.ignoredKeys) > 0 {
		skip = make(map[string]struct{}, len(s.ignoredKeys))
		for _, k := range s.ignoredKeys {
			skip[k] = struct{}{}
		}
	}

	s.data.Range(func(internalKey string, _ U) bool {
		// strip store prefix so we expose the “user key”
		userKey := internalKey
		if s.prefix != "" {
			userKey = strings.TrimPrefix(internalKey, s.prefix)
		}
		if _, ok := skip[userKey]; ok {
			return true // ignore
		}
		v, err := s.Get(internalKey)
		if err != nil {
			return true // ignore errors, just skip this entry
		}

		// Clone the value to ensure no parallel access issues
		return fn(userKey, proto.Clone(v).(U))
	})
}

func (s *Store[T, U]) Clear(ctx context.Context) error {
	keys := s.Keys("")

	errs := multierr.Combine()
	for _, key := range keys {
		if s.ignoredKeys != nil && slices.Contains(s.ignoredKeys, key) {
			continue
		}

		if err := s.Delete(ctx, key); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (s *Store[T, U]) ReadOnly() StoreRO[T, U] {
	return s
}

func (s *Store[T, U]) WatchAll(ctx context.Context) (chan *KeyValueEntry[T, U], error) {
	watcher, err := s.kv.Watch(ctx, s.prefix+">", jetstream.UpdatesOnly())
	if err != nil {
		return nil, fmt.Errorf("failed to start kv watch in store. %w", err)
	}

	ch := make(chan *KeyValueEntry[T, U], 100)

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
				key := entry.Key()
				if s.ignoredKeys != nil && slices.Contains(s.ignoredKeys, key) {
					continue
				}

				ch <- &KeyValueEntry[T, U]{
					key:       key,
					operation: entry.Operation(),
					value:     entry.Value(),
					revision:  entry.Revision(),
				}
			}
		}
	}()

	return ch, nil
}

type KeyValueEntry[T any, U protoutils.ProtoMessageWithMerge[T]] struct {
	key       string
	operation jetstream.KeyValueOp
	value     []byte
	revision  uint64
}

// Key is the name of the key that was retrieved.
func (k *KeyValueEntry[T, U]) Key() string {
	return k.key
}

// Operation is the operation that was performed on the key.
func (k *KeyValueEntry[T, U]) Operation() jetstream.KeyValueOp {
	return k.operation
}

// Value is the value of the key that was retrieved.
func (k *KeyValueEntry[T, U]) Value() (U, error) {
	data := U(new(T))
	if err := proto.Unmarshal(k.value, data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal store watcher update. %w", err)
	}

	return data, nil
}
