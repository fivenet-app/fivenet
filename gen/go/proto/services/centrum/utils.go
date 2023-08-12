package centrum

import (
	"context"
	"errors"
	"fmt"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

var (
	tCentrumSettings = table.FivenetCentrumSettings
)

func (s *Server) getDispatchFromDB(ctx context.Context, tx qrm.DB, id uint64) (*dispatch.Dispatch, error) {
	condition := tDispatch.ID.EQ(jet.Uint64(id)).AND(jet.OR(
		tDispatchStatus.ID.IS_NULL(),
		tDispatchStatus.ID.EQ(
			jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
		),
	))

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Anon,
			tDispatch.UserID,
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchUnit.UnitID,
			tDispatchUnit.DispatchID,
			tDispatchUnit.CreatedAt,
			tDispatchUnit.ExpiresAt,
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				).
				LEFT_JOIN(tDispatchUnit,
					tDispatchUnit.DispatchID.EQ(tDispatch.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tDispatch.CreatedAt.ASC(),
			tDispatch.Job.ASC(),
			tDispatchStatus.Status.ASC(),
		).
		LIMIT(24)

	dispatch := dispatch.Dispatch{}
	if err := stmt.QueryContext(ctx, tx, &dispatch); err != nil {
		return nil, err
	}

	return &dispatch, nil
}

func (s *Server) getDispatchStatus(ctx context.Context, id uint64) (*dispatch.DispatchStatus, error) {
	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tUser.ID,
			tUser.Identifier,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(
			tDispatchStatus.
				LEFT_JOIN(
					tUser,
					tUser.ID.EQ(tDispatchStatus.UserID),
				),
		).
		WHERE(
			tDispatchStatus.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	var dest dispatch.DispatchStatus
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) getUnit(ctx context.Context, userInfo *userinfo.UserInfo, id uint64) (*dispatch.Unit, bool) {
	unit := &dispatch.Unit{}
	if err := s.units.Get(fmt.Sprintf("%s/%d", userInfo.Job, id), unit); err != nil {
		return nil, false
	}

	return unit, true
}

func (s *Server) getUnitFromDB(ctx context.Context, tx qrm.DB, id uint64) (*dispatch.Unit, error) {
	condition := tUnitStatus.ID.IS_NULL().OR(
		tUnitStatus.ID.EQ(
			jet.RawInt("SELECT MAX(`unitstatus`.`id`) FROM `fivenet_centrum_units_status` AS `unitstatus` WHERE `unitstatus`.`unit_id` = `unit`.`id`"),
		),
	)

	stmt := tUnits.
		SELECT(
			tUnits.ID,
			tUnits.Job,
			tUnits.Name,
			tUnits.Initials,
			tUnits.Color,
			tUnits.Description,
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitUser.UnitID,
			tUnitUser.UserID,
			tUnitUser.Identifier,
		).
		FROM(
			tUnits.
				LEFT_JOIN(tUnitStatus,
					tUnitStatus.UnitID.EQ(tUnits.ID),
				).
				LEFT_JOIN(tUnitUser,
					tUnitUser.UnitID.EQ(tUnits.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tUnits.Job.ASC(),
			tUnits.Name.ASC(),
			tUnitStatus.Status.ASC(),
		).
		LIMIT(1)

	unit := dispatch.Unit{}
	if err := stmt.QueryContext(ctx, tx, &unit); err != nil {
		return nil, err
	}

	return &unit, nil
}

func (s *Server) getUnitStatus(ctx context.Context, id uint64) (*dispatch.UnitStatus, error) {
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
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(
					tUser,
					tUser.ID.EQ(tUnitStatus.UserID),
				),
		).
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

func (s *Server) resolveUserById(ctx context.Context, u int32) (*users.UserShort, error) {
	us, err := s.resolveUsersByIds(ctx, []int32{u})
	if err != nil {
		return nil, err
	}

	return us[0], nil
}

func (s *Server) resolveUsersByIds(ctx context.Context, u []int32) ([]*users.UserShort, error) {
	if len(u) == 0 {
		return nil, nil
	}

	userIds := make([]jet.Expression, len(u))
	for i := 0; i < len(u); i++ {
		userIds[i] = jet.Int32(u[i])
	}

	stmt := tUser.
		SELECT(
			tUser.ID.AS("user_id"),
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
		).
		FROM(tUser).
		WHERE(
			tUser.ID.IN(userIds...),
		)

	resolvedUsers := []*users.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, &resolvedUsers); err != nil {
		return nil, err
	}

	return resolvedUsers, nil
}

func (s *Server) getUnitIDForUserID(ctx context.Context, userId int32) (uint64, error) {
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
		UnitID uint64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return 0, err
		}

		return 0, nil
	}

	return dest.UnitID, nil
}

func (s *Server) checkIfUserIsPartOfDispatch(ctx context.Context, userInfo *userinfo.UserInfo, dsp *dispatch.Dispatch, disponentOkay bool) (bool, error) {
	// TODO check if user is disponent

	for i := 0; i < len(dsp.Units); i++ {
		unit, ok := s.getUnit(ctx, userInfo, dsp.Units[i].UnitId)
		if !ok {
			continue
		}

		if s.checkIfUserPartOfUnit(userInfo.UserId, unit) {
			return true, nil
		}
	}

	return false, nil
}

func (s *Server) checkIfUserPartOfUnit(userId int32, unit *dispatch.Unit) bool {
	for i := 0; i < len(unit.Users); i++ {
		if unit.Users[i].UserId == userId {
			return true
		}
	}

	return false
}

func (s *Server) updateDispatchStatus(ctx context.Context, userInfo *userinfo.UserInfo, dsp *dispatch.Dispatch, in *dispatch.DispatchStatus) (*dispatch.DispatchStatus, error) {
	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt := tDispatchStatus.
		INSERT(
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
		).
		VALUES(
			in.DispatchId,
			in.UnitId,
			in.Status,
			in.Reason,
			in.Code,
			in.UserId,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	status, err := s.getDispatchStatus(ctx, uint64(lastId))
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(status)
	if err != nil {
		return nil, err
	}

	for _, u := range dsp.Units {
		s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchStatus, userInfo, u.UnitId), data)
	}

	return status, nil
}

func (s *Server) updateUnitStatus(ctx context.Context, userInfo *userinfo.UserInfo, unit *dispatch.Unit, in *dispatch.UnitStatus) (*dispatch.UnitStatus, error) {
	unitId := jet.NULL
	if in.UnitId <= 0 {
		unitId = jet.Uint64(in.UnitId)
	}

	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		INSERT(
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
		).
		VALUES(
			unitId,
			in.Status,
			in.Reason,
			in.Code,
			userInfo.UserId,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	status, err := s.getUnitStatus(ctx, uint64(lastId))
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(status)
	if err != nil {
		return nil, err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitStatus, userInfo, status.UnitId), data)

	return status, nil
}

func (s *Server) updateDispatchUnitAssignments(ctx context.Context, userInfo *userinfo.UserInfo, unit *dispatch.Unit, toAdd []int32, toRemove []int32) error {
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
			for k := 0; k < len(toRemove); k++ {
				if unit.Users[i].UserId == toRemove[k] {
					unit.Users = utils.RemoveFromSlice(unit.Users, i)
				}
			}
		}
	}

	if len(toAdd) > 0 {
		addIds := make([]jet.IntegerExpression, len(toAdd))
		for i := 0; i < len(toAdd); i++ {
			_, ok := s.tracker.GetUserById(toAdd[i])
			if !ok {
				continue
			}

			addIds[i] = jet.Int32(toAdd[i])
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
					tUser.
						SELECT(
							tUser.Identifier,
						).
						FROM(tUser).
						WHERE(tUser.ID.EQ(id)).
						LIMIT(1),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}
		}

		found := []int32{}
		for k := 0; k < len(toAdd); k++ {
			for i := 0; i < len(unit.Users); i++ {
				if unit.Users[i].UserId == toAdd[k] {
					found = append(found, toAdd[k])
				}
			}
		}

		users, err := s.resolveUsersByIds(ctx, found)
		if err != nil {
			return err
		}
		assignments := []*dispatch.UnitAssignment{}
		for _, v := range users {
			assignments = append(assignments, &dispatch.UnitAssignment{
				UnitId: unit.Id,
				UserId: v.UserId,
				User:   v,
			})
		}
		unit.Users = assignments
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return ErrFailedQuery
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUpdated, userInfo, unit.Id), data)
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUserAssigned, userInfo, unit.Id), data)

	return nil
}

func (s *Server) getSettings(ctx context.Context, job string) (*dispatch.Settings, error) {
	settings := &dispatch.Settings{}
	if err := s.settings.Get(job, settings); err != nil {
		if !errors.Is(nats.ErrKeyNotFound, err) {
			return nil, err
		}

		// Return default settings
		return &dispatch.Settings{
			Job:          job,
			Enabled:      false,
			Mode:         dispatch.CENTRUM_MODE_MANUAL,
			FallbackMode: dispatch.CENTRUM_MODE_MANUAL,
		}, nil
	}

	return settings, nil
}

func (s *Server) getDisponents(ctx context.Context, job string) ([]*users.UserShort, error) {
	stmt := tCentrumUsers.
		SELECT(
			tUser.ID,
			tUser.Identifier,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.Job,
		).
		FROM(
			tCentrumUsers.
				INNER_JOIN(tUser,
					tCentrumUsers.UserID.EQ(tUser.ID),
				),
		).
		WHERE(
			tCentrumUsers.Job.EQ(jet.String(job)),
		)

	var dest []*users.UserShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return dest, nil
}
