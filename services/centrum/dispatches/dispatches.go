package dispatches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/units"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type dispatchMetrics struct {
	lastID *prometheus.GaugeVec
}

var (
	dispatchMetricsOnce sync.Once
	dispatchMetricsInst *dispatchMetrics
)

func getDispatchMetrics() *dispatchMetrics {
	dispatchMetricsOnce.Do(func() {
		dispatchMetricsInst = &dispatchMetrics{
			lastID: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "centrum",
				Name:      "dispatch_last_id",
				Help:      "Last dispatch ID.",
			}, []string{admin.MetricsJobNameLabel}),
		}

		prometheus.MustRegister(dispatchMetricsInst.lastID)
	})

	return dispatchMetricsInst
}

type DispatchDB struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher mstlystcdata.IEnricher
	tracker  tracker.ITracker
	postals  postals.Postals
	appCfg   appconfig.IConfig

	colleagueshydrator colleagueshydrator.IHydrator
	hydrator           citizenshydrator.IHydrator

	settings *settings.SettingsDB
	units    *units.UnitDB

	dispatchLocationsMutex  *sync.Mutex
	dispatchLocations       map[string]*coords.Coords[*centrumdispatches.Dispatch]
	dispatchLocationEntries map[int64]map[string]*centrumdispatches.Dispatch
	projectionLocksMu       sync.Mutex
	projectionLocks         map[int64]*sync.Mutex

	specStore        *store.Store[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec]
	assignmentsStore *store.Store[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments]
	statusStore      *store.Store[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus]
	jobMapping       *store.Store[common.IDMapping, *common.IDMapping]
	idleKV           jetstream.KeyValue
	metrics          *dispatchMetrics
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	JS        *events.JSWrapper
	DB        *sql.DB
	Cfg       *config.Config
	Enricher  mstlystcdata.IEnricher
	Tracker   tracker.ITracker
	Postals   postals.Postals
	AppConfig appconfig.IConfig

	Colleagueshydrator colleagueshydrator.IHydrator
	Hydrator           citizenshydrator.IHydrator

	Settings *settings.SettingsDB
	Units    *units.UnitDB
}

func New(p Params) *DispatchDB {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.Named("centrum.dispatches")
	d := &DispatchDB{
		logger:   logger,
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		postals:  p.Postals,
		appCfg:   p.AppConfig,

		colleagueshydrator: p.Colleagueshydrator,
		hydrator:           p.Hydrator,

		settings: p.Settings,
		units:    p.Units,

		dispatchLocationsMutex:  &sync.Mutex{},
		dispatchLocations:       map[string]*coords.Coords[*centrumdispatches.Dispatch]{},
		dispatchLocationEntries: map[int64]map[string]*centrumdispatches.Dispatch{},
		projectionLocks:         make(map[int64]*sync.Mutex),
		metrics:                 getDispatchMetrics(),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		idleKV, err := d.js.CreateOrUpdateKeyValue(ctxStartup, jetstream.KeyValueConfig{
			Bucket:         "centrum_dispatches_idle",
			Description:    "Timer keys for dispatch inactivity, cleanup, projections, and assignments",
			Storage:        jetstream.MemoryStorage,
			History:        1,
			MaxBytes:       0,
			TTL:            IdleKVTTL,
			LimitMarkerTTL: InactiveLimitMarkerTTL,
		})
		if err != nil {
			return fmt.Errorf("failed to create kv for idle dispatches. %w", err)
		}
		d.idleKV = idleKV

		jobSt, err := store.New[common.IDMapping, *common.IDMapping](
			ctxCancel,
			logger,
			p.JS,
			"centrum_dispatches",
			store.WithKVPrefix[common.IDMapping, *common.IDMapping]("job"),
			store.WithLocks[common.IDMapping, *common.IDMapping](nil),
			store.WithKVConfig[common.IDMapping, *common.IDMapping](
				jetstream.KeyValueConfig{
					TTL:            DispatchKVTTL,
					LimitMarkerTTL: DispatchKVTTL,
				},
			),
			store.WithPurgeMarkerTTL[common.IDMapping, *common.IDMapping](DispatchPurgeMarkerTTL),
		)
		if err != nil {
			return err
		}

		if err := jobSt.Start(ctxCancel, true); err != nil {
			return err
		}
		d.jobMapping = jobSt
		componentConfig := func() jetstream.KeyValueConfig {
			return jetstream.KeyValueConfig{TTL: DispatchKVTTL, LimitMarkerTTL: DispatchKVTTL}
		}
		specSt, err := store.New[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec](
			ctxCancel,
			logger,
			p.JS,
			"centrum_dispatches",
			store.WithKVPrefix[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec](
				"spec",
			),
			store.WithKVConfig[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec](
				componentConfig(),
			),
			store.WithPurgeMarkerTTL[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec](
				DispatchPurgeMarkerTTL,
			),
		)
		if err != nil {
			return err
		}
		assignmentsSt, err := store.New[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments](
			ctxCancel,
			logger,
			p.JS,
			"centrum_dispatches",
			store.WithKVPrefix[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments](
				"assignments",
			),
			store.WithKVConfig[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments](
				componentConfig(),
			),
			store.WithPurgeMarkerTTL[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments](
				DispatchPurgeMarkerTTL,
			),
		)
		if err != nil {
			return err
		}
		statusSt, err := store.New[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus](
			ctxCancel,
			logger,
			p.JS,
			"centrum_dispatches",
			store.WithKVPrefix[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus](
				"status",
			),
			store.WithKVConfig[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus](
				componentConfig(),
			),
			store.WithPurgeMarkerTTL[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus](
				DispatchPurgeMarkerTTL,
			),
		)
		if err != nil {
			return err
		}
		for _, component := range []interface {
			Start(context.Context, bool) error
		}{specSt, assignmentsSt, statusSt} {
			if err := component.Start(ctxCancel, true); err != nil {
				return err
			}
		}
		d.specStore, d.assignmentsStore, d.statusStore = specSt, assignmentsSt, statusSt

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return d
}

func (s *DispatchDB) LoadFromDB(ctx context.Context, cond mysql.BoolExpression) (int, error) {
	if cond == nil {
		if err := s.removeOrphanedProjections(ctx); err != nil {
			return 0, err
		}
	}

	tDispatch := table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus := table.FivenetCentrumDispatchesStatus.AS("dispatch_status")

	condition := tDispatchStatus.ID.IS_NULL().OR(
		mysql.AND(
			tDispatchStatus.ID.EQ(
				mysql.RawInt(
					"SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`",
				),
			).
				// Don't load archived dispatches into cache
				AND(tDispatchStatus.Status.NOT_IN(
					mysql.Int32(int32(centrumdispatches.StatusDispatch_STATUS_DISPATCH_ARCHIVED)),
					mysql.Int32(int32(centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
					mysql.Int32(int32(centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
					mysql.Int32(int32(centrumdispatches.StatusDispatch_STATUS_DISPATCH_DELETED)),
				)),
		),
	)

	if cond != nil {
		condition = condition.AND(cond)
	}

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Jobs,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tDispatchStatus.CreatorJob,
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(200)

	dsps := []*centrumdispatches.Dispatch{}
	if err := stmt.QueryContext(ctx, s.db, &dsps); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	if len(dsps) == 0 {
		return 0, nil
	}

	publicJobs := s.appCfg.Get().JobInfo.GetPublicJobs()
	creatorTargets := make([]citizenshydrator.BasicTarget, 0, len(dsps))
	for i := range dsps {
		dsp := dsps[i]
		var err error
		dsp.Units, err = s.LoadDispatchAssignments(ctx, dsp.GetId())
		if err != nil {
			return 0, err
		}

		if dsp.GetCreatorId() > 0 {
			creatorTargets = append(creatorTargets, citizenshydrator.BasicTarget{
				UserID: dsp.GetCreatorId(),
				Set:    dsp.SetCreator,
			})
		}

		if dsp.Postal == nil {
			if postal, ok := s.postals.Closest(dsp.GetX(), dsp.GetY()); postal != nil &&
				ok {
				dsp.Postal = postal.Code
			}
		}

		// Ensure dispatch has a status
		if dsp.GetStatus() == nil {
			dsp.Status, err = s.AddDispatchStatus(ctx, s.db, &centrumdispatches.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dsp.GetId(),
				Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW,
				Postal:     dsp.Postal,
				X:          &dsp.X,
				Y:          &dsp.Y,
			})
			if err != nil {
				return 0, fmt.Errorf(
					"failed to add dispatch (id: %d) status. %w",
					dsp.GetId(),
					err,
				)
			}
		}

		s.EnrichJobLabels(dsp)

		// Update dispatch in db and in kv
		if _, err := s.Update(ctx, nil, dsp); err != nil {
			return 0, err
		}
		if err := s.syncComponents(ctx, dsp); err != nil {
			return 0, err
		}

		if err := s.UpdateDispatchLocation(dsp); err != nil {
			return 0, err
		}
	}

	if err := s.hydrator.HydrateBasicTargetsSafeFunc(nil)(
		ctx,
		s.db,
		creatorTargets,
	); err != nil {
		return 0, err
	}

	for i := range dsps {
		if dsps[i].GetCreator() == nil {
			continue
		}

		// Clear dispatch creator's job info if it isn't a public job
		if !slices.Contains(publicJobs, dsps[i].GetCreator().GetJob()) {
			dsps[i].Creator.Job = ""
		}
		dsps[i].Creator.JobGrade = 0
	}

	return len(dsps), nil
}

func (s *DispatchDB) removeOrphanedProjections(ctx context.Context) error {
	tDispatch := table.FivenetCentrumDispatches
	var orphaned []*centrumdispatches.DispatchSpec

	s.specStore.Range(func(_ string, spec *centrumdispatches.DispatchSpec) bool {
		if spec != nil {
			orphaned = append(orphaned, spec)
		}
		return true
	})

	for _, dispatch := range orphaned {
		var existing struct {
			ID int64 `alias:"id"`
		}
		stmt := tDispatch.
			SELECT(tDispatch.ID).
			WHERE(tDispatch.ID.EQ(mysql.Int64(dispatch.GetId())))

		err := stmt.QueryContext(ctx, s.db, &existing)
		if err == nil {
			continue
		}
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to check dispatch projection %d: %w", dispatch.GetId(), err)
		}

		unlock := s.lockProjection(dispatch.GetId())
		if err := s.deleteComponents(ctx, dispatch.GetId()); err != nil &&
			!errors.Is(err, jetstream.ErrKeyNotFound) {
			unlock()
			return fmt.Errorf(
				"failed to remove orphaned dispatch projection %d: %w",
				dispatch.GetId(),
				err,
			)
		}

		for _, job := range dispatch.GetJobs().GetJobStrings() {
			if strings.TrimSpace(job) == "" {
				continue
			}
			if err := s.jobMapping.Delete(
				ctx,
				centrumutils.JobIdKey(job, dispatch.GetId()),
			); err != nil &&
				!errors.Is(err, jetstream.ErrKeyNotFound) {
				unlock()
				return fmt.Errorf(
					"failed to remove orphaned dispatch job mapping %d: %w",
					dispatch.GetId(),
					err,
				)
			}
		}
		unlock()
	}

	return nil
}

func (s *DispatchDB) Create(
	ctx context.Context,
	dsp *centrumdispatches.Dispatch,
) (*centrumdispatches.Dispatch, error) {
	if dsp.GetJobs() == nil || len(dsp.GetJobs().GetJobs()) == 0 {
		return nil, errorscentrum.ErrDispatchNoJobs
	}

	for _, job := range dsp.GetJobs().GetJobs() {
		if strings.TrimSpace(job.GetName()) == "" {
			return nil, errorscentrum.ErrDispatchNoJobs
		}
	}
	s.EnrichJobLabels(dsp)

	if dsp.Postal == nil || dsp.GetPostal() == "" {
		if postal, ok := s.postals.Closest(dsp.GetX(), dsp.GetY()); postal != nil && ok {
			dsp.Postal = postal.Code
		}
	}

	if dsp.GetCreatorId() > 0 {
		getShortByUserID := s.hydrator.GetBasicByUserIDSafeFunc(nil)
		var err error
		creator, err := getShortByUserID(ctx, s.db, dsp.GetCreatorId())
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user for dispatch creator. %w", err)
		}
		// Unset creator in case we don't have a user
		if creator == nil {
			dsp.ClearCreatorId()
		} else {
			dsp.SetCreator(creator)
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDispatch := table.FivenetCentrumDispatches
	stmt := tDispatch.
		INSERT(
			tDispatch.CreatedAt,
			tDispatch.Jobs,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
		).
		VALUES(
			mysql.CURRENT_TIMESTAMP(),
			dsp.GetJobs(),
			dsp.GetMessage(),
			dsp.Description,
			dsp.GetAttributes(),
			dsp.GetReferences(),
			dsp.GetX(),
			dsp.GetY(),
			dsp.Postal,
			dsp.GetAnon(),
			dsp.CreatorId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	dsp.SetId(lastId)

	var userId *int32
	if !dsp.GetAnon() && dsp.CreatorId != nil {
		userId = dsp.CreatorId
	}

	var statusUser *jobscolleagues.Colleague
	if dsp.GetCreator() != nil {
		statusUser = dsp.GetCreator().Colleague()
	}

	if dsp.Status, err = s.AddDispatchStatus(ctx, tx, &centrumdispatches.DispatchStatus{
		CreatedAt:  timestamp.Now(),
		DispatchId: dsp.GetId(),
		UserId:     userId,
		User:       statusUser,
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW,
		X:          &dsp.X,
		Y:          &dsp.Y,
		Postal:     dsp.Postal,
	}); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	for _, job := range dsp.GetJobs().GetJobStrings() {
		s.metrics.lastID.WithLabelValues(job).Set(float64(lastId))
	}

	// Hide user info when dispatch is anonymous
	if dsp.GetAnon() {
		dsp.ClearCreator()
	}

	if err := s.updateInKV(ctx, dsp.GetId(), dsp); err != nil {
		return nil, err
	}

	return dsp, nil
}

func (s *DispatchDB) Update(
	ctx context.Context,
	_ *int32,
	dsp *centrumdispatches.Dispatch,
) (*centrumdispatches.Dispatch, error) {
	// Labels are derived presentation data and are not retained in the SQL job
	// list. Restore them before persisting the projection and publishing it.
	s.EnrichJobLabels(dsp)
	dsp.UpdatedAt = timestamp.Now()

	tDispatch := table.FivenetCentrumDispatches
	stmt := tDispatch.
		UPDATE(
			tDispatch.UpdatedAt,
			tDispatch.Jobs,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
		).
		SET(
			mysql.CURRENT_TIMESTAMP(),
			dsp.GetJobs(),
			dsp.GetMessage(),
			dsp.Description,
			dsp.GetAttributes(),
			dsp.GetReferences(),
			dsp.GetX(),
			dsp.GetY(),
			dsp.Postal,
			dsp.GetAnon(),
			dsp.CreatorId,
		).
		WHERE(mysql.AND(
			tDispatch.ID.EQ(mysql.Int64(dsp.GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	if err := s.syncSpec(ctx, dsp); err != nil {
		return nil, err
	}
	current, err := s.Get(ctx, dsp.GetId())
	if err != nil {
		return nil, err
	}
	if err := s.UpdateDispatchLocation(current); err != nil {
		return nil, err
	}

	return dsp, nil
}

func (s *DispatchDB) AddAttributeToDispatch(
	ctx context.Context,
	dsp *centrumdispatches.Dispatch,
	attribute centrumdispatches.DispatchAttribute,
) error {
	var update bool
	if dsp.GetAttributes() == nil {
		dsp.Attributes = &centrumdispatches.DispatchAttributes{
			List: []centrumdispatches.DispatchAttribute{attribute},
		}

		update = true
	} else {
		update = dsp.GetAttributes().Add(attribute)
	}

	if update {
		if _, err := s.Update(ctx, nil, dsp); err != nil {
			return err
		}
	}

	return nil
}

func (s *DispatchDB) AddReferencesToDispatch(
	ctx context.Context,
	dsp *centrumdispatches.Dispatch,
	refs ...*centrumdispatches.DispatchReference,
) error {
	update := false
	if dsp.GetReferences() == nil {
		dsp.References = &centrumdispatches.DispatchReferences{
			References: refs,
		}

		update = true
	} else {
		for _, ref := range refs {
			upd := dsp.GetReferences().Add(ref)
			if upd {
				update = true
			}
		}
	}

	if update {
		if _, err := s.Update(ctx, nil, dsp); err != nil {
			return err
		}
	}

	return nil
}

func (s *DispatchDB) Delete(ctx context.Context, id int64, removeFromDB bool) error {
	if err := s.deleteInKV(ctx, id); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	if !removeFromDB {
		return nil
	}

	dispatches := table.FivenetCentrumDispatches
	stmt := dispatches.DELETE().WHERE(dispatches.ID.EQ(mysql.Int64(id))).LIMIT(1)
	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return errorscentrum.ErrFailedQuery
	}
	return nil
}
