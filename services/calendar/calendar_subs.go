package calendar

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	calendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2025/services/calendar/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) SubscribeToCalendar(
	ctx context.Context,
	req *pbcalendar.SubscribeToCalendarRequest,
) (*pbcalendar.SubscribeToCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "SubscribeToCalendar",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	req.Sub.UserId = userInfo.GetUserId()

	// Check if user has access to existing calendar
	check, err := s.checkIfUserHasAccessToCalendar(
		ctx,
		req.GetSub().GetCalendarId(),
		userInfo,
		calendar.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	if err := s.createOrDeleteSubscription(ctx, req.GetSub().GetCalendarId(), userInfo.GetUserId(), !req.GetDelete(), true, req.GetSub().GetMuted()); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	sub, err := s.getCalendarSub(ctx, userInfo.GetUserId(), req.GetSub().GetCalendarId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbcalendar.SubscribeToCalendarResponse{
		Sub: sub,
	}, nil
}

func (s *Server) createOrDeleteSubscription(
	ctx context.Context,
	calendarId int64,
	userId int32,
	subscribe bool,
	confirmed bool,
	muted bool,
) error {
	tCalendarSubs := table.FivenetCalendarSubs

	if subscribe {
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
				tCalendarSubs.CalendarID.EQ(jet.Int64(calendarId)),
				tCalendarSubs.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) getCalendarSub(
	ctx context.Context,
	userId int32,
	calendarId int64,
) (*calendar.CalendarSub, error) {
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
			tCalendarSubs.CalendarID.EQ(jet.Int64(calendarId)),
			tCalendarSubs.UserID.EQ(jet.Int32(userId)),
		))

	var dest calendar.CalendarSub
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetCalendarId() == 0 || dest.GetUserId() == 0 {
		return nil, nil
	}

	return &dest, nil
}
