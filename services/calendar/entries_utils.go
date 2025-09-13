package calendar

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) listCalendarEntriesQuery(
	condition jet.BoolExpression,
	userInfo *userinfo.UserInfo,
	access calendar.AccessLevel,
) jet.SelectStatement {
	tCreator := tables.User().AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	rsvp2 := tCalendarRSVP.AS("r2")

	accessExists := jet.EXISTS(
		jet.SELECT(jet.Int(1)).
			FROM(tCAccess).
			WHERE(jet.AND(
				tCAccess.TargetID.EQ(tCalendarEntry.CalendarID),
				tCAccess.Access.GT_EQ(jet.Int32(int32(access))),
			)),
	)

	rsvpExists := jet.EXISTS(
		jet.SELECT(jet.Int(1)).
			FROM(tCalendarRSVP.AS("r2")).
			WHERE(
				rsvp2.UserID.EQ(jet.Int32(userInfo.GetUserId())).
					AND(rsvp2.Response.GT(jet.Int32(3))). // response > 3
					AND(rsvp2.EntryID.EQ(tCalendarEntry.ID)),
			),
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
				tCalendar.ID.EQ(tCalendarEntry.CalendarID).
					AND(tCalendar.DeletedAt.IS_NULL()),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tCalendarRSVP,
				tCalendarRSVP.UserID.EQ(jet.Int32(userInfo.GetUserId())).
					AND(tCalendarRSVP.EntryID.EQ(tCalendarEntry.ID)),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(jet.AND(
			jet.OR(
				accessExists,
				rsvpExists,
				tCalendarEntry.CreatorID.EQ(jet.Int32(userInfo.GetUserId())),
			),
			condition,
		)).
		LIMIT(100)

	return stmt
}
