package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToCalendar(ctx context.Context, calendarId uint64, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarIDs(ctx, userInfo, access, publicOk, calendarId)
	return len(out) > 0, err
}

type calendarAccessEntry struct {
	ID     uint64 `alias:"calendar.id"`
	Public bool   `alias:"calendar.public"`
}

func (s *Server) checkIfUserHasAccessToCalendarIDs(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool, calendarIds ...uint64) ([]*calendarAccessEntry, error) {
	var dest []*calendarAccessEntry
	if len(calendarIds) == 0 {
		return dest, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		for i := 0; i < len(calendarIds); i++ {
			dest = append(dest, &calendarAccessEntry{
				ID: calendarIds[i],
			})
		}
		return dest, nil
	}

	tCreator := tables.Users().AS("creator")

	ids := make([]jet.Expression, len(calendarIds))
	for i := 0; i < len(calendarIds); i++ {
		ids[i] = jet.Uint64(calendarIds[i])
	}

	condition := jet.Bool(false)
	if publicOk {
		condition = tCalendar.Public.IS_TRUE()
	}

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
		).
		FROM(tCalendar.
			LEFT_JOIN(tCAccess,
				tCAccess.TargetID.EQ(tCalendar.ID).
					AND(tCAccess.Access.GT_EQ(jet.Int32(int32(access)))),
			).
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(jet.AND(
			tCalendar.ID.IN(ids...),
			tCalendar.DeletedAt.IS_NULL(),
			jet.OR(
				tCalendar.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				tCalendar.CreatorJob.EQ(jet.String(userInfo.Job)),
				tCAccess.UserID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tCAccess.Job.EQ(jet.String(userInfo.Job)),
					tCAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade)),
				),
				condition,
			),
		)).
		ORDER_BY(tCalendar.ID.DESC(), tCAccess.MinimumGrade)

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
