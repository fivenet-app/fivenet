package cache

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/laws"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Cache[T any, U protoutils.ProtoMessageWithMerge[T]] struct {
	logger *zap.Logger

	ap   atomic.Pointer[U]
	data *xsync.MapOf[uint64, *laws.LawBook]

	kv  jetstream.KeyValue
	key string
}

func New[T any, U protoutils.ProtoMessageWithMerge[T]](logger *zap.Logger, kv jetstream.KeyValue, prefix string) (*Cache[T, U], error) {
	c := &Cache[T, U]{
		logger: logger.Named("cache").With(zap.String("prefix", prefix)),

		ap: atomic.Pointer[U]{},

		kv:  kv,
		key: prefix,
	}

	return c, nil
}

func (c *Cache[T, U]) Start(ctx context.Context) error {
	watcher, err := c.kv.Watch(ctx, c.key)
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
					c.ap.Store(nil)

				case jetstream.KeyValuePut:
					// Parse and set value "locally"
					data := U(new(T))
					if err := proto.Unmarshal(entry.Value(), data); err != nil {
						c.logger.Error("failed to unmarshal store watcher update", zap.Error(err))
					}

					c.ap.Store(&data)

				default:
					c.logger.Error("unknown key operation received", zap.Uint8("op", uint8(entry.Operation())))
				}
			}
		}
	}()

	return nil
}

func (c *Cache[T, U]) Get() *U {
	return c.ap.Load()
}

func (c *Cache[T, U]) Set(ctx context.Context, val U) error {
	out, err := protoutils.Marshal(val)
	if err != nil {
		return err
	}

	if _, err := c.kv.Put(ctx, c.key, out); err != nil {
		return err
	}

	return nil
}
