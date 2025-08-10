package units

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func (s *UnitDB) UpsertWithTTL(ctx context.Context, kv jetstream.KeyValue, key string, ttl time.Duration) error {
	if _, err := kv.Create(ctx, key, nil, jetstream.KeyTTL(ttl)); err != nil {
		if !errors.Is(err, jetstream.ErrKeyExists) {
			return err
		}
		ent, err := kv.Get(ctx, key)
		if err != nil {
			return err
		}
		_, err = kv.Update(ctx, key, nil, ent.Revision()) // Resets TTL
		return err
	}
	return nil
}
