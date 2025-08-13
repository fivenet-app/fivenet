package centrumbot

import (
	"context"
	"math/rand"
	"sort"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"go.uber.org/zap"
)

const (
	DelayBetweenAssignments = 35 * time.Second
	MinUnitCountForDelay    = 3
	PerUnitDelaySeconds     = 5
	MaxAssignmentDelayCap   = 60 * time.Second
)

type Bot struct {
	ctx    context.Context
	cancel context.CancelFunc

	logger  *zap.Logger
	tracker tracker.ITracker

	helpers    *helpers.Helpers
	settings   *settings.SettingsDB
	units      *units.UnitDB
	dispatches *dispatches.DispatchDB

	job string

	lastAssignedUnits map[uint64]time.Time
}

func NewBot(
	ctx context.Context,
	logger *zap.Logger,
	tracker tracker.ITracker,
	helpers *helpers.Helpers,
	settings *settings.SettingsDB,
	units *units.UnitDB,
	dispatches *dispatches.DispatchDB,
	job string,
) *Bot {
	ctx, cancel := context.WithCancel(ctx)

	return &Bot{
		ctx:     ctx,
		cancel:  cancel,
		logger:  logger.Named("bot").With(zap.String("job", job)),
		tracker: tracker,

		helpers:    helpers,
		settings:   settings,
		units:      units,
		dispatches: dispatches,

		job:               job,
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

		dispatches := b.dispatches.Filter(b.ctx, []string{b.job}, nil, []centrum.StatusDispatch{
			// Dispatch status that mean it is being worked on
			centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
			centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED,
			centrum.StatusDispatch_STATUS_DISPATCH_EN_ROUTE,
			centrum.StatusDispatch_STATUS_DISPATCH_ON_SCENE,
			centrum.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE,
			// "Completed" states
			centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			centrum.StatusDispatch_STATUS_DISPATCH_DELETED,
		})

		b.logger.Debug(
			"trying to auto assign dispatches",
			zap.Int("dispatch_count", len(dispatches)),
		)

		sort.Slice(dispatches, func(i, j int) bool {
			return dispatches[i].GetId() < dispatches[j].GetId()
		})

		for _, dsp := range dispatches {
			// Dispatch should be at least 7 seconds old to ensure deduplication has happened
			if (dsp.GetCreatedAt() != nil && time.Since(dsp.GetCreatedAt().AsTime()) <= 7*time.Second) ||
				!centrumutils.IsDispatchUnassigned(dsp) {
				continue
			}

			b.logger.Debug("trying to auto assign dispatch", zap.Uint64("dispatch_id", dsp.GetId()))

			unit, ok := b.getAvailableUnit(b.ctx, dsp.GetJobs())
			if !ok {
				// No unit available
				b.logger.Warn(
					"no available units for dispatch",
					zap.Uint64("dispatch_id", dsp.GetId()),
				)
				break
			}

			if err := b.dispatches.UpdateAssignments(b.ctx, nil, dsp.GetId(), []uint64{unit.GetId()}, nil,
				b.settings.DispatchAssignmentExpirationTime()); err != nil {
				b.logger.Error(
					"failed to assgin unit to dispatch",
					zap.Uint64("dispatch_id", dsp.GetId()),
					zap.Uint64("unit_id", unit.GetId()),
					zap.Error(err),
				)
				break
			}
		}
	}
}

func (b *Bot) Stop() {
	b.cancel()

	<-b.ctx.Done()
}

func (b *Bot) getAvailableUnit(ctx context.Context, jobs *centrum.JobList) (*centrum.Unit, bool) {
	units := b.units.Filter(
		ctx,
		jobs.GetJobStrings(),
		[]centrum.StatusUnit{centrum.StatusUnit_STATUS_UNIT_AVAILABLE},
		nil,
		func(unit *centrum.Unit) bool {
			return unit.GetAttributes() == nil ||
				!unit.GetAttributes().
					Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_NO_DISPATCH_AUTO_ASSIGN)
		},
	)

	b.logger.Debug("found available units", zap.Int("available_units_count", len(units)))
	if len(units) == 0 {
		return nil, false
	}

	// Randomize unit ids
	for i := range units {
		//nolint:gosec // G404: rand.Intn is not cryptographically secure, but we don't need it to be here.
		j := rand.Intn(i + 1)
		units[i], units[j] = units[j], units[i]
	}

	var selectedUnit *centrum.Unit
	for _, unit := range units {
		t, ok := b.lastAssignedUnits[unit.GetId()]
		if !ok || time.Now().After(t) {
			// Double check if unit is still available
			if unit.GetStatus() == nil ||
				unit.GetStatus().GetStatus() != centrum.StatusUnit_STATUS_UNIT_AVAILABLE {
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
		delay = min(time.Duration(unitCount*PerUnitDelaySeconds)*time.Second, MaxAssignmentDelayCap)
	}

	b.lastAssignedUnits[selectedUnit.GetId()] = time.Now().Add(DelayBetweenAssignments).Add(delay)

	return selectedUnit, true
}
