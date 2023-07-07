package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getDispatch(ctx context.Context, userInfo *userinfo.UserInfo, id uint64) (*dispatch.Dispatch, error) {
	jobDispatches, ok := s.dispatches.Load(userInfo.Job)
	if !ok {
		return nil, nil
	}

	dispatch, ok := jobDispatches.Load(id)
	if !ok {
		return nil, nil
	}

	return dispatch, nil
}

func (s *Server) getUnit(ctx context.Context, userInfo *userinfo.UserInfo, id uint64) (*dispatch.Unit, error) {
	jobUnits, ok := s.units.Load(userInfo.Job)
	if !ok {
		return nil, nil
	}

	unit, ok := jobUnits.Load(id)
	if !ok {
		return nil, nil
	}

	return unit, nil
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
