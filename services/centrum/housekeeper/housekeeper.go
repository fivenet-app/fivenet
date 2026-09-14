package housekeeper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
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
	pkguserinfo "github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	centrummetrics "github.com/fivenet-app/fivenet/v2026/services/centrum/metrics"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/units"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
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

	cancelOldDispatchesSchedule      = "0 */5 * * * *"
	auditEmptyUnitDispatchesSchedule = "0 */5 * * * *"
	deleteOldDispatchesKVSchedule    = "0 15 2,14 * * *"
	reconcileUserInfoStateSchedule   = "0 0 4,16 * * *"
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

	tracer                    trace.Tracer
	metrics                   *centrummetrics.Metrics
	db                        *sql.DB
	tracker                   tracker.ITracker
	userinfo                  pkguserinfo.UserInfoRetriever
	js                        *events.JSWrapper
	userInfoReconcileConsumer jetstream.Consumer
	le                        *leaderelection.LeaderElector

	settings                  *settings.SettingsDB
	dispatchers               *dispatchers.DispatchersDB
	units                     *units.UnitDB
	unitAssignments           unitAssignments
	unitAssignmentWatchSource unitAssignmentWatchSource
	dispatches                *dispatches.DispatchDB
	unitUserState             unitUserState
	dispatcherUserState       dispatcherUserState
	dispatchLifecycle         dispatchLifecycle
	emptyUnitCleaner          emptyUnitCleaner
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

type unitUserState interface {
	SyncUserUnitMapping(context.Context, int32) error
	ReconcileUserJobChange(context.Context, int32, string) (bool, error)
	Range(func(string, *centrumunits.Unit) bool)
}

type dispatcherUserState interface {
	SetUserState(context.Context, string, int32, bool) error
	Range(func(string, *centrumdispatchers.Dispatchers) bool)
}

type dispatchLifecycle interface {
	Get(context.Context, int64) (*centrumdispatches.Dispatch, error)
	UpdateStatus(
		context.Context,
		int64,
		*centrumdispatches.DispatchStatus,
	) (*centrumdispatches.DispatchStatus, error)
	AddAttributeToDispatch(
		context.Context,
		*centrumdispatches.Dispatch,
		centrumdispatches.DispatchAttribute,
	) error
	Delete(context.Context, int64, bool) error
	ScheduleProjectionCleanup(context.Context, int64, *timestamp.Timestamp) error
}

type emptyUnitCleaner interface {
	Remove(context.Context, *centrumunits.Unit) (int, error)
}

type dispatchAssignmentCleaner struct {
	housekeeper *Housekeeper
}

func (c dispatchAssignmentCleaner) Remove(
	ctx context.Context,
	unit *centrumunits.Unit,
) (int, error) {
	return c.housekeeper.removeDispatchesFromEmptyUnit(ctx, unit)
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	TP       *tracesdk.TracerProvider
	DB       *sql.DB
	JS       *events.JSWrapper
	Config   *config.Config
	Tracker  tracker.ITracker
	UserInfo pkguserinfo.UserInfoRetriever

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

		tracer:   p.TP.Tracer("centrum.manager.housekeeper"),
		metrics:  centrummetrics.Get(),
		db:       p.DB,
		tracker:  p.Tracker,
		userinfo: p.UserInfo,
		js:       p.JS,

		settings:        p.Settings,
		dispatchers:     p.Dispatchers,
		units:           p.Units,
		unitAssignments: p.Units,
		dispatches:      p.Dispatches,
	}
	// For testing, we can override these functions to use mocks or fakes.
	s.unitAssignmentWatchSource = p.Units.Store()
	s.unitUserState = p.Units
	s.dispatcherUserState = p.Dispatchers
	s.dispatchLifecycle = p.Dispatches
	s.emptyUnitCleaner = dispatchAssignmentCleaner{housekeeper: s}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := s.ensureUserInfoReconcileConsumer(ctxStartup); err != nil {
			return fmt.Errorf("failed to register user info reconciliation consumer: %w", err)
		}

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

	// Start durable delivery before the recovery sweep. An event arriving while
	// the sweep runs is safe because the per-user correction is idempotent.
	s.wg.Go(func() {
		s.maintainUserInfoReconcileConsumer(ctx, s.runLeadershipRecovery)
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
}

// runLeadershipRecovery repairs projections that updates-only watchers could
// miss during a leader handoff. It starts only after durable user-info delivery
// is attached, so an overlapping event cannot be lost before the sweep.
func (s *Housekeeper) runLeadershipRecovery(ctx context.Context) {
	if _, _, _, _, err := s.cleanupDispatchers(ctx); err != nil {
		s.logger.Error("failed to reconcile dispatchers on leadership start", zap.Error(err))
	}
	if _, _, err := s.checkUnitUsers(ctx); err != nil {
		s.logger.Error(
			"failed to reconcile unit membership on leadership start",
			zap.Error(err),
		)
	}
	if _, _, err := s.removeDispatchesFromEmptyUnits(ctx); err != nil {
		s.logger.Error(
			"failed to reconcile empty unit dispatch assignments on leadership start",
			zap.Error(err),
		)
	}
	if _, _, err := s.reconcileUserInfoState(ctx); err != nil {
		s.logger.Error("failed to reconcile user info state on leadership start", zap.Error(err))
	}
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
		"centrum.manager_housekeeper.reconcile_userinfo_state",
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
		Name:     "centrum.housekeeper.reconcile_userinfo_state",
		Schedule: reconcileUserInfoStateSchedule, // Twelve-hour authoritative recovery audit; also manually runnable.
		Timeout:  durationpb.New(2 * time.Minute),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.housekeeper.audit_empty_unit_dispatches",
		Schedule: auditEmptyUnitDispatchesSchedule, // Five-minute recovery audit for the targeted unit watcher.
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
	h.Add("centrum.housekeeper.reconcile_userinfo_state", s.runReconcileUserInfoState)
	h.Add("centrum.housekeeper.cleanup_dispatchers", s.runCleanupDispatchers)

	return nil
}
