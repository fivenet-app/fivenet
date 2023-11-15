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

const (
	DelayBetweenDispatchAssignment = 35 * time.Second
	MinUnitCountForDelay           = 3
	PerUnitDelaySeconds            = 6
	MaxDelayCap                    = 80 * time.Second
)

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

		case <-time.After(4 * time.Second):
		}

		dispatches := b.state.FilterDispatches(b.job, nil, []dispatch.StatusDispatch{
			// Dispatch status that mean it is being worked on
			dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
			dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED,
			dispatch.StatusDispatch_STATUS_DISPATCH_EN_ROUTE,
			dispatch.StatusDispatch_STATUS_DISPATCH_ON_SCENE,
			dispatch.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE,
			// Completed states
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
				b.logger.Warn("no available units for dispatch", zap.Uint64("dispatch_id", dsp.Id))
				break
			}

			if err := b.state.UpdateDispatchAssignments(
				ctx, b.job, nil, dsp,
				[]uint64{unit.Id}, nil,
				b.state.DispatchAssignmentExpirationTime(),
			); err != nil {
				b.logger.Warn("failed to assgin unit to dispatch", zap.Uint64("dispatch_id", dsp.Id), zap.Uint64("unit_id", unit.Id))
				break
			}
		}
	}
}

func (b *Bot) getAvailableUnit(ctx context.Context) (*dispatch.Unit, bool) {
	units := b.state.FilterUnits(b.job, []dispatch.StatusUnit{dispatch.StatusUnit_STATUS_UNIT_AVAILABLE}, nil)
	if len(units) == 0 {
		return nil, false
	}

	// Randomize unit ids
	for i := range units {
		j := rand.Intn(i + 1)
		units[i], units[j] = units[j], units[i]
	}

	var unit *dispatch.Unit
	for _, u := range units {
		t, ok := b.lastAssignedUnits[u.Id]
		if !ok || time.Now().After(t) {
			// Double check if unit is still available
			if u.Status != nil && u.Status.Status != dispatch.StatusUnit_STATUS_UNIT_AVAILABLE {
				continue
			}

			unit = u
			break
		}
	}

	if unit == nil {
		return nil, false
	}

	delay := 0 * time.Second
	unitCount := len(units)
	if unitCount > MinUnitCountForDelay {
		delay = time.Duration(unitCount*PerUnitDelaySeconds) * time.Second
		if delay >= MaxDelayCap {
			delay = MaxDelayCap
		}
	}

	b.lastAssignedUnits[unit.Id] = time.Now().Add(DelayBetweenDispatchAssignment).Add(delay)

	return unit, true
}
