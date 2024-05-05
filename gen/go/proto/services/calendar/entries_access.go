package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToCalendarEntry(ctx context.Context, calendarId uint64, entryId uint64, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarEntryIDs(ctx, userInfo, access, publicOk, entryId)
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

func (s *Server) checkIfUserHasAccessToCalendarEntries(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool, entryIds ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarEntryIDs(ctx, userInfo, access, publicOk, entryIds...)
	return len(out) == len(entryIds), err
}

func (s *Server) checkIfUserHasAccessToCalendarEntryIDs(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool, entryIds ...uint64) ([]uint64, error) {
	var dest []uint64
	if len(entryIds) == 0 {
		return dest, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		for i := 0; i < len(entryIds); i++ {
			dest = append(dest, entryIds[i])
		}
		return dest, nil
	}

	ids := make([]jet.Expression, len(entryIds))
	for i := 0; i < len(entryIds); i++ {
		ids[i] = jet.Uint64(entryIds[i])
	}

	condition := jet.Bool(false)
	if publicOk {
		condition = tCalendarEntry.Public.IS_TRUE()
	}

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
		).
		FROM(tCalendarEntry.
			LEFT_JOIN(tCUserAccess,
				tCUserAccess.CalendarID.EQ(tCalendarEntry.CalendarID).
					AND(tCUserAccess.EntryID.EQ(tCalendarEntry.ID)).
					AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tCJobAccess,
				tCJobAccess.CalendarID.EQ(tCalendarEntry.CalendarID).
					AND(tCJobAccess.EntryID.EQ(tCalendarEntry.ID)).
					AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendarEntry.ID).
		WHERE(jet.AND(
			tCalendarEntry.ID.IN(ids...),
			tCalendarEntry.DeletedAt.IS_NULL(),
			jet.OR(
				condition,
				tCalendarEntry.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				tCalendarEntry.CreatorJob.EQ(jet.String(userInfo.Job)),
				jet.AND(
					tCUserAccess.Access.IS_NOT_NULL(),
					tCUserAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
				jet.AND(
					tCUserAccess.Access.IS_NULL(),
					tCJobAccess.Access.IS_NOT_NULL(),
					tCJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
			),
		)).
		ORDER_BY(tCalendarEntry.ID.DESC(), tCJobAccess.MinimumGrade)

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
