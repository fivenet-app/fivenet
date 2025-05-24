package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToCalendarEntry(ctx context.Context, calendarId uint64, entryId uint64, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool) (bool, error) {
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

func (s *Server) checkIfUserHasAccessToCalendarEntryIDs(ctx context.Context, userInfo *userinfo.UserInfo, publicOk bool, entryIds ...uint64) ([]uint64, error) {
	var dest []uint64
	if len(entryIds) == 0 {
		return dest, nil
	}

	// Allow superusers access to any docs
	if userInfo.Superuser {
		dest = append(dest, entryIds...)
		return dest, nil
	}

	ids := make([]jet.Expression, len(entryIds))
	for i := range entryIds {
		ids[i] = jet.Uint64(entryIds[i])
	}

	condition := jet.Bool(false)
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
				tCalendarRSVP.UserID.EQ(jet.Int32(userInfo.UserId)),
			).
			INNER_JOIN(tCalendar,
				tCalendar.ID.EQ(tCalendarEntry.CalendarID).
					AND(tCalendar.DeletedAt.IS_NULL()),
			),
		).
		GROUP_BY(tCalendarEntry.ID).
		WHERE(jet.AND(
			tCalendarEntry.DeletedAt.IS_NULL(),
			tCalendarRSVP.EntryID.IN(ids...),
			jet.OR(
				jet.AND(
					tCalendarEntry.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					tCalendarEntry.CreatorJob.EQ(jet.String(userInfo.Job)),
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
