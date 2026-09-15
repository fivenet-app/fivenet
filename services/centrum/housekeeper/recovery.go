package housekeeper

import (
	"context"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"go.uber.org/zap"
)

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

// capWorkItems reserves one selected row to signal that another recovery run
// is needed, while keeping the work performed in this run bounded.
func capWorkItems[T any](items []T, limit int) ([]T, bool) {
	if len(items) <= limit {
		return items, false
	}

	return items[:limit], true
}

func setBacklogAttribute(data *cron.GenericCronData, key string, backlogRemaining bool) {
	data.SetAttribute(key, strconv.FormatBool(backlogRemaining))
}

// runLeadershipRecovery repairs projections that updates-only watchers could
// miss during a leader handoff. It starts only after durable user-info delivery
// is attached, so an overlapping event cannot be lost before the sweep.
func (s *Housekeeper) runLeadershipRecovery(ctx context.Context) {
	if _, _, _, _, err := s.cleanupDispatchers(ctx); err != nil {
		s.logger.Error("failed to reconcile dispatchers on leadership start", zap.Error(err))
	}
	if _, _, err := s.checkUnitUsers(ctx); err != nil {
		s.logger.Error("failed to reconcile unit membership on leadership start", zap.Error(err))
	}
	if _, _, err := s.removeDispatchesFromEmptyUnits(ctx); err != nil {
		s.logger.Error(
			"failed to reconcile empty unit dispatch assignments on leadership start",
			zap.Error(err),
		)
	}
	if _, _, err := s.reconcileUserInfoState(ctx); err != nil {
		s.logger.Error("failed to reconcile user info state on leadership start", zap.Error(err))
	}
}
