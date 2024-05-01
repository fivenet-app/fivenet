package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToCalendarEntry(ctx context.Context, calendarId uint64, entryId uint64, userInfo *userinfo.UserInfo, access calendar.AccessLevel) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarEntryIDs(ctx, userInfo, access, entryId)
	if err != nil {
		return false, err
	}

	if len(out) > 0 {
		return true, nil
	}

	check, err := s.checkIfUserHasAccessToCalendar(ctx, calendarId, userInfo, access)
	if err != nil {
		return false, err
	}

	return check, err
}

func (s *Server) checkIfUserHasAccessToCalendarEntries(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, entryIds ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarEntryIDs(ctx, userInfo, access, entryIds...)
	return len(out) == len(entryIds), err
}

func (s *Server) checkIfUserHasAccessToCalendarEntryIDs(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, entryIds ...uint64) ([]uint64, error) {
	if len(entryIds) == 0 {
		return entryIds, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		return entryIds, nil
	}

	ids := make([]jet.Expression, len(entryIds))
	for i := 0; i < len(entryIds); i++ {
		ids[i] = jet.Uint64(entryIds[i])
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

	var dest struct {
		IDs []uint64 `alias:"calendar_entry.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}
