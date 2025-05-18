package centrummanager

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/paulmach/orb"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	MaxCancelledDispatchesPerRun = 6

	DeleteDispatchDays = 14
	DeleteUnitDays     = 14
)

var HousekeeperModule = fx.Module("centrum_manager_housekeeper",
	fx.Provide(
		NewHousekeeper,
	))

type Housekeeper struct {
	ctx    context.Context
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer trace.Tracer
	db     *sql.DB

	converterType string
	convertJobs   []string

	*Manager
}

type HousekeeperParams struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	DB      *sql.DB
	Manager *Manager
	Config  *config.Config
}

type HousekeeperResult struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewHousekeeper(p HousekeeperParams) HousekeeperResult {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Housekeeper{
		ctx:           ctxCancel,
		logger:        p.Logger.Named("centrum.manager.housekeeper"),
		wg:            sync.WaitGroup{},
		tracer:        p.TP.Tracer("centrum-manager-housekeeper"),
		db:            p.DB,
		converterType: p.Config.DispatchCenter.Type,
		convertJobs:   p.Config.DispatchCenter.ConvertJobs,
		Manager:       p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.watchUserChanges()
		}()

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.ConvertPhoneJobMsgToDispatch()
		}()

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			if err := s.runDeleteOldDispatches(ctxCancel, nil); err != nil {
				s.logger.Error("failed to delete old dispatches on startup")
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		return nil
	}))

	return HousekeeperResult{
		Housekeeper:  s,
		CronRegister: s,
	}
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.manager_housekeeper.dispatch_assignment_expiration",
		Schedule: "*/2 * * * * * *", // Every 2 seconds
		Timeout:  durationpb.New(3 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.manager_housekeeper.dispatch_deduplication",
		Schedule: "*/2 * * * * * *", // Every 2 seconds
		Timeout:  durationpb.New(5 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.manager_housekeeper.cleanup_units",
		Schedule: "*/5 * * * * * *", // Every 5 seconds
		Timeout:  durationpb.New(10 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.manager_housekeeper.cancel_old_dispatches",
		Schedule: "*/15 * * * * * *", // Every 15 seconds
		Timeout:  durationpb.New(20 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.manager_housekeeper.load_new_dispatches",
		Schedule: "*/2 * * * * * *", // Every 2 seconds
		Timeout:  durationpb.New(5 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.manager_housekeeper.delete_old_dispatches",
		Schedule: "*/2 * * * *", // Every 2 minutes
		Timeout:  durationpb.New(15 * time.Second),
	}); err != nil {
		return err
	}

	return nil
}

func (s *Housekeeper) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add("centrum.manager_housekeeper.dispatch_assignment_expiration", s.runHandleDispatchAssignmentExpiration)
	h.Add("centrum.manager_housekeeper.dispatch_deduplication", s.runDispatchDeduplication)
	h.Add("centrum.manager_housekeeper.cleanup_units", s.runCleanupUnits)
	h.Add("centrum.manager_housekeeper.cancel_old_dispatches", s.runCancelOldDispatches)
	h.Add("centrum.manager_housekeeper.load_new_dispatches", s.loadNewDispatches)
	h.Add("centrum.manager_housekeeper.delete_old_dispatches", s.runDeleteOldDispatches)

	return nil
}

func (s *Housekeeper) loadNewDispatches(ctx context.Context, data *cron.CronjobData) error {
	// Load dispatches with null postal field (they are considered "new")
	if err := s.LoadDispatchesFromDB(ctx, tDispatch.Postal.IS_NULL()); err != nil {
		s.logger.Error("failed loading new dispatches from DB", zap.Error(err))
	}

	return nil
}

func (s *Housekeeper) runHandleDispatchAssignmentExpiration(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum-dispatch-assignment-expiration")
	defer span.End()

	if err := s.handleDispatchAssignmentExpiration(ctx); err != nil {
		s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
	}

	return nil
}

// Handle expired dispatch unit assignments
func (s *Housekeeper) handleDispatchAssignmentExpiration(ctx context.Context) error {
	stmt := tDispatchUnit.
		SELECT(
			tDispatchUnit.DispatchID.AS("dispatch_id"),
			tDispatchUnit.UnitID.AS("unit_id"),
			tUnits.Job.AS("job"),
		).
		FROM(
			tDispatchUnit.
				INNER_JOIN(tUnits,
					tUnits.ID.EQ(tDispatchUnit.UnitID),
				),
		).
		WHERE(jet.AND(
			tDispatchUnit.ExpiresAt.IS_NOT_NULL(),
			tDispatchUnit.ExpiresAt.LT_EQ(jet.CURRENT_TIMESTAMP()),
		))

	var dest []*struct {
		DispatchID uint64
		UnitID     uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	assignments := map[string]map[uint64][]uint64{}
	for _, ua := range dest {
		if _, ok := assignments[ua.Job]; !ok {
			assignments[ua.Job] = map[uint64][]uint64{}
		}
		if _, ok := assignments[ua.Job][ua.DispatchID]; !ok {
			assignments[ua.Job][ua.DispatchID] = []uint64{}
		}

		assignments[ua.Job][ua.DispatchID] = append(assignments[ua.Job][ua.DispatchID], ua.UnitID)
	}

	for job, dsps := range assignments {
		s.logger.Debug("handling dispatch assignment expiration", zap.String("job", job), zap.Int("expired_assignments", len(dsps)))
		for dispatchId, units := range dsps {
			if err := s.UpdateDispatchAssignments(ctx, job, nil, dispatchId, nil, units, time.Time{}); err != nil {
				return fmt.Errorf("failed to update dispatch %d assignments. %w", dispatchId, err)
			}
		}
	}

	return nil
}

func (s *Housekeeper) runCancelOldDispatches(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum-dispatch-cancel")
	defer span.End()

	if err := s.cancelOldDispatches(ctx); err != nil {
		s.logger.Error("failed to archive dispatches", zap.Error(err))
	}

	return nil
}

// Cancel dispatches that haven't been worked on for some time
func (s *Housekeeper) cancelOldDispatches(ctx context.Context) error {
	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.DispatchID.AS("dispatch_id"),
			tDispatch.Job.AS("job"),
			tDispatchStatus.Status.AS("status"),
		).
		FROM(
			tDispatchStatus.
				INNER_JOIN(tDispatch,
					tDispatch.ID.EQ(tDispatchStatus.DispatchID),
				),
		).
		// Dispatches that are older than time X and are not in a completed/cancelled/archived state, or have no status at all
		WHERE(jet.AND(
			tDispatchStatus.ID.EQ(
				jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
			),
			tDispatchStatus.Status.NOT_IN(
				jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
				jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
				jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED)),
			),
			tDispatch.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(60, jet.MINUTE)),
			),
		)).
		GROUP_BY(
			tDispatchStatus.DispatchID,
		).
		ORDER_BY(
			tDispatchStatus.DispatchID.ASC(),
		).
		LIMIT(200)

	var dest []*struct {
		DispatchID uint64
		Job        string
		Status     centrum.StatusDispatch
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	s.logger.Debug("canceling expired dispatches", zap.Int("dispatch_count", len(dest)))
	for _, ds := range dest {
		// Ignore already cancelled dispatches
		if ds.Status == centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED {
			continue
		}

		// Add "too old" attribute when we are able to retrieve the dispatch
		if dsp, err := s.GetDispatch(ctx, ds.Job, ds.DispatchID); err == nil && dsp != nil {
			if err := s.AddAttributeToDispatch(ctx, dsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_TOO_OLD); err != nil {
				s.logger.Error("failed to add too old attribute to cancelled dispatch", zap.Uint64("dispatch_id", ds.DispatchID), zap.Error(err))
			}
		}

		if _, err := s.UpdateDispatchStatus(ctx, ds.Job, ds.DispatchID, &centrum.DispatchStatus{
			CreatedAt:  timestamp.Now(),
			DispatchId: ds.DispatchID,
			Status:     centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		}); err != nil {
			s.logger.Error("failed to cancel dispatch", zap.Uint64("dispatch_id", ds.DispatchID), zap.Error(err))
			continue
		}

		// Remove dispatch from state and publish event so clients remove it
		if err := s.DeleteDispatch(ctx, ds.Job, ds.DispatchID, false); err != nil {
			s.logger.Error("failed to delete cancelled dispatch", zap.Uint64("dispatch_id", ds.DispatchID), zap.Error(err))
			continue
		}
	}

	return nil
}

func (s *Housekeeper) runDeleteOldDispatches(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum-dispatch-old-delete")
	defer span.End()

	errs := multierr.Combine()

	if err := s.deleteOldDispatches(ctx); err != nil {
		s.logger.Error("failed to remove old dispatches", zap.Error(err))
		errs = multierr.Append(errs, err)
	}

	if err := s.deleteOldDispatchesFromKV(ctx); err != nil {
		s.logger.Error("failed to remove old dispatches from kv", zap.Error(err))
		errs = multierr.Append(errs, err)
	}

	if err := s.deleteOldUnitStatus(ctx); err != nil {
		s.logger.Error("failed to remove old unit status", zap.Error(err))
		errs = multierr.Append(errs, err)
	}

	return errs
}

func (s *Housekeeper) deleteOldDispatches(ctx context.Context) error {
	stmt := tDispatch.
		SELECT(
			tDispatch.ID.AS("dispatch_id"),
			tDispatch.Job.AS("job"),
		).
		FROM(
			tDispatch,
		).
		WHERE(jet.AND(
			tDispatch.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(DeleteDispatchDays, jet.DAY)),
			),
		)).
		// Get 75 at a time
		LIMIT(75)

	var dest []*struct {
		DispatchID uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	errs := multierr.Combine()
	for _, ds := range dest {
		if err := s.DeleteDispatch(ctx, ds.Job, ds.DispatchID, true); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
	}

	return errs
}

func (s *Housekeeper) deleteOldDispatchesFromKV(ctx context.Context) error {
	errs := multierr.Combine()

	dsps := s.State.DispatchesStore().List()
	for _, dsp := range dsps {
		if dsp == nil {
			continue
		}

		if (
		// Created dispatches older than the delete dispatch days amount
		dsp.CreatedAt != nil && time.Since(dsp.CreatedAt.AsTime()) > DeleteDispatchDays*24*time.Hour) ||
			// Remove nil status
			dsp.Status == nil ||
			// "Completed" dispatches with their status being older than 15 minutes
			(centrumutils.IsStatusDispatchComplete(dsp.Status.Status) &&
				time.Since(dsp.Status.CreatedAt.AsTime()) > 15*time.Minute) {
			s.logger.Debug("old dispatch deleted from kv", zap.Uint64("dispatch_id", dsp.Id))
			if err := s.DeleteDispatch(ctx, dsp.Job, dsp.Id, false); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		}
	}

	return errs
}

func (s *Housekeeper) deleteOldUnitStatus(ctx context.Context) error {
	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		DELETE().
		WHERE(jet.AND(
			tUnitStatus.CreatedAt.LT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(DeleteUnitDays, jet.DAY))),
		)).
		LIMIT(1500)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return fmt.Errorf("failed to delete old unit status. %w", err)
	}

	return nil
}

func (s *Housekeeper) runDispatchDeduplication(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum-dispatch-deduplicatation")
	defer span.End()

	if err := s.deduplicateDispatches(ctx); err != nil {
		s.logger.Error("failed to deduplicate dispatches", zap.Error(err))
	}

	return nil
}

func (s *Housekeeper) deduplicateDispatches(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for _, settings := range s.ListSettings(ctx) {
		job := settings.Job
		if locs, ok := s.State.GetDispatchLocations(job); locs == nil || !ok {
			continue
		}

		wg.Add(1)
		go func(job string) {
			defer wg.Done()

			dsps := s.State.FilterDispatches(ctx, job, nil, []centrum.StatusDispatch{
				centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
				centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
				centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			})

			if len(dsps) <= 1 {
				return
			}

			removedCount := 0
			dispatchIds := map[uint64]any{}
			for _, dsp := range dsps {
				// Skip handled dispatches
				if _, ok := dispatchIds[dsp.Id]; ok {
					continue
				}

				// Add the handled dispatch to the list
				dispatchIds[dsp.Id] = nil

				if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
					continue
				}

				// Iterate over close by dispatches and collect the active ones (if locations are available)
				locs, ok := s.State.GetDispatchLocations(dsp.Job)
				if locs != nil || !ok {
					continue
				}
				closestsDsp := locs.KNearest(dsp.Point(), 8, func(p orb.Pointer) bool {
					return p.(*centrum.Dispatch).Id != dsp.Id
				}, 45.0)
				s.logger.Debug("deduplicating dispatches", zap.String("job", dsp.Job), zap.Uint64("dispatch_id", dsp.Id), zap.Int("closeby_dsps", len(closestsDsp)))

				var refs *centrum.DispatchReferences
				if dsp.References != nil {
					refs = dsp.References
				} else {
					refs = &centrum.DispatchReferences{}
				}

				activeDispatchesCloseBy := []*centrum.Dispatch{}
				for _, dest := range closestsDsp {
					if dest == nil {
						continue
					}

					closeByDsp := dest.(*centrum.Dispatch)
					if closeByDsp.Status != nil && centrumutils.IsStatusDispatchComplete(closeByDsp.Status.Status) {
						continue
					}

					if closeByDsp.CreatedAt != nil && time.Since(closeByDsp.CreatedAt.AsTime()) >= 3*time.Minute {
						continue
					}

					// Skip dispatches that are marked as multiple or duplicate dispatches they have already been handled
					if closeByDsp.Attributes != nil && (closeByDsp.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE) || closeByDsp.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE)) {
						continue
					}

					// Add close by dispatch as a reference
					refs.Add(&centrum.DispatchReference{
						TargetDispatchId: closeByDsp.Id,
						ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATED_BY,
					})

					activeDispatchesCloseBy = append(activeDispatchesCloseBy, closeByDsp)
				}

				// Prevent unnecessary updates to the dispatch
				if len(activeDispatchesCloseBy) == 0 {
					continue
				}

				// Add "multiple" attribute when multiple dispatches close by
				if err := s.AddAttributeToDispatch(ctx, dsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE); err != nil {
					s.logger.Error("failed to update original dispatch attribute", zap.Error(err))
				}

				// Set dispatch references on dispatch
				if err := s.AddReferencesToDispatch(ctx, dsp, refs.References...); err != nil {
					s.logger.Error("failed to update duplicate dispatch references", zap.Error(err))
				}

				sourceDspRef := &centrum.DispatchReference{
					TargetDispatchId: dsp.Id,
					ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATE_OF,
				}

				for _, closeByDsp := range activeDispatchesCloseBy {
					// Already took care of the dispatch
					if _, ok := dispatchIds[closeByDsp.Id]; ok {
						continue
					}
					dispatchIds[closeByDsp.Id] = nil

					if closeByDsp.Status != nil && centrumutils.IsStatusDispatchComplete(closeByDsp.Status.Status) {
						continue
					}

					if err := s.AddAttributeToDispatch(ctx, closeByDsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE); err != nil {
						s.logger.Error("failed to update duplicate dispatch attribute", zap.Error(err))
					}

					if err := s.AddReferencesToDispatch(ctx, closeByDsp, sourceDspRef); err != nil {
						s.logger.Error("failed to update duplicate dispatch references", zap.Error(err))
					}

					if _, err := s.UpdateDispatchStatus(ctx, closeByDsp.Job, closeByDsp.Id, &centrum.DispatchStatus{
						CreatedAt:  timestamp.Now(),
						DispatchId: closeByDsp.Id,
						Status:     centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
					}); err != nil {
						s.logger.Error("failed to update duplicate dispatch status", zap.Error(err))
						return
					}

					toRemove := []uint64{}
					for _, ua := range closeByDsp.Units {
						toRemove = append(toRemove, ua.UnitId)
					}
					if err := s.UpdateDispatchAssignments(ctx, closeByDsp.Job, nil, closeByDsp.Id, nil, toRemove, time.Time{}); err != nil {
						s.logger.Error("failed to remove assigned units from duplicate dispatch", zap.Error(err))
						return
					}

					removedCount++

					if removedCount >= MaxCancelledDispatchesPerRun {
						break
					}
				}
			}
		}(job)
	}

	wg.Wait()

	return nil
}

func (s *Housekeeper) runCleanupUnits(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum-units-cleanup")
	defer span.End()

	if err := s.removeDispatchesFromEmptyUnits(ctx); err != nil {
		s.logger.Error("failed to clean empty units from dispatches", zap.Error(err))
	}

	if err := s.cleanupUnitStatus(ctx); err != nil {
		s.logger.Error("failed to clean up unit status", zap.Error(err))
	}

	if err := s.checkUnitUsers(ctx); err != nil {
		s.logger.Error("failed to check duty state of unit users", zap.Error(err))
	}

	return nil
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty
func (s *Housekeeper) removeDispatchesFromEmptyUnits(ctx context.Context) error {
	for _, settings := range s.ListSettings(ctx) {
		job := settings.Job

		dsps := s.State.FilterDispatches(ctx, job, nil, []centrum.StatusDispatch{
			centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
		})

		for _, dsp := range dsps {
			// Make sure unassigned dispatch has the unassigned status
			if len(dsp.Units) == 0 && dsp.Status != nil && !centrumutils.IsStatusDispatchUnassigned(dsp.Status.Status) {
				s.logger.Debug("updating dispatch status to unassigned because it has no assignments",
					zap.String("job", job), zap.Uint64("dispatch_id", dsp.Id))
				if _, err := s.UpdateDispatchStatus(ctx, dsp.Job, dsp.Id, &centrum.DispatchStatus{
					CreatedAt:  timestamp.Now(),
					DispatchId: dsp.Id,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
				}); err != nil {
					return err
				}

				continue
			}

			for i := range slices.Backward(dsp.Units) {
				if i > (len(dsp.Units) - 1) {
					break
				}

				unitId := dsp.Units[i].UnitId
				// If unit isn't empty, continue with the loop
				if unitId <= 0 {
					continue
				}

				unit, err := s.GetUnit(ctx, job, unitId)
				if err != nil {
					continue
				}

				if len(unit.Users) > 0 {
					continue
				}

				s.logger.Debug("removing empty unit from dispatch",
					zap.String("job", job), zap.Uint64("unit_id", unitId), zap.Uint64("dispatch_id", dsp.Id))

				if err := s.UpdateDispatchAssignments(ctx, job, nil, dsp.Id, nil, []uint64{unitId}, time.Time{}); err != nil {
					s.logger.Error("failed to remove empty unit from dispatch",
						zap.String("job", job), zap.Uint64("unit_id", unitId), zap.Uint64("dispatch_id", dsp.Id), zap.Error(err))
					continue
				}
			}
		}
	}

	return nil
}

// Iterate over units to ensure that, e.g., an empty unit status is set to `unavailable`
func (s *Housekeeper) cleanupUnitStatus(ctx context.Context) error {
	for _, settings := range s.ListSettings(ctx) {
		job := settings.Job

		units, ok := s.ListUnits(ctx, job)
		if !ok {
			continue
		}

		for _, unit := range units {
			// Either unit has users but is static and in a wrong status
			if len(unit.Users) > 0 {
				if unit.Attributes == nil || !unit.Attributes.Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
					continue
				}

				if unit.Status != nil &&
					(unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_BUSY ||
						unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_ON_BREAK ||
						unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE) {
					continue
				}
			} else {
				// Or the unit is not already set to be unavailable (because it is empty)
				if unit.Status != nil &&
					unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE {
					continue
				}
			}

			var userId *int32
			if unit.Status != nil && unit.Status.UserId != nil {
				userId = unit.Status.UserId
			}

			s.logger.Debug("setting unit status to unavailable it is empty or static attribute (wrong status)",
				zap.String("job", job), zap.Uint64("unit_id", unit.Id), zap.Int32p("user_id", userId))
			if _, err := s.UpdateUnitStatus(ctx, job, unit.Id, &centrum.UnitStatus{
				CreatedAt: timestamp.Now(),
				UnitId:    unit.Id,
				Status:    centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:    userId,
			}); err != nil {
				s.logger.Error("failed to update empty unit status to unavailable",
					zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Error(err))
				continue
			}
		}
	}

	return nil
}

// Make sure that all users in units are still on duty
func (s *Housekeeper) checkUnitUsers(ctx context.Context) error {
	foundUserIds := []int32{}

	for _, settings := range s.ListSettings(ctx) {
		job := settings.Job

		units, ok := s.ListUnits(ctx, job)
		if !ok {
			continue
		}

		for _, u := range units {
			unit, err := s.GetUnit(ctx, job, u.Id)
			if err != nil {
				continue
			}

			if len(unit.Users) == 0 {
				continue
			}

			toRemove := []int32{}
			for i := range slices.Backward(unit.Users) {
				if i > (len(unit.Users) - 1) {
					break
				}

				userId := unit.Users[i].UserId
				if userId == 0 {
					s.logger.Warn("zero user id found during unit user checkup", zap.Uint64("unit_id", unit.Id))
					continue
				}

				unitId, found := s.GetUserUnitID(ctx, userId)
				// If user is in that unit and still on duty, nothing to do, otherwise remove the user from the unit
				if found && unit.Id == unitId && s.tracker.IsUserOnDuty(userId) {
					foundUserIds = append(foundUserIds, userId)
					continue
				}

				toRemove = append(toRemove, userId)
			}

			if len(toRemove) > 0 {
				s.logger.Debug("removing off-duty users from unit",
					zap.String("job", job), zap.Uint64("unit_id", unit.Id), zap.Int32s("to_remove", toRemove))

				if err := s.UpdateUnitAssignments(ctx, job, nil, unit.Id, nil, toRemove); err != nil {
					s.logger.Error("failed to remove off-duty users from unit",
						zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Int32s("user_ids", toRemove), zap.Error(err))
				}
			}
		}
	}

	userUnitIds, err := s.State.ListUserUnitMappings(ctx)
	if err != nil {
		return err
	}

	errs := multierr.Combine()
	for _, userUnit := range userUnitIds {
		// Check if user id is part of an unit
		if slices.Contains(foundUserIds, userUnit.UserId) {
			continue
		}

		s.logger.Warn("found user id with unit mapping that isn't in any unit", zap.Int32("user_id", userUnit.UserId), zap.Int32s("users_in_units", foundUserIds), zap.Any("mapping", userUnit))

		// TODO this isn't working as intended at the moment..
		/*
			// Unset unit id for user when user is not in any unit
			if err := s.UnsetUnitIDForUser(ctx, userId); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		*/
	}

	return errs
}

func (s *Housekeeper) watchUserChanges() {
	userCh := s.tracker.Subscribe()
	defer s.tracker.Unsubscribe(userCh)

	for {
		select {
		case <-s.ctx.Done():
			return

		case event := <-userCh:
			if event == nil {
				s.logger.Error("received nil user changes event, skipping")
				continue
			}

			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-watch-users")
				defer span.End()

				s.logger.Debug("received user changes", zap.Int("added", len(event.Added)), zap.Int("removed", len(event.Removed)))
				for _, userMarker := range event.Added {
					if _, ok := s.GetUserUnitID(ctx, userMarker.UserId); ok {
						break
					}

					unitId, err := s.LoadUnitIDForUserID(ctx, userMarker.UserId)
					if err != nil {
						s.logger.Error("failed to load user unit id", zap.Error(err))
						continue
					}

					if unitId == 0 {
						if err := s.UnsetUnitIDForUser(ctx, userMarker.UserId); err != nil {
							s.logger.Error("failed to unset user's unit id", zap.Error(err))
						}
						continue
					}

					if err := s.SetUnitForUser(ctx, userMarker.Job, userMarker.UserId, unitId); err != nil {
						s.logger.Error("failed to update user unit id mapping in kv", zap.Error(err))
						continue
					}
				}

				for _, userMarker := range event.Removed {
					s.handleRemovedUserDisponents(ctx, userMarker.Job, userMarker.UserId)
					s.handleRemovedUserUnit(ctx, userMarker.Job, userMarker.UserId)
				}
			}()
		}
	}
}

func (s *Housekeeper) handleRemovedUserDisponents(ctx context.Context, job string, userId int32) {
	if s.CheckIfUserIsDisponent(ctx, job, userId) {
		if err := s.DisponentSignOn(ctx, job, userId, false); err != nil {
			s.logger.Error("failed to remove user from disponents", zap.Error(err))
			return
		}
	}
}

func (s *Housekeeper) handleRemovedUserUnit(ctx context.Context, job string, userId int32) bool {
	unitId, ok := s.GetUserUnitID(ctx, userId)
	if !ok {
		// Nothing to do
		return false
	}

	unit, err := s.GetUnit(ctx, job, unitId)
	if err != nil {
		if err := s.UnsetUnitIDForUser(ctx, userId); err != nil {
			s.logger.Error("failed to unset user's unit id", zap.Error(err))
		}
		return false
	}

	if err := s.UpdateUnitAssignments(ctx, job, &userId, unit.Id, nil, []int32{userId}); err != nil {
		s.logger.Error("failed to remove user from unit", zap.Error(err))
		return false
	}

	return true
}
