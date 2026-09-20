package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/units"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *Housekeeper) runTTLWatcher(ctx context.Context) error {
	for {
		if err := s.unitKVPing(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.recordWatcherRestart("unit_ping", err)
				s.logger.Error("unit ping watcher stopped", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			return nil

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) unitKVPing(ctx context.Context) error {
	if err := s.reconcileUnitPings(ctx); err != nil {
		return err
	}

	watch, err := s.units.KVPing.Watch(ctx, "ping.*")
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case e, ok := <-watch.Updates():
			if !ok {
				return watcherUpdatesClosedError(ctx)
			}
			// Ignore nil event
			if e == nil {
				continue
			}
			// We only care about expiry events
			if e.Operation() != jetstream.KeyValueDelete &&
				e.Operation() != jetstream.KeyValuePurge {
				continue
			}

			id, err := parseUnitPingKey(e.Key())
			if err != nil {
				s.logger.Warn("ignoring malformed unit ping key", zap.String("key", e.Key()), zap.Error(err))
				continue
			}
			if err := s.handleUnitKVPing(ctx, id); err != nil {
				s.logger.Error(
					"failed to handle TTL event",
					zap.Int64("unit_id", id),
					zap.Error(err),
				)
			}
		}
	}
}

// reconcileUnitPings repairs the local in-memory timer bucket after leadership
// changes. A unit may have expired from this process while another instance
// was leader, so relying only on replayed ping keys would leave a supervision
// gap after failover.
func (s *Housekeeper) reconcileUnitPings(ctx context.Context) error {
	var errs error
	s.units.Range(func(_ string, unit *centrumunits.Unit) bool {
		if err := s.units.SyncUnitPing(ctx, unit); err != nil {
			errs = errors.Join(errs, err)
		}
		return true
	})

	return errs
}

func parseUnitPingKey(key string) (int64, error) {
	const prefix = "ping."
	if !strings.HasPrefix(key, prefix) {
		return 0, fmt.Errorf("invalid unit ping key %q", key)
	}

	id, err := strconv.ParseInt(strings.TrimPrefix(key, prefix), 10, 64)
	if err != nil || id <= 0 {
		if err == nil {
			err = fmt.Errorf("unit ID must be positive")
		}
		return 0, fmt.Errorf("invalid unit ping key %q: %w", key, err)
	}

	return id, nil
}

// handleUnitKVPing checks if a unit is empty or has the static attribute and sets its status to unavailable if so.
func (s *Housekeeper) handleUnitKVPing(ctx context.Context, unitId int64) error {
	unit, err := s.units.Get(ctx, unitId)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			if err := s.units.KVPing.Delete(ctx, fmt.Sprintf("ping.%d", unitId)); err != nil {
				return fmt.Errorf(
					"failed to delete ghost unit %d ping from kv ping store. %w",
					unitId,
					err,
				)
			}
		}
		return fmt.Errorf("failed to get unit %d from store. %w", unitId, err)
	}

	// Fast check if unit is already unavailable and empty
	if unit.GetStatus() != nil &&
		unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE &&
		len(unit.GetUsers()) == 0 {
		return nil
	}

	// Check and verify the users in the unit and in case of changes, "schedule" another ping check
	toRemove, removed, err := s.checkAndUpdateUnitUsers(ctx, unit)
	if err != nil {
		s.logger.Error("failed to check users in unit", zap.Error(err))
	}
	stillAvailable := removed > 0

	// If all users are still valid (toRemove is empty) and status checks pass, keep the unit available
	if len(unit.GetUsers()) > 0 && len(toRemove) == 0 {
		if unit.GetAttributes() == nil ||
			!unit.GetAttributes().Has(centrumunits.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
			stillAvailable = true
		} else if unit.GetStatus() != nil &&
			(unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_BUSY ||
				unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_ON_BREAK ||
				unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE) {
			stillAvailable = true
		}
	}

	if stillAvailable {
		// Reset the ping timer
		return s.resetUnitPing(ctx, unitId)
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
	if _, _, err := s.units.UpdateStatus(ctx, unit.GetId(), &centrumunits.UnitStatus{
		CreatedAt:  timestamp.Now(),
		UnitId:     unit.GetId(),
		Status:     centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE,
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
	if err := s.units.UpsertWithTTL(
		ctx,
		s.units.KVPing,
		fmt.Sprintf("ping.%d", id),
		units.PingTTL,
	); err != nil {
		return fmt.Errorf("failed to upsert ping unit timer. %w", err)
	}
	return nil
}
