package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToCalendarEntry(
	ctx context.Context,
	calendarId int64,
	entryId int64,
	userInfo *userinfo.UserInfo,
	access calendar.AccessLevel,
	publicOk bool,
) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarEntryIDs(ctx, userInfo, publicOk, entryId)
	if err != nil {
		return false, err
	}

	if len(out) > 0 {
		return true, nil
	}

	check, err := s.checkIfUserHasAccessToCalendar(ctx, calendarId, userInfo, access, publicOk)
	if err != nil {
		return false, err
	}

	return check, err
}

func (s *Server) checkIfUserHasAccessToCalendarEntryIDs(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	publicOk bool,
	entryIds ...int64,
) ([]int64, error) {
	var dest []int64
	if len(entryIds) == 0 {
		return dest, nil
	}

	// Allow superusers access to any docs
	if userInfo.GetSuperuser() {
		dest = append(dest, entryIds...)
		return dest, nil
	}

	ids := make([]mysql.Expression, len(entryIds))
	for i := range entryIds {
		ids[i] = mysql.Int64(entryIds[i])
	}

	condition := mysql.Bool(false)
	if publicOk {
		condition = tCalendar.Public.IS_TRUE()
	}

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarRSVP.EntryID,
		).
		FROM(tCalendarEntry.
			LEFT_JOIN(tCalendarRSVP,
				tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			).
			INNER_JOIN(tCalendar,
				mysql.AND(
					tCalendar.ID.EQ(tCalendarEntry.CalendarID),
					tCalendar.DeletedAt.IS_NULL(),
				),
			),
		).
		WHERE(mysql.AND(
			tCalendarEntry.DeletedAt.IS_NULL(),
			tCalendarRSVP.EntryID.IN(ids...),
			mysql.OR(
				mysql.AND(
					tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarEntry.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
				),
				tCalendarRSVP.EntryID.IS_NOT_NULL(),
				condition,
			),
		)).
		ORDER_BY(tCalendarEntry.ID.DESC())

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
