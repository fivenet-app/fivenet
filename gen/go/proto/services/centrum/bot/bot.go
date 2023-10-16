package bot

import (
	"context"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
)

type Bot struct {
	job   string
	state *manager.Manager

	lastAssignedUnits map[uint64]time.Time
}

func NewBot(job string, state *manager.Manager) *Bot {
	return &Bot{
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

		case <-time.After(7 * time.Second):
		}

		dispatches := b.state.GetDispatchesMap(b.job)
		dispatches.Range(func(key uint64, dsp *dispatch.Dispatch) bool {
			if !centrumutils.IsDispatchUnassigned(dsp) {
				return true
			}

			unit, ok := b.getAvailableUnit(ctx)
			if !ok {
				// No unit available
				return false
			}

			if err := b.state.UpdateDispatchAssignments(
				ctx, b.job, nil, dsp,
				[]uint64{unit.Id}, nil,
				b.state.DispatchAssignmentExpirationTime(),
			); err != nil {
				return false
			}

			return false
		})
	}
}

func (b *Bot) getAvailableUnit(ctx context.Context) (*dispatch.Unit, bool) {
	units := b.state.GetUnitsMap(b.job)

	var unit *dispatch.Unit
	units.Range(func(key uint64, value *dispatch.Unit) bool {
		if value.Status != nil && value.Status.Status == dispatch.StatusUnit_STATUS_UNIT_AVAILABLE {
			unit = value
			return false
		}

		return true
	})

	if unit == nil {
		return nil, false
	}

	t, ok := b.lastAssignedUnits[unit.Id]
	if ok {
		if !time.Now().After(t) {
			return nil, false
		}
	}

	b.lastAssignedUnits[unit.Id] = time.Now().Add(35 * time.Second)

	return unit, true
}
