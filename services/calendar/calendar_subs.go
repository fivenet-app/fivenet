package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2025/services/calendar/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) SubscribeToCalendar(ctx context.Context, req *pbcalendar.SubscribeToCalendarRequest) (*pbcalendar.SubscribeToCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "SubscribeToCalendar",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	req.Sub.UserId = userInfo.UserId

	// Check if user has access to existing calendar
	check, err := s.checkIfUserHasAccessToCalendar(ctx, req.Sub.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	if err := s.createOrDeleteSubscription(ctx, req.Sub.CalendarId, userInfo.UserId, !req.Delete, true, req.Sub.Muted); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	sub, err := s.getCalendarSub(ctx, userInfo.UserId, req.Sub.CalendarId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &pbcalendar.SubscribeToCalendarResponse{
		Sub: sub,
	}, nil
}

func (s *Server) createOrDeleteSubscription(ctx context.Context, calendarId uint64, userId int32, subscribe bool, confirmed bool, muted bool) error {
	if subscribe {
		tCalendarSubs := table.FivenetCalendarSubs
		stmt := tCalendarSubs.
			INSERT(
				tCalendarSubs.CalendarID,
				tCalendarSubs.UserID,
				tCalendarSubs.Confirmed,
				tCalendarSubs.Muted,
			).
			VALUES(
				calendarId,
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
		stmt := tCalendarSubs.
			DELETE().
			WHERE(jet.AND(
				tCalendarSubs.CalendarID.EQ(jet.Uint64(calendarId)),
				tCalendarSubs.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) getCalendarSub(ctx context.Context, userId int32, calendarId uint64) (*calendar.CalendarSub, error) {
	stmt := tCalendarSubs.
		SELECT(
			tCalendarSubs.CalendarID,
			tCalendarSubs.UserID,
			tCalendarSubs.CreatedAt,
			tCalendarSubs.Confirmed,
			tCalendarSubs.Muted,
		).
		FROM(tCalendarSubs).
		WHERE(jet.AND(
			tCalendarSubs.CalendarID.EQ(jet.Uint64(calendarId)),
			tCalendarSubs.UserID.EQ(jet.Int32(userId)),
		))

	var dest calendar.CalendarSub
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.CalendarId == 0 || dest.UserId == 0 {
		return nil, nil
	}

	return &dest, nil
}
