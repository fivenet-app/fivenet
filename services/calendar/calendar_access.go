package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToCalendar(
	ctx context.Context,
	calendarId int64,
	userInfo *userinfo.UserInfo,
	access calendar.AccessLevel,
	publicOk bool,
) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarIDs(ctx, userInfo, access, publicOk, calendarId)
	return len(out) > 0, err
}

type calendarAccessEntry struct {
	ID     int64 `alias:"calendar.id"`
	Public bool  `alias:"calendar.public"`
}

func (s *Server) checkIfUserHasAccessToCalendarIDs(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	access calendar.AccessLevel,
	publicOk bool,
	calendarIds ...int64,
) ([]*calendarAccessEntry, error) {
	var dest []*calendarAccessEntry
	if len(calendarIds) == 0 {
		return dest, nil
	}

	// Allow superusers access to any docs
	if userInfo.GetSuperuser() {
		for i := range calendarIds {
			dest = append(dest, &calendarAccessEntry{
				ID: calendarIds[i],
			})
		}
		return dest, nil
	}

	tCreator := tables.User().AS("creator")

	ids := make([]jet.Expression, len(calendarIds))
	for i := range calendarIds {
		ids[i] = jet.Int64(calendarIds[i])
	}

	condition := jet.Bool(false)
	if publicOk {
		condition = tCalendar.Public.IS_TRUE()
	}

	var accessExists jet.BoolExpression
	if !userInfo.GetSuperuser() {
		accessExists = jet.EXISTS(
			jet.
SELECT(jet.Int(1)).
				FROM(tCAccess).
				WHERE(
					jet.AND(
						tCAccess.TargetID.EQ(tCalendar.ID),
						tCAccess.Access.GT_EQ(jet.Int32(int32(access))),
						jet.OR(
							tCAccess.UserID.EQ(jet.Int32(userInfo.GetUserId())),
							jet.AND(
								tCAccess.Job.EQ(jet.String(userInfo.GetJob())),
								tCAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.GetJobGrade())),
							),
						),
					),
				),
		)
	} else {
		accessExists = jet.Bool(true)
	}

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
		).
		FROM(tCalendar.
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		WHERE(jet.AND(
			tCalendar.ID.IN(ids...),
			tCalendar.DeletedAt.IS_NULL(),
			jet.OR(
				accessExists,
				condition,
			),
		)).
		ORDER_BY(tCalendar.ID.DESC())

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
