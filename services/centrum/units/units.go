package units

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/users"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	PingTTL  = 8 * time.Second  // OFFLINE
	EmptyTTL = 30 * time.Second // No crew → delete
)

type UnitDB struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	tracker  tracker.ITracker
	postals  postals.Postals

	store      *store.Store[centrum.Unit, *centrum.Unit]
	jobMapping *store.Store[common.IDMapping, *common.IDMapping]

	unitAccess *access.Grouped[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitUserAccess, *centrum.UnitUserAccess, centrum.UnitQualificationAccess, *centrum.UnitQualificationAccess, centrum.UnitAccessLevel]

	KVPing jetstream.KeyValue
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	JS       *events.JSWrapper
	DB       *sql.DB
	Cfg      *config.Config
	Enricher *mstlystcdata.Enricher
	Tracker  tracker.ITracker
	Postals  postals.Postals
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

		unitAccess: access.NewGrouped[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitUserAccess](
			p.DB,
			table.FivenetCentrumUnits,
			&access.TargetTableColumns{
				ID:        table.FivenetCentrumUnits.ID,
				DeletedAt: table.FivenetCentrumUnits.DeletedAt,
			},
			access.NewJobs[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitAccessLevel](
				table.FivenetCentrumUnitsAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCentrumUnitsAccess.ID,
						TargetID: table.FivenetCentrumUnitsAccess.TargetID,
						Access:   table.FivenetCentrumUnitsAccess.Access,
					},
					Job:          table.FivenetCentrumUnitsAccess.Job,
					MinimumGrade: table.FivenetCentrumUnitsAccess.MinimumGrade,
				},
				table.FivenetCentrumUnitsAccess.AS("unit_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCentrumUnitsAccess.AS("unit_job_access").ID,
						TargetID: table.FivenetCentrumUnitsAccess.AS("unit_job_access").TargetID,
						Access:   table.FivenetCentrumUnitsAccess.AS("unit_job_access").Access,
					},
					Job: table.FivenetCentrumUnitsAccess.AS("unit_job_access").Job,
					MinimumGrade: table.FivenetCentrumUnitsAccess.AS(
						"unit_job_access",
					).MinimumGrade,
				},
			),
			nil,
			access.NewQualifications[centrum.UnitQualificationAccess, *centrum.UnitQualificationAccess, centrum.UnitAccessLevel](
				table.FivenetCentrumUnitsAccess,
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCentrumUnitsAccess.ID,
						TargetID: table.FivenetCentrumUnitsAccess.TargetID,
						Access:   table.FivenetCentrumUnitsAccess.Access,
					},
					QualificationId: table.FivenetCentrumUnitsAccess.QualificationID,
				},
				table.FivenetCentrumUnitsAccess.AS("unit_qualification_access"),
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID: table.FivenetCentrumUnitsAccess.AS(
							"unit_qualification_access",
						).ID,
						TargetID: table.FivenetCentrumUnitsAccess.AS(
							"unit_qualification_access",
						).TargetID,
						Access: table.FivenetCentrumUnitsAccess.AS(
							"unit_qualification_access",
						).Access,
					},
					QualificationId: table.FivenetCentrumUnitsAccess.AS(
						"unit_qualification_access",
					).QualificationID,
				},
			),
		),
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

		st, err := store.New[centrum.Unit, *centrum.Unit](
			ctxCancel,
			storeLogger,
			p.JS,
			"centrum_units",
			store.WithKVPrefix[centrum.Unit, *centrum.Unit]("id"),
			store.WithOnUpdateFn[centrum.Unit, *centrum.Unit](
				func(ctx context.Context, _ *centrum.Unit, unit *centrum.Unit) (*centrum.Unit, error) {
					if unit == nil {
						return nil, nil
					}

					if err := jobSt.Put(ctx, centrumutils.JobIdKey(unit.GetJob(), unit.GetId()), &common.IDMapping{
						Id: unit.GetId(),
					}); err != nil {
						return nil, fmt.Errorf(
							"failed to update job %s mapping for unit %d. %w",
							unit.GetJob(),
							unit.GetId(),
							err,
						)
					}

					// Reset unit ping timer
					if err := d.UpsertWithTTL(ctx, d.KVPing, fmt.Sprintf("ping.%d", unit.GetId()), PingTTL); err != nil {
						return nil, fmt.Errorf("failed to upsert ping unit timer. %w", err)
					}

					return unit, nil
				},
			),
			store.WithOnDeleteFn(func(ctx context.Context, _ string, unit *centrum.Unit) error {
				if unit == nil {
					return nil
				}

				if err := jobSt.Delete(ctx, centrumutils.JobIdKey(unit.GetJob(), unit.GetId())); err != nil {
					return fmt.Errorf(
						"failed to delete job %s mapping for unit %d. %w",
						unit.GetJob(),
						unit.GetId(),
						err,
					)
				}

				if err := d.KVPing.Delete(ctx, fmt.Sprintf("ping.%d", unit.GetId())); err != nil {
					d.logger.Error(
						"failed to delete ping timer for unit",
						zap.Int64("unit_id", unit.GetId()),
						zap.Error(err),
					)
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

func (s *UnitDB) LoadFromDB(ctx context.Context, id int64) error {
	tUnits := table.FivenetCentrumUnits.AS("unit")
	tUnitUser := table.FivenetCentrumUnitsUsers.AS("unit_assignment")

	condition := tUnits.DeletedAt.IS_NULL()

	if id > 0 {
		condition = condition.AND(
			tUnits.ID.EQ(jet.Int64(id)),
		)
	}

	stmt := tUnits.
		SELECT(
			tUnits.ID,
			tUnits.CreatedAt,
			tUnits.UpdatedAt,
			tUnits.Job,
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
			tUnits.Name.ASC(),
		)

	units := []*centrum.Unit{}
	if err := stmt.QueryContext(ctx, s.db, &units); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for i := range units {
		access, err := s.ListAccess(ctx, units[i].GetId())
		if err != nil {
			return err
		}
		units[i].Access = access

		status, err := s.GetLastStatus(ctx, s.db, units[i].GetId())
		if err != nil {
			return err
		}
		units[i].Status = status

		if err := users.RetrieveUsersForUnit(ctx, s.db, s.enricher, &units[i].Users); err != nil {
			return err
		}

		s.enricher.EnrichJobName(units[i])

		if err := s.updateInKV(ctx, units[i].GetId(), units[i]); err != nil {
			return err
		}

		for _, user := range units[i].GetUsers() {
			if err := s.tracker.SetUserMappingForUser(ctx, user.GetUserId(), &units[i].Id); err != nil {
				s.logger.Error("failed to set user's unit id", zap.Error(err))
			}
		}
	}

	return nil
}

func (s *UnitDB) LoadUnitIDForUserID(ctx context.Context, userId int32) (int64, error) {
	tUnitUser := table.FivenetCentrumUnitsUsers.AS("unit_assignment")

	stmt := tUnitUser.
		SELECT(
			tUnitUser.UnitID.AS("unit_id"),
		).
		FROM(tUnitUser).
		WHERE(
			tUnitUser.UserID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	var dest struct {
		UnitID int64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}

		return 0, nil
	}

	return dest.UnitID, nil
}

func (s *UnitDB) UpdateStatus(
	ctx context.Context,
	unitId int64,
	in *centrum.UnitStatus,
) (*centrum.UnitStatus, error) {
	unit, err := s.Get(ctx, unitId)
	if err != nil {
		return nil, err
	}

	// If the unit status is the same and is a status that shouldn't be duplicated, don't update the status again
	if unit.GetStatus() != nil &&
		unit.GetStatus().GetStatus() == in.GetStatus() &&
		(in.GetStatus() == centrum.StatusUnit_STATUS_UNIT_ON_BREAK ||
			in.GetStatus() == centrum.StatusUnit_STATUS_UNIT_BUSY ||
			in.GetStatus() == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE ||
			in.GetStatus() == centrum.StatusUnit_STATUS_UNIT_AVAILABLE) &&
		// Additionally if the status is under 2 minutes disallow the same status update
		(unit.GetStatus().GetCreatedAt() == nil || time.Since(unit.GetStatus().GetCreatedAt().AsTime()) < 2*time.Minute) {
		s.logger.Debug(
			"skipping unit status update due to same status or time",
			zap.Int64("unit_id", unitId),
			zap.String("status", in.GetStatus().String()),
		)
		return nil, nil
	}

	if unit.GetAttributes() != nil &&
		unit.GetAttributes().Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
		// Only allow a static unit to be set busy, on break or unavailable
		if in.GetStatus() != centrum.StatusUnit_STATUS_UNIT_BUSY &&
			in.GetStatus() != centrum.StatusUnit_STATUS_UNIT_ON_BREAK &&
			in.GetStatus() != centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE {
			return nil, nil
		}
	}

	s.logger.Debug(
		"updating unit status",
		zap.Int64("unit_id", unitId),
		zap.String("status", in.GetStatus().String()),
	)

	if in.UserId != nil {
		var err error
		in.User, err = users.RetrieveUserShortById(ctx, s.db, s.enricher, in.GetUserId())
		if err != nil {
			return nil, err
		}

		if um, ok := s.tracker.GetUserMarkerById(in.GetUserId()); ok {
			in.X = &um.X
			in.Y = &um.Y
			in.Postal = um.Postal
		}
	}
	if in.CreatorId != nil {
		// If the creator of the status is the same as the user, no need to query the db
		if in.UserId != nil && in.GetCreatorId() == in.GetUserId() {
			in.Creator = in.GetUser()
		} else {
			var err error
			in.Creator, err = users.RetrieveUserShortById(ctx, s.db, s.enricher, in.GetCreatorId())
			if err != nil {
				return nil, err
			}
		}
	}

	if in.GetCreatedAt() == nil {
		in.CreatedAt = timestamp.Now()
	}

	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		INSERT(
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorID,
			tUnitStatus.CreatorJob,
		).
		VALUES(
			in.GetCreatedAt(),
			in.GetUnitId(),
			in.GetStatus(),
			in.Reason,
			in.Code,
			in.UserId,
			in.X,
			in.Y,
			in.Postal,
			in.CreatorId,
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
	in.Id = lastId

	if err := s.updateStatusInKV(ctx, in.GetUnitId(), in); err != nil {
		return nil, err
	}

	return in, nil
}

func (s *UnitDB) UpdateUnitAssignments(
	ctx context.Context,
	_ string,
	userId *int32,
	unitId int64,
	toAdd []int32,
	toRemove []int32,
) error {
	s.logger.Debug(
		"updating unit assignments",
		zap.Int64("unit_id", unitId),
		zap.Int32s("toAdd", toAdd),
		zap.Int32s("toRemove", toRemove),
	)

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

	tUnitUser := table.FivenetCentrumUnitsUsers

	toAddUsers := []*jobs.Colleague{}

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
			removeIds[i] = jet.Int32(toRemove[i])
		}

		stmt := tUnitUser.
			DELETE().
			WHERE(jet.AND(
				tUnitUser.UnitID.EQ(jet.Int64(unitId)),
				tUnitUser.UserID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		addIds := []int32{}
		for i := range toAdd {
			if um, ok := s.tracker.GetUserMarkerById(toAdd[i]); !ok || um.GetHidden() {
				continue
			}

			unit, err := s.Get(ctx, unitId)
			if err != nil {
				return err
			}
			// Skip already added units
			if slices.ContainsFunc(unit.GetUsers(), func(in *centrum.UnitAssignment) bool {
				return in.GetUserId() == toAdd[i]
			}) {
				continue
			}

			addIds = append(addIds, toAdd[i])
		}

		if len(addIds) > 0 {
			stmt := tUnitUser.
				INSERT(
					tUnitUser.UnitID,
					tUnitUser.UserID,
				)

			for _, id := range addIds {
				stmt = stmt.
					VALUES(
						unitId,
						id,
					)
			}

			stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
				tUnitUser.UnitID.SET(jet.IntExp(jet.Raw("VALUES(`unit_id`)"))),
			)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}

			var err error
			toAddUsers, err = users.RetrieveColleagueById(ctx, s.db, s.enricher, addIds...)
			if err != nil {
				return err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	key := centrumutils.IdKey(unitId)
	if err := s.store.ComputeUpdate(ctx, key, func(key string, unit *centrum.Unit) (*centrum.Unit, bool, error) {
		if len(toRemove) > 0 {
			toAnnounce := []int32{}

			if len(unit.GetUsers()) == 0 {
				// No users in unit? Make sure to announce all users that should be removed just in case
				toAnnounce = append(toAnnounce, toRemove...)
			} else {
				unit.Users = slices.DeleteFunc(unit.GetUsers(), func(in *centrum.UnitAssignment) bool {
					for k := range toRemove {
						if in.GetUserId() != toRemove[k] {
							continue
						}

						toAnnounce = append(toAnnounce, toRemove[k])
						return true
					}

					return false
				})
			}

			// Send updates
			for _, user := range toAnnounce {
				if _, err := s.AddStatus(ctx, s.db, &centrum.UnitStatus{
					CreatedAt: timestamp.Now(),
					UnitId:    unit.GetId(),
					Status:    centrum.StatusUnit_STATUS_UNIT_USER_REMOVED,
					UserId:    &user,
					CreatorId: userId,
					X:         x,
					Y:         y,
					Postal:    postal,
				}, true, unit.GetJob()); err != nil {
					return nil, false, err
				}

				if err := s.tracker.UnsetUnitIDForUser(ctx, user); err != nil {
					return nil, false, err
				}
			}
		}

		if len(toAddUsers) > 0 {
			for _, user := range toAddUsers {
				unit.Users = append(unit.Users, &centrum.UnitAssignment{
					UnitId: unit.GetId(),
					UserId: user.GetUserId(),
					User:   user,
				})

				if _, err := s.AddStatus(ctx, s.db, &centrum.UnitStatus{
					CreatedAt:  timestamp.Now(),
					UnitId:     unit.GetId(),
					Status:     centrum.StatusUnit_STATUS_UNIT_USER_ADDED,
					UserId:     &user.UserId,
					CreatorId:  userId,
					X:          x,
					Y:          y,
					Postal:     postal,
					CreatorJob: &user.Job,
				}, true, unit.GetJob()); err != nil {
					return nil, false, err
				}

				if err := s.tracker.SetUserMappingForUser(ctx, user.GetUserId(), &unit.Id); err != nil {
					return nil, false, err
				}
			}
		}

		// Unit is empty now, set unit status to be unavailable automatically
		if len(unit.GetUsers()) == 0 {
			if unit.Status, err = s.AddStatus(ctx, s.db, &centrum.UnitStatus{
				CreatedAt: timestamp.Now(),
				UnitId:    unit.GetId(),
				Status:    centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:    userId,
				CreatorId: userId,
				X:         x,
				Y:         y,
				Postal:    postal,
			}, true, unit.GetJob()); err != nil {
				return nil, false, err
			}
		}

		return unit, true, nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *UnitDB) CreateUnit(
	ctx context.Context,
	creatorJob string,
	unit *centrum.Unit,
) (*centrum.Unit, error) {
	if unit.GetAccess() == nil {
		unit.Access = &centrum.UnitAccess{}
	}
	unit.GetAccess().ClearQualificationResults()

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tUnits := table.FivenetCentrumUnits
	stmt := tUnits.
		INSERT(
			tUnits.Job,
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
	unit.Job = creatorJob
	unit.Id = lastId

	// A new unit shouldn't have a status, so we make sure we add one
	if unit.Status, err = s.AddStatus(ctx, tx, &centrum.UnitStatus{
		CreatedAt: timestamp.Now(),
		UnitId:    unit.GetId(),
		Status:    centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
	}, false, ""); err != nil {
		return nil, err
	}

	if _, err := s.unitAccess.HandleAccessChanges(ctx, tx, unit.GetId(), unit.GetAccess().GetJobs(), nil, unit.GetAccess().GetQualifications()); err != nil {
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

func (s *UnitDB) Update(ctx context.Context, unit *centrum.Unit) (*centrum.Unit, error) {
	if unit.GetAccess() == nil {
		unit.Access = &centrum.UnitAccess{}
	}
	unit.GetAccess().ClearQualificationResults()

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
		WHERE(jet.AND(
			tUnits.ID.EQ(jet.Int64(unit.GetId())),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	if _, err := s.unitAccess.HandleAccessChanges(ctx, tx, unit.GetId(), unit.GetAccess().GetJobs(), nil, unit.GetAccess().GetQualifications()); err != nil {
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

func (s *UnitDB) AddStatus(
	ctx context.Context,
	tx qrm.DB,
	status *centrum.UnitStatus,
	publish bool,
	job string,
) (*centrum.UnitStatus, error) {
	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		INSERT(
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorID,
			tUnitStatus.CreatorJob,
		).
		VALUES(
			jet.CURRENT_TIMESTAMP(),
			status.GetUnitId(),
			status.GetStatus(),
			status.Reason,
			status.Code,
			status.UserId,
			status.X,
			status.Y,
			status.Postal,
			status.CreatorId,
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

	newStatus, err := s.GetStatusByID(ctx, tx, lastId)
	if err != nil {
		return nil, err
	}

	if publish {
		data, err := proto.Marshal(status)
		if err != nil {
			return nil, err
		}

		if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, job), data); err != nil {
			return nil, err
		}
	}

	return newStatus, nil
}

func (s *UnitDB) GetStatusByID(
	ctx context.Context,
	tx qrm.DB,
	id int64,
) (*centrum.UnitStatus, error) {
	tUnitStatus := table.FivenetCentrumUnitsStatus.AS("unit_status")
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")
	tUsers := tables.User().AS("colleague")
	tUserProps := table.FivenetUserProps.AS("user_props")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.CreatorID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorJob,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tUsers.ID).
						AND(tColleagueProps.Job.EQ(tUsers.Job)),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tUnitStatus.ID.EQ(jet.Int64(id)),
		).
		ORDER_BY(tUnitStatus.ID.DESC()).
		LIMIT(1)

	var dest centrum.UnitStatus
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return &dest, nil
}

func (s *UnitDB) GetLastStatus(
	ctx context.Context,
	tx qrm.DB,
	unitId int64,
) (*centrum.UnitStatus, error) {
	tUnitStatus := table.FivenetCentrumUnitsStatus.AS("unit_status")
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")
	tUsers := tables.User().AS("colleague")
	tUserProps := table.FivenetUserProps.AS("user_props")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.CreatorID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorJob,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tUsers.ID).
						AND(tColleagueProps.Job.EQ(tUsers.Job)),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(jet.AND(
			tUnitStatus.UnitID.EQ(jet.Int64(unitId)),
			tUnitStatus.Status.NOT_IN(
				jet.Int32(int32(centrum.StatusUnit_STATUS_UNIT_USER_ADDED)),
				jet.Int32(int32(centrum.StatusUnit_STATUS_UNIT_USER_REMOVED)),
			),
		)).
		ORDER_BY(tUnitStatus.ID.DESC()).
		LIMIT(1)

	var dest centrum.UnitStatus
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return &dest, nil
}

func (s *UnitDB) Delete(ctx context.Context, id int64) error {
	tUnits := table.FivenetCentrumUnits

	stmt := tUnits.
		DELETE().
		WHERE(jet.AND(
			tUnits.ID.EQ(jet.Int64(id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	if err := s.deleteInKV(ctx, id); err != nil {
		return fmt.Errorf("failed to delete unit in KV store. %w", err)
	}

	return nil
}

func (s *UnitDB) ListAccess(ctx context.Context, id int64) (*centrum.UnitAccess, error) {
	access := &centrum.UnitAccess{}

	jobsAccess, err := s.unitAccess.Jobs.List(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	access.Jobs = jobsAccess

	for i := range access.GetJobs() {
		s.enricher.EnrichJobInfo(access.GetJobs()[i])
	}

	qualificationsccess, err := s.unitAccess.Qualifications.List(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	access.Qualifications = qualificationsccess

	return access, nil
}

func (s *UnitDB) GetAccess() *access.Grouped[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitUserAccess, *centrum.UnitUserAccess, centrum.UnitQualificationAccess, *centrum.UnitQualificationAccess, centrum.UnitAccessLevel] {
	return s.unitAccess
}
