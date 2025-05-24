package calendar

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) listCalendarEntriesQuery(condition jet.BoolExpression, userInfo *userinfo.UserInfo, access calendar.AccessLevel) jet.SelectStatement {
	tCreator := tables.User().AS("creator")

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarEntry.CreatedAt,
			tCalendarEntry.UpdatedAt,
			tCalendarEntry.DeletedAt,
			tCalendarEntry.CalendarID,
			tCalendar.ID,
			tCalendar.Name,
			tCalendar.Color,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendarEntry.Job,
			tCalendarEntry.StartTime,
			tCalendarEntry.EndTime,
			tCalendarEntry.Title,
			tCalendarEntry.Closed,
			tCalendarEntry.RsvpOpen,
			tCalendarEntry.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
			tCalendarEntry.Recurring,
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		).
		FROM(tCalendarEntry.
			INNER_JOIN(tCalendar,
				tCalendar.ID.EQ(tCalendarEntry.CalendarID).
					AND(tCalendar.DeletedAt.IS_NULL()),
			).
			LEFT_JOIN(tCAccess,
				tCAccess.TargetID.EQ(tCalendarEntry.CalendarID).
					AND(tCAccess.Access.GT_EQ(jet.Int32(int32(access)))),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tCalendarRSVP,
				tCalendarRSVP.UserID.EQ(jet.Int32(userInfo.UserId)).
					AND(tCalendarRSVP.EntryID.EQ(tCalendarEntry.ID)),
			),
		).
		GROUP_BY(tCalendarEntry.ID).
		WHERE(condition).
		LIMIT(100)

	return stmt
}
