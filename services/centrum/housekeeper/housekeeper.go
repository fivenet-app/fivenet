package housekeeper

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/leaderelection"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	MaxCancelledDispatchesPerRun = 6

	DeleteDispatchDays = 14
	DeleteUnitDays     = 14
)

var Module = fx.Module("centrum_housekeeper",
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
	le      *leaderelection.LeaderElector

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
	JS      *events.JSWrapper
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
		nodeName := instance.ID() + "_centrum_housekeeper"

		var err error
		s.le, err = leaderelection.New(
			ctxCancel, s.logger, p.JS,
			"leader_election",     // Bucket
			"centrum_housekeeper", // Key
			12*time.Second,        // TTL for the lock
			6*time.Second,         // Heartbeat interval
			func(ctx context.Context) {
				s.logger.Info("housekeeper started", zap.String("node_name", nodeName))

				s.start(ctx)
			},
			nil, // No on stopped function, context cancels the centrum housekeeper
		)
		if err != nil {
			return fmt.Errorf("failed to create leader elector. %w", err)
		}

		s.le.Start()

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
		s.le.Stop()
		cancel()

		s.wg.Wait()

		return nil
	}))

	return Result{
		Housekeeper:  s,
		CronRegister: s,
	}
}

func (s *Housekeeper) start(ctx context.Context) {
	s.wg.Go(func() {
		s.runUserChangesWatch(ctx)
	})

	s.wg.Go(func() {
		s.runDispatchWatch(ctx)
	})

	s.wg.Go(func() {
		s.runIdleWatcher(ctx)
	})

	s.wg.Go(func() {
		s.runTTLWatcher(ctx)
	})
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	for _, c := range []string{
		"centrum.manager_housekeeper.dispatch_deduplication",
		"centrum.manager_housekeeper.load_new_dispatches",
		"centrum.manager_housekeeper.dispatch_assignment_expiration",
		"centrum.manager_housekeeper.cleanup_units",
		"centrum.manager_housekeeper.cancel_old_dispatches",
		"centrum.manager_housekeeper.delete_old_dispatches",
		"centrum.manager_housekeeper.delete_old_dispatches_from_kv",
	} {
		if err := registry.UnregisterCronjob(ctx, c); err != nil {
			return err
		}
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.load_new_dispatches",
		Schedule: "*/4 * * * * * *", // Every 4 seconds
		Timeout:  durationpb.New(3 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.dispatch_assignment_expiration",
		Schedule: "*/2 * * * * * *", // Every 2 seconds
		Timeout:  durationpb.New(3 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.cleanup_units",
		Schedule: "15 * * * *", // Every hour at 15 minutes past the hour
		Timeout:  durationpb.New(6 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.cancel_old_dispatches",
		Schedule: "*/12 * * * * * *", // Every 12 hours
		Timeout:  durationpb.New(30 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.delete_old_dispatches",
		Schedule: "@hourly", // Hourly
		Timeout:  durationpb.New(30 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.delete_old_dispatches_from_kv",
		Schedule: "@always", // Every minute
		Timeout:  durationpb.New(30 * time.Second),
	}); err != nil {
		return err
	}

	return nil
}

func (s *Housekeeper) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add("centrum.housekeeper.load_new_dispatches", s.loadNewDispatches)
	h.Add(
		"centrum.housekeeper.dispatch_assignment_expiration",
		s.runHandleDispatchAssignmentExpiration,
	)
	h.Add("centrum.housekeeper.cleanup_units", s.runCleanupUnits)
	h.Add("centrum.housekeeper.cancel_old_dispatches", s.runCancelOldDispatches)
	h.Add("centrum.housekeeper.delete_old_dispatches", s.runDeleteOldDispatches)
	h.Add("centrum.housekeeper.delete_old_dispatches_from_kv", s.runDeleteOldDispatchesFromKV)

	return nil
}
