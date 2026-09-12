package housekeeper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/leaderelection"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	centrummetrics "github.com/fivenet-app/fivenet/v2026/services/centrum/metrics"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/units"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/go-jet/jet/v2/qrm"
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

	cancelOldDispatchesSchedule   = "0 */5 * * * *"
	deleteOldDispatchesKVSchedule = "0 15 2,14 * * *"
)

var Module = fx.Module("centrum_housekeeper",
	fx.Provide(
		New,
	))

var errWatcherUpdatesClosed = errors.New("watcher updates channel closed")

type Housekeeper struct {
	ctx    context.Context //nolint:containedctx // Housekeeper retains lifecycle context for watcher loops across files.
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer  trace.Tracer
	metrics *centrummetrics.Metrics
	db      *sql.DB
	tracker tracker.ITracker
	le      *leaderelection.LeaderElector

	settings                  *settings.SettingsDB
	dispatchers               *dispatchers.DispatchersDB
	units                     *units.UnitDB
	unitAssignments           unitAssignments
	unitAssignmentWatchSource unitAssignmentWatchSource
	removeEmptyUnit           func(context.Context, *centrumunits.Unit) (int, error)
	syncUserUnitMapping       func(context.Context, int32) error
	setDispatcherState        func(context.Context, string, int32, bool) error
	getDispatchProjection     func(context.Context, int64) (*centrumdispatches.Dispatch, error)
	updateDispatchStatus      func(context.Context, int64, *centrumdispatches.DispatchStatus) (*centrumdispatches.DispatchStatus, error)
	addDispatchAttribute      func(context.Context, *centrumdispatches.Dispatch, centrumdispatches.DispatchAttribute) error
	deleteDispatch            func(context.Context, int64, bool) error
	scheduleProjectionCleanup func(context.Context, int64, *timestamp.Timestamp) error
	dispatches                *dispatches.DispatchDB
}

type unitAssignments interface {
	UserInJob(ctx context.Context, db qrm.DB, job string, userID int32) (bool, error)
	UpdateUnitAssignments(
		ctx context.Context,
		creatorJob string,
		creatorId *int32,
		unitId int64,
		toAdd []int32,
		toRemove []int32,
	) error
}

type unitAssignmentWatchSource interface {
	WatchAll(
		ctx context.Context,
	) (store.IKVWatcher[centrumunits.Unit, *centrumunits.Unit], error)
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
		metrics: centrummetrics.Get(),
		db:      p.DB,
		tracker: p.Tracker,

		settings:        p.Settings,
		dispatchers:     p.Dispatchers,
		units:           p.Units,
		unitAssignments: p.Units,
		dispatches:      p.Dispatches,
	}
	s.unitAssignmentWatchSource = p.Units.Store()
	s.removeEmptyUnit = s.removeDispatchesFromEmptyUnit
	s.syncUserUnitMapping = p.Units.SyncUserUnitMapping
	s.setDispatcherState = p.Dispatchers.SetUserState
	s.getDispatchProjection = func(ctx context.Context, id int64) (*centrumdispatches.Dispatch, error) {
		return p.Dispatches.Store().Get(centrumutils.IdKey(id))
	}
	s.updateDispatchStatus = p.Dispatches.UpdateStatus
	s.addDispatchAttribute = p.Dispatches.AddAttributeToDispatch
	s.deleteDispatch = p.Dispatches.Delete
	s.scheduleProjectionCleanup = p.Dispatches.ScheduleProjectionCleanup

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
	// This sweep mutates durable state and must run only while elected leader.
	s.wg.Go(func() {
		if err := s.runDeleteOldDispatches(ctx, nil); err != nil {
			s.logger.Error("failed to delete old dispatches on leadership start", zap.Error(err))
		}
	})

	s.wg.Go(func() {
		s.runUserChangesWatch(ctx)
	})

	s.wg.Go(func() {
		s.runUnitAssignmentWatch(ctx)
	})

	s.wg.Go(func() {
		s.runDispatchWatch(ctx)
	})

	s.wg.Go(func() {
		s.runIdleWatcher(ctx)
	})

	s.wg.Go(func() {
		s.runDispatchCleanupWatcher(ctx)
	})

	s.wg.Go(func() {
		s.runProjectionCleanupWatcher(ctx)
	})

	s.wg.Go(func() {
		s.runTTLWatcher(ctx)
	})

	// Watchers use updates-only subscriptions, so a leadership handoff can miss
	// changes that occurred while this process was not leader. Reconcile the
	// session-bound projections immediately; concurrent watcher updates are
	// idempotent and keep the result current while this sweep runs.
	s.wg.Go(func() {
		if _, _, _, _, err := s.cleanupDispatchers(ctx); err != nil {
			s.logger.Error("failed to reconcile dispatchers on leadership start", zap.Error(err))
		}
		if _, _, err := s.checkUnitUsers(ctx); err != nil {
			s.logger.Error(
				"failed to reconcile unit membership on leadership start",
				zap.Error(err),
			)
		}
	})
}

func (s *Housekeeper) recordWatcherRestart(watcher string, err error) {
	outcome := "restart"
	if errors.Is(err, errWatcherUpdatesClosed) {
		outcome = "updates_closed"
	}
	s.metrics.IncHousekeeperEvent(watcher, outcome)
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	for _, c := range []string{
		"centrum.manager_housekeeper.dispatch_deduplication",
		"centrum.manager_housekeeper.load_new_dispatches",
		"centrum.manager_housekeeper.dispatch_assignment_expiration",
		"centrum.manager_housekeeper.cleanup_units",
		"centrum.manager_housekeeper.audit_unit_membership",
		"centrum.manager_housekeeper.audit_empty_unit_dispatches",
		"centrum.manager_housekeeper.cancel_old_dispatches",
		"centrum.manager_housekeeper.delete_old_dispatches",
		"centrum.manager_housekeeper.delete_old_dispatches_from_kv",
	} {
		if err := registry.UnregisterCronjob(ctx, c); err != nil {
			return err
		}
	}

	// Once legacy FiveM plugins create dispatches through the Centrum API, remove
	// the registration below and enable this replacement instead.
	// if err := registry.UnregisterCronjob(ctx, "centrum.housekeeper.load_new_dispatches"); err != nil {
	// 	return err
	// }
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
		Name:     "centrum.housekeeper.audit_empty_unit_dispatches",
		Schedule: "15 4 * * *", // Recovery audit for the targeted unit watcher.
		Timeout:  durationpb.New(30 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.audit_unit_membership",
		Schedule: "30 3 * * *", // Daily at 03:30
		Timeout:  durationpb.New(30 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.cancel_old_dispatches",
		Schedule: cancelOldDispatchesSchedule, // Every five minutes; targeted SQL recovery query.
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
		Schedule: deleteOldDispatchesKVSchedule, // Twelve-hour recovery audit; targeted timers handle normal cleanup.
		Timeout:  durationpb.New(30 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.cleanup_dispatchers",
		Schedule: "45 3 * * *", // Daily audit; tracker changes handle normal cleanup.
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
	h.Add("centrum.housekeeper.audit_unit_membership", s.runAuditUnitMembership)
	h.Add("centrum.housekeeper.audit_empty_unit_dispatches", s.runAuditEmptyUnitDispatches)
	h.Add("centrum.housekeeper.cancel_old_dispatches", s.runCancelOldDispatches)
	h.Add("centrum.housekeeper.delete_old_dispatches", s.runDeleteOldDispatches)
	h.Add("centrum.housekeeper.delete_old_dispatches_from_kv", s.runDeleteOldDispatchesFromKV)
	h.Add("centrum.housekeeper.cleanup_dispatchers", s.runCleanupDispatchers)

	return nil
}
