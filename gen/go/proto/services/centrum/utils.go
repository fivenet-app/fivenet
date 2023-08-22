package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/puzpuzpuz/xsync/v2"
)

var (
	tCentrumSettings = table.FivenetCentrumSettings
)

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

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.Job,
		).
		FROM(tUsers).
		WHERE(
			tUsers.ID.IN(userIds...),
		).
		LIMIT(int64(len(u)))

	dest := []*users.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Server) checkIfUserIsDisponent(job string, userId int32) bool {
	ds, ok := s.disponents.Load(job)
	if !ok {
		return false
	}

	for i := 0; i < len(ds); i++ {
		if userId == ds[i].UserId {
			return true
		}
	}

	return false
}

func (s *Server) checkIfUserIsPartOfDispatch(userInfo *userinfo.UserInfo, dsp *dispatch.Dispatch, disponentOkay bool) bool {
	// Check if user is a disponent
	if s.checkIfUserIsDisponent(userInfo.Job, userInfo.UserId) {
		return true
	}

	// Iterate over units of dispatch and check if the user is in one of the units
	for i := 0; i < len(dsp.Units); i++ {
		unit, ok := s.getUnit(userInfo.Job, dsp.Units[i].UnitId)
		if !ok {
			continue
		}

		if s.checkIfUserPartOfUnit(userInfo.UserId, unit) {
			return true
		}
	}

	return false
}

func (s *Server) checkIfUserPartOfUnit(userId int32, unit *dispatch.Unit) bool {
	for i := 0; i < len(unit.Users); i++ {
		if unit.Users[i].UserId == userId {
			return true
		}
	}

	return false
}

func (s *Server) getSettings(job string) *dispatch.Settings {
	settings, ok := s.settings.Load(job)
	if !ok {
		// Return default settings
		return &dispatch.Settings{
			Job:          job,
			Enabled:      false,
			Mode:         dispatch.CENTRUM_MODE_MANUAL,
			FallbackMode: dispatch.CENTRUM_MODE_MANUAL,
		}
	}

	return settings
}

func (s *Server) getDisponents(job string) []*users.UserShort {
	disponents, ok := s.disponents.Load(job)
	if !ok {
		return nil
	}

	return disponents
}

func (s *Server) getDispatchesMap(job string) *xsync.MapOf[uint64, *dispatch.Dispatch] {
	store, _ := s.dispatches.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Dispatch] {
		return xsync.NewIntegerMapOf[uint64, *dispatch.Dispatch]()
	})

	return store
}

func (s *Server) getUnitsMap(job string) *xsync.MapOf[uint64, *dispatch.Unit] {
	store, _ := s.units.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Unit] {
		return xsync.NewIntegerMapOf[uint64, *dispatch.Unit]()
	})

	return store
}
