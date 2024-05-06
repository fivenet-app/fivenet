package calendar

import (
	"context"

	calendar "github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	errorscalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) createOrDeleteSubscription(ctx context.Context, calendarId uint64, entryId *uint64, userId int32, subscribe bool, confirmed bool, muted bool) error {
	if subscribe {
		tCalendarSubs := table.FivenetCalendarSubs
		stmt := tCalendarSubs.
			INSERT(
				tCalendarSubs.CalendarID,
				tCalendarSubs.EntryID,
				tCalendarSubs.UserID,
				tCalendarSubs.Confirmed,
				tCalendarSubs.Muted,
			).
			VALUES(
				calendarId,
				entryId,
				userId,
				confirmed,
				muted,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCalendarSubs.Confirmed.SET(jet.Bool(confirmed)),
				tCalendarSubs.Muted.SET(jet.Bool(muted)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	} else {
		entryIdColumn := jet.IntExp(jet.NULL)
		if entryId != nil {
			entryIdColumn = jet.Uint64(*entryId)
		}

		stmt := tCalendarSubs.
			DELETE().
			WHERE(jet.AND(
				tCalendarSubs.CalendarID.EQ(jet.Uint64(calendarId)),
				tCalendarSubs.EntryID.EQ(entryIdColumn),
				tCalendarSubs.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) SubscribeToCalendar(ctx context.Context, req *SubscribeToCalendarRequest) (*SubscribeToCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user has access to existing calendar
	var check bool
	var err error
	if req.Sub.EntryId == nil {
		check, err = s.checkIfUserHasAccessToCalendar(ctx, req.Sub.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	} else {
		check, err = s.checkIfUserHasAccessToCalendarEntry(ctx, req.Sub.CalendarId, *req.Sub.EntryId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	if err := s.createOrDeleteSubscription(ctx, req.Sub.CalendarId, req.Sub.EntryId, userInfo.UserId, req.Delete, true, req.Sub.Muted); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &SubscribeToCalendarResponse{}, nil
}
