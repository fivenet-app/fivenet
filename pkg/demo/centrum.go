package demo

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"go.uber.org/zap"
)

var (
	dispatchDescriptions = []string{
		"A person was seen acting suspiciously in the area.",
		"A vehicle was reported speeding.",
		"Loud noises were heard coming from a building.",
		"Possible altercation in progress.",
		"Unattended package found.",
	}
	dispatchMessages = []string{
		"Suspicious activity reported",
		"Speeding vehicle reported",
		"Noise complaint received",
		"Possible fight reported",
		"Suspicious package reported",
	}

	dispatchStatusProgression = []centrumdispatches.StatusDispatch{
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_EN_ROUTE,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_ON_SCENE,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	}
)

// generateDispatches creates up to 2 random dispatches per run.
func (d *Demo) generateDispatches(ctx context.Context) {
	numDispatches := d.randIntN(2) // 0-1 dispatches

	for range numDispatches {
		x := d.randFloat64()*(xBounds[1]-xBounds[0]) + xBounds[0]
		y := d.randFloat64()*(yBounds[1]-yBounds[0]) + yBounds[0]
		desc := dispatchDescriptions[d.randIntN(len(dispatchDescriptions))]
		msg := dispatchMessages[d.randIntN(len(dispatchMessages))]
		if _, err := d.dispatches.Create(ctx, &centrumdispatches.Dispatch{
			Jobs: &centrum.JobList{
				Jobs: []*centrum.JobListEntry{
					{Name: d.cfg.Demo.TargetJob},
				},
			},
			Message:     msg,
			Description: &desc,
			X:           x,
			Y:           y,
			Anon:        true,
		}); err != nil {
			d.logger.Error("failed to create dispatch", zap.Error(err))
		}
	}
}

// updateDispatches randomly updates the status and position of up to 2 existing dispatches.
func (d *Demo) updateDispatches(ctx context.Context) error {
	dsps := d.dispatches.List(ctx, []string{d.cfg.Demo.TargetJob})
	if len(dsps) == 0 {
		return nil
	}

	perm := d.randPerm(len(dsps))
	numToUpdate := min(len(dsps), 2)

	for i := range numToUpdate {
		dsp := dsps[perm[i]]

		x := dsp.GetX() + d.randFloat64()*700 - 350
		y := dsp.GetY() + d.randFloat64()*700 - 350

		currStatus := dsp.GetStatus().GetStatus()
		newStatusValue := currStatus
		for i, s := range dispatchStatusProgression {
			if s == currStatus && i+1 < len(dispatchStatusProgression) {
				newStatusValue = dispatchStatusProgression[i+1]
				break
			}
		}

		newStatus := &centrumdispatches.DispatchStatus{
			DispatchId: dsp.GetId(),
			Status:     newStatusValue,
			X:          &x,
			Y:          &y,
			CreatorJob: &d.cfg.Demo.TargetJob,
		}
		if _, err := d.dispatches.UpdateStatus(ctx, dsp.GetId(), newStatus); err != nil {
			d.logger.Error("failed to update dispatch status", zap.Error(err))
		}
	}

	return nil
}
