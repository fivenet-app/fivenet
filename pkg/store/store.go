package store

import (
	"errors"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Store[T proto.Message] struct {
	kv nats.KeyValue
}

func (s *Store[T]) Start(js nats.JetStreamContext) error {
	kv, err := js.KeyValue("fivenet_store")
	if errors.Is(nats.ErrBucketNotFound, err) {
		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "fivenet_store",
		})
		if err != nil {
			return err
		}
	}

	s.kv = kv

	return nil
}

func (s *Store[T]) Put(key string, msg T) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := s.kv.Put(key, data); err != nil {
		return err
	}

	return nil
}

func (s *Store[T]) Get(key string) (T, error) {
	var dest T
	entry, err := s.kv.Get(key)
	if err != nil {
		return dest, err
	}

	if err := proto.Unmarshal(entry.Value(), dest); err != nil {
		return dest, err
	}

	return dest, nil
}

func (s *Store[T]) Keys() ([]string, error) {
	keys, err := s.kv.Keys(nats.MetaOnly())
	if err != nil {
		return nil, err
	}

	return keys, nil
}
