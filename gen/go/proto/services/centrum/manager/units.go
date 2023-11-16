package manager

import (
	"context"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) GetUnitStatusFromDB(ctx context.Context, job string, id uint64) (*dispatch.UnitStatus, error) {
	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
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
		FROM(tUnitStatus).
		WHERE(
			tUnitStatus.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	var dest dispatch.UnitStatus
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	if dest.UserId != nil {
		var err error
		dest.User, err = s.resolveUserShortById(ctx, *dest.UserId)
		if err != nil {
			return nil, err
		}
	}
	if dest.CreatorId != nil {
		var err error
		dest.Creator, err = s.resolveUserShortById(ctx, *dest.CreatorId)
		if err != nil {
			return nil, err
		}
	}

	if dest.UnitId > 0 {
		unit, ok := s.GetUnit(job, dest.UnitId)
		if ok && unit != nil {
			newUnit := proto.Clone(unit)
			dest.Unit = newUnit.(*dispatch.Unit)
		}
	}

	return &dest, nil
}

func (s *Manager) UpdateUnitStatus(ctx context.Context, job string, unit *dispatch.Unit, in *dispatch.UnitStatus) error {
	// If the unit status is the same and is a status that shouldn't be duplicated, don't update the status again
	if unit.Status != nil &&
		unit.Status.Status == in.Status &&
		(in.Status == dispatch.StatusUnit_STATUS_UNIT_ON_BREAK ||
			in.Status == dispatch.StatusUnit_STATUS_UNIT_BUSY ||
			in.Status == dispatch.StatusUnit_STATUS_UNIT_UNAVAILABLE ||
			in.Status == dispatch.StatusUnit_STATUS_UNIT_AVAILABLE) &&
		// Additionally if the status is older than 2 minutes allow it to be updated again
		(unit.Status.CreatedAt == nil || time.Since(unit.Status.CreatedAt.AsTime()) < 3*time.Minute) {
		s.logger.Debug("skipping unit status update due to same status or time", zap.Uint64("unit_id", unit.Id), zap.String("status", in.Status.String()))
		return nil
	}

	s.logger.Debug("updating unit status", zap.Uint64("unit_id", unit.Id))

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
		return err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	status, err := s.GetUnitStatusFromDB(ctx, job, uint64(lastId))
	if err != nil {
		return err
	}

	unit.Status = status
	newUnit, ok := s.GetUnit(job, unit.Id)
	if !ok {
		return nil
	}

	newUnit.Status = status
	data, err := proto.Marshal(newUnit)
	if err != nil {
		return err
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, job, status.UnitId), data)

	return nil
}

func (s *Manager) UpdateUnitAssignments(ctx context.Context, job string, userId *int32, unit *dispatch.Unit, toAdd []int32, toRemove []int32) error {
	s.logger.Debug("updating unit assignments", zap.String("job", job), zap.Uint64("unit_id", unit.Id), zap.Int32s("toAdd", toAdd), zap.Int32s("toRemove", toRemove))

	var x, y *float64
	var postal *string
	if userId != nil {
		if marker, ok := s.tracker.GetUserById(*userId); ok {
			x = &marker.Info.X
			y = &marker.Info.Y
			postal = marker.Info.Postal
		}
	}

	var previousStatus proto.Message
	if unit.Status != nil && (len(toAdd) > 0 || len(toRemove) > 0) {
		previousStatus = proto.Clone(unit.Status)
	}

	lock := s.State.GetUnitLock(unit.Id)
	lock.Lock()
	defer lock.Unlock()

	tUnitUser := table.FivenetCentrumUnitsUsers
	if len(toRemove) > 0 {
		removeIds := make([]jet.Expression, len(toRemove))
		for i := 0; i < len(toRemove); i++ {
			removeIds[i] = jet.Int32(toRemove[i])
		}

		stmt := tUnitUser.
			DELETE().
			WHERE(jet.AND(
				tUnitUser.UnitID.EQ(jet.Uint64(unit.Id)),
				tUnitUser.UserID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}

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
				unit.Users = utils.RemoveFromSlice(unit.Users, i)
				s.UnsetUnitIDForUser(toRemove[k])
			}
		}

		// Send updates
		for _, user := range toAnnounce {
			if err := s.UpdateUnitStatus(ctx, job, unit, &dispatch.UnitStatus{
				UnitId:    unit.Id,
				Status:    dispatch.StatusUnit_STATUS_UNIT_USER_REMOVED,
				UserId:    &user,
				CreatorId: userId,
				X:         x,
				Y:         y,
				Postal:    postal,
			}); err != nil {
				return err
			}
		}
	}

	if len(toAdd) > 0 {
		notFound := []int32{}
		addIds := []jet.IntegerExpression{}
		for i := 0; i < len(toAdd); i++ {
			if _, ok := s.tracker.GetUserById(toAdd[i]); !ok {
				continue
			}

			// Skip already added units
			if utils.InSliceFunc(unit.Users, func(in *dispatch.UnitAssignment) bool {
				return in.UserId == toAdd[i]
			}) {
				continue
			}

			addIds = append(addIds, jet.Int32(toAdd[i]))
			notFound = append(notFound, toAdd[i])
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
						unit.Id,
						id,
					)
			}

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}
		}

		users, err := s.resolveUserShortsByIds(ctx, notFound)
		if err != nil {
			return err
		}

		for _, user := range users {
			s.SetUnitForUser(user.UserId, unit.Id)

			unit.Users = append(unit.Users, &dispatch.UnitAssignment{
				UnitId: unit.Id,
				UserId: user.UserId,
				User:   user,
			})
		}

		for _, user := range users {
			if err := s.UpdateUnitStatus(ctx, job, unit, &dispatch.UnitStatus{
				UnitId:    unit.Id,
				Status:    dispatch.StatusUnit_STATUS_UNIT_USER_ADDED,
				UserId:    &user.UserId,
				CreatorId: userId,
				X:         x,
				Y:         y,
				Postal:    postal,
			}); err != nil {
				return err
			}
		}
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitUpdated, job, unit.Id), data)

	// Unit is empty, set unit status to be unavailable automatically
	if len(unit.Users) == 0 {
		if err := s.UpdateUnitStatus(ctx, job, unit, &dispatch.UnitStatus{
			UnitId:    unit.Id,
			Status:    dispatch.StatusUnit_STATUS_UNIT_UNAVAILABLE,
			UserId:    userId,
			CreatorId: userId,
			X:         x,
			Y:         y,
			Postal:    postal,
		}); err != nil {
			return err
		}
	} else {
		// F*ck it, just set the unit status to what it was before (again)
		if previousStatus != nil {
			if err := s.UpdateUnitStatus(ctx, job, unit, previousStatus.(*dispatch.UnitStatus)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Manager) CreateUnit(ctx context.Context, job string, unit *dispatch.Unit) (*dispatch.Unit, error) {
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
		).
		VALUES(
			job,
			unit.Name,
			unit.Initials,
			unit.Color,
			unit.Description,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	// A new unit shouldn't have a status, so we make sure it has one
	if err := s.AddUnitStatus(ctx, tx, &dispatch.UnitStatus{
		UnitId: uint64(lastId),
		Status: dispatch.StatusUnit_STATUS_UNIT_UNKNOWN,
	}); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Load new/updated unit from database
	if err := s.LoadUnits(ctx, uint64(lastId)); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	unit, ok := s.GetUnit(job, uint64(lastId))
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitCreated, job, unit.Id), data)

	return unit, nil
}

func (s *Manager) UpdateUnit(ctx context.Context, job string, unit *dispatch.Unit) (*dispatch.Unit, error) {
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
		).
		SET(
			tUnits.Name.SET(jet.String(unit.Name)),
			tUnits.Initials.SET(jet.String(unit.Initials)),
			tUnits.Color.SET(jet.String(unit.Color)),
			tUnits.Description.SET(jet.String(description)),
		).
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(job)),
			tUnits.ID.EQ(jet.Uint64(unit.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	// Load new/updated unit from database
	if err := s.LoadUnits(ctx, unit.Id); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	unit, ok := s.GetUnit(job, unit.Id)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitUpdated, job, unit.Id), data)

	return unit, nil
}

func (s *Manager) AddUnitStatus(ctx context.Context, tx qrm.DB, status *dispatch.UnitStatus) error {
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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Manager) DeleteUnit(ctx context.Context, job string, id uint64) error {
	unit, ok := s.State.GetUnit(job, id)
	if !ok {
		return nil
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
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitDeleted, job, id), data)

	s.State.DeleteUnit(job, id)

	return nil
}
