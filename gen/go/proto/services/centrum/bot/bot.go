package bot

import (
	"context"
	"math/rand"
	"sort"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"go.uber.org/zap"
)

const DelayBetweenDispatchAssignment = 45 * time.Second

type Bot struct {
	logger *zap.Logger

	job   string
	state *manager.Manager

	lastAssignedUnits map[uint64]time.Time
}

func NewBot(logger *zap.Logger, job string, state *manager.Manager) *Bot {
	return &Bot{
		logger:            logger.Named("bot"),
		job:               job,
		state:             state,
		lastAssignedUnits: map[uint64]time.Time{},
	}
}

func (b *Bot) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-time.After(6 * time.Second):
		}

		dispatches := b.state.FilterDispatches(b.job, nil, []dispatch.StatusDispatch{
			dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		})
		sort.Slice(dispatches, func(i, j int) bool {
			return dispatches[i].Id < dispatches[j].Id
		})

		for _, dsp := range dispatches {
			if !centrumutils.IsDispatchUnassigned(dsp) {
				continue
			}

			unit, ok := b.getAvailableUnit(ctx)
			if !ok {
				// No unit available
				b.logger.Warn("No available units for dispatch", zap.Uint64("dispatch_id", dsp.Id))
				break
			}

			if err := b.state.UpdateDispatchAssignments(
				ctx, b.job, nil, dsp,
				[]uint64{unit.Id}, nil,
				b.state.DispatchAssignmentExpirationTime(),
			); err != nil {
				b.logger.Warn("Failed to assgin unit to dispatch", zap.Uint64("dispatch_id", dsp.Id), zap.Uint64("unit_id", unit.Id))
				break
			}
		}
	}
}

func (b *Bot) getAvailableUnit(ctx context.Context) (*dispatch.Unit, bool) {
	units := b.state.GetUnitsMap(b.job)

	unitIds := []uint64{}
	units.Range(func(unitId uint64, value *dispatch.Unit) bool {
		if value.Status != nil && value.Status.Status == dispatch.StatusUnit_STATUS_UNIT_AVAILABLE {
			unitIds = append(unitIds, unitId)
		}
		return true
	})

	// Randomize unit ids
	for i := range unitIds {
		j := rand.Intn(i + 1)
		unitIds[i], unitIds[j] = unitIds[j], unitIds[i]
	}

	var unit *dispatch.Unit
	for _, unitId := range unitIds {
		t, ok := b.lastAssignedUnits[unitId]
		if !ok || time.Now().After(t) {
			unit, _ = units.Load(unitId)
			break
		}
	}

	if unit == nil {
		return nil, false
	}

	b.lastAssignedUnits[unit.Id] = time.Now().Add(DelayBetweenDispatchAssignment)

	return unit, true
}
