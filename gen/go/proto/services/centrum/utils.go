package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getDispatch(ctx context.Context, userInfo *userinfo.UserInfo, id uint64) (*dispatch.Dispatch, bool) {
	jobDispatches, ok := s.dispatches.Load(userInfo.Job)
	if !ok {
		return nil, false
	}

	dispatch, ok := jobDispatches.Load(id)
	if !ok {
		return nil, false
	}

	return dispatch, true
}

func (s *Server) getDispatchFromDB(ctx context.Context, tx qrm.Queryable, id uint64) (*dispatch.Dispatch, error) {
	condition := tDispatch.ID.EQ(jet.Uint64(id)).AND(jet.OR(
		tDispatchStatus.ID.IS_NULL(),
		tDispatchStatus.ID.EQ(
			jet.RawInt(`SELECT MAX(dispatchstatus.id) FROM fivenet_centrum_dispatches_status AS dispatchstatus WHERE dispatchstatus.dispatch_id = dispatch.id`),
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
		LIMIT(1)

	dispatch := dispatch.Dispatch{}
	if err := stmt.QueryContext(ctx, tx, &dispatch); err != nil {
		return nil, err
	}

	return &dispatch, nil
}

func (s *Server) getUnit(ctx context.Context, userInfo *userinfo.UserInfo, id uint64) (*dispatch.Unit, bool) {
	jobUnits, ok := s.units.Load(userInfo.Job)
	if !ok {
		return nil, false
	}

	unit, ok := jobUnits.Load(id)
	if !ok {
		return nil, false
	}

	return unit, true
}

func (s *Server) getUnitFromDB(ctx context.Context, tx qrm.Queryable, id uint64) (*dispatch.Unit, error) {
	condition := tUnitStatus.ID.IS_NULL().OR(
		tUnitStatus.ID.EQ(
			jet.RawInt(`SELECT MAX(unitstatus.id) FROM fivenet_centrum_units_status AS unitstatus WHERE unitstatus.unit_id = unit.id`),
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
			tUnitStatus.InSquad,
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

func (s *Server) resolveUsersForUnit(ctx context.Context, u []*dispatch.UnitAssignment) ([]*dispatch.UnitAssignment, error) {
	userIds := make([]int32, len(u))
	for i := 0; i < len(u); i++ {
		userIds[i] = u[i].UserId
	}

	if len(userIds) == 0 {
		return nil, nil
	}

	us, err := s.resolveUsersById(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(u); i++ {
		u[i].User = us[i]
	}

	return u, nil
}

func (s *Server) resolveUsersById(ctx context.Context, u []int32) ([]*users.UserShort, error) {
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
