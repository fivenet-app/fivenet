package housekeeper

import (
	"context"
	"errors"
	"fmt"
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

	// Dispatches can belong to more than one job, only deduplicate when the dispatch belongs to a single job
	if dsp.Jobs.IsEmpty() {
		return nil // no jobs assigned, nothing to deduplicate
	}

	// More than one job assigned, skip dispatch
	if len(dsp.Jobs.Jobs) > 1 {
		return nil
	}

	job := dsp.Jobs.Jobs[0].Name
	locs, ok := s.dispatches.GetLocations(job)
	if !ok || locs == nil {
		return nil // no spatial index yet
	}

	// Search the spatial index for nearby active dispatches (same logic as before)
	closeBy := locs.KNearest(dsp.Point(), 8, func(p orb.Pointer) bool {
		return p.(*centrum.Dispatch).Id != dsp.Id
	}, 45.0) // metres

	if len(closeBy) == 0 {
		return nil
	}

	refs := &centrum.DispatchReferences{}
	active := []*centrum.Dispatch{}

	for _, dest := range closeBy {
		other := dest.(*centrum.Dispatch)
		if other.Status != nil && centrumutils.IsStatusDispatchComplete(other.Status.Status) {
			continue
		}
		if other.CreatedAt != nil && time.Since(other.CreatedAt.AsTime()) >= 3*time.Minute {
			continue
		}
		if other.Attributes != nil &&
			(other.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE) ||
				other.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE)) {
			continue
		}

		refs.Add(&centrum.DispatchReference{
			TargetDispatchId: other.Id,
			ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATED_BY,
		})
		active = append(active, other)
	}

	if len(active) == 0 {
		return nil
	}

	// mark the current dispatch as “multiple” and add the references
	if err := s.dispatches.AddAttributeToDispatch(ctx, dsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE); err != nil {
		return err
	}
	if err := s.dispatches.AddReferencesToDispatch(ctx, dsp, refs.References...); err != nil {
		return err
	}

	// Mark the close-by ones as duplicates & cancel them (same as original)
	sourceRef := &centrum.DispatchReference{
		TargetDispatchId: dsp.Id,
		ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATE_OF,
	}

	for _, dup := range active {
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
