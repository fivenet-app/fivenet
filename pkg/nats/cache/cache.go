package cache

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Cache[T any, U protoutils.ProtoMessage[T]] struct {
	logger *zap.Logger

	data *xsync.Map[string, *EntryWrapper[T, U]]

	kv     jetstream.KeyValue
	prefix string
	ttl    *time.Duration
}

type EntryWrapper[T any, U protoutils.ProtoMessage[T]] struct {
	Data    U
	Created time.Time
}

type Option[T any, U protoutils.ProtoMessage[T]] func(s *Cache[T, U])

func New[T any, U protoutils.ProtoMessage[T]](logger *zap.Logger, kv jetstream.KeyValue, opts ...Option[T, U]) (*Cache[T, U], error) {
	c := &Cache[T, U]{
		logger: logger.Named("cache").With(zap.String("bucket", kv.Bucket())),

		data: xsync.NewMap[string, *EntryWrapper[T, U]](),

		kv: kv,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.prefix != "" {
		c.prefix = c.prefix + "."
	}

	return c, nil
}

func (c *Cache[T, U]) Start(ctx context.Context) error {
	watcher, err := c.kv.Watch(ctx, c.prefix)
	if err != nil {
		return err
	}

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
					continue
				}

				switch entry.Operation() {
				case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
					// Value is deleted
					c.data.Delete(entry.Key())

				case jetstream.KeyValuePut:
					// Parse and set value "locally"
					data := U(new(T))
					if err := proto.Unmarshal(entry.Value(), data); err != nil {
						c.logger.Error("failed to unmarshal store watcher update", zap.String("key", entry.Key()), zap.Error(err))
					}

					c.data.Store(entry.Key(), &EntryWrapper[T, U]{
						Data:    data,
						Created: entry.Created(),
					})

				default:
					c.logger.Error("unknown key operation received", zap.String("key", entry.Key()), zap.Uint8("op", uint8(entry.Operation())))
				}
			}
		}
	}()

	return nil
}

func (c *Cache[T, U]) Get(key string) (U, bool) {
	key = c.prefix + events.SanitizeKey(key)

	value, ok := c.data.Load(c.prefix + key)
	if !ok {
		return nil, false
	}

	// TTL expired
	if c.ttl != nil && time.Since(value.Created) > *c.ttl {
		return nil, false
	}

	return value.Data, ok
}

func (c *Cache[T, U]) Set(ctx context.Context, key string, val U) error {
	key = c.prefix + events.SanitizeKey(key)

	out, err := protoutils.Marshal(val)
	if err != nil {
		return err
	}

	if _, err := c.kv.Put(ctx, key, out); err != nil {
		return err
	}

	return nil
}

// Local deletion will happen via the watcher which might be delayed but this is just a cache
func (c *Cache[T, U]) Delete(ctx context.Context, key string) error {
	if err := c.kv.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
