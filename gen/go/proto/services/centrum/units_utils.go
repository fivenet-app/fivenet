package centrum

import (
	"context"
	"strings"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
)

func (s *Server) getUnit(job string, id uint64) (*dispatch.Unit, bool) {
	units, ok := s.units.Load(job)
	if !ok {
		return nil, false
	}

	return units.Load(id)
}

func (s *Server) listUnits(job string) ([]*dispatch.Unit, error) {
	us := []*dispatch.Unit{}

	units, ok := s.units.Load(job)
	if !ok {
		return nil, nil
	}

	units.Range(func(key uint64, unit *dispatch.Unit) bool {
		us = append(us, unit)
		return true
	})

	slices.SortFunc(us, func(a, b *dispatch.Unit) int {
		return strings.Compare(a.Name, b.Name)
	})

	return us, nil
}

func (s *Server) getUnitStatusFromDB(ctx context.Context, id uint64) (*dispatch.UnitStatus, error) {
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
		dest.User, err = s.resolveUserById(ctx, *dest.UserId)
		if err != nil {
			return nil, err
		}
	}
	if dest.CreatorId != nil {
		var err error
		dest.Creator, err = s.resolveUserById(ctx, *dest.CreatorId)
		if err != nil {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) resolveUsersForUnit(ctx context.Context, u []*dispatch.UnitAssignment) ([]*dispatch.UnitAssignment, error) {
	userIds := make([]int32, len(u))
	for i := 0; i < len(u); i++ {
		userIds[i] = u[i].UserId
	}

	if len(userIds) == 0 {
		return nil, nil
	}

	us, err := s.resolveUsersByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(u); i++ {
		u[i].User = us[i]
	}

	return u, nil
}

func (s *Server) getUnitIDForUserID(userId int32) (uint64, bool) {
	return s.userIDToUnitID.Load(userId)
}

func (s *Server) updateUnitStatus(ctx context.Context, job string, unit *dispatch.Unit, in *dispatch.UnitStatus) error {
	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		INSERT(
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.CreatorID,
		).
		VALUES(
			in.UnitId,
			in.Status,
			in.Reason,
			in.Code,
			in.UserId,
			in.X,
			in.Y,
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

	status, err := s.getUnitStatusFromDB(ctx, uint64(lastId))
	if err != nil {
		return err
	}
	unit.Status = status

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitStatus, job, status.UnitId), data)

	return nil
}

func (s *Server) updateUnitAssignments(ctx context.Context, userInfo *userinfo.UserInfo, unit *dispatch.Unit, toAdd []int32, toRemove []int32) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

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

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}

		for i := 0; i < len(unit.Users); i++ {
			for k := len(toRemove) - 1; k >= 0; k-- {
				if unit.Users[i].UserId == toRemove[k] {
					if err := s.updateUnitStatus(ctx, userInfo.Job, unit, &dispatch.UnitStatus{
						UnitId:    unit.Id,
						Status:    dispatch.UNIT_STATUS_USER_REMOVED,
						UserId:    &toRemove[k],
						CreatorId: &userInfo.UserId,
					}); err != nil {
						return err
					}

					unit.Users = utils.RemoveFromSlice(unit.Users, i)
					s.userIDToUnitID.Delete(toRemove[k])

					break
				}
			}
		}
	}

	if len(toAdd) > 0 {
		found := []int32{}
		addIds := []jet.IntegerExpression{}
		for i := 0; i < len(toAdd); i++ {
			_, ok := s.tracker.GetUserById(toAdd[i])
			if !ok {
				continue
			}

			// Skip already added units
			if utils.InSliceFunc(unit.Users, func(in *dispatch.UnitAssignment) bool {
				return in.UserId == toAdd[i]
			}) {
				continue
			}

			addIds = append(addIds, jet.Int32(toAdd[i]))
			found = append(found, toAdd[i])
		}

		for _, id := range addIds {
			stmt := tUnitUser.
				INSERT(
					tUnitUser.UnitID,
					tUnitUser.UserID,
					tUnitUser.Identifier,
				).
				VALUES(
					unit.Id,
					id,
					tUsers.
						SELECT(
							tUsers.Identifier,
						).
						FROM(tUsers).
						WHERE(tUsers.ID.EQ(id)).
						LIMIT(1),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if dbutils.IsDuplicateError(err) {
					return ErrAlreadyInUnit
				} else {
					return err
				}
			}
		}

		users, err := s.resolveUsersByIds(ctx, found)
		if err != nil {
			return err
		}

		for _, user := range users {
			unit.Users = append(unit.Users, &dispatch.UnitAssignment{
				UnitId: unit.Id,
				UserId: user.UserId,
				User:   user,
			})

			if err := s.updateUnitStatus(ctx, userInfo.Job, unit, &dispatch.UnitStatus{
				UnitId:    unit.Id,
				Status:    dispatch.UNIT_STATUS_USER_ADDED,
				UserId:    &user.UserId,
				CreatorId: &userInfo.UserId,
			}); err != nil {
				return err
			}

			s.userIDToUnitID.Store(user.UserId, unit.Id)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return ErrFailedQuery
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}

	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUpdated, userInfo.Job, unit.Id), data)

	// Send unit user assigned message when needed
	if len(toAdd) > 0 {
		s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUserAssigned, userInfo.Job, 0), data)
	}
	if len(toRemove) > 0 {
		s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUserAssigned, userInfo.Job, unit.Id), data)
	}

	// Unit is empty, set unit status to be unavailable automatically
	if len(unit.Users) == 0 {
		if err := s.updateUnitStatus(ctx, userInfo.Job, unit, &dispatch.UnitStatus{
			UnitId:    unit.Id,
			Status:    dispatch.UNIT_STATUS_UNAVAILABLE,
			UserId:    &userInfo.UserId,
			CreatorId: &userInfo.UserId,
		}); err != nil {
			return err
		}
	}

	return nil
}
