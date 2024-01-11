package manager

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) UpdateUnitStatus(ctx context.Context, job string, unitId uint64, in *centrum.UnitStatus) (*centrum.UnitStatus, error) {
	unit, err := s.GetUnit(job, unitId)
	if err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	// If the unit status is the same and is a status that shouldn't be duplicated, don't update the status again
	if unit.Status != nil &&
		unit.Status.Status == in.Status &&
		(in.Status == centrum.StatusUnit_STATUS_UNIT_ON_BREAK ||
			in.Status == centrum.StatusUnit_STATUS_UNIT_BUSY ||
			in.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE ||
			in.Status == centrum.StatusUnit_STATUS_UNIT_AVAILABLE) &&
		// Additionally if the status is under 2 minutes disallow the same status update
		(unit.Status.CreatedAt == nil || time.Since(unit.Status.CreatedAt.AsTime()) < 3*time.Minute) {
		s.logger.Debug("skipping unit status update due to same status or time", zap.Uint64("unit_id", unitId), zap.String("status", in.Status.String()))
		return nil, nil
	}

	if unit.Attributes != nil && unit.Attributes.Has(centrum.UnitAttributeStatic) {
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
		in.User, err = s.resolveUserShortById(ctx, *in.UserId)
		if err != nil {
			return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
		}

		if marker, ok := s.tracker.GetUserById(*in.UserId); ok {
			in.X = &marker.Info.X
			in.Y = &marker.Info.Y
			in.Postal = marker.Info.Postal
		}
	}
	if in.CreatorId != nil {
		var err error
		in.Creator, err = s.resolveUserShortById(ctx, *in.CreatorId)
		if err != nil {
			return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
		}
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
		).
		VALUES(
			jet.CURRENT_TIMESTAMP(),
			in.UnitId,
			in.Status,
			in.Reason,
			in.Code,
			in.UserId,
			in.X,
			in.Y,
			in.Postal,
			in.CreatorId,
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

	if err := s.State.UpdateUnitStatus(ctx, job, in.UnitId, in); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, job), data); err != nil {
		return nil, fmt.Errorf("failed to publish unit status event (size: %d, message: '%+v'): %w", len(data), in, err)
	}

	return in, nil
}

func (s *Manager) UpdateUnitAssignments(ctx context.Context, job string, userId *int32, unitId uint64, toAdd []int32, toRemove []int32) error {
	s.logger.Debug("updating unit assignments", zap.String("job", job), zap.Uint64("unit_id", unitId), zap.Int32s("toAdd", toAdd), zap.Int32s("toRemove", toRemove))

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return nil
	}

	var x, y *float64
	var postal *string
	if userId != nil {
		if marker, ok := s.tracker.GetUserById(*userId); ok {
			x = &marker.Info.X
			y = &marker.Info.Y
			postal = marker.Info.Postal
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

	if len(toRemove) > 0 {
		removeIds := make([]jet.Expression, len(toRemove))
		for i := 0; i < len(toRemove); i++ {
			removeIds[i] = jet.Int32(toRemove[i])
		}

		stmt := tUnitUser.
			DELETE().
			WHERE(jet.AND(
				tUnitUser.UnitID.EQ(jet.Uint64(unitId)),
				tUnitUser.UserID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		addIds := []jet.IntegerExpression{}
		for i := 0; i < len(toAdd); i++ {
			if _, ok := s.tracker.GetUserById(toAdd[i]); !ok {
				continue
			}

			unit, err := s.GetUnit(job, unitId)
			if err != nil {
				return err
			}
			// Skip already added units
			if slices.ContainsFunc(unit.Users, func(in *centrum.UnitAssignment) bool {
				return in.UserId == toAdd[i]
			}) {
				continue
			}

			addIds = append(addIds, jet.Int32(toAdd[i]))
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

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
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

	key := state.JobIdKey(job, unitId)
	if err := store.ComputeUpdate(key, true, func(key string, unit *centrum.Unit) (*centrum.Unit, error) {
		if len(toRemove) > 0 {
			toAnnounce := []int32{}
			for i := len(unit.Users) - 1; i >= 0; i-- {
				if i > (len(unit.Users) - 1) {
					break
				}

				for k := 0; k < len(toRemove); k++ {
					if unit.Users[i].UserId != toRemove[k] {
						continue
					}

					toAnnounce = append(toAnnounce, toRemove[k])
					unit.Users = slices.Delete(unit.Users, i, i+1)
				}
			}

			// Send updates
			for _, user := range toAnnounce {
				if err := s.AddUnitStatus(ctx, s.db, job, &centrum.UnitStatus{
					CreatedAt: timestamp.Now(),
					UnitId:    unit.Id,
					Status:    centrum.StatusUnit_STATUS_UNIT_USER_REMOVED,
					UserId:    &user,
					CreatorId: userId,
					X:         x,
					Y:         y,
					Postal:    postal,
				}); err != nil {
					return nil, err
				}

				if err := s.UnsetUnitIDForUser(user); err != nil {
					return nil, err
				}
			}
		}

		if len(toAdd) > 0 {
			notFound := []int32{}
			for i := 0; i < len(toAdd); i++ {
				if _, ok := s.tracker.GetUserById(toAdd[i]); !ok {
					continue
				}

				// Skip already added units
				if slices.ContainsFunc(unit.Users, func(in *centrum.UnitAssignment) bool {
					return in.UserId == toAdd[i]
				}) {
					continue
				}

				notFound = append(notFound, toAdd[i])
			}

			users, err := s.resolveUserShortsByIds(ctx, notFound)
			if err != nil {
				return nil, err
			}

			for _, user := range users {
				unit.Users = append(unit.Users, &centrum.UnitAssignment{
					UnitId: unit.Id,
					UserId: user.UserId,
					User:   user,
				})
			}

			for _, user := range users {
				if err := s.AddUnitStatus(ctx, s.db, job, &centrum.UnitStatus{
					CreatedAt: timestamp.Now(),
					UnitId:    unit.Id,
					Status:    centrum.StatusUnit_STATUS_UNIT_USER_ADDED,
					UserId:    &user.UserId,
					CreatorId: userId,
					X:         x,
					Y:         y,
					Postal:    postal,
				}); err != nil {
					return nil, err
				}

				if err := s.SetUnitForUser(user.Job, user.UserId, unit.Id); err != nil {
					return nil, err
				}
			}
		}

		// Unit is empty, set unit status to be unavailable automatically
		if len(unit.Users) == 0 {
			unit.Status = &centrum.UnitStatus{
				CreatedAt: timestamp.Now(),
				UnitId:    unit.Id,
				Status:    centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:    userId,
				CreatorId: userId,
				X:         x,
				Y:         y,
				Postal:    postal,
			}
			if err := s.AddUnitStatus(ctx, s.db, job, unit.Status); err != nil {
				return nil, err
			}
		}

		return unit, nil
	}); err != nil {
		return err
	}

	unit, err := s.GetUnit(job, unitId)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitUpdated, job), data); err != nil {
		return err
	}

	return nil
}

func (s *Manager) CreateUnit(ctx context.Context, job string, unit *centrum.Unit) (*centrum.Unit, error) {
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
		).
		VALUES(
			job,
			unit.Name,
			unit.Initials,
			unit.Color,
			unit.Description,
			unit.Attributes,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	// A new unit shouldn't have a status, so we make sure it has one
	if err := s.AddUnitStatus(ctx, tx, job, &centrum.UnitStatus{
		CreatedAt: timestamp.Now(),
		UnitId:    uint64(lastId),
		Status:    centrum.StatusUnit_STATUS_UNIT_UNKNOWN,
	}); err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Load new/updated unit from database
	if err := s.LoadUnitsFromDB(ctx, uint64(lastId)); err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	unit.Id = uint64(lastId)

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitCreated, job), data); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *Manager) UpdateUnit(ctx context.Context, job string, unit *centrum.Unit) (*centrum.Unit, error) {
	description := ""
	if unit.Description != nil {
		description = *unit.Description
	}

	stmt := tUnits.
		UPDATE(
			tUnits.Name,
			tUnits.Initials,
			tUnits.Color,
			tUnits.Description,
			tUnits.Attributes,
		).
		SET(
			unit.Name,
			unit.Initials,
			unit.Color,
			description,
			unit.Attributes,
		).
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(job)),
			tUnits.ID.EQ(jet.Uint64(unit.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	// Load new/updated unit from database
	if err := s.LoadUnitsFromDB(ctx, unit.Id); err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	if err := s.State.UpdateUnit(ctx, unit.Job, unit.Id, unit); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitUpdated, job), data); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *Manager) AddUnitStatus(ctx context.Context, tx qrm.DB, job string, status *centrum.UnitStatus) error {
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
		).
		VALUES(
			jet.CURRENT_TIMESTAMP(),
			status.UnitId,
			status.Status,
			status.Reason,
			status.Code,
			status.UserId,
			status.X,
			status.Y,
			status.Postal,
			status.CreatorId,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	status.Id = uint64(lastId)

	data, err := proto.Marshal(status)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, job), data); err != nil {
		return err
	}

	return nil
}

func (s *Manager) DeleteUnit(ctx context.Context, job string, id uint64) error {
	unit, err := s.State.GetUnit(job, id)
	if err != nil {
		return nil
	}

	if unit.Job != job {
		return errorscentrum.ErrFailedQuery
	}

	stmt := tUnits.
		DELETE().
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(job)),
			tUnits.ID.EQ(jet.Uint64(id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}

	if err := s.State.DeleteUnit(job, id); err != nil {
		return err
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitDeleted, job), data); err != nil {
		return err
	}

	return nil
}