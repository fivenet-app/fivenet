package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/puzpuzpuz/xsync/v2"
	"google.golang.org/protobuf/proto"
)

var (
	tCentrumSettings = table.FivenetCentrumSettings
)

func (s *Server) resolveUserById(ctx context.Context, u int32) (*users.User, error) {
	tUsers := tUsers.AS("user")
	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(tUsers).
		WHERE(
			tUsers.ID.EQ(jet.Int32(u)),
		).
		LIMIT(1)

	dest := users.User{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) resolveUserShortById(ctx context.Context, u int32) (*users.UserShort, error) {
	us, err := s.resolveUserShortsByIds(ctx, []int32{u})
	if err != nil {
		return nil, err
	}

	return us[0], nil
}

func (s *Server) resolveUserShortsByIds(ctx context.Context, u []int32) ([]*users.UserShort, error) {
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
			tUsers.Sex,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
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
	ds, ok := s.state.Disponents.Load(job)
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

func (s *Server) dispatchCenterSignOn(ctx context.Context, job string, userId int32, signon bool) error {
	if signon {
		if _, ok := s.tracker.GetUserByJobAndID(job, userId); !ok {
			return ErrNotOnDuty
		}

		stmt := tCentrumUsers.
			INSERT(
				tCentrumUsers.Job,
				tCentrumUsers.UserID,
				tCentrumUsers.Identifier,
			).
			VALUES(
				job,
				userId,
				tUsers.
					SELECT(
						tUsers.Identifier.AS("identifier"),
					).
					FROM(tUsers).
					WHERE(
						tUsers.ID.EQ(jet.Int32(userId)),
					).
					LIMIT(1),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return ErrFailedQuery
			}
		}
	} else {
		stmt := tCentrumUsers.
			DELETE().
			WHERE(jet.AND(
				tCentrumUsers.Job.EQ(jet.String(job)),
				tCentrumUsers.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return ErrFailedQuery
		}
	}

	// Load updated disponents into state
	if err := s.loadDisponents(ctx, job); err != nil {
		return ErrFailedQuery
	}

	disponents := s.getDisponents(job)
	change := &DisponentsChange{
		Job:        job,
		Disponents: disponents,
	}
	data, err := proto.Marshal(change)
	if err != nil {
		return ErrFailedQuery
	}
	s.events.JS.PublishAsync(buildSubject(TopicGeneral, TypeGeneralDisponents, job, 0), data)

	return nil
}

func (s *Server) getSettings(job string) *dispatch.Settings {
	settings, ok := s.state.Settings.Load(job)
	if !ok {
		// Return default settings
		return &dispatch.Settings{
			Job:          job,
			Enabled:      false,
			Mode:         dispatch.CentrumMode_CENTRUM_MODE_MANUAL,
			FallbackMode: dispatch.CentrumMode_CENTRUM_MODE_MANUAL,
		}
	}

	return settings
}

func (s *Server) getDisponents(job string) []*users.UserShort {
	disponents, ok := s.state.Disponents.Load(job)
	if !ok {
		return nil
	}

	return disponents
}

func (s *Server) getDispatchesMap(job string) *xsync.MapOf[uint64, *dispatch.Dispatch] {
	store, _ := s.state.Dispatches.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Dispatch] {
		return xsync.NewIntegerMapOf[uint64, *dispatch.Dispatch]()
	})

	return store
}

func (s *Server) getUnitsMap(job string) *xsync.MapOf[uint64, *dispatch.Unit] {
	store, _ := s.state.Units.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Unit] {
		return xsync.NewIntegerMapOf[uint64, *dispatch.Unit]()
	})

	return store
}

func (s *Server) isStatusDispatchComplete(in dispatch.StatusDispatch) bool {
	return in == dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED
}

func (s *Server) checkIfBotNeeded(job string) bool {
	settings := s.getSettings(job)

	if settings.Mode == dispatch.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
		return true
	}

	disponents := s.getDisponents(job)
	if len(disponents) == 0 {
		if settings.FallbackMode == dispatch.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
			return true
		}
	}

	return false
}
