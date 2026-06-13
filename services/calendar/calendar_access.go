package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToCalendar(
	ctx context.Context,
	calendarId int64,
	userInfo *userinfo.UserInfo,
	access calendaraccess.AccessLevel,
	publicOk bool,
) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarIDs(ctx, userInfo, access, publicOk, calendarId)
	return len(out) > 0, err
}

var calendarSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(calendaraccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_SHARE),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
	},
}

type calendarAccessEntry struct {
	ID     int64 `alias:"calendar.id"`
	Public bool  `alias:"calendar.public"`
}

func (s *Server) checkIfUserHasAccessToCalendarIDs(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	access calendaraccess.AccessLevel,
	publicOk bool,
	calendarIds ...int64,
) ([]*calendarAccessEntry, error) {
	var dest []*calendarAccessEntry
	if len(calendarIds) == 0 {
		return dest, nil
	}

	tCreator := table.FivenetUser.AS("creator")

	ids := make([]mysql.Expression, len(calendarIds))
	for i := range calendarIds {
		ids[i] = mysql.Int64(calendarIds[i])
	}

	condition := mysql.AND(
		tCalendar.Job.IS_NULL(),
		tCalendar.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
	)
	if publicOk {
		condition = condition.OR(tCalendar.Public.IS_TRUE())
	}

	var accessExists mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		accessExists = s.access.ACLAccessExistsCondition(tCalendar.ID, userInfo, int32(access))
	}

	visibleCondition := mysql.OR(accessExists, condition)
	if userInfo.GetSuperuser() {
		visibleCondition = mysql.OR(
			tCalendar.SystemKind.IS_NULL(),
			tCalendar.SystemKind.NOT_EQ(
				mysql.Int32(
					int32(calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
				),
			),
			s.birthdayCalendarVisible(tCalendar.ID, access, userInfo),
		)
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
		WHERE(mysql.AND(
			tCalendar.ID.IN(ids...),
			tCalendar.DeletedAt.IS_NULL(),
			visibleCondition,
		)).
		ORDER_BY(tCalendar.ID.DESC())

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
