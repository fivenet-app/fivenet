package units

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

// createMemoryBucket creates/updates a memory KV bucket with a suitable LimitMarkerTTL.
func createMemoryBucket(ctx context.Context, js *events.JSWrapper, name string, ttl time.Duration) (jetstream.KeyValue, error) {
	cfg := jetstream.KeyValueConfig{
		Bucket:         name,
		Storage:        jetstream.MemoryStorage,
		History:        1,
		TTL:            0,       // Exclusively per-key TTL
		LimitMarkerTTL: 2 * ttl, // Tombstones live 2× rule
	}
	kv, err := js.CreateOrUpdateKeyValue(ctx, cfg)
	return kv, err
}

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
