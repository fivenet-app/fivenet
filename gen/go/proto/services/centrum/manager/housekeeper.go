package manager

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/paulmach/orb"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	MaxCancelledDispatchesPerRun = 4
)

var HousekeeperModule = fx.Module("centrum_manager_housekeeper", fx.Provide(
	NewHousekeeper,
))

type Housekeeper struct {
	ctx    context.Context
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer trace.Tracer
	db     *sql.DB

	convertJobs []string

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

func NewHousekeeper(p HousekeeperParams) *Housekeeper {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Housekeeper{
		ctx:         ctx,
		logger:      p.Logger.Named("centrum_manager_housekeeper"),
		wg:          sync.WaitGroup{},
		tracer:      p.TP.Tracer("centrum-manager-housekeeper"),
		db:          p.DB,
		convertJobs: p.Config.Game.DispatchCenter.ConvertJobs,
		Manager:     p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runHandleDispatchAssignmentExpiration()
		}()

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runArchiveDispatches()
		}()

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runDispatchDeduplication()
		}()

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runCleanupUnits()
		}()

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

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		return nil
	}))

	return s
}

func (s *Housekeeper) runHandleDispatchAssignmentExpiration() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(1 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-assignment-expiration")
				defer span.End()

				if err := s.handleDispatchAssignmentExpiration(ctx); err != nil {
					s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
				}
			}()
		}
	}
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
			tDispatchUnit.ExpiresAt.LT_EQ(jet.NOW()),
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
				return err
			}
		}
	}

	return nil
}

func (s *Housekeeper) runArchiveDispatches() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(4 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-archival")
				defer span.End()

				if err := s.archiveDispatches(ctx); err != nil {
					s.logger.Error("failed to archive dispatches", zap.Error(err))
				}

				if err := s.deleteOldDispatches(ctx); err != nil {
					s.logger.Error("failed to cleanup dispatches", zap.Error(err))
				}
			}()
		}
	}
}

// Set `COMPLETED`/`CANCELLED` dispatches to status `ARCHIVED` when the status is older than 5 minutes
func (s *Housekeeper) archiveDispatches(ctx context.Context) error {
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
		// Dispatches that are at 5 minutes or older and are in status completed/cancelled or have no status at all
		WHERE(jet.AND(
			tDispatchStatus.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(5, jet.MINUTE)),
			),
			tDispatchStatus.ID.IS_NULL().OR(
				jet.AND(
					tDispatchStatus.ID.EQ(
						jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
					),
					tDispatchStatus.Status.IN(
						jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
						jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
					),
				),
			),
		))

	var dest []*struct {
		DispatchID uint64
		Job        string
		Status     centrum.StatusDispatch
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for _, ds := range dest {
		// Ignore already archived dispatches
		if ds.Status == centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED {
			continue
		}

		if _, err := s.UpdateDispatchStatus(ctx, ds.Job, ds.DispatchID, &centrum.DispatchStatus{
			DispatchId: ds.DispatchID,
			Status:     centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		}); err != nil {
			return err
		}

		if err := s.DeleteDispatch(ctx, ds.Job, ds.DispatchID, false); err != nil {
			return err
		}
	}

	return nil
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
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(14, jet.DAY)),
			),
		)).
		LIMIT(50)

	var dest []*struct {
		DispatchID uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for _, ds := range dest {
		if err := s.DeleteDispatch(ctx, ds.Job, ds.DispatchID, true); err != nil {
			return err
		}
	}

	return nil
}
func (s *Housekeeper) runDispatchDeduplication() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(2 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-deduplicatation")
				defer span.End()

				if err := s.deduplicateDispatches(ctx); err != nil {
					s.logger.Error("failed to deduplicate dispatches", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Housekeeper) deduplicateDispatches(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for _, job := range s.trackedJobs {
		wg.Add(1)
		go func(job string) {
			defer wg.Done()

			dsps := s.State.FilterDispatches(job, nil, []centrum.StatusDispatch{
				centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
				centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
				centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			})

			if len(dsps) <= 1 {
				return
			}

			removedCount := 0
			dispatchIds := map[uint64]interface{}{}
			for _, dsp := range dsps {
				// Add the handled dispatch to the list
				dispatchIds[dsp.Id] = nil

				if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
					continue
				}

				closestsDsp := s.State.GetDispatchLocations(dsp.Job).KNearest(dsp.Point(), 8, func(p orb.Pointer) bool {
					return p.(*centrum.Dispatch).Id != dsp.Id
				}, 45.0)
				// Add "multiple" attribute when multiple dispatches close by
				if len(closestsDsp) > 0 {
					if err := s.addAttributeToDispatch(ctx, dsp, "multiple"); err != nil {
						s.logger.Error("failed to update original dispatch attribute", zap.Error(err))
					}
				}

				for _, dest := range closestsDsp {
					if dest == nil {
						continue
					}

					// Already took care of the dispatch
					closeByDsp := dest.(*centrum.Dispatch)
					if _, ok := dispatchIds[closeByDsp.Id]; ok {
						continue
					}
					dispatchIds[closeByDsp.Id] = nil

					if closeByDsp.Status != nil && centrumutils.IsStatusDispatchComplete(closeByDsp.Status.Status) {
						continue
					}

					if closeByDsp.CreatedAt != nil && time.Since(closeByDsp.CreatedAt.AsTime()) >= 3*time.Minute {
						continue
					}

					s.State.GetDispatchLocations(closeByDsp.Job).Remove(closeByDsp, func(p orb.Pointer) bool {
						return p.(*centrum.Dispatch).Id == closeByDsp.Id
					})

					if _, err := s.UpdateDispatchStatus(ctx, closeByDsp.Job, closeByDsp.Id, &centrum.DispatchStatus{
						DispatchId: closeByDsp.Id,
						Status:     centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
					}); err != nil {
						s.logger.Error("failed to update duplicate dispatch status", zap.Error(err))
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

func (s *Housekeeper) addAttributeToDispatch(ctx context.Context, dsp *centrum.Dispatch, attribute string) error {
	update := false
	if dsp.Attributes == nil {
		dsp.Attributes = &centrum.Attributes{
			List: []string{attribute},
		}
		update = true
	} else {
		if !utils.InSlice(dsp.Attributes.List, attribute) {
			dsp.Attributes.List = append(dsp.Attributes.List, attribute)
			update = true
		}
	}

	if update {
		if _, err := s.UpdateDispatch(ctx, dsp.Job, dsp.CreatorId, dsp, true); err != nil {
			return err
		}
	}

	return nil
}

func (s *Housekeeper) runCleanupUnits() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(5 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-units-empty")
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
			}()
		}
	}
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty
func (s *Housekeeper) removeDispatchesFromEmptyUnits(ctx context.Context) error {
	for _, job := range s.trackedJobs {
		dsps := s.State.FilterDispatches(job, nil, []centrum.StatusDispatch{
			centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
		})

		for _, dsp := range dsps {
			for i := len(dsp.Units) - 1; i >= 0; i-- {
				if i > (len(dsp.Units) - 1) {
					break
				}

				unitId := dsp.Units[i].UnitId
				// If unit isn't empty, continue with the loop
				if unitId > 0 {
					continue
				}

				unit, err := s.GetUnit(job, unitId)
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
	for _, job := range s.trackedJobs {
		units, ok := s.ListUnits(job)
		if !ok {
			continue
		}

		for _, unit := range units {
			if len(unit.Users) > 0 {
				continue
			}

			if unit.Status != nil && unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE {
				continue
			}

			var userId *int32
			if unit.Status != nil && unit.Status.UserId != nil {
				userId = unit.Status.UserId
			}

			s.logger.Debug("cleaning up unit status to unavailable because it is empty",
				zap.String("job", job), zap.Uint64("unit_id", unit.Id))
			if _, err := s.UpdateUnitStatus(ctx, job, unit.Id, &centrum.UnitStatus{
				UnitId: unit.Id,
				Status: centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId: userId,
			}); err != nil {
				s.logger.Error("failed to update empty unit status to unavailable",
					zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Error(err))
				continue
			}
		}
	}

	return nil
}

func (s *Housekeeper) checkUnitUsers(ctx context.Context) error {
	for _, job := range s.trackedJobs {
		units, ok := s.ListUnits(job)
		if !ok {
			continue
		}

		for _, u := range units {
			unit, err := s.GetUnit(job, u.Id)
			if err != nil {
				continue
			}

			if len(unit.Users) == 0 {
				continue
			}

			toRemove := []int32{}
			for i := len(unit.Users) - 1; i >= 0; i-- {
				if i > (len(unit.Users) - 1) {
					break
				}

				userId := unit.Users[i].UserId
				if userId == 0 {
					s.logger.Warn("zero user id found during unit user checkup", zap.Uint64("unit_id", unit.Id))
					continue
				}

				unitId, _ := s.GetUserUnitID(userId)
				// If user is in that unit and still on duty, nothing to do, otherwise remove the user from the unit
				if unit.Id == unitId && s.tracker.IsUserOnDuty(job, userId) {
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

	return nil
}

func (s *Housekeeper) watchUserChanges() {
	userCh := s.tracker.Subscribe()

	for {
		select {
		case <-s.ctx.Done():
			return

		case event := <-userCh:
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-watch-users")
				defer span.End()

				s.logger.Debug("received user changes", zap.Int("added", len(event.Added)), zap.String("added_users", fmt.Sprintf("%+v", event.Added)),
					zap.Int("removed", len(event.Removed)), zap.String("removed_users", fmt.Sprintf("%+v", event.Removed)))
				for _, userInfo := range event.Added {
					if _, ok := s.GetUserUnitID(userInfo.UserID); ok {
						break
					}

					unitId, err := s.LoadUnitIDForUserID(ctx, userInfo.UserID)
					if err != nil {
						s.logger.Error("failed to load user unit id", zap.Error(err))
						continue
					}

					if unitId == 0 {
						s.UnsetUnitIDForUser(userInfo.UserID)
						continue
					}

					if err := s.SetUnitForUser(userInfo.UserID, unitId); err != nil {
						s.logger.Error("failed to update user unit id mapping in kv", zap.Error(err))
						continue
					}
				}

				for _, userInfo := range event.Removed {
					s.handleRemoveUserFromDisponents(ctx, userInfo.Job, userInfo.UserID)
					s.handleRemoveUserFromUnit(ctx, userInfo.Job, userInfo.UserID)
				}
			}()
		}
	}
}

func (s *Housekeeper) handleRemoveUserFromDisponents(ctx context.Context, job string, userId int32) {
	if s.CheckIfUserIsDisponent(job, userId) {
		if err := s.DisponentSignOn(ctx, job, userId, false); err != nil {
			s.logger.Error("failed to remove user from disponents", zap.Error(err))
			return
		}
	}
}

func (s *Housekeeper) handleRemoveUserFromUnit(ctx context.Context, job string, userId int32) bool {
	unitId, ok := s.GetUserUnitID(userId)
	if !ok {
		// Nothing to do
		return false
	}

	unit, err := s.GetUnit(job, unitId)
	if err != nil {
		s.UnsetUnitIDForUser(userId)
		return false
	}

	if err := s.UpdateUnitAssignments(ctx, job, &userId, unit.Id, nil, []int32{userId}); err != nil {
		s.logger.Error("failed to remove user from unit", zap.Error(err))
		return false
	}

	return true
}
