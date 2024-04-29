package calendar

import (
	"context"
	"errors"

	calendar "github.com/galexrt/fivenet/gen/go/proto/resources/calendar"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getAccess(ctx context.Context, calendarId uint64, entryId *uint64) (*calendar.CalendarAccess, error) {
	dest := &calendar.CalendarAccess{}

	jobCondition := tCJobAccess.CalendarID.EQ(jet.Uint64(calendarId))
	if entryId != nil {
		jobCondition = jobCondition.AND(tCJobAccess.EntryID.EQ(jet.Uint64(*entryId)))
	}

	jobStmt := tCJobAccess.
		SELECT(
			tCJobAccess.ID,
			tCJobAccess.CreatedAt,
			tCJobAccess.CalendarID,
			tCJobAccess.EntryID,
			tCJobAccess.Job,
			tCJobAccess.MinimumGrade,
			tCJobAccess.Access,
		).
		FROM(tCJobAccess).
		WHERE(jobCondition).
		ORDER_BY(tCJobAccess.ID.ASC())

	if err := jobStmt.QueryContext(ctx, s.db, &dest.Jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	userCondition := tCUserAccess.CalendarID.EQ(jet.Uint64(calendarId))
	if entryId != nil {
		userCondition = userCondition.AND(tCUserAccess.EntryID.EQ(jet.Uint64(*entryId)))
	}

	userStmt := tCUserAccess.
		SELECT(
			tCUserAccess.ID,
			tCUserAccess.CreatedAt,
			tCUserAccess.CalendarID,
			tCUserAccess.EntryID,
			tCUserAccess.UserID,
			tCUserAccess.Access,
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tCUserAccess.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tCUserAccess.UserID),
				),
		).
		WHERE(userCondition).
		ORDER_BY(tCUserAccess.ID.ASC())

	if err := userStmt.QueryContext(ctx, s.db, &dest.Jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Server) checkIfUserHasAccessToCalendar(ctx context.Context, calendarId uint64, userInfo *userinfo.UserInfo, access calendar.AccessLevel) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarIDs(ctx, userInfo, access, calendarId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToCalendars(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, calendarIds ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarIDs(ctx, userInfo, access, calendarIds...)
	return len(out) == len(calendarIds), err
}

func (s *Server) checkIfUserHasAccessToCalendarIDs(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, calendarIds ...uint64) ([]uint64, error) {
	if len(calendarIds) == 0 {
		return calendarIds, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		return calendarIds, nil
	}

	ids := make([]jet.Expression, len(calendarIds))
	for i := 0; i < len(calendarIds); i++ {
		ids[i] = jet.Uint64(calendarIds[i])
	}

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
		).
		FROM(tCalendar.
			LEFT_JOIN(tCUserAccess,
				tCUserAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCUserAccess.EntryID.IS_NULL()).
					AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tCJobAccess,
				tCJobAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCJobAccess.EntryID.IS_NULL()).
					AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(jet.AND(
			tCalendar.ID.IN(ids...),
			tCalendar.DeletedAt.IS_NULL(),
			jet.OR(
				tCalendar.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				tCalendar.CreatorJob.EQ(jet.String(userInfo.Job)),
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
		ORDER_BY(tCalendar.ID.DESC(), tCJobAccess.MinimumGrade)

	var dest struct {
		IDs []uint64 `alias:"calendar.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}
