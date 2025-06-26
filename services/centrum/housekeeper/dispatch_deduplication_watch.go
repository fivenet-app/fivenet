package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
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
	store := s.dispatches.Store() // Helper that returns our nats store wrapper

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eventCh, err := store.WatchAll(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case e := <-eventCh:
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
					s.logger.Error("dispatch deduplication failed", zap.Uint64("dispatch_id", dsp.Id), zap.Error(err))
				}
			}
		}
	}
}

// tryDeduplicate single-dispatch dedup logic
func (s *Housekeeper) tryDeduplicate(ctx context.Context, dsp *centrum.Dispatch) error {
	// Check if the dispatch has already been cancelled, completed, etc.
	if dsp.Status == nil || centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
		return nil // already completed or cancelled, nothing to do
	}

	// Dispatches can belong to more than one job, don't duplicate when the dispatch belongs to no job or more than one job
	if dsp.Jobs.IsEmpty() || len(dsp.Jobs.Jobs) > 1 {
		return nil // No jobs assigned, nothing to deduplicate
	}

	if dsp.Attributes != nil && (dsp.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE) || dsp.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE)) {
		return nil // Already marked as multiple or duplicate, no need to deduplicate
	}

	job := dsp.Jobs.Jobs[0].Name
	locs := s.dispatches.GetLocations(job)
	if locs == nil {
		return nil // no spatial index yet
	}

	settings, err := s.settings.Get(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to get settings for job %s: %w", job, err)
	}
	if !settings.Configuration.DeduplicationEnabled {
		return nil
	}

	radius := float64(settings.Configuration.DeduplicationRadius)
	duration := settings.Configuration.DeduplicationDuration.AsDuration()

	// Search the spatial index for nearby active dispatches (same logic as before)
	closeBy := locs.KNearest(dsp.Point(), 8, func(p orb.Pointer) bool {
		return p.(*centrum.Dispatch).Id != dsp.Id
	}, radius) // metres
	if len(closeBy) == 0 {
		return nil
	}

	active := []*centrum.Dispatch{}
	for _, dest := range closeBy {
		other := dest.(*centrum.Dispatch)
		if other.Status != nil && centrumutils.IsStatusDispatchComplete(other.Status.Status) {
			continue
		}
		if other.CreatedAt != nil && time.Since(other.CreatedAt.AsTime()) >= duration {
			continue
		}
		if other.Attributes != nil {
			if other.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE) {
			} else if other.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE) {
				continue // Already marked as duplicate, skip it
			}
		}

		active = append(active, other)
	}

	if len(active) == 0 {
		return nil
	}

	slices.SortFunc(active, func(a, b *centrum.Dispatch) int {
		return int(a.Id - b.Id)
	})
	mainDsp := dsp
	if len(active) > 0 && active[0].Id < mainDsp.Id {
		mainDsp = active[0]
		if len(active) > 1 {
			active = active[1:]          // Remove the new main dispatch from the list of duplicates
			active = append(active, dsp) // Add the original dispatch to the list of duplicates
		} else {
			active = []*centrum.Dispatch{dsp}
		}
	}

	refs := &centrum.DispatchReferences{}
	for _, dup := range active {
		if dup.Id == mainDsp.Id {
			continue // Skip the main dispatch itself
		}

		refs.Add(&centrum.DispatchReference{
			TargetDispatchId: dup.Id,
			ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATED_BY,
		})
	}

	// Mark the current dispatch as "multiple" and add the references
	if err := s.dispatches.AddAttributeToDispatch(ctx, mainDsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE); err != nil {
		return err
	}
	if err := s.dispatches.AddReferencesToDispatch(ctx, mainDsp, refs.References...); err != nil {
		return err
	}

	// Mark the close-by ones as duplicates & cancel them (same as original)
	sourceRef := &centrum.DispatchReference{
		TargetDispatchId: mainDsp.Id,
		ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATE_OF,
	}

	for _, dup := range active {
		if dup.Id == mainDsp.Id {
			continue // Skip the main dispatch itself
		}

		if err := s.dispatches.AddAttributeToDispatch(ctx, dup, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE); err != nil {
			return err
		}
		if err := s.dispatches.AddReferencesToDispatch(ctx, dup, sourceRef); err != nil {
			return err
		}

		if _, err := s.dispatches.UpdateStatus(ctx, dup.Id, &centrum.DispatchStatus{
			CreatedAt:  timestamp.Now(),
			DispatchId: dup.Id,
			Status:     centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		}); err != nil {
			return err
		}

		toRemove := []uint64{}
		for _, ua := range dup.Units {
			toRemove = append(toRemove, ua.UnitId)
		}
		if err := s.dispatches.UpdateAssignments(ctx, nil, dup.Id, nil, toRemove, time.Time{}); err != nil {
			return fmt.Errorf("failed to remove assigned units from duplicate dispatch. %w", err)
		}
	}

	return nil
}
