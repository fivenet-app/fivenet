package cache

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// metricDataMapCount is a Prometheus gauge for tracking the number of entries in the cache data map.
var metricDataMapCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "nats_cache",
	Name:      "datamap_count",
	Help:      "Count of data map entries.",
}, []string{"bucket"})

// Cache provides a generic, in-memory and NATS-backed cache for protobuf messages.
type Cache[T any, U protoutils.ProtoMessage[T]] struct {
	// logger for cache operations
	logger *zap.Logger
	// bucket name for the NATS KeyValue store
	bucket string
	// NATS JetStream KeyValue store instance
	kv jetstream.KeyValue

	// prefix for all cache keys (optional)
	prefix string
	// list of keys to ignore in cache operations
	ignoredKeys []string
	// optional time-to-live for cache entries
	ttl *time.Duration

	// concurrent map holding cached entries
	data *xsync.Map[string, *EntryWrapper[T, U]]
}

// EntryWrapper wraps a cached protobuf message and its creation time.
type EntryWrapper[T any, U protoutils.ProtoMessage[T]] struct {
	// cached protobuf message
	Data U
	// time when the entry was created
	Created time.Time
}

// Option is a functional option for configuring the Cache.
type Option[T any, U protoutils.ProtoMessage[T]] func(s *Cache[T, U])

// New creates a new Cache instance with the given options and bucket.
func New[T any, U protoutils.ProtoMessage[T]](ctx context.Context, logger *zap.Logger, js *events.JSWrapper, bucket string, opts ...Option[T, U]) (*Cache[T, U], error) {
	c := &Cache[T, U]{
		logger: logger.Named("cache").With(zap.String("bucket", bucket)),
		bucket: bucket,
		data:   xsync.NewMap[string, *EntryWrapper[T, U]](),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.kv == nil {
		storeKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:      bucket,
			Description: fmt.Sprintf("%s Cache", bucket),
			History:     1,
			Storage:     jetstream.MemoryStorage,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create kv (bucket %s) for cache. %w", bucket, err)
		}
		c.kv = storeKV
	}

	if c.prefix != "" {
		c.prefix = strings.TrimSuffix(c.prefix, ".") + "."
	}

	return c, nil
}

// Start begins watching the NATS KeyValue store for changes and updates the cache accordingly.
// If wait is true, it blocks until the initial sync is complete.
func (c *Cache[T, U]) Start(ctx context.Context, wait bool) error {
	watcher, err := c.kv.Watch(ctx, c.prefix+">")
	if err != nil {
		return fmt.Errorf("failed to start cache kv. %w", err)
	}

	var ready atomic.Bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer watcher.Stop()
		updateCh := watcher.Updates()

		for {
			select {
			case <-ctx.Done():
				return
			case entry := <-updateCh:
				if entry == nil {
					if !ready.Swap(true) {
						wg.Done()
					}
					continue
				}

				key := entry.Key()
				if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, key) {
					continue
				}

				switch entry.Operation() {
				case jetstream.KeyValuePut:
					c.handleWatcherPut(key, entry)
				case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
					c.handleWatcherDelete(key)
				default:
					c.logger.Error("unknown key operation received", zap.String("key", key), zap.Uint8("op", uint8(entry.Operation())))
				}
			}
		}
	}()

	go func() {
		for {
			// Update Prometheus metric with current cache size
			metricDataMapCount.WithLabelValues(c.bucket).Set(float64(c.data.Size()))

			select {
			case <-ctx.Done():
				return
			case <-time.After(15 * time.Second):
			}
		}
	}()

	if wait {
		wg.Wait()
	}

	return nil
}

// handleWatcherPut handles a put/update event from the NATS KeyValue watcher.
func (c *Cache[T, U]) handleWatcherPut(key string, entry jetstream.KeyValueEntry) {
	// Unmarshal the value and store it in the cache with its creation time
	data := U(new(T))
	if err := proto.Unmarshal(entry.Value(), data); err != nil {
		c.logger.Error("failed to unmarshal cache update", zap.String("key", key), zap.Error(err))
		return
	}

	c.data.Store(key, &EntryWrapper[T, U]{
		Data:    data,
		Created: entry.Created(),
	})
}

// handleWatcherDelete handles a delete event from the NATS KeyValue watcher.
func (c *Cache[T, U]) handleWatcherDelete(key string) {
	c.data.Delete(key)
}

// prefixed returns the full cache key with prefix and sanitized user key.
func (c *Cache[T, U]) prefixed(key string) string {
	return c.prefix + events.SanitizeKey(key)
}

// Has checks if a key exists in the cache and is not expired.
func (c *Cache[T, U]) Has(key string) bool {
	key = c.prefixed(key)
	v, ok := c.data.Load(key)
	if !ok || v == nil {
		return false
	}
	if c.ttl != nil && time.Since(v.Created) > *c.ttl {
		return false
	}
	return true
}

// Get retrieves a value from the cache by key, returning an error if not found or expired.
func (c *Cache[T, U]) Get(key string) (U, error) {
	key = c.prefixed(key)

	value, ok := c.data.Load(key)
	if !ok || value == nil {
		return nil, jetstream.ErrKeyNotFound
	}

	if c.ttl != nil && time.Since(value.Created) > *c.ttl {
		return nil, jetstream.ErrKeyNotFound
	}

	return proto.Clone(value.Data).(U), nil
}

// Keys returns a list of user keys in the cache, optionally filtered by prefix.
func (c *Cache[T, U]) Keys(prefix string) []string {
	fullPrefix := c.prefix
	if prefix != "" {
		fullPrefix += events.SanitizeKey(prefix) + "."
	}

	var keys []string
	c.data.Range(func(key string, _ *EntryWrapper[T, U]) bool {
		if fullPrefix != "" && !strings.HasPrefix(key, fullPrefix) {
			return true
		}

		userKey := strings.TrimPrefix(key, c.prefix)
		if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, userKey) {
			return true
		}

		keys = append(keys, userKey)
		return true
	})

	return keys
}

// List returns all values in the cache as a slice.
func (c *Cache[T, U]) List() []U {
	var list []U
	c.data.Range(func(internalKey string, value *EntryWrapper[T, U]) bool {
		if value == nil {
			return true
		}

		userKey := strings.TrimPrefix(internalKey, c.prefix)
		if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, userKey) {
			return true
		}

		if c.ttl != nil && time.Since(value.Created) > *c.ttl {
			return true
		}

		list = append(list, proto.Clone(value.Data).(U))
		return true
	})

	return list
}

// Range executes the given function for each key-value pair in the cache.
func (c *Cache[T, U]) Range(fn func(key string, value U) bool) {
	c.data.Range(func(internalKey string, wrapper *EntryWrapper[T, U]) bool {
		if wrapper == nil {
			return true
		}

		userKey := strings.TrimPrefix(internalKey, c.prefix)
		if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, userKey) {
			return true
		}

		if c.ttl != nil && time.Since(wrapper.Created) > *c.ttl {
			return true
		}

		return fn(userKey, proto.Clone(wrapper.Data).(U))
	})
}

// Put adds or updates a value in the cache, storing it in the NATS KeyValue store.
func (c *Cache[T, U]) Put(ctx context.Context, key string, val U) error {
	key = c.prefixed(key)

	data, err := proto.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value %s for cache. %w", key, err)
	}

	if _, err := c.kv.Put(ctx, key, data); err != nil {
		return fmt.Errorf("failed to put value %s in cache. %w", key, err)
	}

	return nil
}

// Delete removes a value from the cache and the NATS KeyValue store.
func (c *Cache[T, U]) Delete(ctx context.Context, key string) error {
	if err := c.kv.Delete(ctx, c.prefixed(key)); err != nil {
		return fmt.Errorf("failed to delete value %s from cache. %w", key, err)
	}
	return nil
}

// Clear removes all values from the cache and the NATS KeyValue store.
func (c *Cache[T, U]) Clear(ctx context.Context) error {
	var errs error
	for _, key := range c.Keys("") {
		if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, key) {
			continue
		}
		if err := c.Delete(ctx, key); err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}
