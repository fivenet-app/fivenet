package store

import (
	"errors"
	"strings"

	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
)

type Store[T proto.Message] struct {
	kv         nats.KeyValue
	prefix     string
	pathPrefix string
}

func New[T proto.Message](prefix string) *Store[T] {
	return &Store[T]{
		prefix: prefix,
	}
}

func NewWithPathPrefix[T proto.Message](prefix string, pathPrefix string) *Store[T] {
	return &Store[T]{
		prefix:     prefix,
		pathPrefix: pathPrefix,
	}
}

func (s *Store[T]) Start(js nats.JetStreamContext) error {
	bucket := "fivenet"
	if s.prefix != "" {
		bucket += "_" + s.prefix
	}

	kv, err := js.KeyValue(bucket)
	if err != nil {
		if !errors.Is(nats.ErrBucketNotFound, err) {
			return err
		}

		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      bucket,
			Description: "FiveNet",
		})
		if err != nil {
			return err
		}
	}

	s.kv = kv

	return nil
}

func (s *Store[T]) prefixKey(key string) string {
	if s.pathPrefix != "" {
		return s.pathPrefix + "/" + key
	}

	return key
}

func (s *Store[T]) Get(key string, msg T) error {
	entry, err := s.kv.Get(s.prefixKey(key))
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(entry.Value(), msg); err != nil {
		return err
	}

	return nil
}

func (s *Store[T]) Put(key string, msg T) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := s.kv.Put(s.prefixKey(key), data); err != nil {
		return err
	}

	return nil
}

func (s *Store[T]) Delete(key string) error {
	return s.kv.Delete(s.prefixKey(key))
}

func (s *Store[T]) Keys() ([]string, error) {
	keys, err := s.kv.Keys(nats.MetaOnly())
	if err != nil {
		return nil, err
	}

	slices.Sort(keys)

	return keys, nil
}

func (s *Store[T]) KeysWithPrefix(prefix string) ([]string, error) {
	keys, err := s.Keys()
	if err != nil {
		return nil, err
	}

	if s.pathPrefix != "" {
		prefix = s.pathPrefix + "/" + prefix
	}

	for i := len(keys) - 1; i >= 0; i-- {
		if !strings.HasPrefix(keys[i], prefix) {
			keys = append(keys[:i], keys[i+1:]...)
		}
	}

	slices.Sort(keys)

	return keys, nil
}
