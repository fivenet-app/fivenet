package centrummanager

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumstate"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) UpdateUnitStatus(ctx context.Context, unitJob string, unitId uint64, in *centrum.UnitStatus) (*centrum.UnitStatus, error) {
	unit, err := s.GetUnit(ctx, unitJob, unitId)
	if err != nil {
		return nil, err
	}

	// If the unit status is the same and is a status that shouldn't be duplicated, don't update the status again
	if unit.Status != nil &&
		unit.Status.Status == in.Status &&
		(in.Status == centrum.StatusUnit_STATUS_UNIT_ON_BREAK ||
			in.Status == centrum.StatusUnit_STATUS_UNIT_BUSY ||
			in.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE ||
			in.Status == centrum.StatusUnit_STATUS_UNIT_AVAILABLE) &&
		// Additionally if the status is under 2 minutes disallow the same status update
		(unit.Status.CreatedAt == nil || time.Since(unit.Status.CreatedAt.AsTime()) < 2*time.Minute) {
		s.logger.Debug("skipping unit status update due to same status or time", zap.Uint64("unit_id", unitId), zap.String("status", in.Status.String()))
		return nil, nil
	}

	if unit.Attributes != nil && unit.Attributes.Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
		// Only allow a static unit to be set busy, on break or unavailable
		if in.Status != centrum.StatusUnit_STATUS_UNIT_BUSY &&
			in.Status != centrum.StatusUnit_STATUS_UNIT_ON_BREAK &&
			in.Status != centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE {
			return nil, nil
		}
	}

	s.logger.Debug("updating unit status", zap.Uint64("unit_id", unitId), zap.String("status", in.Status.String()))

	if in.UserId != nil {
		var err error
		in.User, err = s.retrieveUserShortById(ctx, *in.UserId)
		if err != nil {
			return nil, err
		}

		if um, ok := s.tracker.GetUserById(*in.UserId); ok {
			in.X = &um.X
			in.Y = &um.Y
			in.Postal = um.Postal
		}
	}
	if in.CreatorId != nil {
		// If the creator of the status is the same as the user, no need to query the db
		if in.UserId != nil && *in.CreatorId == *in.UserId {
			in.Creator = in.User
		} else {
			var err error
			in.Creator, err = s.retrieveUserShortById(ctx, *in.CreatorId)
			if err != nil {
				return nil, err
			}
		}
	}

	status, err := s.AddUnitStatus(ctx, s.db, unitJob, in, false)
	if err != nil {
		return nil, err
	}
	in.Id = uint64(status.Id)

	if err := s.State.UpdateUnitStatus(ctx, unitJob, in.UnitId, in); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, unitJob), data); err != nil {
		return nil, fmt.Errorf("failed to publish unit status event (size: %d, message: '%+v'): %w", len(data), in, err)
	}

	return in, nil
}

func (s *Manager) UpdateUnitAssignments(ctx context.Context, userJob *string, userId *int32, unitJob string, unitId uint64, toAdd []int32, toRemove []int32, allowedJobs []string) error {
	s.logger.Debug("updating unit assignments", zap.String("job", unitJob), zap.Uint64("unit_id", unitId), zap.Int32s("toAdd", toAdd), zap.Int32s("toRemove", toRemove))

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return nil
	}

	var x, y *float64
	var postal *string
	if userId != nil {
		if um, ok := s.tracker.GetUserById(*userId); ok {
			x = &um.X
			y = &um.Y
			postal = um.Postal
		}
	}

	tUnitUser := table.FivenetCentrumUnitsUsers

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	usersToRemove := []*jobs.Colleague{}
	if len(toRemove) > 0 {
		users, err := s.retrieveColleagueById(ctx, toRemove...)
		if err != nil {
			return err
		}
		usersToRemove = append(usersToRemove, users...)

		removeIds := make([]jet.Expression, len(users))
		for i := range users {
			removeIds[i] = jet.Int32(users[i].UserId)
		}

		stmt := tUnitUser.
			DELETE().
			WHERE(jet.AND(
				tUnitUser.UnitID.EQ(jet.Uint64(unitId)),
				tUnitUser.UserID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	usersToAdd := []*jobs.Colleague{}
	if len(toAdd) > 0 {
		for i := range toAdd {
			um, ok := s.tracker.GetUserById(toAdd[i])
			if !ok || um.Hidden {
				continue
			}

			unit, err := s.GetUnit(ctx, unitJob, unitId)
			if err != nil {
				return err
			}

			// Skip already added users
			if slices.ContainsFunc(unit.Users, func(in *centrum.UnitAssignment) bool {
				return in.UserId == toAdd[i]
			}) {
				continue
			}

			users, err := s.retrieveColleagueById(ctx, toAdd...)
			if err != nil {
				return err
			}

			if allowedJobs == nil {
				usersToAdd = append(usersToAdd, users...)
			} else {
				// Check if users are in a job that the user is allowed to access
				for _, user := range users {
					if slices.ContainsFunc(allowedJobs, func(in string) bool {
						return in == user.Job
					}) {
						usersToAdd = append(usersToAdd, user)
					}
				}
			}
		}

		if len(usersToAdd) > 0 {
			stmt := tUnitUser.
				INSERT(
					tUnitUser.UnitID,
					tUnitUser.UnitJob,
					tUnitUser.UserID,
					tUnitUser.UserJob,
				)

			for _, user := range usersToAdd {
				stmt = stmt.
					VALUES(
						unitId,
						unitJob,
						user.UserId,
						user.Job,
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
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	store := s.State.UnitsStore()

	key := centrumstate.JobIdKey(unitJob, unitId)
	if err := store.ComputeUpdate(ctx, key, true, func(key string, unit *centrum.Unit) (*centrum.Unit, bool, error) {
		if len(usersToRemove) > 0 {
			toAnnounce := []*jobs.Colleague{}

			unit.Users = slices.DeleteFunc(unit.Users, func(in *centrum.UnitAssignment) bool {
				for k := range usersToRemove {
					if in.UserId != usersToRemove[k].UserId {
						continue
					}

					toAnnounce = append(toAnnounce, usersToRemove[k])
					return true
				}

				return false
			})

			// Send updates
			for _, user := range toAnnounce {
				if _, err := s.AddUnitStatus(ctx, s.db, unitJob, &centrum.UnitStatus{
					CreatedAt:  timestamp.Now(),
					UnitId:     unit.Id,
					UnitJob:    unit.Job,
					Status:     centrum.StatusUnit_STATUS_UNIT_USER_REMOVED,
					UserId:     &user.UserId,
					UserJob:    &user.Job,
					CreatorId:  userId,
					CreatorJob: userJob,
					X:          x,
					Y:          y,
					Postal:     postal,
				}, true); err != nil {
					return nil, false, err
				}

				if err := s.UnsetUnitIDForUser(ctx, user.UserId); err != nil {
					return nil, false, err
				}
			}
		}

		if len(usersToAdd) > 0 {
			notFound := []*jobs.Colleague{}
			for _, user := range usersToAdd {
				if um, ok := s.tracker.GetUserById(user.UserId); !ok || um.Hidden {
					continue
				}

				// Skip already added units
				if slices.ContainsFunc(unit.Users, func(in *centrum.UnitAssignment) bool {
					return in.UserId == user.UserId
				}) {
					continue
				}

				notFound = append(notFound, user)
			}

			for _, user := range notFound {
				unit.Users = append(unit.Users, &centrum.UnitAssignment{
					UnitId:  unit.Id,
					UnitJob: unit.Job,
					UserId:  user.UserId,
					UserJob: user.Job,
					User:    user,
				})

				if _, err := s.AddUnitStatus(ctx, s.db, unitJob, &centrum.UnitStatus{
					CreatedAt:  timestamp.Now(),
					UnitId:     unit.Id,
					UnitJob:    unit.Job,
					Status:     centrum.StatusUnit_STATUS_UNIT_USER_ADDED,
					UserId:     &user.UserId,
					UserJob:    &user.Job,
					CreatorId:  userId,
					CreatorJob: userJob,
					X:          x,
					Y:          y,
					Postal:     postal,
				}, true); err != nil {
					return nil, false, err
				}

				if err := s.SetUnitForUser(ctx, user.Job, user.UserId, unit.Job, unit.Id); err != nil {
					return nil, false, err
				}
			}
		}

		// Unit is empty, set unit status to be unavailable automatically
		if len(unit.Users) == 0 {
			if unit.Status, err = s.AddUnitStatus(ctx, s.db, unitJob, &centrum.UnitStatus{
				CreatedAt:  timestamp.Now(),
				UnitId:     unit.Id,
				UnitJob:    unit.Job,
				Status:     centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:     userId,
				UserJob:    userJob,
				CreatorId:  userId,
				CreatorJob: userJob,
				X:          x,
				Y:          y,
				Postal:     postal,
			}, true); err != nil {
				return nil, false, err
			}
		}

		return unit, true, nil
	}); err != nil {
		return err
	}

	unit, err := s.GetUnit(ctx, unitJob, unitId)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitUpdated, unitJob), data); err != nil {
		return err
	}

	return nil
}

func (s *Manager) CreateUnit(ctx context.Context, job string, unit *centrum.Unit) (*centrum.Unit, error) {
	if unit.Access == nil {
		unit.Access = &centrum.UnitAccess{}
	}
	unit.Access.ClearQualificationResults()

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
			tUnits.Description,
			tUnits.Attributes,
			tUnits.HomePostal,
		).
		VALUES(
			job,
			unit.Name,
			unit.Initials,
			unit.Color,
			unit.Description,
			unit.Attributes,
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
	unit.Id = uint64(lastId)

	// A new unit doesn't have a status, so we make sure to add one
	if unit.Status, err = s.AddUnitStatus(ctx, tx, job, &centrum.UnitStatus{
		CreatedAt: timestamp.Now(),
		UnitId:    unit.Id,
		UnitJob:   unit.Job,
		Status:    centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
	}, false); err != nil {
		return nil, err
	}

	if _, err := s.GetUnitAccess().HandleAccessChanges(ctx, tx, unit.Id, unit.Access.Jobs, nil, unit.Access.Qualifications); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Load new/updated unit from database
	if err := s.LoadUnitsFromDB(ctx, unit.Id); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, err
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitCreated, job), data); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *Manager) UpdateUnit(ctx context.Context, unit *centrum.Unit) (*centrum.Unit, error) {
	description := ""
	if unit.Description != nil {
		description = *unit.Description
	}

	if unit.Access == nil {
		unit.Access = &centrum.UnitAccess{}
	}
	unit.Access.ClearQualificationResults()

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
			tUnits.Description,
			tUnits.Attributes,
			tUnits.HomePostal,
		).
		SET(
			unit.Name,
			unit.Initials,
			unit.Color,
			description,
			unit.Attributes,
			unit.HomePostal,
		).
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(unit.Job)),
			tUnits.ID.EQ(jet.Uint64(unit.Id)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	if _, err := s.GetUnitAccess().HandleAccessChanges(ctx, tx, unit.Id, unit.Access.Jobs, nil, unit.Access.Qualifications); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Load new/updated unit from database
	if err := s.LoadUnitsFromDB(ctx, unit.Id); err != nil {
		return nil, err
	}

	if err := s.State.UpdateUnit(ctx, unit.Job, unit.Id, unit); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, err
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitUpdated, unit.Job), data); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *Manager) AddUnitStatus(ctx context.Context, tx qrm.DB, unitJob string, status *centrum.UnitStatus, publish bool) (*centrum.UnitStatus, error) {
	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		INSERT(
			tUnitStatus.UnitID,
			tUnitStatus.UnitJob,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.UserJob,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorID,
			tUnitStatus.CreatorJob,
		).
		VALUES(
			status.UnitId,
			status.UnitJob,
			status.Status,
			status.Reason,
			status.Code,
			status.UserId,
			status.UserJob,
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

	newStatus, err := s.GetUnitStatusByID(ctx, tx, unitJob, uint64(lastId))
	if err != nil {
		return nil, err
	}

	if newStatus != nil && newStatus.CreatedAt == nil {
		newStatus.CreatedAt = timestamp.Now()
	}

	if publish {
		data, err := proto.Marshal(status)
		if err != nil {
			return nil, err
		}

		if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, unitJob), data); err != nil {
			return nil, err
		}
	}

	return newStatus, nil
}

func (s *Manager) GetUnitStatusByID(ctx context.Context, tx qrm.DB, unitJob string, unitId uint64) (*centrum.UnitStatus, error) {
	tUsers := tables.Users().AS("colleague")

	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.UnitJob,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.UserJob,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorID,
			tUnitStatus.CreatorJob,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.NamePrefix,
			tJobsUserProps.NameSuffix,
			tUserProps.Avatar.AS("colleague.avatar"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUsers.ID).
						AND(tJobsUserProps.Job.EQ(tUsers.Job)),
				),
		).
		WHERE(
			tUnitStatus.ID.EQ(jet.Uint64(unitId)),
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

func (s *Manager) GetLastUnitStatus(ctx context.Context, tx qrm.DB, unitId uint64) (*centrum.UnitStatus, error) {
	tUsers := tables.Users().AS("colleague")

	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.UnitJob,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.UserJob,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorID,
			tUnitStatus.CreatorJob,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.NamePrefix,
			tJobsUserProps.NameSuffix,
			tUserProps.Avatar.AS("colleague.avatar"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUsers.ID).
						AND(tJobsUserProps.Job.EQ(tUsers.Job)),
				),
		).
		WHERE(
			tUnitStatus.UnitID.EQ(jet.Uint64(unitId)),
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

func (s *Manager) DeleteUnit(ctx context.Context, unitJob string, unitId uint64) error {
	unit, err := s.State.GetUnit(ctx, unitJob, unitId)
	if err != nil {
		return nil
	}

	if unit.Job != unitJob {
		return errorscentrum.ErrFailedQuery
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}

	if err := s.State.DeleteUnit(ctx, unitJob, unitId); err != nil {
		return err
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitDeleted, unitJob), data); err != nil {
		return err
	}

	stmt := tUnits.
		DELETE().
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(unitJob)),
			tUnits.ID.EQ(jet.Uint64(unitId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Manager) ListUnitAccess(ctx context.Context, unitId uint64) (*centrum.UnitAccess, error) {
	access := &centrum.UnitAccess{}

	jobsAccess, err := s.unitAccess.Jobs.List(ctx, s.db, unitId)
	if err != nil {
		return nil, err
	}
	access.Jobs = jobsAccess

	for i := range access.Jobs {
		s.enricher.EnrichJobInfo(access.Jobs[i])
	}

	qualificationsccess, err := s.unitAccess.Qualifications.List(ctx, s.db, unitId)
	if err != nil {
		return nil, err
	}
	access.Qualifications = qualificationsccess

	return access, nil
}
