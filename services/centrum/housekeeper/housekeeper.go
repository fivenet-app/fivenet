package housekeeper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
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

var Module = fx.Module("centrum_manager_housekeeper",
	fx.Provide(
		New,
	))

type Housekeeper struct {
	ctx    context.Context
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer  trace.Tracer
	db      *sql.DB
	tracker tracker.ITracker

	helpers     *helpers.Helpers
	settings    *settings.SettingsDB
	dispatchers *dispatchers.DispatchersDB
	units       *units.UnitDB
	dispatches  *dispatches.DispatchDB
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	DB      *sql.DB
	Config  *config.Config
	Tracker tracker.ITracker

	Helpers     *helpers.Helpers
	Settings    *settings.SettingsDB
	Dispatchers *dispatchers.DispatchersDB
	Units       *units.UnitDB
	Dispatches  *dispatches.DispatchDB
}

type Result struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func New(p Params) Result {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Housekeeper{
		ctx:    ctxCancel,
		logger: p.Logger.Named("centrum.manager.housekeeper"),
		wg:     sync.WaitGroup{},

		tracer:  p.TP.Tracer("centrum.manager.housekeeper"),
		db:      p.DB,
		tracker: p.Tracker,

		helpers:     p.Helpers,
		settings:    p.Settings,
		dispatchers: p.Dispatchers,
		units:       p.Units,
		dispatches:  p.Dispatches,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runUserChangesWatch()
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

	return Result{
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
		Schedule: "*/7 * * * * * *", // Every 7 seconds
		Timeout:  durationpb.New(6 * time.Second),
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

func (s *Housekeeper) runUserChangesWatch() {
	for {
		if err := s.watchUserChanges(s.ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.logger.Error("failed to watch user changes", zap.Error(err))
			}
		}

		select {
		case <-s.ctx.Done():
			s.logger.Info("stopping user changes watcher")
			return

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) watchUserChanges(ctx context.Context) error {
	userCh, err := s.tracker.Subscribe(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event := <-userCh:
			if event == nil {
				s.logger.Error("received nil user changes event, skipping")
				continue
			}

			func() {
				ctx, span := s.tracer.Start(ctx, "centrum.watch-users")
				defer span.End()

				if event.Operation() == jetstream.KeyValuePut {
					userMarker, err := event.Value()
					if err != nil {
						s.logger.Error("failed to get user marker from event", zap.Error(err))
						return
					}

					if _, err := s.tracker.GetUserMapping(userMarker.UserId); err != nil {
						return
					}

					unitId, err := s.units.LoadUnitIDForUserID(ctx, userMarker.UserId)
					if err != nil {
						s.logger.Error("failed to load user unit id", zap.Error(err))
						return
					}

					if err := s.tracker.SetUserMappingForUser(ctx, userMarker.UserId, &unitId); err != nil {
						s.logger.Error("failed to update user unit id mapping in kv", zap.Error(err))
						return
					}
				} else if event.Operation() == jetstream.KeyValueDelete || event.Operation() == jetstream.KeyValuePurge {
					userId, job, _, err := tracker.DecodeUserMarkerKey(event.Key())
					if err != nil {
						s.logger.Error("failed to decode user marker key", zap.Error(err), zap.String("key", string(event.Key())))
						return
					}

					if err := s.handleRemovedUser(ctx, job, userId); err != nil {
						s.logger.Error("failed to handle removed user", zap.Int32("user_id", userId), zap.Error(err))
					}
				}
			}()
		}
	}
}

func (s *Housekeeper) handleRemovedUser(ctx context.Context, job string, userId int32) error {
	var errs error
	if s.helpers.CheckIfUserIsDispatcher(ctx, job, userId) {
		if err := s.dispatchers.SetUserState(ctx, job, userId, false); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to remove user from disponents. %w", err))
		}
	}

	um, err := s.tracker.GetUserMapping(userId)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to get user unit mapping. %w", err))
		// User not in any unit, nothing to do
		return errs
	}

	if um != nil && um.UnitId != nil && *um.UnitId > 0 {
		if err := s.units.UpdateUnitAssignments(ctx, job, &userId, *um.UnitId, nil, []int32{userId}); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to remove user from unit. %w", err))
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

func (s *Housekeeper) runCleanupUnits(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.units-cleanup")
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
	for _, settings := range s.settings.List(ctx) {
		job := settings.Job

		dsps := s.dispatches.Filter(ctx, []string{job}, nil, []centrum.StatusDispatch{
			centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			centrum.StatusDispatch_STATUS_DISPATCH_DELETED,
		})

		for _, dsp := range dsps {
			// Make sure unassigned dispatch has the unassigned status
			if len(dsp.Units) == 0 && dsp.Status != nil && !centrumutils.IsStatusDispatchUnassigned(dsp.Status.Status) {
				s.logger.Debug("updating dispatch status to unassigned because it has no assignments",
					zap.String("job", job), zap.Uint64("dispatch_id", dsp.Id))
				if _, err := s.dispatches.UpdateStatus(ctx, dsp.Id, &centrum.DispatchStatus{
					CreatedAt:  timestamp.Now(),
					DispatchId: dsp.Id,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
					CreatorJob: &job,
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

				unit, err := s.units.Get(ctx, unitId)
				if err != nil {
					continue
				}

				if len(unit.Users) > 0 {
					continue
				}

				s.logger.Debug("removing empty unit from dispatch",
					zap.String("job", job), zap.Uint64("unit_id", unitId), zap.Uint64("dispatch_id", dsp.Id))

				if err := s.dispatches.UpdateAssignments(ctx, nil, dsp.Id, nil, []uint64{unitId}, time.Time{}); err != nil {
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
	for _, settings := range s.settings.List(ctx) {
		job := settings.Job

		units := s.units.List(ctx, []string{job})
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
			if _, err := s.units.UpdateStatus(ctx, unit.Id, &centrum.UnitStatus{
				CreatedAt:  timestamp.Now(),
				UnitId:     unit.Id,
				Status:     centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:     userId,
				CreatorJob: &job,
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

	for _, settings := range s.settings.List(ctx) {
		job := settings.Job

		units := s.units.List(ctx, []string{job})
		for _, u := range units {
			unit, err := s.units.Get(ctx, u.Id)
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

				unitMapping, err := s.tracker.GetUserMapping(userId)
				// If user is in that unit and still on duty, nothing to do, otherwise remove the user from the unit
				if err == nil && unitMapping.UnitId != nil && unit.Id == *unitMapping.UnitId && s.tracker.IsUserOnDuty(userId) {
					foundUserIds = append(foundUserIds, userId)
					continue
				}

				toRemove = append(toRemove, userId)
			}

			if len(toRemove) > 0 {
				s.logger.Debug("removing off-duty users from unit",
					zap.String("job", job), zap.Uint64("unit_id", unit.Id), zap.Int32s("to_remove", toRemove))

				if err := s.units.UpdateUnitAssignments(ctx, job, nil, unit.Id, nil, toRemove); err != nil {
					s.logger.Error("failed to remove off-duty users from unit",
						zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Int32s("user_ids", toRemove), zap.Error(err))
				}
			}
		}
	}

	userUnitIds, err := s.tracker.ListUserMappings(ctx)
	if err != nil {
		return err
	}

	errs := multierr.Combine()
	for _, userUnit := range userUnitIds {
		// Check if user id is part of an unit
		if slices.Contains(foundUserIds, userUnit.UserId) {
			continue
		}

		s.logger.Warn("found user id with unit mapping that isn't in any unit anymore", zap.Int32("user_id", userUnit.UserId), zap.Int32s("users_in_units", foundUserIds), zap.Any("mapping", userUnit))

		// TODO this isn't working as intended at the moment..
		/*
			// Unset unit id for user when user is not in any unit
			if err := s.tracker.UnsetUnitIDForUser(ctx, userId); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		*/
	}

	return errs
}
