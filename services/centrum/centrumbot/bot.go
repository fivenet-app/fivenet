package centrumbot

import (
	"context"
	"math/rand"
	"sort"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrummanager"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"go.uber.org/zap"
)

const (
	DelayBetweenDispatchAssignment = 35 * time.Second
	MinUnitCountForDelay           = 3
	PerUnitDelaySeconds            = 5
	MaxDelayCap                    = 60 * time.Second
)

type Bot struct {
	ctx    context.Context
	cancel context.CancelFunc

	logger *zap.Logger

	job     string
	state   *centrummanager.Manager
	tracker tracker.ITracker

	lastAssignedUnits map[uint64]time.Time
}

func NewBot(ctx context.Context, logger *zap.Logger, job string, state *centrummanager.Manager, tracker tracker.ITracker) *Bot {
	ctx, cancel := context.WithCancel(ctx)

	return &Bot{
		ctx:               ctx,
		cancel:            cancel,
		logger:            logger.Named("bot").With(zap.String("job", job)),
		job:               job,
		state:             state,
		tracker:           tracker,
		lastAssignedUnits: map[uint64]time.Time{},
	}
}

func (b *Bot) Run() {
	for {
		select {
		case <-b.ctx.Done():
			return

		case <-time.After(4 * time.Second):
		}

		dispatches := b.state.FilterDispatches(b.ctx, b.job, nil, []centrum.StatusDispatch{
			// Dispatch status that mean it is being worked on
			centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
			centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED,
			centrum.StatusDispatch_STATUS_DISPATCH_EN_ROUTE,
			centrum.StatusDispatch_STATUS_DISPATCH_ON_SCENE,
			centrum.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE,
			// Completed states
			centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		})

		b.logger.Debug("trying to auto assign dispatches", zap.Int("dispatch_count", len(dispatches)))

		sort.Slice(dispatches, func(i, j int) bool {
			return dispatches[i].Id < dispatches[j].Id
		})

		for _, dsp := range dispatches {
			// Dispatch should be at least 5 seconds old to ensure deduplication has happened
			if (dsp.CreatedAt != nil && time.Since(dsp.CreatedAt.AsTime()) <= 5*time.Second) ||
				!centrumutils.IsDispatchUnassigned(dsp) {
				continue
			}

			b.logger.Debug("trying to auto assign dispatch", zap.Uint64("dispatch_id", dsp.Id))

			unit, ok := b.getAvailableUnit(b.ctx)
			if !ok {
				// No unit available
				b.logger.Warn("no available units for dispatch", zap.Uint64("dispatch_id", dsp.Id))
				break
			}

			if err := b.state.UpdateDispatchAssignments(b.ctx, b.job, nil, dsp.Id, []uint64{unit.Id}, nil,
				b.state.DispatchAssignmentExpirationTime()); err != nil {
				b.logger.Error("failed to assgin unit to dispatch", zap.Uint64("dispatch_id", dsp.Id), zap.Uint64("unit_id", unit.Id), zap.Error(err))
				break
			}
		}
	}
}

func (b *Bot) Stop() {
	b.cancel()

	<-b.ctx.Done()
}

func (b *Bot) getAvailableUnit(ctx context.Context) (*centrum.Unit, bool) {
	units := b.state.FilterUnits(ctx, b.job, []centrum.StatusUnit{centrum.StatusUnit_STATUS_UNIT_AVAILABLE}, nil,
		func(unit *centrum.Unit) bool {
			return unit.Attributes == nil || !unit.Attributes.Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_NO_DISPATCH_AUTO_ASSIGN)
		},
	)

	b.logger.Debug("found available units", zap.Int("available_units_count", len(units)))
	if len(units) == 0 {
		return nil, false
	}

	// Randomize unit ids
	for i := range units {
		j := rand.Intn(i + 1)
		units[i], units[j] = units[j], units[i]
	}

	var selectedUnit *centrum.Unit
	for _, unit := range units {
		t, ok := b.lastAssignedUnits[unit.Id]
		if !ok || time.Now().After(t) {
			// Double check if unit is still available
			if unit.Status == nil || unit.Status.Status != centrum.StatusUnit_STATUS_UNIT_AVAILABLE {
				continue
			}

			selectedUnit = unit
			break
		}
	}

	if selectedUnit == nil {
		return nil, false
	}

	delay := 0 * time.Second
	unitCount := len(units)
	if unitCount > MinUnitCountForDelay {
		delay = min(time.Duration(unitCount*PerUnitDelaySeconds)*time.Second, MaxDelayCap)
	}

	b.lastAssignedUnits[selectedUnit.Id] = time.Now().Add(DelayBetweenDispatchAssignment).Add(delay)

	return selectedUnit, true
}
