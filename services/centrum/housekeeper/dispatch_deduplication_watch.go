package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/paulmach/orb"
	"go.uber.org/zap"
)

func (s *Housekeeper) runDispatchWatch(ctx context.Context) {
	for {
		if err := s.watchDispatches(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.logger.Error("dispatch watcher stopped", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) watchDispatches(ctx context.Context) error {
	watch, err := s.dispatches.Store().WatchAll(ctx)
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case e := <-watch.Updates():
			if e == nil { // heartbeat
				continue
			}

			switch e.Operation() {
			case jetstream.KeyValuePut:
				dsp, err := e.Value()
				if err != nil {
					s.logger.Warn("cannot read dispatch value", zap.Error(err))
					continue
				}

				if err := s.tryDeduplicate(ctx, dsp); err != nil {
					s.logger.Error(
						"dispatch deduplication failed",
						zap.Int64("dispatch_id", dsp.GetId()),
						zap.Error(err),
					)
				}
			}
		}
	}
}

// tryDeduplicate single-dispatch dedup logic.
func (s *Housekeeper) tryDeduplicate(
	ctx context.Context,
	mainDsp *centrumdispatches.Dispatch,
) error {
	// Check if the dispatch has already been cancelled, completed, etc.
	if mainDsp.GetStatus() == nil ||
		centrumutils.IsStatusDispatchComplete(mainDsp.GetStatus().GetStatus()) {
		return nil // already completed or cancelled, nothing to do
	}

	// Dispatches can belong to more than one job, don't duplicate when the dispatch belongs to no job or more than one job
	if mainDsp.GetJobs().IsEmpty() || len(mainDsp.GetJobs().GetJobs()) > 1 {
		return nil // No jobs assigned, nothing to deduplicate
	}

	if mainDsp.GetAttributes() != nil &&
		(mainDsp.GetAttributes().Has(centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE) || mainDsp.GetAttributes().Has(centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE)) {
		return nil // Already marked as multiple or duplicate, no need to deduplicate
	}

	activeDsps, err := s.findCloseByActiveDispatches(ctx, mainDsp)
	if err != nil {
		return err
	}
	if len(activeDsps) == 0 {
		return nil
	}

	slices.SortFunc(activeDsps, func(a, b *centrumdispatches.Dispatch) int {
		return int(a.GetId() - b.GetId())
	})

	if len(activeDsps) > 0 && activeDsps[0].GetId() < mainDsp.GetId() {
		mainDsp = activeDsps[0]
		if len(activeDsps) > 1 {
			activeDsps = activeDsps[1:] // Remove the new main dispatch from the list of duplicates
			activeDsps = append(
				activeDsps,
				mainDsp,
			) // Add the original dispatch to the list of duplicates
		} else {
			activeDsps = []*centrumdispatches.Dispatch{mainDsp}
		}
	}

	refs := &centrumdispatches.DispatchReferences{}
	for _, dup := range activeDsps {
		if dup.GetId() == mainDsp.GetId() {
			continue // Skip the main dispatch itself
		}

		refs.Add(&centrumdispatches.DispatchReference{
			TargetDispatchId: dup.GetId(),
			ReferenceType:    centrumdispatches.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATED_BY,
		})
	}

	// Mark the current dispatch as "multiple" and add the references
	if err := s.dispatches.AddAttributeToDispatch(
		ctx,
		mainDsp,
		centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE,
	); err != nil {
		return err
	}
	if err := s.dispatches.AddReferencesToDispatch(
		ctx,
		mainDsp,
		refs.GetReferences()...); err != nil {
		return err
	}

	if err := s.createReferences(ctx, mainDsp.GetId(), activeDsps); err != nil {
		return err
	}

	return nil
}

func (s *Housekeeper) findCloseByActiveDispatches(
	ctx context.Context,
	sourceDsp *centrumdispatches.Dispatch,
) ([]*centrumdispatches.Dispatch, error) {
	job := sourceDsp.GetJobs().GetJobs()[0].GetName()
	locs := s.dispatches.GetLocations(job)
	if locs == nil {
		return nil, nil // no spatial index yet
	}

	settings, err := s.settings.Get(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings for job %s. %w", job, err)
	}
	if !settings.GetConfiguration().GetDeduplicationEnabled() {
		return nil, nil
	}

	radius := float64(settings.GetConfiguration().GetDeduplicationRadius())
	duration := settings.GetConfiguration().GetDeduplicationDuration().AsDuration()

	// Search the spatial index for nearby active dispatches (same logic as before)
	closeBy := locs.KNearest(sourceDsp.Point(), 8, func(p orb.Pointer) bool {
		//nolint:forcetypeassert // We know that p is a *centrumdispatches.Dispatch because locs is a generics spatial index
		return p.(*centrumdispatches.Dispatch).GetId() != sourceDsp.GetId()
	}, radius) // meters
	if len(closeBy) == 0 {
		return nil, nil
	}

	active := []*centrumdispatches.Dispatch{}
	for _, dest := range closeBy {
		//nolint:forcetypeassert // We know that p is a *centrumdispatches.Dispatch because closeBy is a list of dispatches from a generics spatial index
		other := dest.(*centrumdispatches.Dispatch)
		if other.GetStatus() != nil &&
			centrumutils.IsStatusDispatchComplete(other.GetStatus().GetStatus()) {
			continue
		}
		if other.GetCreatedAt() != nil && time.Since(other.GetCreatedAt().AsTime()) >= duration {
			continue
		}
		if other.GetAttributes() != nil {
			if other.GetAttributes().
				Has(centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE) {
			} else if other.GetAttributes().Has(centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE) {
				continue // Already marked as duplicate, skip it
			}
		}

		active = append(active, other)
	}

	return active, nil
}

func (s *Housekeeper) createReferences(
	ctx context.Context,
	dispatchId int64,
	activeDispatches []*centrumdispatches.Dispatch,
) error {
	// Mark the close-by ones as duplicates & cancel them (same as original)
	sourceRef := &centrumdispatches.DispatchReference{
		TargetDispatchId: dispatchId,
		ReferenceType:    centrumdispatches.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATE_OF,
	}

	for _, dup := range activeDispatches {
		if dup.GetId() == dispatchId {
			continue // Skip the main dispatch itself
		}

		if err := s.dispatches.AddAttributeToDispatch(
			ctx,
			dup,
			centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE,
		); err != nil {
			return err
		}
		if err := s.dispatches.AddReferencesToDispatch(ctx, dup, sourceRef); err != nil {
			return err
		}

		if _, err := s.dispatches.UpdateStatus(ctx, dup.GetId(), &centrumdispatches.DispatchStatus{
			CreatedAt:  timestamp.Now(),
			DispatchId: dup.GetId(),
			Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		}); err != nil {
			return err
		}

		toRemove := []int64{}
		for _, ua := range dup.GetUnits() {
			toRemove = append(toRemove, ua.GetUnitId())
		}
		if err := s.dispatches.UpdateAssignments(
			ctx,
			nil,
			nil,
			dup.GetId(),
			nil,
			toRemove,
			time.Time{},
		); err != nil {
			return fmt.Errorf("failed to remove assigned units from duplicate dispatch. %w", err)
		}
	}

	return nil
}
