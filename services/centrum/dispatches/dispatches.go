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

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/users"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/paulmach/orb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var metricDispatchLastID = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "centrum",
	Name:      "dispatch_last_id",
	Help:      "Last dispatch ID.",
}, []string{"job_name"})

type DispatchDB struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	tracker  tracker.ITracker
	postals  postals.Postals
	appCfg   appconfig.IConfig

	settings *settings.SettingsDB
	units    *units.UnitDB

	dispatchLocationsMutex *sync.Mutex
	dispatchLocations      map[string]*coords.Coords[*centrum.Dispatch]

	store      *store.Store[centrum.Dispatch, *centrum.Dispatch]
	jobMapping *store.Store[common.IDMapping, *common.IDMapping]
	idleKV     jetstream.KeyValue
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	JS        *events.JSWrapper
	DB        *sql.DB
	Cfg       *config.Config
	Enricher  *mstlystcdata.Enricher
	Tracker   tracker.ITracker
	Postals   postals.Postals
	AppConfig appconfig.IConfig

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

		settings: p.Settings,
		units:    p.Units,

		dispatchLocationsMutex: &sync.Mutex{},
		dispatchLocations:      map[string]*coords.Coords[*centrum.Dispatch]{},
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		idleKV, err := d.js.CreateOrUpdateKeyValue(ctxStartup, jetstream.KeyValueConfig{
			Bucket:         "centrum_dispatches_idle",
			Description:    "Timer keys that expire when a dispatch is inactive",
			Storage:        jetstream.MemoryStorage,
			History:        1,
			MaxBytes:       0,
			TTL:            0,
			LimitMarkerTTL: InactiveLimitMarkerTTL,
		})
		d.idleKV = idleKV

		jobSt, err := store.New[common.IDMapping, *common.IDMapping](ctxCancel, logger, p.JS, "centrum_dispatches",
			store.WithKVPrefix[common.IDMapping, *common.IDMapping]("job"),
			store.WithLocks[common.IDMapping, *common.IDMapping](nil),
			store.WithKVConfig[common.IDMapping, *common.IDMapping](jetstream.KeyValueConfig{TTL: 7 * 24 * time.Hour}),
		)
		if err != nil {
			return err
		}

		if err := jobSt.Start(ctxCancel, false); err != nil {
			return err
		}
		d.jobMapping = jobSt

		st, err := store.New[centrum.Dispatch, *centrum.Dispatch](ctxCancel, logger, p.JS, "centrum_dispatches",
			store.WithKVPrefix[centrum.Dispatch, *centrum.Dispatch]("id"),
			// Make sure dispatches are removed from the store after 7 days of inactivity (if all other cleanup mechanisms fail)
			store.WithKVConfig[centrum.Dispatch, *centrum.Dispatch](jetstream.KeyValueConfig{TTL: 7 * 24 * time.Hour}),
			store.WithOnUpdateFn[centrum.Dispatch, *centrum.Dispatch](func(ctx context.Context, _ *store.Store[centrum.Dispatch, *centrum.Dispatch], dispatch *centrum.Dispatch) (*centrum.Dispatch, error) {
				if dispatch == nil {
					return nil, nil
				}

				var errs error
				if err := d.TouchActivity(ctx, dispatch.Id); err != nil {
					errs = multierr.Append(errs, fmt.Errorf("failed to touch activity for dispatch %d. %w", dispatch.Id, err))
				}

				for _, job := range dispatch.Jobs.GetJobStrings() {
					if err := jobSt.Put(ctx, centrumutils.JobIdKey(job, dispatch.Id), &common.IDMapping{
						Id: dispatch.Id,
					}); err != nil {
						errs = multierr.Append(errs, fmt.Errorf("failed to update job %s mapping for dispatch %d. %w", job, dispatch.Id, err))
						continue
					}
				}

				return dispatch, nil
			}),
			store.WithOnDeleteFn[centrum.Dispatch, *centrum.Dispatch](func(ctx context.Context, s *store.Store[centrum.Dispatch, *centrum.Dispatch], _ string, dispatch *centrum.Dispatch) error {
				if dispatch == nil {
					return nil
				}

				var errs error
				for _, job := range dispatch.Jobs.GetJobStrings() {
					if err := jobSt.Delete(ctx, centrumutils.JobIdKey(job, dispatch.Id)); err != nil {
						errs = multierr.Append(errs, fmt.Errorf("failed to delete job %s mapping for dispatch %d. %w", job, dispatch.Id, err))
						continue
					}
				}

				return errs
			}),
			store.WithOnRemoteUpdatedFn[centrum.Dispatch, *centrum.Dispatch](func(ctx context.Context, _ *store.Store[centrum.Dispatch, *centrum.Dispatch], dispatch *centrum.Dispatch) (*centrum.Dispatch, error) {
				if dispatch == nil || dispatch.Jobs == nil {
					return dispatch, nil
				}

				addLocation := dispatch.Status != nil && centrumutils.IsStatusDispatchComplete(dispatch.Status.Status)
				// Ensure the dispatch has a valid ID
				for _, job := range dispatch.Jobs.GetJobStrings() {
					locs, ok := d.GetLocations(job)
					if locs == nil || !ok {
						continue
					}

					if addLocation {
						if locs.Has(dispatch, centrum.DispatchPointMatchFn(dispatch.Id)) {
							locs.Remove(dispatch, centrum.DispatchPointMatchFn(dispatch.Id))
						}
					} else {
						if err := locs.Replace(dispatch, func(p orb.Pointer) bool {
							return p.(*centrum.Dispatch).Id == dispatch.Id
						}, func(p1, p2 orb.Pointer) bool {
							return p1.Point().Equal(p2.Point())
						}); err != nil {
							d.logger.Error("failed to add non-existant dispatch to locations", zap.Uint64("dispatch_id", dispatch.Id))
						}
					}
				}

				return dispatch, nil
			}),
			store.WithOnRemoteDeletedFn[centrum.Dispatch, *centrum.Dispatch](func(ctx context.Context, _ *store.Store[centrum.Dispatch, *centrum.Dispatch], key string, dispatch *centrum.Dispatch) error {
				if dispatch != nil {
					for _, job := range dispatch.Jobs.GetJobStrings() {
						if locs, ok := d.GetLocations(job); ok && locs != nil {
							locs.Remove(nil, centrum.DispatchPointMatchFn(dispatch.Id))
						}
					}

					return nil
				}

				// Fallback to iterating over each job's locations map and delete the dispatch from the map by id
				split := strings.Split(key, ".")
				if len(split) < 1 {
					d.logger.Warn("unable to delete dispatch location, invalid key", zap.String("store_dispatch_key", key))
					return fmt.Errorf("invalid key format for dispatch remote delete. %s", key)
				}

				idKey := split[1]
				dspId, err := strconv.ParseUint(idKey, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse dispatch id from key %s. %w", key, err)
				}

				for _, job := range d.GetLocationsJob() {
					if locs, ok := d.GetLocations(job); ok && locs != nil {
						locs.Remove(nil, centrum.DispatchPointMatchFn(dspId))
					}
				}

				return nil
			}),
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

func (s *DispatchDB) LoadFromDB(ctx context.Context, cond jet.BoolExpression) error {
	tDispatch := table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus := table.FivenetCentrumDispatchesStatus.AS("dispatch_status")

	condition := tDispatchStatus.ID.IS_NULL().OR(
		jet.AND(
			tDispatchStatus.ID.EQ(
				jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
			).
				// Don't load archived dispatches into cache
				AND(tDispatchStatus.Status.NOT_IN(
					jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED)),
					jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
					jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
					jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_DELETED)),
				)),
		),
	)

	if cond != nil {
		condition = condition.AND(cond)
	}

	tUsers := tables.User().AS("user")

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
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				).
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tDispatchStatus.UserID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(200)

	dsps := []*centrum.Dispatch{}
	if err := stmt.QueryContext(ctx, s.db, &dsps); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	if len(dsps) == 0 {
		return nil
	}

	publicJobs := s.appCfg.Get().JobInfo.PublicJobs
	for i := range dsps {
		var err error
		dsps[i].Units, err = s.LoadDispatchAssignments(ctx, dsps[i].Id)
		if err != nil {
			return err
		}

		if dsps[i].CreatorId != nil && *dsps[i].CreatorId > 0 {
			dsps[i].Creator, err = users.RetrieveUserById(ctx, s.db, *dsps[i].CreatorId)
			if err != nil {
				return err
			}

			if dsps[i].Creator != nil {
				// Clear dispatch creator's job info if not a visible job
				if !slices.Contains(publicJobs, dsps[i].Creator.Job) {
					dsps[i].Creator.Job = ""
				}
				dsps[i].Creator.JobGrade = 0
			}
		}

		if dsps[i].Postal == nil {
			if postal, ok := s.postals.Closest(dsps[i].X, dsps[i].Y); postal != nil && ok {
				dsps[i].Postal = postal.Code
			}
		}

		// Ensure dispatch has a status
		if dsps[i].Status == nil {
			dsps[i].Status, err = s.AddDispatchStatus(ctx, s.db, &centrum.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dsps[i].Id,
				Status:     centrum.StatusDispatch_STATUS_DISPATCH_NEW,
				Postal:     dsps[i].Postal,
				X:          &dsps[i].X,
				Y:          &dsps[i].Y,
			}, false, nil)
			if err != nil {
				return err
			}
		}

		// Ensure dispatch has a valid job list (fallback to deprecated Jobs field for old dispatches)
		if dsps[i].Jobs == nil || len(dsps[i].Jobs.GetJobs()) == 0 {
			dsps[i].Jobs = &centrum.JobList{
				Jobs: []*centrum.Job{
					{
						Name: dsps[i].Job,
					},
				},
			}
			dsps[i].Job = ""
		}
		for _, job := range dsps[i].Jobs.GetJobs() {
			s.enricher.EnrichJobName(job)
		}

		// Update dispatch in db and in kv
		if _, err := s.Update(ctx, nil, dsps[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *DispatchDB) LoadDispatchAssignments(ctx context.Context, dispatchId uint64) ([]*centrum.DispatchAssignment, error) {
	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts.AS("dispatch_assignment")

	stmt := tDispatchUnit.
		SELECT(
			tDispatchUnit.DispatchID,
			tDispatchUnit.UnitID,
			tDispatchUnit.CreatedAt,
			tDispatchUnit.ExpiresAt,
		).
		FROM(tDispatchUnit).
		ORDER_BY(
			tDispatchUnit.CreatedAt.ASC(),
		).
		WHERE(
			tDispatchUnit.DispatchID.EQ(jet.Uint64(dispatchId)),
		)

	dest := []*centrum.DispatchAssignment{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	// Retrieve units based on the dispatch unit assignments
	for i := range dest {
		unit, err := s.units.Get(ctx, dest[i].UnitId)
		if unit == nil || err != nil {
			return nil, fmt.Errorf("no unit found for dispatch id %d with id %d", dispatchId, dest[i].UnitId)
		}

		dest[i].Unit = unit
	}

	return dest, nil
}

func (s *DispatchDB) GetLocations(job string) (*coords.Coords[*centrum.Dispatch], bool) {
	s.dispatchLocationsMutex.Lock()
	defer s.dispatchLocationsMutex.Unlock()

	locations, ok := s.dispatchLocations[job]
	if !ok {
		locations = coords.New[*centrum.Dispatch]()
		s.dispatchLocations[job] = locations
	}
	return locations, ok
}

func (s *DispatchDB) GetLocationsJob() []string {
	s.dispatchLocationsMutex.Lock()
	defer s.dispatchLocationsMutex.Unlock()

	var jobs []string
	for job := range s.dispatchLocations {
		jobs = append(jobs, job)
	}

	return jobs
}

func (s *DispatchDB) Delete(ctx context.Context, id uint64, removeFromDB bool) error {
	if err := s.deleteInKV(ctx, id); err != nil {
		return err
	}

	if removeFromDB {
		tDispatch := table.FivenetCentrumDispatches

		stmt := tDispatch.
			DELETE().
			WHERE(jet.AND(
				tDispatch.ID.EQ(jet.Uint64(id)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return errorscentrum.ErrFailedQuery
		}
	}

	return nil
}

func (s *DispatchDB) UpdateStatus(ctx context.Context, dspId uint64, in *centrum.DispatchStatus) (*centrum.DispatchStatus, error) {
	dsp, err := s.Get(ctx, dspId)
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, err
		}
	}

	if dsp != nil && dsp.Status != nil {
		// If the dispatch status is the same and is a status that shouldn't be duplicated, don't update the status again
		if dsp.Status.Status == in.Status &&
			(in.Status == centrum.StatusDispatch_STATUS_DISPATCH_NEW ||
				in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED) {
			s.logger.Debug("skipping dispatch status update due to being new or same status", zap.Uint64("dispatch_id", dsp.Id), zap.String("status", in.Status.String()))
			return in, nil
		}

		// If the dispatch is complete, we ignore any unit unassignments/accepts/declines
		if centrumutils.IsStatusDispatchComplete(dsp.Status.Status) &&
			(in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED ||
				in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED ||
				in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED ||
				in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED) {
			return in, nil
		}
	}

	s.logger.Debug("updating dispatch status", zap.Uint64("dispatch_id", dspId), zap.String("status", in.Status.String()))

	if in.UserId != nil {
		var err error
		in.User, err = users.RetrieveUserShortById(ctx, s.db, s.enricher, *in.UserId)
		if err != nil {
			return nil, err
		}

		if um, ok := s.tracker.GetUserMarkerById(*in.UserId); ok {
			in.X = &um.X
			in.Y = &um.Y
			in.Postal = um.Postal
		}
	}

	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt := tDispatchStatus.
		INSERT(
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
		VALUES(
			jet.CURRENT_TIMESTAMP(),
			in.DispatchId,
			in.UnitId,
			in.Status,
			in.Reason,
			in.Code,
			in.UserId,
			in.X,
			in.Y,
			in.Postal,
			in.CreatorJob,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	in.Id = uint64(lastId)

	if err := s.updateStatusInKV(ctx, in.DispatchId, in); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	for _, job := range dsp.Jobs.GetJobStrings() {
		if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchStatus, job), data); err != nil {
			return nil, fmt.Errorf("failed to publish dispatch status event (size: %d, message: '%+v'). %w", len(data), in, err)
		}
	}

	return in, nil
}

func (s *DispatchDB) UpdateAssignments(ctx context.Context, userId *int32, dspId uint64, toAdd []uint64, toRemove []uint64, expiresAt time.Time) error {
	s.logger.Debug("updating dispatch assignments", zap.Int32p("user_id", userId), zap.Uint64("dispatch_id", dspId), zap.Uint64s("toAdd", toAdd), zap.Uint64s("toRemove", toRemove))

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return nil
	}

	var x, y *float64
	var postal *string
	if userId != nil {
		if um, ok := s.tracker.GetUserMarkerById(*userId); ok {
			x = &um.X
			y = &um.Y
			postal = um.Postal
		}
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if len(toRemove) > 0 {
		removeIds := make([]jet.Expression, len(toRemove))
		for i := range toRemove {
			removeIds[i] = jet.Uint64(toRemove[i])
		}

		stmt := tDispatchUnit.
			DELETE().
			WHERE(jet.AND(
				tDispatchUnit.DispatchID.EQ(jet.Uint64(dspId)),
				tDispatchUnit.UnitID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	var expiresAtTS *timestamp.Timestamp
	// If expires at time is not zero
	expiresAtVal := jet.NULL
	if !expiresAt.IsZero() {
		expiresAtTS = timestamp.New(expiresAt)
		expiresAtVal = jet.TimeT(expiresAt)
	}

	if len(toAdd) > 0 {
		units := []uint64{}
		dsp, err := s.Get(ctx, dspId)
		if err != nil {
			return err
		}
		for i := range toAdd {

			// Skip already added units
			if slices.ContainsFunc(dsp.Units, func(in *centrum.DispatchAssignment) bool {
				return in.UnitId == toAdd[i]
			}) {
				continue
			}

			unit, err := s.units.Get(ctx, toAdd[i])
			if err != nil {
				continue
			}

			// Skip empty units
			if len(unit.Users) == 0 {
				continue
			}

			// Only add unit to dispatch if not already assigned/in list
			units = append(units, toAdd[i])
		}

		if len(units) > 0 {
			stmt := tDispatchUnit.
				INSERT(
					tDispatchUnit.DispatchID,
					tDispatchUnit.UnitID,
					tDispatchUnit.ExpiresAt,
				)

			for _, unitId := range units {
				stmt = stmt.
					VALUES(
						dspId,
						unitId,
						expiresAtVal,
					)
			}

			stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
				tDispatchUnit.ExpiresAt.SET(jet.RawTimestamp("VALUES(`expires_at`)")),
			)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	key := centrumutils.IdKey(dspId)
	if err := s.store.ComputeUpdate(ctx, key, true, func(key string, dsp *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
		if dsp == nil {
			s.logger.Error("nil dispatch in computing dispatch assignment logic", zap.String("key", key), zap.Any("dsp", dsp))
			return dsp, false, nil
		}

		if len(toRemove) > 0 {
			toAnnounce := []uint64{}
			dsp.Units = slices.DeleteFunc(dsp.Units, func(in *centrum.DispatchAssignment) bool {
				for k := range toRemove {
					if in.UnitId != toRemove[k] {
						continue
					}

					toAnnounce = append(toAnnounce, toRemove[k])
					return true
				}

				return false
			})

			// Send updates
			for _, unitId := range toAnnounce {
				if _, err := s.AddDispatchStatus(ctx, s.db, &centrum.DispatchStatus{
					CreatedAt:  timestamp.Now(),
					DispatchId: dsp.Id,
					UnitId:     &unitId,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED,
					UserId:     userId,
					X:          x,
					Y:          y,
					Postal:     postal,
				}, true, dsp.Jobs.GetJobStrings()); err != nil {
					return nil, false, err
				}
			}
		}

		if len(toAdd) > 0 {
			units := []uint64{}
			for i := range toAdd {
				// Skip already added units
				if slices.ContainsFunc(dsp.Units, func(in *centrum.DispatchAssignment) bool {
					return in.UnitId == toAdd[i]
				}) {
					continue
				}

				unit, err := s.units.Get(ctx, toAdd[i])
				if err != nil {
					continue
				}

				// Skip empty units
				if len(unit.Users) == 0 {
					continue
				}

				// Only add unit to dispatch if not already assigned/in list
				units = append(units, toAdd[i])
			}

			for _, unitId := range units {
				unit, err := s.units.Get(ctx, unitId)
				if err != nil {
					continue
				}

				dsp.Units = append(dsp.Units, &centrum.DispatchAssignment{
					DispatchId: dsp.Id,
					UnitId:     unit.Id,
					Unit:       unit,
					ExpiresAt:  expiresAtTS,
				})
			}

			for _, unitId := range units {
				if _, err := s.AddDispatchStatus(ctx, s.db, &centrum.DispatchStatus{
					CreatedAt:  timestamp.Now(),
					DispatchId: dsp.Id,
					UnitId:     &unitId,
					UserId:     userId,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
					X:          x,
					Y:          y,
					Postal:     postal,
				}, true, dsp.Jobs.GetJobStrings()); err != nil {
					return nil, false, err
				}
			}
		}

		return dsp, len(toRemove) > 0 || len(toAdd) > 0, nil
	}); err != nil {
		return err
	}

	dsp, err := s.Get(ctx, dspId)
	if err != nil {
		return err
	}

	// Dispatch has no units assigned anymore
	if len(dsp.Units) == 0 {
		// Check dispatch status to not be completed/archived, etc.
		if dsp.Status != nil && !centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
			if _, err := s.UpdateStatus(ctx, dspId, &centrum.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dspId,
				Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
				UserId:     userId,
				X:          x,
				Y:          y,
				Postal:     postal,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *DispatchDB) Create(ctx context.Context, dsp *centrum.Dispatch) (*centrum.Dispatch, error) {
	// Check if the dispatch has at least one job, till the Job field is removed keep using it as a fallback
	if (dsp.Jobs == nil || len(dsp.Jobs.GetJobs()) == 0) && dsp.Job == "" {
		return nil, errorscentrum.ErrDispatchNoJobs
	}

	// If the deprecated Job field is used, convert it to Jobs but only if the jobs list is empty
	if dsp.Jobs == nil || len(dsp.Jobs.GetJobs()) == 0 {
		dsp.Jobs = &centrum.JobList{
			Jobs: []*centrum.Job{
				{
					Name: dsp.Job,
				},
			},
		}
		dsp.Job = ""
	}

	for _, job := range dsp.Jobs.GetJobs() {
		s.enricher.EnrichJobName(job)
	}

	if dsp.Postal == nil || *dsp.Postal == "" {
		if postal, ok := s.postals.Closest(dsp.X, dsp.Y); postal != nil && ok {
			dsp.Postal = postal.Code
		}
	}

	if dsp.CreatorId != nil {
		var err error
		dsp.Creator, err = users.RetrieveUserById(ctx, s.db, *dsp.CreatorId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user for dispatch creator. %w", err)
		}
		// Unset creator in case we don't have a user
		if dsp.Creator == nil {
			dsp.Creator = nil
			dsp.CreatorId = nil
		} else if !slices.Contains(dsp.Jobs.GetJobStrings(), dsp.Creator.Job) {
			// Remove creator props when job isn't equal
			dsp.Creator.Props = nil
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
			jet.CURRENT_TIMESTAMP(),
			dsp.Jobs,
			dsp.Message,
			dsp.Description,
			dsp.Attributes,
			dsp.References,
			dsp.X,
			dsp.Y,
			dsp.Postal,
			dsp.Anon,
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
	dsp.Id = uint64(lastId)

	var userId *int32
	if !dsp.Anon && dsp.CreatorId != nil {
		userId = dsp.CreatorId
	}

	var statusUser *jobs.Colleague
	if dsp.Creator != nil {
		statusUser = dsp.Creator.Colleague()
	}

	if dsp.Status, err = s.AddDispatchStatus(ctx, tx, &centrum.DispatchStatus{
		CreatedAt:  timestamp.Now(),
		DispatchId: dsp.Id,
		UserId:     userId,
		User:       statusUser,
		Status:     centrum.StatusDispatch_STATUS_DISPATCH_NEW,
		X:          &dsp.X,
		Y:          &dsp.Y,
		Postal:     dsp.Postal,
	}, false, nil); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	for _, job := range dsp.Jobs.GetJobStrings() {
		metricDispatchLastID.WithLabelValues(job).Set(float64(lastId))
	}

	// Hide user info when dispatch is anonymous
	if dsp.Anon {
		dsp.Creator = nil
	}

	if err := s.updateInKV(ctx, dsp.Id, dsp); err != nil {
		return nil, err
	}

	return dsp, nil
}

func (s *DispatchDB) Update(ctx context.Context, userId *int32, dsp *centrum.Dispatch) (*centrum.Dispatch, error) {
	tDispatch := table.FivenetCentrumDispatches.AS("dispatch")

	dsp.UpdatedAt = timestamp.Now()
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
			jet.CURRENT_TIMESTAMP(),
			dsp.Jobs,
			dsp.Message,
			dsp.Description,
			dsp.Attributes,
			dsp.References,
			dsp.X,
			dsp.Y,
			dsp.Postal,
			dsp.Anon,
			dsp.CreatorId,
		).
		WHERE(jet.AND(
			tDispatch.ID.EQ(jet.Uint64(dsp.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	if err := s.updateInKV(ctx, dsp.Id, dsp); err != nil {
		return nil, err
	}

	return dsp, nil
}

func (s *DispatchDB) AddDispatchStatus(ctx context.Context, tx qrm.DB, status *centrum.DispatchStatus, publish bool, jobs []string) (*centrum.DispatchStatus, error) {
	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt := tDispatchStatus.
		INSERT(
			tDispatchStatus.CreatedAt,
			tDispatchStatus.DispatchID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UnitID,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tDispatchStatus.CreatorJob,
		).
		VALUES(
			status.CreatedAt,
			status.DispatchId,
			status.Status,
			status.Reason,
			status.Code,
			status.UnitId,
			status.UserId,
			status.X,
			status.Y,
			status.Postal,
			status.CreatorJob,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	newStatus, err := s.GetStatus(ctx, tx, uint64(lastId))
	if err != nil {
		return nil, err
	}

	if publish {
		data, err := proto.Marshal(newStatus)
		if err != nil {
			return nil, err
		}

		for _, job := range jobs {
			if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchStatus, job), data); err != nil {
				return nil, err
			}
		}
	}

	return newStatus, nil
}

func (s *DispatchDB) GetStatus(ctx context.Context, tx qrm.DB, id uint64) (*centrum.DispatchStatus, error) {
	tDispatchStatus := table.FivenetCentrumDispatchesStatus.AS("dispatch_status")
	tUsers := tables.User().AS("colleague")

	stmt := tDispatchStatus.
		SELECT(
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
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tDispatchStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tDispatchStatus.UserID),
				),
		).
		WHERE(
			tDispatchStatus.ID.EQ(jet.Uint64(id)),
		).
		ORDER_BY(tDispatchStatus.ID.DESC()).
		LIMIT(1)

	var dest centrum.DispatchStatus
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	if dest.UnitId != nil && *dest.UnitId > 0 && dest.User != nil {
		unit, err := s.units.Get(ctx, *dest.UnitId)
		if err != nil {
			return nil, err
		}

		dest.Unit = unit
	}

	return &dest, nil
}

func (s *DispatchDB) TakeDispatch(ctx context.Context, userJob string, userId int32, unitId uint64, resp centrum.TakeDispatchResp, dispatchIds []uint64) error {
	settings, err := s.settings.Get(ctx, userJob)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	// If the dispatch center is in central command mode, units can't self assign dispatches
	if settings.Mode == centrum.CentrumMode_CENTRUM_MODE_CENTRAL_COMMAND {
		return errorscentrum.ErrModeForbidsAction
	}

	unit, err := s.units.Get(ctx, unitId)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	var x, y *float64
	var postal *string
	if um, ok := s.tracker.GetUserMarkerById(userId); ok {
		x = &um.X
		y = &um.Y
		postal = um.Postal
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts

	for _, dspId := range dispatchIds {
		if resp == centrum.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
			stmt := tDispatchUnit.
				INSERT(
					tDispatchUnit.DispatchID,
					tDispatchUnit.UnitID,
					tDispatchUnit.ExpiresAt,
				).
				VALUES(
					dspId,
					unit.Id,
					jet.NULL,
				).
				ON_DUPLICATE_KEY_UPDATE(
					tDispatchUnit.ExpiresAt.SET(jet.TimestampExp(jet.NULL)),
				)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
			}
		} else {
			stmt := tDispatchUnit.
				DELETE().
				WHERE(jet.AND(
					tDispatchUnit.DispatchID.EQ(jet.Uint64(dspId)),
					tDispatchUnit.UnitID.EQ(jet.Uint64(unit.Id)),
				)).
				LIMIT(1)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
			}
		}

		key := centrumutils.IdKey(dspId)
		if err := s.store.ComputeUpdate(ctx, key, true, func(key string, dsp *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
			// If dispatch is nil or completed, disallow to accept the dispatch
			if dsp == nil || (dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status)) {
				return nil, false, errorscentrum.ErrDispatchAlreadyCompleted
			}

			status := centrum.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED

			// Dispatch accepted
			if resp == centrum.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
				status = centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED

				found := false
				accepted := true
				// Set unit expires at to nil
				for _, ua := range dsp.Units {
					if ua.UnitId == unit.Id {
						found = true
						// If there's no expiration time the unit has been directly assigned
						if ua.ExpiresAt == nil {
							accepted = false
						}
						ua.ExpiresAt = nil
						break
					}
				}

				if !found {
					dsp.Units = append(dsp.Units, &centrum.DispatchAssignment{
						DispatchId: dsp.Id,
						UnitId:     unit.Id,
						Unit:       unit,
						CreatedAt:  timestamp.Now(),
					})
				}

				if accepted {
					// Set unit to busy when unit accepts a dispatch
					if unit.Status == nil || unit.Status.Status != centrum.StatusUnit_STATUS_UNIT_BUSY {
						if _, err := s.units.UpdateStatus(ctx, unit.Id, &centrum.UnitStatus{
							CreatedAt:  timestamp.Now(),
							UnitId:     unit.Id,
							Status:     centrum.StatusUnit_STATUS_UNIT_BUSY,
							UserId:     &userId,
							CreatorId:  &userId,
							X:          x,
							Y:          y,
							Postal:     postal,
							CreatorJob: &userJob,
						}); err != nil {
							return nil, false, err
						}
					}
				}
			} else {
				// Dispatch declined
				status = centrum.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED

				// Remove the unit's assignment
				dsp.Units = slices.DeleteFunc(dsp.Units, func(in *centrum.DispatchAssignment) bool {
					return in.UnitId == unit.Id
				})
			}

			if dsp.Status, err = s.AddDispatchStatus(ctx, s.db, &centrum.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dspId,
				Status:     status,
				UnitId:     &unitId,
				UserId:     &userId,
				X:          x,
				Y:          y,
				Postal:     postal,
				CreatorJob: &userJob,
			}, true, dsp.Jobs.GetJobStrings()); err != nil {
				return nil, false, err
			}

			return dsp, true, nil
		}); err != nil {
			// Ignore errors that are "okay" to encounter
			if !errors.Is(err, errorscentrum.ErrDispatchAlreadyCompleted) {
				return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}
	}

	return nil
}

func (s *DispatchDB) AddAttributeToDispatch(ctx context.Context, dsp *centrum.Dispatch, attribute centrum.DispatchAttribute) error {
	update := false
	if dsp.Attributes == nil {
		dsp.Attributes = &centrum.DispatchAttributes{
			List: []centrum.DispatchAttribute{attribute},
		}

		update = true
	} else {
		update = dsp.Attributes.Add(attribute)
	}

	if update {
		if _, err := s.Update(ctx, nil, dsp); err != nil {
			return err
		}
	}

	return nil
}

func (s *DispatchDB) AddReferencesToDispatch(ctx context.Context, dsp *centrum.Dispatch, refs ...*centrum.DispatchReference) error {
	update := false
	if dsp.References == nil {
		dsp.References = &centrum.DispatchReferences{
			References: refs,
		}

		update = true
	} else {
		for _, ref := range refs {
			upd := dsp.References.Add(ref)
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
