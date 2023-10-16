package bot

import (
	"context"
	"fmt"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
)

type Bot struct {
	job   string
	state *manager.Manager

	lastAssignedUnits []uint64
}

func (b *Bot) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-time.After(2 * time.Second):
		}

		dispatches := b.state.GetDispatchesMap(b.job)
		dispatches.Range(func(key uint64, dsp *dispatch.Dispatch) bool {
			if !centrumutils.IsDispatchUnassigned(dsp) {
				return true
			}

			fmt.Println("DSP IS UNASSIGNED", dsp.Id)
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

			return true
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

	// TODO make sure units are not double assigned

	return unit, true
}
