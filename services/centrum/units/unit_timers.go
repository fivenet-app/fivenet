package units

import (
	"context"
	"errors"
	"fmt"
	"time"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/nats-io/nats.go/jetstream"
)

// unitNeedsPing reports whether a unit needs another liveness/membership
// check. The ping KV is a short-lived supervision queue, not an index of all
// units.
func unitNeedsPing(unit *centrumunits.Unit) bool {
	if unit == nil {
		return false
	}

	// Units with members must continue to be checked so stale/off-duty members
	// are removed. An empty unit only needs one check while it can still be in
	// an actionable state; stable empty/unavailable units need no timer.
	return len(unit.GetUsers()) > 0 ||
		unit.GetStatus() == nil ||
		unit.GetStatus().GetStatus() != centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE
}

// SyncUnitPing reconciles a unit's ephemeral supervision timer. It is also
// used by the elected housekeeper during leadership recovery because timers
// are intentionally stored in per-process memory.
func (s *UnitDB) SyncUnitPing(ctx context.Context, unit *centrumunits.Unit) error {
	if s.KVPing == nil || unit == nil || unit.GetId() <= 0 {
		return nil
	}

	key := fmt.Sprintf("ping.%d", unit.GetId())
	if !unitNeedsPing(unit) {
		if err := s.KVPing.Delete(ctx, key); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			return err
		}
		return nil
	}

	return s.UpsertWithTTL(ctx, s.KVPing, key, PingTTL)
}

func (s *UnitDB) UpsertWithTTL(
	ctx context.Context,
	kv jetstream.KeyValue,
	key string,
	ttl time.Duration,
) error {
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
