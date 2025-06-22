package housekeeper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
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
