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
	"github.com/nats-io/nats.go/jetstream"
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
	kv     jetstream.KeyValue

	bucket string

	data *xsync.MapOf[string, U]

	l *locks.Locks

	OnUpdate OnUpdateFn[T, U]
	OnDelete OnDeleteFn[T, U]
}

type StoreOption[T any, U protoMessage[T]] func(s *Store[T, U]) error

type OnUpdateFn[T any, U protoMessage[T]] func(U) (U, error)
type OnDeleteFn[T any, U protoMessage[T]] func(jetstream.KeyValueEntry, U) error

func NewWithLocks[T any, U protoMessage[T]](ctx context.Context, logger *zap.Logger, js jetstream.JetStream, bucket string, l *locks.Locks, opts ...StoreOption[T, U]) (*Store[T, U], error) {
	kv, err := natsutils.CreateKeyValue(ctx, js, jetstream.KeyValueConfig{
		Bucket:      bucket,
		Description: natsutils.Description,
		History:     3,
		Storage:     jetstream.MemoryStorage,
	})
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

func New[T any, U protoMessage[T]](ctx context.Context, logger *zap.Logger, js jetstream.JetStream, bucket string, opts ...StoreOption[T, U]) (*Store[T, U], error) {
	lBucket := fmt.Sprintf("%s_locks", bucket)
	lkv, err := natsutils.CreateKeyValue(ctx, js, jetstream.KeyValueConfig{
		Bucket:      lBucket,
		Description: natsutils.Description,
		History:     3,
		Storage:     jetstream.MemoryStorage,
		TTL:         3 * LockTimeout,
	})
	if err != nil {
		return nil, err
	}
	l, err := locks.New(logger, lkv, lBucket, 6*LockTimeout)
	if err != nil {
		return nil, err
	}

	return NewWithLocks[T, U](ctx, logger, js, bucket, l, opts...)
}

// Put upload the message to kv and local
func (s *Store[T, U]) Put(ctx context.Context, key string, msg U) error {
	ctx, cancel := context.WithTimeout(ctx, LockTimeout)
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

func (s *Store[T, U]) ComputeUpdate(ctx context.Context, key string, load bool, fn func(key string, existing U) (U, error)) error {
	ctx, cancel := context.WithTimeout(ctx, LockTimeout)
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
		existing, err = s.Load(ctx, key)
		if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
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

	if err := s.put(ctx, key, computed); err != nil {
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
func (s *Store[T, U]) Load(ctx context.Context, key string) (U, error) {
	entry, err := s.kv.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	return s.update(entry)
}

func (s *Store[T, U]) update(entry jetstream.KeyValueEntry) (U, error) {
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

func (s *Store[T, U]) GetOrLoad(ctx context.Context, key string) (U, error) {
	i, ok := s.Get(key)
	if !ok || i == nil {
		var err error
		i, err = s.Load(ctx, key)
		if err != nil {
			return nil, err
		}

		return s.updateFromType(key, i), nil
	}

	return i, nil
}

func (s *Store[T, U]) Delete(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, LockTimeout)
	defer cancel()

	if s.l != nil {
		if err := s.l.Lock(ctx, key); err != nil {
			return err
		}
		defer s.l.Unlock(ctx, key)
	}

	if err := s.kv.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}

func (s *Store[T, U]) Keys(ctx context.Context, prefix string) ([]string, error) {
	lister, err := s.kv.ListKeys(ctx, jetstream.MetaOnly())
	if err != nil {
		return nil, err
	}

	hasPrefix := prefix != ""

	keys := []string{}
	for key := range lister.Keys() {
		if hasPrefix {
			if strings.HasPrefix(key, prefix+".") {
				keys = append(keys, key)
			}
		} else {
			keys = append(keys, key)
		}
	}

	return keys, nil
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

				if entry.Operation() == jetstream.KeyValueDelete || entry.Operation() == jetstream.KeyValuePurge {
					item, _ := s.data.LoadAndDelete(entry.Key())
					if s.OnDelete != nil {
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

func (s *Store[T, U]) List() ([]U, error) {
	list := []U{}

	s.data.Range(func(key string, value U) bool {
		if value == nil {
			return true
		}

		list = append(list, value)
		return true
	})

	return list, nil
}
