package dispatches

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// WatchSpecs exposes the independently stored dispatch specification updates
// to maintenance workers. The composed dispatch can be loaded by ID when an
// update arrives.
func (s *DispatchDB) WatchSpecs(
	ctx context.Context,
) (store.IKVWatcher[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec], error) {
	if s.specStore == nil {
		return nil, errors.New("dispatch specification projection is not initialized")
	}
	return s.specStore.WatchAll(ctx)
}

// WatchStatuses exposes status changes so consumers that require a fully
// composed dispatch can react after the spec-triggered initial write.
func (s *DispatchDB) WatchStatuses(
	ctx context.Context,
) (store.IKVWatcher[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus], error) {
	if s.statusStore == nil {
		return nil, errors.New("dispatch status projection is not initialized")
	}
	return s.statusStore.WatchAll(ctx)
}

func (s *DispatchDB) SpecStore() *store.Store[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec] {
	return s.specStore
}

func (s *DispatchDB) syncComponents(
	ctx context.Context,
	dispatch *centrumdispatches.Dispatch,
) error {
	if dispatch == nil {
		return nil
	}
	key := centrumutils.IdKey(dispatch.GetId())
	oldSpec, err := s.specStore.Get(key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	spec := &centrumdispatches.DispatchSpec{
		Id:          dispatch.GetId(),
		CreatedAt:   dispatch.CreatedAt,
		UpdatedAt:   dispatch.UpdatedAt,
		Jobs:        dispatch.GetJobs(),
		Message:     dispatch.GetMessage(),
		Description: dispatch.Description,
		Attributes:  dispatch.GetAttributes(),
		X:           dispatch.GetX(),
		Y:           dispatch.GetY(),
		Postal:      dispatch.Postal,
		Anon:        dispatch.GetAnon(),
		CreatorId:   dispatch.CreatorId,
		Creator:     dispatch.Creator,
		References:  dispatch.GetReferences(),
	}
	job := ""
	if jobs := dispatch.GetJobs().GetJobStrings(); len(jobs) > 0 {
		job = jobs[0]
	}
	if _, err := s.specStore.PutIfChanged(ctx, key, spec); err != nil {
		return err
	}
	newJobs := make(map[string]struct{}, len(dispatch.GetJobs().GetJobStrings()))
	for _, jobEntry := range dispatch.GetJobs().GetJobStrings() {
		if strings.TrimSpace(jobEntry) == "" {
			continue
		}
		newJobs[jobEntry] = struct{}{}
		if err := s.jobMapping.Put(
			ctx,
			centrumutils.JobIdKey(jobEntry, dispatch.GetId()),
			&common.IDMapping{Id: dispatch.GetId()},
		); err != nil {
			return err
		}
	}
	if oldSpec != nil {
		for _, oldJob := range oldSpec.GetJobs().GetJobStrings() {
			if strings.TrimSpace(oldJob) == "" {
				continue
			}
			if _, exists := newJobs[oldJob]; exists {
				continue
			}
			if err := s.jobMapping.Delete(
				ctx,
				centrumutils.JobIdKey(oldJob, dispatch.GetId()),
			); err != nil &&
				!errors.Is(err, jetstream.ErrKeyNotFound) {
				return err
			}
		}
	}
	if _, err := s.assignmentsStore.PutIfChanged(
		ctx,
		key,
		&centrumdispatches.DispatchAssignments{
			DispatchId: dispatch.GetId(),
			Job:        job,
			Units:      dispatch.GetUnits(),
		},
	); err != nil {
		return err
	}
	if dispatch.GetStatus() != nil {
		if _, err := s.statusStore.PutIfChanged(ctx, key, dispatch.GetStatus()); err != nil {
			return err
		}
	} else if err := s.statusStore.Delete(ctx, key); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	return nil
}

func (s *DispatchDB) deleteComponents(ctx context.Context, id int64) error {
	key := centrumutils.IdKey(id)
	var errs error
	// The job index is also the delete-feed source. Delete every mapping for
	// this ID, including stale mappings left behind by an interrupted update.
	for _, jobKey := range s.jobMapping.KeysFiltered("", func(jobKey string) bool {
		mappedID, err := centrumutils.ExtractID(jobKey)
		return err == nil && mappedID == id
	}) {
		if err := s.jobMapping.Delete(
			ctx,
			jobKey,
		); err != nil &&
			!errors.Is(err, jetstream.ErrKeyNotFound) {
			errs = errors.Join(errs, err)
		}
	}
	for _, deleteFn := range []func(context.Context, string) error{s.specStore.Delete, s.assignmentsStore.Delete, s.statusStore.Delete} {
		if err := deleteFn(ctx, key); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (s *DispatchDB) updateTakeDispatchProjection(
	ctx context.Context,
	dispatchID int64,
	unit *centrumunits.Unit,
	resp centrumdispatches.TakeDispatchResp,
) (*centrumdispatches.Dispatch, bool, bool, error) {
	dsp, err := s.Get(ctx, dispatchID)
	if err != nil {
		return nil, false, false, err
	}
	if dsp == nil ||
		(dsp.GetStatus() != nil && centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus())) {
		return nil, false, false, errorscentrum.ErrDispatchAlreadyCompleted
	}
	var preserveOperationalStatus bool
	if resp == centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
		found := false
		for _, assignment := range dsp.GetUnits() {
			if assignment.GetUnitId() != unit.GetId() {
				continue
			}
			found = true
			if assignment.GetExpiresAt() == nil {
				return dsp, false, false, nil
			}
			assignment.ExpiresAt = nil
			break
		}
		if !found {
			dsp.Units = append(dsp.Units, &centrumdispatches.DispatchAssignment{
				DispatchId: dsp.GetId(),
				UnitId:     unit.GetId(),
				Unit:       unit,
				CreatedAt:  timestamp.Now(),
			})
		}
		preserveOperationalStatus = centrumutils.IsStatusDispatchOperational(
			dsp.GetStatus().GetStatus(),
		)
	} else {
		dsp.Units = slices.DeleteFunc(
			dsp.GetUnits(),
			func(in *centrumdispatches.DispatchAssignment) bool {
				return in.GetUnitId() == unit.GetId()
			},
		)
		preserveOperationalStatus = len(dsp.GetUnits()) > 0 &&
			centrumutils.IsStatusDispatchOperational(dsp.GetStatus().GetStatus())
	}
	job := ""
	if jobs := dsp.GetJobs().GetJobStrings(); len(jobs) > 0 {
		job = jobs[0]
	}
	if _, err := s.assignmentsStore.PutIfChanged(
		ctx,
		centrumutils.IdKey(dispatchID),
		&centrumdispatches.DispatchAssignments{
			DispatchId: dispatchID, Job: job, Units: dsp.GetUnits(),
		},
	); err != nil {
		return nil, false, false, err
	}
	return dsp, true, preserveOperationalStatus, nil
}

func (s *DispatchDB) IdleStore() jetstream.KeyValue {
	return s.idleKV
}

func (s *DispatchDB) updateInKV(
	ctx context.Context,
	id int64,
	dsp *centrumdispatches.Dispatch,
) error {
	if err := s.syncComponents(ctx, dsp); err != nil {
		return err
	}
	if err := s.ScheduleProjectionCleanup(ctx, id, dsp.GetCreatedAt()); err != nil {
		// The weekly recovery audit covers failed timer creation. Do not turn a
		// successful database/KV projection write into a failed API operation.
		s.logger.Warn(
			"failed to schedule dispatch projection cleanup",
			zap.Int64("dispatch_id", id),
			zap.Error(err),
		)
	}

	return nil
}

func (s *DispatchDB) deleteInKV(ctx context.Context, id int64) error {
	return s.deleteComponents(ctx, id)
}

func (s *DispatchDB) Get(ctx context.Context, id int64) (*centrumdispatches.Dispatch, error) {
	key := centrumutils.IdKey(id)
	if s.specStore != nil {
		spec, specErr := s.specStore.GetOrLoad(ctx, key)
		if specErr == nil && spec != nil {
			assignments, err := s.assignmentsStore.GetOrLoad(ctx, key)
			if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
				return nil, err
			}
			status, err := s.statusStore.GetOrLoad(ctx, key)
			if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
				return nil, err
			}
			dsp := &centrumdispatches.Dispatch{
				Id:          spec.GetId(),
				CreatedAt:   spec.CreatedAt,
				UpdatedAt:   spec.UpdatedAt,
				Jobs:        spec.GetJobs(),
				Message:     spec.GetMessage(),
				Description: spec.Description,
				Attributes:  spec.GetAttributes(),
				X:           spec.GetX(),
				Y:           spec.GetY(),
				Postal:      spec.Postal,
				Anon:        spec.GetAnon(),
				CreatorId:   spec.CreatorId,
				Creator:     spec.Creator,
				References:  spec.GetReferences(),
			}
			if assignments != nil {
				dsp.Units = assignments.GetUnits()
			}
			dsp.Status = status
			s.EnrichJobLabels(dsp)
			return dsp, nil
		}
		if !errors.Is(specErr, jetstream.ErrKeyNotFound) {
			return nil, specErr
		}
	}
	return nil, jetstream.ErrKeyNotFound
}

// EnrichJobLabels adds the presentation-only job labels that are intentionally
// omitted from the database representation of a dispatch's job list.
func (s *DispatchDB) EnrichJobLabels(dsp *centrumdispatches.Dispatch) {
	if dsp == nil {
		return
	}

	for _, job := range dsp.GetJobs().GetJobs() {
		s.enricher.EnrichJobName(job)
	}
}

// List returns all dispatches that match the given job prefixes.
// If jobs is nil, it returns all dispatches.
func (s *DispatchDB) List(ctx context.Context, jobs []string) []*centrumdispatches.Dispatch {
	if len(jobs) == 0 {
		jobs = []string{""}
	}

	keys := s.jobMapping.KeysFiltered("", func(key string) bool {
		for _, job := range jobs {
			if strings.HasPrefix(key, job+".") {
				return true
			}
		}

		return false
	})

	ds := []*centrumdispatches.Dispatch{}
	for _, key := range keys {
		uid, err := centrumutils.ExtractIDString(key)
		if err != nil {
			continue
		}

		id, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			continue
		}
		dsp, err := s.Get(ctx, id)
		if err != nil {
			continue
		}
		s.EnrichJobLabels(dsp)
		// Skip any broken dispatches (e.g. missing job or ID)
		if dsp == nil || dsp.GetId() == 0 || len(dsp.GetJobs().GetJobs()) == 0 {
			continue
		}
		ds = append(ds, dsp)
	}

	slices.SortFunc(ds, func(a, b *centrumdispatches.Dispatch) int {
		return int(a.GetId() - b.GetId())
	})

	return ds
}

func (s *DispatchDB) Filter(
	ctx context.Context,
	jobs []string,
	statuses []centrumdispatches.StatusDispatch,
	notStatuses []centrumdispatches.StatusDispatch,
) []*centrumdispatches.Dispatch {
	ds := s.List(ctx, jobs)

	// List returns cache-owned objects. Redaction is response-specific and must
	// not mutate the shared dispatch projection.
	for i, dispatch := range ds {
		if dispatch.GetAnon() {
			ds[i] = proto.Clone(dispatch).(*centrumdispatches.Dispatch)
			ds[i].ClearCreator()
		}
	}

	ds = slices.DeleteFunc(ds, func(dispatch *centrumdispatches.Dispatch) bool {
		// Include statuses that should be listed
		if len(statuses) > 0 && !slices.Contains(statuses, dispatch.GetStatus().GetStatus()) {
			return true
		} else if len(notStatuses) > 0 && dispatch.GetStatus() != nil {
			// Which statuses to ignore
			if slices.Contains(notStatuses, dispatch.GetStatus().GetStatus()) {
				return true
			}
		}

		return false
	})

	return ds
}

func (s *DispatchDB) updateStatusInKV(
	ctx context.Context,
	id int64,
	status *centrumdispatches.DispatchStatus,
) error {
	if status == nil || s.statusStore == nil {
		return nil
	}
	_, err := s.statusStore.PutIfChanged(ctx, centrumutils.IdKey(id), status)
	return err
}

const (
	InactiveTTL            = 60 * time.Minute // How long a Dispatch may stay idle/quiet
	InactiveLimitMarkerTTL = 2 * time.Hour    // How long TTL-expiration markers remain replayable
	CompletedTTL           = 15 * time.Minute // How long completed dispatches remain visible
	ProjectionTTL          = 4 * time.Hour    // Maximum time a dispatch projection remains in KV
	DispatchKVTTL          = 24 * time.Hour   // Maximum time dispatch KV entries remain as a final safety net
	DispatchPurgeMarkerTTL = 4 * time.Hour    // How long removed dispatch/index tombstones remain replayable
	IdleKVTTL              = 6 * time.Hour    // Projection timer lifetime plus recovery headroom
)

func cleanupKey(id int64) string {
	return "cleanup." + centrumutils.IdKey(id)
}

func projectionCleanupKey(id int64) string {
	return "projection." + centrumutils.IdKey(id)
}

func assignmentExpirationKey(dispatchID, unitID int64) string {
	return "assignment." + centrumutils.IdKey(dispatchID) + "." + centrumutils.IdKey(unitID)
}

// ScheduleAssignmentExpiration creates the short-lived timer which wakes the
// elected housekeeper when a pending unit assignment reaches its deadline.
// MySQL remains authoritative; the periodic expiration job recovers from a
// timer that could not be created or observed.
func (d *DispatchDB) ScheduleAssignmentExpiration(
	ctx context.Context,
	dispatchID, unitID int64,
	expiresAt time.Time,
) error {
	if expiresAt.IsZero() {
		return d.CancelAssignmentExpiration(ctx, dispatchID, unitID)
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = time.Second
	}

	key := assignmentExpirationKey(dispatchID, unitID)
	for {
		_, err := d.idleKV.Create(ctx, key, nil, jetstream.KeyTTL(ttl))
		if !errors.Is(err, jetstream.ErrKeyExists) {
			return err
		}

		// Updating a KV entry refreshes the current entry TTL; it cannot apply
		// the new per-entry KeyTTL. Remove the old entry conditionally, then
		// retry Create with the requested deadline.
		entry, err := d.idleKV.Get(ctx, key)
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		if err := d.idleKV.Delete(ctx, key, jetstream.LastRevision(entry.Revision())); err != nil {
			if errors.Is(err, jetstream.ErrKeyNotFound) || errors.Is(err, jetstream.ErrKeyExists) {
				// A concurrent schedule or expiry changed the key. Start again from
				// Create so the latest successful operation installs its TTL.
				continue
			}
			return err
		}
	}
}

// CancelAssignmentExpiration cancels an outstanding assignment timer. A
// regular delete is intentionally distinct from the TTL expiry marker watched
// by the housekeeper.
func (d *DispatchDB) CancelAssignmentExpiration(
	ctx context.Context,
	dispatchID, unitID int64,
) error {
	err := d.idleKV.Delete(ctx, assignmentExpirationKey(dispatchID, unitID))
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil
	}

	return err
}

// TouchActivity updates the idle timer for a Dispatch.
func (d *DispatchDB) TouchActivity(ctx context.Context, id int64) error {
	key := "idle." + centrumutils.IdKey(id)

	// First try to create the timer key with TTL.
	if _, err := d.idleKV.Create(
		ctx, key, nil,
		jetstream.KeyTTL(InactiveTTL),
	); err != nil {
		if !errors.Is(err, jetstream.ErrKeyExists) {
			return err // a real failure
		}

		// Key already exists -> refresh its TTL with Update.
		// We need the current revision to comply with optimistic locking.
		entry, err := d.idleKV.Get(ctx, key)
		if err != nil {
			return err
		}

		// Update resets the previously-set TTL back to the full 60 min.
		if _, err = d.idleKV.Update(ctx, key, nil, entry.Revision()); err != nil {
			return err
		}
	}

	return nil
}

// ScheduleCleanup removes a completed dispatch from the live projection after CompletedTTL.
func (d *DispatchDB) ScheduleCleanup(ctx context.Context, id int64) error {
	_, err := d.idleKV.Create(ctx, cleanupKey(id), nil, jetstream.KeyTTL(CompletedTTL))
	if errors.Is(err, jetstream.ErrKeyExists) {
		return nil
	}

	return err
}

func (d *DispatchDB) CancelScheduledCleanup(ctx context.Context, id int64) error {
	key := cleanupKey(id)
	if _, err := d.idleKV.Get(ctx, key); err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil
		}

		return err
	}

	err := d.idleKV.Delete(ctx, key)
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil
	}

	return err
}

// ScheduleProjectionCleanup bounds the lifetime of a dispatch in the live KV
// projection. It deliberately does not refresh on updates: retention is based
// on the original dispatch creation time.
func (d *DispatchDB) ScheduleProjectionCleanup(
	ctx context.Context,
	id int64,
	createdAt *timestamp.Timestamp,
) error {
	if createdAt == nil {
		return nil
	}

	ttl := time.Until(createdAt.AsTime().Add(ProjectionTTL))
	if ttl <= 0 {
		ttl = time.Second
	}

	_, err := d.idleKV.Create(ctx, projectionCleanupKey(id), nil, jetstream.KeyTTL(ttl))
	if errors.Is(err, jetstream.ErrKeyExists) {
		return nil
	}

	return err
}
