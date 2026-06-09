package calendar

import (
	"context"
	"errors"
	"slices"
	"time"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func calendarEntryVisibility(
	userInfo *userinfo.UserInfo,
	access calendaraccess.AccessLevel,
	rsvpResponse calendarentries.RsvpResponses,
) mysql.BoolExpression {
	tCalendarEntryRsvp := tCalendarRSVP.AS("r2")

	return mysql.OR(
		mysql.EXISTS(
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
		),
		mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tCalendarEntryRsvp).
				WHERE(mysql.AND(
					tCalendarEntryRsvp.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarEntryRsvp.Response.GT(mysql.Int32(int32(rsvpResponse))),
					tCalendarEntryRsvp.EntryID.EQ(tCalendarEntry.ID),
				)),
		),
		tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
	)
}

func (s *Server) birthdayCalendarVisible(userInfo *userinfo.UserInfo) mysql.BoolExpression {
	return mysql.AND(
		tCalendar.SystemKind.EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
		mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tCAccess).
				WHERE(mysql.AND(
					tCAccess.TargetID.EQ(tCalendarEntry.CalendarID),
					tCAccess.Access.GT_EQ(
						mysql.Int32(int32(calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
					mysql.OR(
						tCAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						mysql.AND(
							tCAccess.Job.EQ(mysql.String(userInfo.GetJob())),
							tCAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
						),
					),
				)),
		),
	)
}

func calendarEntriesQuery(
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
	visibility mysql.BoolExpression,
	limit *int64,
) mysql.SelectStatement {
	tCreator := table.FivenetUser.AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarEntry.CreatedAt,
			tCalendarEntry.UpdatedAt,
			tCalendarEntry.DeletedAt,
			tCalendarEntry.CalendarID,
			tCalendar.ID,
			tCalendar.Job,
			tCalendar.SystemKind,
			tCalendar.Name,
			tCalendar.Color,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
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
			visibility,
			condition,
		)).
		ORDER_BY(
			tCalendarEntry.StartTime.ASC(),
			tCalendarEntry.ID.ASC(),
		)

	if limit != nil && *limit > 0 {
		stmt = stmt.LIMIT(*limit)
	}

	return stmt
}

func (s *Server) loadExpandedCalendarEntries(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
	visibility mysql.BoolExpression,
	rangeStart, rangeEnd time.Time,
	limit *int64,
) ([]*calendarentries.CalendarEntry, error) {
	stmt := calendarEntriesQuery(userInfo, condition, visibility, limit)

	entries := []*calendarentries.CalendarEntry{}
	if err := stmt.QueryContext(ctx, s.db, &entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	expanded, err := s.expandCalendarEntries(ctx, userInfo, entries, rangeStart, rangeEnd)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return expanded, nil
}

func (s *Server) finalizeCalendarEntries(
	entries []*calendarentries.CalendarEntry,
	userInfo *userinfo.UserInfo,
) []*calendarentries.CalendarEntry {
	slices.SortFunc(entries, func(left, right *calendarentries.CalendarEntry) int {
		l := left.GetStartTime().AsTime()
		r := right.GetStartTime().AsTime()
		if l.Before(r) {
			return -1
		}
		if l.After(r) {
			return 1
		}
		if left.GetCalendarId() < right.GetCalendarId() {
			return -1
		}
		if left.GetCalendarId() > right.GetCalendarId() {
			return 1
		}
		if left.GetId() < right.GetId() {
			return -1
		}
		if left.GetId() > right.GetId() {
			return 1
		}
		return 0
	})

	if s.enricher == nil {
		return entries
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range entries {
		if entries[i].GetCreator() != nil {
			jobInfoFn(entries[i].GetCreator())
		}
	}

	return entries
}

func int64Ptr(v int64) *int64 {
	return &v
}
