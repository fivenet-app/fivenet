package cache

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

var metricDataMapCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "nats_cache",
	Name:      "datamap_count",
	Help:      "Count of data map entries.",
}, []string{"bucket"})

type Cache[T any, U protoutils.ProtoMessage[T]] struct {
	logger *zap.Logger
	bucket string
	kv     jetstream.KeyValue

	prefix      string
	ignoredKeys []string
	ttl         *time.Duration

	data *xsync.Map[string, *EntryWrapper[T, U]]
}

type EntryWrapper[T any, U protoutils.ProtoMessage[T]] struct {
	Data    U
	Created time.Time
}

type Option[T any, U protoutils.ProtoMessage[T]] func(s *Cache[T, U])

func New[T any, U protoutils.ProtoMessage[T]](ctx context.Context, logger *zap.Logger, js *events.JSWrapper, bucket string, opts ...Option[T, U]) (*Cache[T, U], error) {
	c := &Cache[T, U]{
		logger: logger.Named("cache").With(zap.String("bucket", bucket)),
		bucket: bucket,

		data: xsync.NewMap[string, *EntryWrapper[T, U]](),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.kv == nil {
		storeKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:      bucket,
			Description: fmt.Sprintf("%s Store", bucket),
			History:     1,
			Storage:     jetstream.MemoryStorage,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create kv (bucket %s) for store. %w", bucket, err)
		}

		c.kv = storeKV
	}

	if c.prefix != "" {
		c.prefix = strings.TrimSuffix(c.prefix, ".") + "."
	}

	return c, nil
}

func (c *Cache[T, U]) Start(ctx context.Context, wait bool) error {
	watcher, err := c.kv.Watch(ctx, c.prefix+">")
	if err != nil {
		return fmt.Errorf("failed to start cache kv. %w", err)
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
						c.logger.Error("error while stopping watcher", zap.Error(err))
					}
				} else {
					c.logger.Debug("store watcher done")
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
				if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, key) {
					continue
				}

				switch entry.Operation() {
				case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
					// Handle delete and purge operations
					c.data.Delete(key)

				case jetstream.KeyValuePut:
					// Parse and set value "locally"
					data := U(new(T))
					if err := proto.Unmarshal(entry.Value(), data); err != nil {
						c.logger.Error("failed to unmarshal store watcher update", zap.String("key", key), zap.Error(err))
					}

					c.data.Store(key, &EntryWrapper[T, U]{
						Data:    data,
						Created: entry.Created(),
					})

				default:
					c.logger.Error("unknown key operation received", zap.String("key", key), zap.Uint8("op", uint8(entry.Operation())))
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
				metricDataMapCount.WithLabelValues(c.bucket).Set(float64(c.data.Size()))
			}
		}
	}()

	if wait {
		wg.Wait()
	}

	return nil
}

func (c *Cache[T, U]) Get(key string) (U, error) {
	key = c.prefix + events.SanitizeKey(key)

	value, ok := c.data.Load(key)
	if !ok {
		return nil, jetstream.ErrKeyNotFound
	}

	// TTL expired
	if c.ttl != nil && time.Since(value.Created) > *c.ttl {
		return nil, jetstream.ErrKeyNotFound
	}

	return value.Data, nil
}

func (c *Cache[T, U]) Keys(prefix string) []string {
	hasPrefix := (c.prefix + prefix) != ""
	if hasPrefix {
		if prefix != "" {
			prefix += "."
		}
		prefix = c.prefix + events.SanitizeKey(prefix)
	}

	keys := []string{}
	c.data.Range(func(key string, _ *EntryWrapper[T, U]) bool {
		if hasPrefix {
			if strings.HasPrefix(key, prefix) {
				if c.prefix != "" {
					after, _ := strings.CutPrefix(key, c.prefix)
					key = after
				}

				if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, key) {
					return true
				}
			}
		} else {
			keys = append(keys, key)
		}

		return true
	})

	return keys
}

func (c *Cache[T, U]) List() []U {
	list := []U{}

	c.data.Range(func(key string, value *EntryWrapper[T, U]) bool {
		if value == nil {
			return true
		}

		if c.prefix != "" {
			key, _ = strings.CutPrefix(key, c.prefix)
		}

		if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, key) {
			return true
		}

		item, err := c.Get(key)
		if err != nil {
			return true
		}

		list = append(list, item)
		return true
	})

	return list
}

func (c *Cache[T, U]) Range(fn func(key string, value U) bool) {
	keys := c.Keys("")

	for _, key := range keys {
		if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, key) {
			continue
		}

		v, err := c.Get(key)
		if err != nil {
			continue
		}

		if !fn(key, v) {
			break
		}
	}
}

func (c *Cache[T, U]) Put(ctx context.Context, key string, val U) error {
	key = c.prefix + events.SanitizeKey(key)

	out, err := proto.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value %s for cache. %w", key, err)
	}

	if _, err := c.kv.Put(ctx, key, out); err != nil {
		return fmt.Errorf("failed to put value %s in cache. %w", key, err)
	}

	return nil
}

// Local deletion will happen via the watcher which might be delayed but this is just a cache
func (c *Cache[T, U]) Delete(ctx context.Context, key string) error {
	if err := c.kv.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete value %s from cache. %w", key, err)
	}

	return nil
}

func (c *Cache[T, U]) Clear(ctx context.Context) error {
	keys := c.Keys("")

	errs := multierr.Combine()
	for _, key := range keys {
		if c.ignoredKeys != nil && slices.Contains(c.ignoredKeys, key) {
			continue
		}

		if err := c.Delete(ctx, key); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}
