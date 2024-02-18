package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	natsutils "github.com/galexrt/fivenet/pkg/nats"
	"github.com/galexrt/fivenet/pkg/nats/locks"
	"github.com/nats-io/nats.go"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const LockTimeout = 350 * time.Millisecond

type protoMessage[T any] interface {
	*T
	proto.Message

	Merge(in *T) *T
}

type Store[T any, U protoMessage[T]] struct {
	logger *zap.Logger
	kv     nats.KeyValue

	bucket string

	data *xsync.MapOf[string, U]

	l *locks.Locks

	OnUpdate OnUpdateFn[T, U]
	OnDelete OnDeleteFn[T, U]
}

type StoreOption[T any, U protoMessage[T]] func(s *Store[T, U]) error

type OnUpdateFn[T any, U protoMessage[T]] func(U) (U, error)
type OnDeleteFn[T any, U protoMessage[T]] func(nats.KeyValueEntry, U) error

func New[T any, U protoMessage[T]](logger *zap.Logger, js nats.JetStreamContext, bucket string, opts ...StoreOption[T, U]) (*Store[T, U], error) {
	kv, err := natsutils.CreateKeyValue(js, bucket, &nats.KeyValueConfig{
		Bucket:      bucket,
		Description: natsutils.Description,
		History:     3,
		Storage:     nats.MemoryStorage,
	})
	if err != nil {
		return nil, err
	}

	lBucket := fmt.Sprintf("%s_locks", bucket)
	lkv, err := natsutils.CreateKeyValue(js, lBucket, &nats.KeyValueConfig{
		Bucket:      lBucket,
		Description: natsutils.Description,
		History:     3,
		Storage:     nats.MemoryStorage,
		TTL:         3 * LockTimeout,
	})
	if err != nil {
		return nil, err
	}
	l, err := locks.New(logger, lkv, lBucket, 6*LockTimeout)
	if err != nil {
		return nil, err
	}

	s := &Store[T, U]{
		logger: logger.Named("store").With(zap.String("bucket", bucket)),
		bucket: bucket,
		kv:     kv,

		data: xsync.NewMapOf[string, U](),

		l: l,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Put upload the message to kv and local
func (s *Store[T, U]) Put(key string, msg U) error {
	ctx, cancel := context.WithTimeout(context.Background(), LockTimeout)
	defer cancel()

	if err := s.l.Lock(ctx, key); err != nil {
		return err
	}
	defer s.l.Unlock(ctx, key)

	if err := s.put(key, msg); err != nil {
		return err
	}

	return nil
}

func (s *Store[T, U]) put(key string, msg U) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := s.kv.Put(key, data); err != nil {
		return err
	}

	s.updateFromType(key, msg)

	return nil
}

func (s *Store[T, U]) ComputeUpdate(key string, load bool, fn func(key string, existing U) (U, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), LockTimeout)
	defer cancel()

	if err := s.l.Lock(ctx, key); err != nil {
		return err
	}
	defer s.l.Unlock(ctx, key)

	var existing U
	if load {
		var err error
		existing, err = s.Load(key)
		if err != nil && !errors.Is(err, nats.ErrKeyNotFound) {
			return err
		}
	} else {
		var ok bool
		existing, ok = s.Get(key)
		if !ok {
			return fmt.Errorf("no item for key %s found in local store", key)
		}
	}

	var cloned U
	if existing != nil {
		cloned = proto.Clone(existing).(U)
	}
	computed, err := fn(key, cloned)
	if err != nil {
		return err
	}

	// Skip nil updates for now
	if computed == nil {
		s.logger.Error("compute update returned nil, skipping store put", zap.String("key", key))
		return nil
	}

	if err := s.put(key, computed); err != nil {
		return err
	}

	return nil
}

// Get data from local data
func (s *Store[T, U]) Get(key string) (U, bool) {
	i, ok := s.data.Load(key)
	if !ok {
		return nil, ok
	}

	return i, true
}

// Load data from kv store (this will add/update any existing local data entry)
// No record will be returned as nil and not stored
func (s *Store[T, U]) Load(key string) (U, error) {
	entry, err := s.kv.Get(key)
	if err != nil {
		return nil, err
	}

	return s.update(entry)
}

func (s *Store[T, U]) update(entry nats.KeyValueEntry) (U, error) {
	data, err := s.unmarshal(entry.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal store watcher update: %w", err)
	}

	item := s.updateFromType(entry.Key(), data)
	if s.OnUpdate != nil {
		return s.OnUpdate(item)
	}

	return item, nil
}

func (s *Store[T, U]) updateFromType(key string, updated U) U {
	current, ok := s.data.LoadOrStore(key, updated)
	if ok && current != nil {
		// Compare using protobuf magic and merge if not equal
		if !proto.Equal(current, updated) {
			current.Merge(updated)
		}

		return current
	}

	return updated
}

func (s *Store[T, U]) GetOrLoad(key string) (U, error) {
	i, ok := s.Get(key)
	if !ok || i == nil {
		var err error
		i, err = s.Load(key)
		if err != nil {
			return nil, err
		}

		return s.updateFromType(key, i), nil
	}

	return i, nil
}

func (s *Store[T, U]) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), LockTimeout)
	defer cancel()

	if err := s.l.Lock(ctx, key); err != nil {
		return err
	}
	defer s.l.Unlock(ctx, key)

	if err := s.kv.Delete(key); err != nil {
		return err
	}

	return nil
}

func (s *Store[T, U]) Keys(prefix string) ([]string, error) {
	keys, err := s.kv.Keys(nats.MetaOnly())
	if err != nil {
		return nil, err
	}

	if prefix == "" {
		return keys, nil
	}

	filtered := []string{}
	for i := 0; i < len(keys); i++ {
		if strings.HasPrefix(keys[i], prefix+".") {
			filtered = append(filtered, keys[i])
		}
	}

	return filtered, nil
}

func (s *Store[T, U]) unmarshal(data []byte) (U, error) {
	msg := U(new(T))
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *Store[T, U]) Start(ctx context.Context) error {
	watcher, err := s.kv.WatchAll(nats.Context(ctx))
	if err != nil {
		return err
	}

	go func() {
		updateCh := watcher.Updates()
		for {
			select {
			case <-ctx.Done():
				watcher.Stop()
				return

			case entry := <-updateCh:
				// After all initial keys have been received, a nil entry is returned
				if entry == nil {
					continue
				}

				s.logger.Debug("key update received via watcher", zap.String("key", entry.Key()))

				if entry.Operation() == nats.KeyValueDelete || entry.Operation() == nats.KeyValuePurge {
					if s.OnDelete != nil {
						item, _ := s.data.Load(entry.Key())
						if err := s.OnDelete(entry, item); err != nil {
							s.logger.Error("failed to run on update logic in store watcher", zap.Error(err))
							continue
						}
					}
				} else {
					if _, err := s.update(entry); err != nil {
						s.logger.Error("failed to run on update logic in store watcher", zap.Error(err))
						continue
					}
				}
			}
		}
	}()

	return nil
}

func (s *Store[T, U]) GetBucket() string {
	return s.bucket
}
