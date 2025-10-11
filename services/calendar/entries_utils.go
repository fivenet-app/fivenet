package calendar

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) listCalendarEntriesQuery(
	condition mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
	access calendar.AccessLevel,
) mysql.SelectStatement {
	tCreator := tables.User().AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	rsvp2 := tCalendarRSVP.AS("r2")

	accessExists := mysql.EXISTS(
		mysql.
			SELECT(mysql.Int(1)).
			FROM(tCAccess).
			WHERE(mysql.AND(
				tCAccess.TargetID.EQ(tCalendarEntry.CalendarID),
				tCAccess.Access.GT_EQ(mysql.Int32(int32(access))),
				mysql.OR(
					tCAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					mysql.AND(
						tCAccess.Job.EQ(mysql.String(userInfo.GetJob())),
						tCAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					),
				),
			)),
	)

	rsvpExists := mysql.EXISTS(
		mysql.
			SELECT(mysql.Int(1)).
			FROM(tCalendarRSVP.AS("r2")).
			WHERE(mysql.AND(
				rsvp2.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				rsvp2.Response.GT(mysql.Int32(3)), // response > 3
				rsvp2.EntryID.EQ(tCalendarEntry.ID),
			)),
	)

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
			tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
			tAvatar.FilePath.AS("creator.profile_picture"),
			tCalendarEntry.Recurring,
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		).
		FROM(tCalendarEntry.
			INNER_JOIN(tCalendar,
				mysql.AND(
					tCalendar.ID.EQ(tCalendarEntry.CalendarID),
					tCalendar.DeletedAt.IS_NULL(),
				),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tCalendarRSVP,
				mysql.AND(
					tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarRSVP.EntryID.EQ(tCalendarEntry.ID),
				),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(mysql.AND(
			mysql.OR(
				accessExists,
				rsvpExists,
				tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
			),
			condition,
		)).
		LIMIT(100)

	return stmt
}
