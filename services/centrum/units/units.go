package units

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	unitsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	PingTTL  = 8 * time.Second  // OFFLINE
	EmptyTTL = 30 * time.Second // No crew -> delete
)

var unitSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(unitsaccess.UnitAccessLevel_UNIT_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(unitsaccess.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN),
	},
}

func UnitSubjectAccessOptions() access.SubjectAccessOptions {
	return unitSubjectAccessOptions
}

type sortOrderResult struct {
	SortOrder int32 `alias:"sort_order"`
}

type UnitDB struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher mstlystcdata.IEnricher
	tracker  tracker.ITracker
	postals  postals.Postals
	jobs     jobsstore.IStore

	store      *store.Store[centrumunits.Unit, *centrumunits.Unit]
	jobMapping *store.Store[common.IDMapping, *common.IDMapping]

	unitAccess         *access.CentrumUnitsObjectAccess
	unitAccessResolver *access.SubjectResolver
	colleagueHydrator  colleagueshydrator.IHydrator

	KVPing jetstream.KeyValue
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger            *zap.Logger
	JS                *events.JSWrapper
	DB                *sql.DB
	Cfg               *config.Config
	Enricher          mstlystcdata.IEnricher
	Tracker           tracker.ITracker
	Postals           postals.Postals
	Jobs              jobsstore.IStore
	UnitAccess        *access.CentrumUnitsObjectAccess
	ColleagueHydrator colleagueshydrator.IHydrator
}

func New(p Params) *UnitDB {
	logger := p.Logger.Named("centrum.units")

	ctxCancel, cancel := context.WithCancel(context.Background())

	d := &UnitDB{
		logger:   logger,
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		postals:  p.Postals,
		jobs:     p.Jobs,

		unitAccess:         p.UnitAccess,
		unitAccessResolver: access.NewSubjectResolver(p.DB),
		colleagueHydrator:  p.ColleagueHydrator,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		storeLogger := logger.WithOptions(
			zap.IncreaseLevel(
				p.Cfg.Log.LevelOverrides.Get(config.LoggingComponentKVStore, p.Cfg.LogLevel),
			),
		)

		kvPing, err := d.js.CreateOrUpdateKeyValue(ctxStartup, jetstream.KeyValueConfig{
			Bucket:         "unit_ping",
			Description:    "Centrum Unit Ping Timers",
			Storage:        jetstream.MemoryStorage,
			History:        1,
			TTL:            0,
			LimitMarkerTTL: 2 * PingTTL, // Tombstones live 2× rule
		})
		if err != nil {
			return fmt.Errorf("failed to create memory kv buckets. %w", err)
		}
		d.KVPing = kvPing

		jobSt, err := store.New[common.IDMapping, *common.IDMapping](
			ctxCancel,
			storeLogger,
			p.JS,
			"centrum_units",
			store.WithKVPrefix[common.IDMapping, *common.IDMapping]("job"),
			store.WithLocks[common.IDMapping, *common.IDMapping](nil),
		)
		if err != nil {
			return err
		}

		if err := jobSt.Start(ctxCancel, false); err != nil {
			return err
		}
		d.jobMapping = jobSt
		ensureJobMapping := func(ctx context.Context, job string, unitID int64, source string, refresh bool) error {
			key := centrumutils.JobIdKey(job, unitID)
			mapping := &common.IDMapping{Id: unitID}
			var changed bool
			var err error
			if refresh {
				// The job mapping is the job-scoped stream trigger. Rewrite it for
				// every local unit projection update, even when its ID is unchanged.
				err = jobSt.Put(ctx, key, mapping)
				changed = err == nil
			} else {
				// Remote watchers only repair a missing mapping. Rewriting it from
				// every server would duplicate the stream trigger.
				changed, err = jobSt.PutIfChanged(ctx, key, mapping)
			}
			logger.Debug(
				"ensured unit job mapping",
				zap.String("source", source),
				zap.String("key", key),
				zap.Int64("unit_id", unitID),
				zap.Bool("changed", changed),
				zap.Error(err),
			)
			return err
		}

		st, err := store.New[centrumunits.Unit, *centrumunits.Unit](
			ctxCancel,
			storeLogger,
			p.JS,
			"centrum_units",
			store.WithKVPrefix[centrumunits.Unit, *centrumunits.Unit]("id"),
			store.WithOnUpdateFn[centrumunits.Unit, *centrumunits.Unit](
				func(ctx context.Context, old *centrumunits.Unit, unit *centrumunits.Unit) (*centrumunits.Unit, error) {
					if unit == nil {
						return nil, nil
					}

					if err := ensureJobMapping(
						ctx,
						unit.GetJob(),
						unit.GetId(),
						"local_update",
						true,
					); err != nil {
						return nil, fmt.Errorf(
							"failed to update job %s mapping for unit %d. %w",
							unit.GetJob(),
							unit.GetId(),
							err,
						)
					}

					if old != nil && old.GetJob() != "" && old.GetJob() != unit.GetJob() {
						if err := jobSt.Delete(
							ctx,
							centrumutils.JobIdKey(old.GetJob(), unit.GetId()),
						); err != nil {
							return nil, fmt.Errorf(
								"failed to delete stale job %s mapping for unit %d. %w",
								old.GetJob(),
								unit.GetId(),
								err,
							)
						}
					}

					// Reset unit ping timer
					if err := d.UpsertWithTTL(
						ctx,
						d.KVPing,
						fmt.Sprintf("ping.%d", unit.GetId()),
						PingTTL,
					); err != nil {
						return nil, fmt.Errorf("failed to upsert ping unit timer. %w", err)
					}

					return unit, nil
				},
			),
			store.WithOnRemoteUpdatedFn[centrumunits.Unit, *centrumunits.Unit](
				func(ctx context.Context, old *centrumunits.Unit, unit *centrumunits.Unit) (*centrumunits.Unit, error) {
					if unit == nil {
						return nil, nil
					}

					if err := ensureJobMapping(
						ctx,
						unit.GetJob(),
						unit.GetId(),
						"remote_update",
						false,
					); err != nil {
						return nil, fmt.Errorf(
							"failed to update job %s mapping for unit %d. %w",
							unit.GetJob(),
							unit.GetId(),
							err,
						)
					}

					if old != nil && old.GetJob() != "" && old.GetJob() != unit.GetJob() {
						if err := jobSt.Delete(
							ctx,
							centrumutils.JobIdKey(old.GetJob(), unit.GetId()),
						); err != nil {
							return nil, fmt.Errorf(
								"failed to delete stale job %s mapping for unit %d. %w",
								old.GetJob(),
								unit.GetId(),
								err,
							)
						}
					}

					return unit, nil
				},
			),
			store.WithOnDeleteFn(
				func(ctx context.Context, _ string, unit *centrumunits.Unit) error {
					if unit == nil {
						return nil
					}

					if err := jobSt.Delete(
						ctx,
						centrumutils.JobIdKey(unit.GetJob(), unit.GetId()),
					); err != nil {
						return fmt.Errorf(
							"failed to delete job %s mapping for unit %d. %w",
							unit.GetJob(),
							unit.GetId(),
							err,
						)
					}

					if err := d.KVPing.Delete(
						ctx,
						fmt.Sprintf("ping.%d", unit.GetId()),
					); err != nil {
						d.logger.Error(
							"failed to delete ping timer for unit",
							zap.Int64("unit_id", unit.GetId()),
							zap.Error(err),
						)
					}

					return nil
				},
			),
			store.WithOnRemoteDeletedFn(
				func(ctx context.Context, key string, unit *centrumunits.Unit) error {
					if unit != nil {
						return jobSt.Delete(
							ctx,
							centrumutils.JobIdKey(unit.GetJob(), unit.GetId()),
						)
					}

					idKey, err := centrumutils.ExtractIDString(key)
					if err != nil {
						return fmt.Errorf("failed to parse unit id from key %s. %w", key, err)
					}

					var errs error
					for _, jobKey := range jobSt.KeysFiltered("", func(jobKey string) bool {
						uid, err := centrumutils.ExtractIDString(jobKey)
						return err == nil && uid == idKey
					}) {
						if err := jobSt.Delete(ctx, jobKey); err != nil {
							errs = fmt.Errorf(
								"failed to delete unit job mapping %s. %w",
								jobKey,
								err,
							)
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

func (s *UnitDB) LoadFromDB(ctx context.Context, id int64) error {
	units, err := s.loadUnitsFromDB(ctx, id)
	if err != nil {
		return err
	}

	if id > 0 && len(units) == 0 {
		return s.syncMissingUnitMembership(ctx, id)
	}

	for i := range units {
		if err := s.syncLoadedUnitMembership(ctx, units[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *UnitDB) loadUnitsFromDB(ctx context.Context, id int64) ([]*centrumunits.Unit, error) {
	tUnits := table.FivenetCentrumUnits.AS("unit")
	tUnitUser := table.FivenetCentrumUnitsUsers.AS("unit_assignment")

	condition := tUnits.DeletedAt.IS_NULL()

	if id > 0 {
		condition = condition.AND(
			tUnits.ID.EQ(mysql.Int64(id)),
		)
	}

	stmt := tUnits.
		SELECT(
			tUnits.ID,
			tUnits.CreatedAt,
			tUnits.UpdatedAt,
			tUnits.Job,
			tUnits.SortOrder,
			tUnits.Name,
			tUnits.Initials,
			tUnits.Color,
			tUnits.Icon,
			tUnits.Description,
			tUnits.Attributes,
			tUnits.HomePostal,
			tUnitUser.UnitID,
			tUnitUser.UserID,
		).
		FROM(
			tUnits.
				LEFT_JOIN(tUnitUser,
					tUnitUser.UnitID.EQ(tUnits.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tUnits.Job.ASC(),
			tUnits.SortOrder.ASC(),
			tUnits.Name.ASC(),
			tUnits.ID.ASC(),
		)

	units := []*centrumunits.Unit{}
	if err := stmt.QueryContext(ctx, s.db, &units); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	for i := range units {
		access, err := s.ListAccess(ctx, units[i].GetId())
		if err != nil {
			return nil, err
		}
		units[i].Access = access

		status, err := s.GetLastStatus(ctx, s.db, units[i].GetId())
		if err != nil {
			return nil, err
		}
		units[i].Status = status

		if len(units[i].Users) > 0 {
			targets := make([]colleagueshydrator.Target, 0, len(units[i].Users))
			for j := range units[i].Users {
				if units[i].Users[j] == nil {
					continue
				}

				targets = append(targets, colleagueshydrator.Target{
					UserID: units[i].Users[j].GetUserId(),
					Set:    units[i].Users[j].SetUser,
				})
			}

			if err := s.colleagueHydrator.HydrateTargets(
				ctx,
				s.db,
				nil,
				targets,
				colleagueshydrator.ResolveOpts{
					Scope: colleagueshydrator.JobScope{
						Mode: colleagueshydrator.JobScopePrimary,
						Job:  units[i].GetJob(),
					},
				},
			); err != nil {
				return nil, err
			}
		}

		s.enricher.EnrichJobName(units[i])
	}

	return units, nil
}

func (s *UnitDB) refreshUnitCacheFromDB(ctx context.Context, id int64) error {
	units, err := s.loadUnitsFromDB(ctx, id)
	if err != nil {
		return err
	}
	if len(units) == 0 {
		if err := s.deleteInKV(ctx, id); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			return err
		}
		return nil
	}

	return s.updateInKV(ctx, id, units[0])
}

func (s *UnitDB) nextSortOrder(ctx context.Context, q qrm.Queryable, job string) (int32, error) {
	tUnits := table.FivenetCentrumUnits.AS("unit")

	stmt := tUnits.
		SELECT(
			mysql.COALESCE(mysql.MAX(tUnits.SortOrder), mysql.Int32(-1)).AS("sort_order"),
		).
		FROM(tUnits).
		WHERE(mysql.AND(
			tUnits.Job.EQ(mysql.String(job)),
			tUnits.DeletedAt.IS_NULL(),
		))

	var dest sortOrderResult
	if err := stmt.QueryContext(ctx, q, &dest); err != nil {
		return 0, err
	}

	if dest.SortOrder == math.MaxInt32 {
		return 0, errors.New("unit sort order overflow")
	}

	return dest.SortOrder + 1, nil
}

func sameUnitStatusContent(current, next *centrumunits.UnitStatus) bool {
	return current.GetStatus() == next.GetStatus() &&
		current.GetReason() == next.GetReason() &&
		current.GetCode() == next.GetCode() &&
		current.GetUserId() == next.GetUserId() &&
		current.GetCreatorId() == next.GetCreatorId() &&
		current.GetCreatorJob() == next.GetCreatorJob()
}

func (s *UnitDB) CreateUnit(
	ctx context.Context,
	creatorJob string,
	creatorGrade int32,
	unit *centrumunits.Unit,
) (*centrumunits.Unit, error) {
	if unit.GetAccess() == nil {
		unit.Access = &unitsaccess.UnitAccess{}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tUnits := table.FivenetCentrumUnits

	sortOrder, err := s.nextSortOrder(ctx, tx, creatorJob)
	if err != nil {
		return nil, err
	}
	unit.SortOrder = sortOrder

	stmt := tUnits.
		INSERT(
			tUnits.Job,
			tUnits.SortOrder,
			tUnits.Name,
			tUnits.Initials,
			tUnits.Color,
			tUnits.Icon,
			tUnits.Description,
			tUnits.Attributes,
			tUnits.HomePostal,
		).
		VALUES(
			creatorJob,
			mysql.Int32(sortOrder),
			unit.GetName(),
			unit.GetInitials(),
			unit.GetColor(),
			unit.Icon,
			unit.Description,
			unit.GetAttributes(),
			unit.HomePostal,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	unit.SetJob(creatorJob)
	unit.SetId(lastId)

	// A new unit should have a status, so we make sure we add one
	if unit.Status, err = s.AddStatus(ctx, tx, &centrumunits.UnitStatus{
		CreatedAt:  timestamp.Now(),
		UnitId:     unit.GetId(),
		Status:     centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE,
		CreatorJob: &creatorJob,
	}); err != nil {
		return nil, err
	}

	highestGrade := creatorGrade
	if grade, ok := s.enricher.GetHighestJobGrade(creatorJob); ok {
		highestGrade = grade
	}
	unit.Access = access.EnsureJobAccessEntries(
		unit.GetAccess(),
		&resourcesaccess.JobAccess{
			Job:          creatorJob,
			MinimumGrade: creatorGrade,
			Access:       int32(unitsaccess.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN),
		},
		&resourcesaccess.JobAccess{
			Job:          creatorJob,
			MinimumGrade: highestGrade,
			Access:       int32(unitsaccess.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN),
			Required:     new(true),
		},
	)
	normalizedAccess, err := access.NormalizeAccess(unit.GetAccess(), nil, &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{{
			Job:          creatorJob,
			MinimumGrade: highestGrade,
			Access:       int32(unitsaccess.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN),
		}},
	}, 15)
	if err != nil {
		return nil, err
	}

	if _, err := s.unitAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.unitAccessResolver,
		unit.GetId(),
		normalizedAccess,
		unitSubjectAccessOptions,
	); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Load new/updated unit from database
	if err := s.LoadFromDB(ctx, unit.GetId()); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *UnitDB) Update(
	ctx context.Context,
	userGrade int32,
	unit *centrumunits.Unit,
) (*centrumunits.Unit, error) {
	if unit.GetAccess() == nil {
		unit.Access = &unitsaccess.UnitAccess{}
	}

	tUnits := table.FivenetCentrumUnits

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	stmt := tUnits.
		UPDATE(
			tUnits.Name,
			tUnits.Initials,
			tUnits.Color,
			tUnits.Icon,
			tUnits.Description,
			tUnits.Attributes,
			tUnits.HomePostal,
		).
		SET(
			unit.GetName(),
			unit.GetInitials(),
			unit.GetColor(),
			unit.GetIcon(),
			unit.GetDescription(),
			unit.GetAttributes(),
			unit.GetHomePostal(),
		).
		WHERE(mysql.AND(
			tUnits.ID.EQ(mysql.Int64(unit.GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	highestGrade := userGrade
	if grade, ok := s.enricher.GetHighestJobGrade(unit.GetJob()); ok {
		highestGrade = grade
	}
	normalizedAccess, err := access.NormalizeAccess(unit.GetAccess(), nil, &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{{
			Job:          unit.GetJob(),
			MinimumGrade: highestGrade,
			Access:       int32(unitsaccess.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN),
		}},
	}, 15)
	if err != nil {
		return nil, err
	}

	if _, err := s.unitAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.unitAccessResolver,
		unit.GetId(),
		normalizedAccess,
		unitSubjectAccessOptions,
	); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Load new/updated unit from database
	if err := s.LoadFromDB(ctx, unit.GetId()); err != nil {
		return nil, err
	}

	if err := s.updateInKV(ctx, unit.GetId(), unit); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *UnitDB) Delete(ctx context.Context, id int64) error {
	tUnits := table.FivenetCentrumUnits

	stmt := tUnits.
		DELETE().
		WHERE(mysql.AND(
			tUnits.ID.EQ(mysql.Int64(id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return s.syncMissingUnitMembership(ctx, id)
}

func (s *UnitDB) ListAccess(ctx context.Context, id int64) (*unitsaccess.UnitAccess, error) {
	access, err := s.unitAccess.ListTargetAccess(ctx, s.db, id, unitSubjectAccessOptions)
	if err != nil {
		return nil, err
	}

	for i := range access.GetJobs() {
		s.enricher.EnrichJobInfo(access.GetJobs()[i])
	}

	return access, nil
}

func (s *UnitDB) GetAccess() *access.CentrumUnitsObjectAccess {
	return s.unitAccess
}
