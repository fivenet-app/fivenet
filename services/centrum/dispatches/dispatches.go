package dispatches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
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
	"github.com/paulmach/orb"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"
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

	dispatchLocationsMutex *sync.Mutex
	dispatchLocations      map[string]*coords.Coords[*centrumdispatches.Dispatch]

	store      *store.Store[centrumdispatches.Dispatch, *centrumdispatches.Dispatch]
	jobMapping *store.Store[common.IDMapping, *common.IDMapping]
	idleKV     jetstream.KeyValue
	metrics    *dispatchMetrics
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

		dispatchLocationsMutex: &sync.Mutex{},
		dispatchLocations:      map[string]*coords.Coords[*centrumdispatches.Dispatch]{},
		metrics:                getDispatchMetrics(),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		idleKV, err := d.js.CreateOrUpdateKeyValue(ctxStartup, jetstream.KeyValueConfig{
			Bucket:         "centrum_dispatches_idle",
			Description:    "Timer keys for dispatch inactivity, cleanup, projections, and assignments",
			Storage:        jetstream.MemoryStorage,
			History:        1,
			MaxBytes:       0,
			TTL:            0,
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
				jetstream.KeyValueConfig{TTL: 7 * 24 * time.Hour},
			),
		)
		if err != nil {
			return err
		}

		if err := jobSt.Start(ctxCancel, false); err != nil {
			return err
		}
		d.jobMapping = jobSt
		ensureJobMapping := func(ctx context.Context, job string, dispatchID int64, source string, refresh bool) error {
			if strings.TrimSpace(job) == "" {
				return nil
			}
			key := centrumutils.JobIdKey(job, dispatchID)
			mapping := &common.IDMapping{Id: dispatchID}
			var changed bool
			var err error
			if refresh {
				// The job mapping is the job-scoped stream trigger. Rewrite it for
				// every local dispatch projection update, even when its ID is unchanged.
				err = jobSt.Put(ctx, key, mapping)
				changed = err == nil
			} else {
				// Remote watchers only repair a missing mapping. Rewriting it from
				// every server would duplicate the stream trigger.
				changed, err = jobSt.PutIfChanged(ctx, key, mapping)
			}
			logger.Debug(
				"ensured dispatch job mapping",
				zap.String("source", source),
				zap.String("key", key),
				zap.Int64("dispatch_id", dispatchID),
				zap.Bool("changed", changed),
				zap.Error(err),
			)
			return err
		}

		st, err := store.New[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](
			ctxCancel,
			logger,
			p.JS,
			"centrum_dispatches",
			store.WithKVPrefix[centrumdispatches.Dispatch, *centrumdispatches.Dispatch]("id"),
			// Make sure dispatches are removed from the store after 7 days of inactivity (if all other cleanup mechanisms fail)
			store.WithKVConfig[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](
				jetstream.KeyValueConfig{TTL: 7 * 24 * time.Hour},
			),
			store.WithOnUpdateFn[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](
				func(ctx context.Context, old *centrumdispatches.Dispatch, dispatch *centrumdispatches.Dispatch) (*centrumdispatches.Dispatch, error) {
					if dispatch == nil {
						return nil, nil
					}

					var errs error
					if err := d.TouchActivity(ctx, dispatch.GetId()); err != nil {
						errs = multierr.Append(
							errs,
							fmt.Errorf(
								"failed to touch activity for dispatch %d. %w",
								dispatch.GetId(),
								err,
							),
						)
					}

					newJobSet := make(map[string]struct{}, len(dispatch.GetJobs().GetJobStrings()))
					for _, job := range dispatch.GetJobs().GetJobStrings() {
						if strings.TrimSpace(job) == "" {
							continue
						}
						newJobSet[job] = struct{}{}
						if err := ensureJobMapping(
							ctx,
							job,
							dispatch.GetId(),
							"local_update",
							true,
						); err != nil {
							errs = multierr.Append(
								errs,
								fmt.Errorf(
									"failed to update job %s mapping for dispatch %d. %w",
									job,
									dispatch.GetId(),
									err,
								),
							)
							continue
						}
					}

					if old != nil && old.GetJobs() != nil {
						for _, oldJob := range old.GetJobs().GetJobStrings() {
							if strings.TrimSpace(oldJob) == "" {
								continue
							}
							if _, ok := newJobSet[oldJob]; ok {
								continue
							}

							if err := jobSt.Delete(
								ctx,
								centrumutils.JobIdKey(oldJob, dispatch.GetId()),
							); err != nil {
								errs = multierr.Append(
									errs,
									fmt.Errorf(
										"failed to delete stale job %s mapping for dispatch %d. %w",
										oldJob,
										dispatch.GetId(),
										err,
									),
								)
							}
						}
					}

					if errs != nil {
						return nil, fmt.Errorf(
							"failed to update dispatch %d in kv store. %w",
							dispatch.GetId(),
							errs,
						)
					}

					return dispatch, nil
				},
			),
			store.WithOnDeleteFn[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](
				func(ctx context.Context, _ string, dispatch *centrumdispatches.Dispatch) error {
					if dispatch == nil {
						return nil
					}

					var errs error
					for _, job := range dispatch.GetJobs().GetJobStrings() {
						if strings.TrimSpace(job) == "" {
							continue
						}
						if err := jobSt.Delete(
							ctx,
							centrumutils.JobIdKey(job, dispatch.GetId()),
						); err != nil {
							errs = multierr.Append(
								errs,
								fmt.Errorf(
									"failed to delete job %s mapping for dispatch %d. %w",
									job,
									dispatch.GetId(),
									err,
								),
							)
							continue
						}
					}

					if err := d.idleKV.Delete(
						ctx,
						"idle."+centrumutils.IdKey(dispatch.GetId()),
					); err != nil {
						errs = multierr.Append(
							errs,
							fmt.Errorf(
								"failed to delete idle key for dispatch %d. %w",
								dispatch.GetId(),
								err,
							),
						)
					}

					return errs
				},
			),
			store.WithOnRemoteUpdatedFn[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](
				func(ctx context.Context, old *centrumdispatches.Dispatch, dispatch *centrumdispatches.Dispatch) (*centrumdispatches.Dispatch, error) {
					if dispatch == nil {
						return dispatch, nil
					}

					newJobSet := make(map[string]struct{}, len(dispatch.GetJobs().GetJobStrings()))
					for _, job := range dispatch.GetJobs().GetJobStrings() {
						if strings.TrimSpace(job) == "" {
							continue
						}
						newJobSet[job] = struct{}{}
						if err := ensureJobMapping(
							ctx,
							job,
							dispatch.GetId(),
							"remote_update",
							false,
						); err != nil {
							return nil, fmt.Errorf(
								"failed to update job %s mapping for dispatch %d. %w",
								job,
								dispatch.GetId(),
								err,
							)
						}
					}

					if old != nil && old.GetJobs() != nil {
						for _, oldJob := range old.GetJobs().GetJobStrings() {
							if strings.TrimSpace(oldJob) == "" {
								continue
							}
							if _, ok := newJobSet[oldJob]; ok {
								continue
							}

							if err := jobSt.Delete(
								ctx,
								centrumutils.JobIdKey(oldJob, dispatch.GetId()),
							); err != nil {
								return nil, fmt.Errorf(
									"failed to delete stale job %s mapping for dispatch %d. %w",
									oldJob,
									dispatch.GetId(),
									err,
								)
							}
						}
					}

					removeLoc := dispatch.GetStatus() != nil &&
						centrumutils.IsStatusDispatchComplete(dispatch.GetStatus().GetStatus())
					// Ensure the dispatch has a valid ID
					for _, job := range dispatch.GetJobs().GetJobStrings() {
						if strings.TrimSpace(job) == "" {
							continue
						}
						locs := d.GetLocations(job)
						if locs == nil {
							continue
						}

						if removeLoc {
							if locs.Has(
								dispatch,
								centrumdispatches.DispatchPointMatchFn(dispatch.GetId()),
							) {
								locs.Remove(
									dispatch,
									centrumdispatches.DispatchPointMatchFn(dispatch.GetId()),
								)
							}
						} else {
							if err := locs.Replace(
								dispatch,
								centrumdispatches.DispatchPointMatchFn(dispatch.GetId()),
								func(p1, p2 orb.Pointer) bool {
									return p1.Point().Equal(p2.Point())
								},
							); err != nil {
								d.logger.Error(
									"failed to add non-existent dispatch to locations",
									zap.Int64("dispatch_id", dispatch.GetId()),
								)
							}
						}
					}

					return dispatch, nil
				},
			),
			store.WithOnRemoteDeletedFn[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](
				func(ctx context.Context, key string, dispatch *centrumdispatches.Dispatch) error {
					if dispatch != nil {
						var errs error
						for _, job := range dispatch.GetJobs().GetJobStrings() {
							if strings.TrimSpace(job) == "" {
								continue
							}
							if err := jobSt.Delete(
								ctx,
								centrumutils.JobIdKey(job, dispatch.GetId()),
							); err != nil {
								errs = multierr.Append(
									errs,
									fmt.Errorf(
										"failed to delete job %s mapping for dispatch %d. %w",
										job,
										dispatch.GetId(),
										err,
									),
								)
							}
						}

						for _, job := range dispatch.GetJobs().GetJobStrings() {
							if strings.TrimSpace(job) == "" {
								continue
							}
							if locs := d.GetLocations(job); locs != nil {
								locs.Remove(
									nil,
									centrumdispatches.DispatchPointMatchFn(dispatch.GetId()),
								)
							}
						}

						return errs
					}

					// Fallback to iterating over each job's locations map and delete the dispatch from the map by id
					split := strings.Split(key, ".")
					if len(split) < 2 {
						d.logger.Warn(
							"unable to delete dispatch location, invalid key",
							zap.String("store_dispatch_key", key),
						)
						return fmt.Errorf("invalid key format for dispatch remote delete. %s", key)
					}

					idKey := split[1]
					dspId, err := strconv.ParseInt(idKey, 10, 64)
					if err != nil {
						return fmt.Errorf("failed to parse dispatch id from key %s. %w", key, err)
					}

					var errs error
					for _, jobKey := range jobSt.KeysFiltered("", func(jobKey string) bool {
						id, err := centrumutils.ExtractIDString(jobKey)
						return err == nil && id == idKey
					}) {
						if err := jobSt.Delete(ctx, jobKey); err != nil {
							errs = multierr.Append(
								errs,
								fmt.Errorf(
									"failed to delete job mapping %s for dispatch %d. %w",
									jobKey,
									dspId,
									err,
								),
							)
						}
					}

					for _, job := range d.GetLocationsJob() {
						if locs := d.GetLocations(job); locs != nil {
							locs.Remove(nil, centrumdispatches.DispatchPointMatchFn(dspId))
						}
					}

					return errs
				},
			),
		)
		if err != nil {
			return err
		}

		if err := st.Start(ctxCancel, false); err != nil {
			return err
		}
		d.store = st

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return d
}

func (s *DispatchDB) LoadFromDB(ctx context.Context, cond mysql.BoolExpression) (int, error) {
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

		// Ensure dispatch has a valid job list (fallback to deprecated Jobs field for old dispatches)
		if dsp.GetJobs() == nil || len(dsp.GetJobs().GetJobs()) == 0 {
			dsp.Jobs = &centrum.JobList{
				Jobs: []*centrum.JobListEntry{
					{
						//nolint:staticcheck // This is a fallback for old dispatches.
						Name: dsp.GetJob(),
					},
				},
			}
			//nolint:staticcheck // Clear old job info. This is a fallback for old dispatches.
			dsp.Job = ""
		}
		for _, job := range dsp.GetJobs().GetJobs() {
			s.enricher.EnrichJobName(job)
		}

		// Update dispatch in db and in kv
		if _, err := s.Update(ctx, nil, dsp); err != nil {
			return 0, err
		}

		for _, job := range dsp.GetJobs().GetJobStrings() {
			locs := s.GetLocations(job)
			if locs == nil {
				continue
			}

			if !locs.Has(dsp, centrumdispatches.DispatchPointMatchFn(dsp.GetId())) {
				locs.Add(dsp)
			} else {
				err := locs.Replace(
					dsp,
					centrumdispatches.DispatchPointMatchFn(dsp.GetId()),
					func(p1, p2 orb.Pointer) bool {
						return p1.Point().Equal(p2.Point())
					},
				)
				if err != nil {
					s.logger.Error(
						"failed to replace dispatch in locations",
						zap.Int64("dispatch_id", dsp.GetId()),
						zap.Error(err),
					)
				}
			}
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

func (s *DispatchDB) Create(
	ctx context.Context,
	dsp *centrumdispatches.Dispatch,
) (*centrumdispatches.Dispatch, error) {
	// Check if the dispatch has at least one job, till the Job field is removed keep using it as a fallback
	//nolint:staticcheck // Check old job info. This is a fallback for old dispatches.
	if (dsp.GetJobs() == nil || len(dsp.GetJobs().GetJobs()) == 0) && dsp.GetJob() == "" {
		return nil, errorscentrum.ErrDispatchNoJobs
	}

	// If the deprecated Job field is used, convert it to Jobs but only if the jobs list is empty
	if dsp.GetJobs() == nil || len(dsp.GetJobs().GetJobs()) == 0 {
		dsp.SetJobs(&centrum.JobList{
			Jobs: []*centrum.JobListEntry{
				{
					//nolint:staticcheck // This is a fallback for old dispatches.
					Name: dsp.GetJob(),
				},
			},
		})
		//nolint:staticcheck // Clear old job info. This is a fallback for old dispatches.
		dsp.SetJob("")
	}

	for _, job := range dsp.GetJobs().GetJobs() {
		if strings.TrimSpace(job.GetName()) == "" {
			return nil, errorscentrum.ErrDispatchNoJobs
		}
		s.enricher.EnrichJobName(job)
	}

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

	if err := s.updateInKV(ctx, dsp.GetId(), dsp); err != nil {
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
