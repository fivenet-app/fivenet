package store

import (
	"errors"

	"github.com/galexrt/fivenet/pkg/events"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Store struct {
	kv nats.KeyValue
}

func New(e *events.Eventus) (*Store, error) {
	kv, err := e.JS.KeyValue("fivenet_store")
	if errors.Is(nats.ErrBucketNotFound, err) {
		kv, err = e.JS.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "fivenet_store",
		})
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		kv: kv,
	}, nil
}

func (s *Store) Put(key string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := s.kv.Put(key, data); err != nil {
		return err
	}

	return nil
}

func (s *Store) Get(key string, msg proto.Message) (proto.Message, error) {
	entry, err := s.kv.Get(key)
	if err != nil {
		return nil, err
	}

	if err := proto.Unmarshal(entry.Value(), msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *Store) Keys() ([]string, error) {
	keys, err := s.kv.Keys(nats.MetaOnly())
	if err != nil {
		return nil, err
	}

	return keys, nil
}
