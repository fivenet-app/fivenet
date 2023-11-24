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

const LockTimeout = 3 * time.Second

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

	onUpdate func(nats.KeyValueEntry) (U, error)
}

func New[T any, U protoMessage[T]](logger *zap.Logger, js nats.JetStreamContext, bucket string) (*Store[T, U], error) {
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
	})
	if err != nil {
		return nil, err
	}
	l, err := locks.New(logger, lkv, lBucket)
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

	s.onUpdate = s.update

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
	var existing U
	if load {
		var err error
		existing, err = s.Load(key)
		if err != nil && !errors.Is(err, nats.ErrKeyNotFound) {
			return err
		}
	} else {
		existing, _ = s.Get(key)
	}

	ctx, cancel := context.WithTimeout(context.Background(), LockTimeout)
	defer cancel()

	if err := s.l.Lock(ctx, key); err != nil {
		return err
	}
	defer s.l.Unlock(ctx, key)

	computed, err := fn(key, existing)
	if err != nil {
		return err
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

	return s.onUpdate(entry)
}

func (s *Store[T, U]) update(entry nats.KeyValueEntry) (U, error) {
	data, err := s.unmarshal(entry.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal store watcher update: %w", err)
	}

	return s.updateFromType(entry.Key(), data), nil
}

func (s *Store[T, U]) updateFromType(key string, data U) U {
	current, ok := s.data.LoadOrStore(key, data)
	if ok && current != nil {
		// Compare using protobuf magic and merge if not equal
		if !proto.Equal(current, data) {
			current.Merge(data)
		}

		return current
	}

	return data
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

				if _, err := s.onUpdate(entry); err != nil {
					s.logger.Error("failed to run on update logic in store watcher", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

func (s *Store[T, U]) Locks() *locks.Locks {
	return s.l
}

func (s *Store[T, U]) GetBucket() string {
	return s.bucket
}
