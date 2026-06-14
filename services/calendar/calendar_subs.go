package calendar

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
)

func (s *Server) SubscribeToCalendar(
	ctx context.Context,
	req *pbcalendar.SubscribeToCalendarRequest,
) (*pbcalendar.SubscribeToCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	req.Sub.UserId = userInfo.GetUserId()

	// Check if user has access to existing calendar
	check, err := s.store.CheckIfUserHasAccessToCalendar(
		ctx,
		req.GetSub().GetCalendarId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	if err := s.store.SetSubscription(
		ctx,
		req.GetSub().GetCalendarId(),
		userInfo.GetUserId(),
		!req.GetDelete(),
		true,
		req.GetSub().GetMuted(),
	); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	sub, err := s.store.GetCalendarSub(ctx, userInfo.GetUserId(), req.GetSub().GetCalendarId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbcalendar.SubscribeToCalendarResponse{
		Sub: sub,
	}, nil
}
