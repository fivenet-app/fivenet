package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *Housekeeper) runTTLWatcher(ctx context.Context) error {
	for {
		if err := s.unitKVPing(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.logger.Error("unit ping watcher stopped", zap.Error(err))
			}
		}

		select {
		case <-s.ctx.Done():
			return nil

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) unitKVPing(ctx context.Context) error {
	watch, err := s.units.KVPing.Watch(ctx, "ping.*")
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case e := <-watch.Updates():
			if e == nil || (e.Operation() != jetstream.KeyValueDelete &&
				e.Operation() != jetstream.KeyValuePurge) {
				continue
			}

			id, _ := strconv.ParseInt(strings.TrimPrefix(e.Key(), "ping."), 10, 64)
			if err := s.handleUnitKVPing(ctx, id); err != nil {
				s.logger.Error(
					"failed to handle TTL event",
					zap.Error(err),
					zap.Int64("unit_id", id),
				)
			}
		}
	}
}

// handleUnitKVPing checks if a unit is empty or has the static attribute and sets its status to unavailable if so.
func (s *Housekeeper) handleUnitKVPing(ctx context.Context, id int64) error {
	unit, err := s.units.Get(ctx, id)
	if err != nil {
		return err
	}

	// Fast check if unit is already unavailable and empty
	if unit.GetStatus() != nil &&
		unit.GetStatus().GetStatus() == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE &&
		len(unit.GetUsers()) == 0 {
		return nil
	}

	// Check and verify the users in the unit and in case of changes, "schedule" another ping check
	toRemove, changed, err := s.checkAndUpdateUnitUsers(ctx, unit)
	if err != nil {
		s.logger.Error("failed to check users in unit", zap.Error(err))
	}
	stillAvailable := changed

	// If all users are still valid (toRemove is empty) and status checks pass, keep the unit available
	if len(unit.GetUsers()) > 0 && len(toRemove) == 0 {
		if unit.GetAttributes() == nil ||
			!unit.GetAttributes().Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
			stillAvailable = true
		} else if unit.GetStatus() != nil &&
			(unit.GetStatus().GetStatus() == centrum.StatusUnit_STATUS_UNIT_BUSY ||
				unit.GetStatus().GetStatus() == centrum.StatusUnit_STATUS_UNIT_ON_BREAK ||
				unit.GetStatus().GetStatus() == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE) {
			stillAvailable = true
		}
	}

	if stillAvailable {
		// Reset the ping timer
		return s.resetUnitPing(ctx, id)
	}

	var userId *int32
	if unit.GetStatus() != nil && unit.Status.UserId != nil {
		userId = unit.GetStatus().UserId
	}

	s.logger.Debug(
		"setting unit status to unavailable it is empty or static attribute (wrong status)",
		zap.String("job", unit.GetJob()),
		zap.Int64("unit_id", unit.GetId()),
		zap.Int32p("user_id", userId),
	)
	if _, err := s.units.UpdateStatus(ctx, unit.GetId(), &centrum.UnitStatus{
		CreatedAt:  timestamp.Now(),
		UnitId:     unit.GetId(),
		Unit:       proto.Clone(unit).(*centrum.Unit),
		Status:     centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
		UserId:     userId,
		CreatorJob: &unit.Job,
	}); err != nil {
		s.logger.Error("failed to update empty unit status to unavailable",
			zap.String("job", unit.GetJob()), zap.Int64("unit_id", unit.GetId()), zap.Error(err))
		return nil
	}

	return nil
}

func (s *Housekeeper) resetUnitPing(ctx context.Context, id int64) error {
	// Reset unit ping timer
	if err := s.units.UpsertWithTTL(ctx, s.units.KVPing, fmt.Sprintf("ping.%d", id), units.PingTTL); err != nil {
		return fmt.Errorf("failed to upsert ping unit timer. %w", err)
	}
	return nil
}
