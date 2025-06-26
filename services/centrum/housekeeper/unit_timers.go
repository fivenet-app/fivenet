package housekeeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *Housekeeper) runTTLWatcher(
	kv jetstream.KeyValue,
	prefix string,
	handler func(ctx context.Context, unitID uint64) error,
) {
	sub, _ := kv.Watch(s.ctx, prefix+".*")
	for {
		select {
		case <-s.ctx.Done():
			return
		case e := <-sub.Updates():
			if e == nil || (e.Operation() != jetstream.KeyValueDelete &&
				e.Operation() != jetstream.KeyValuePurge) {
				continue
			}
			id, _ := strconv.ParseUint(strings.TrimPrefix(string(e.Key()), prefix+"."), 10, 64)
			if err := handler(s.ctx, id); err != nil {
				s.logger.Error("failed to handle TTL event", zap.Error(err), zap.Uint64("unit_id", id))
			}
		}
	}
}

// handleUnitKVPing checks if a unit is empty or has the static attribute and sets its status to unavailable if so.
func (s *Housekeeper) handleUnitKVPing(ctx context.Context, id uint64) error {
	unit, err := s.units.Get(ctx, id)
	if err != nil {
		return err
	}

	// Fast check if unit is already unavailable and empty
	if unit.Status != nil &&
		unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE &&
		len(unit.Users) == 0 {
		return nil
	}

	// Check and verify the users in the unit and in case of changes, "schedule" another ping check
	toRemove, changed, err := s.checkAndUpdateUnitUsers(ctx, unit)
	if err != nil {
		s.logger.Error("failed to check users in unit", zap.Error(err))
	}
	stillAvailable := changed

	// If all users are still valid (toRemove is empty) and status checks pass, keep the unit available
	if len(unit.Users) > 0 && len(toRemove) == 0 {
		if unit.Attributes == nil || !unit.Attributes.Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
			stillAvailable = true
		} else if unit.Status != nil &&
			(unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_BUSY ||
				unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_ON_BREAK ||
				unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE) {
			stillAvailable = true
		}
	}

	if stillAvailable {
		// Reset the ping timer
		return s.resetUnitPing(ctx, id)
	}

	var userId *int32
	if unit.Status != nil && unit.Status.UserId != nil {
		userId = unit.Status.UserId
	}

	s.logger.Debug("setting unit status to unavailable it is empty or static attribute (wrong status)",
		zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Int32p("user_id", userId))
	if _, err := s.units.UpdateStatus(ctx, unit.Id, &centrum.UnitStatus{
		CreatedAt:  timestamp.Now(),
		UnitId:     unit.Id,
		Status:     centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
		UserId:     userId,
		CreatorJob: &unit.Job,
	}); err != nil {
		s.logger.Error("failed to update empty unit status to unavailable",
			zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Error(err))
		return nil
	}

	return nil
}

func (s *Housekeeper) resetUnitPing(ctx context.Context, id uint64) error {
	// Reset unit ping timer
	if err := s.units.UpsertWithTTL(ctx, s.units.KVPing, fmt.Sprintf("ping.%d", id), units.PingTTL); err != nil {
		return fmt.Errorf("failed to upsert ping unit timer. %w", err)
	}
	return nil
}
