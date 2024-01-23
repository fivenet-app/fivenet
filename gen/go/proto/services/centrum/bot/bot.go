package bot

import (
	"context"
	"math/rand"
	"sort"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/paulmach/orb"
	"go.uber.org/zap"
)

const (
	DelayBetweenDispatchAssignment = 35 * time.Second
	MinUnitCountForDelay           = 3
	PerUnitDelaySeconds            = 5
	MaxDelayCap                    = 80 * time.Second
)

type Bot struct {
	logger *zap.Logger

	job     string
	state   *manager.Manager
	tracker tracker.ITracker

	lastAssignedUnits map[uint64]time.Time
}

func NewBot(logger *zap.Logger, job string, state *manager.Manager, tracker tracker.ITracker) *Bot {
	return &Bot{
		logger:            logger.Named("bot"),
		job:               job,
		state:             state,
		tracker:           tracker,
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

		dispatches := b.state.FilterDispatches(b.job, nil, []centrum.StatusDispatch{
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

		sort.Slice(dispatches, func(i, j int) bool {
			return dispatches[i].Id < dispatches[j].Id
		})

		b.logger.Debug("trying to auto assign dispatches", zap.Int("dispatch_count", len(dispatches)))
		for _, dsp := range dispatches {
			if !centrumutils.IsDispatchUnassigned(dsp) {
				continue
			}

			b.logger.Debug("trying to auto assign dispatch", zap.Uint64("dispatch_id", dsp.Id))

			unit, ok := b.getAvailableUnit(ctx, dsp.Point())
			if !ok {
				// No unit available
				b.logger.Warn("no available units for dispatch", zap.Uint64("dispatch_id", dsp.Id))
				break
			}

			if err := b.state.UpdateDispatchAssignments(
				ctx, b.job, nil, dsp.Id,
				[]uint64{unit.Id}, nil,
				b.state.DispatchAssignmentExpirationTime(),
			); err != nil {
				b.logger.Error("failed to assgin unit to dispatch", zap.Uint64("dispatch_id", dsp.Id), zap.Uint64("unit_id", unit.Id), zap.Error(err))
				break
			}
		}
	}
}

func (b *Bot) getAvailableUnit(ctx context.Context, point orb.Point) (*centrum.Unit, bool) {
	var units []*centrum.Unit

	locs := b.tracker.GetUserJobLocations(b.job)
	if locs != nil {
		points := locs.KNearest(point, 7, nil, 5000.0)
		for _, point := range points {
			user := point.(*livemap.UserMarker)
			if user.UnitId == nil {
				continue
			}

			unit, err := b.state.GetUnit(user.Info.Job, *user.UnitId)
			if err != nil {
				b.logger.Error("failed to get user's unit", zap.String("job", user.Info.Job), zap.Error(err))
				continue
			}
			units = append(units, unit)
		}
	}

	if len(units) == 0 {
		if locs != nil {
			b.logger.Warn("falling back to normal unit selection, no close by user units found", zap.String("job", b.job))
		}

		units := b.state.FilterUnits(b.job, []centrum.StatusUnit{centrum.StatusUnit_STATUS_UNIT_AVAILABLE}, nil)
		if len(units) == 0 {
			return nil, false
		}

		// Randomize unit ids
		for i := range units {
			j := rand.Intn(i + 1)
			units[i], units[j] = units[j], units[i]
		}
	}

	var unit *centrum.Unit
	for _, u := range units {
		t, ok := b.lastAssignedUnits[u.Id]
		if !ok || time.Now().After(t) {
			// Double check if unit is still available
			if u.Status != nil && u.Status.Status != centrum.StatusUnit_STATUS_UNIT_AVAILABLE {
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
