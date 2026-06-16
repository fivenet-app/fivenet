package calendar

import (
	"context"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	calendarstore "github.com/fivenet-app/fivenet/v2026/stores/calendar"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ListCalendarEntryRSVP(
	ctx context.Context,
	req *pbcalendar.ListCalendarEntryRSVPRequest,
) (*pbcalendar.ListCalendarEntryRSVPResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.store.GetEntry(
		ctx,
		userInfo,
		tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrFailedQuery
	}
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() !=
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.store.CheckIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		entry.GetId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	resp, err := s.store.ListCalendarEntryRSVP(ctx, calendarstore.ListCalendarEntryRSVPOptions{
		EntryID:    entry.GetId(),
		Pagination: req.GetPagination(),
	}, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetEntries() {
		if resp.GetEntries()[i].GetUser() != nil {
			jobInfoFn(resp.GetEntries()[i].GetUser())
		}
	}

	return resp, nil
}

func (s *Server) RSVPCalendarEntry(
	ctx context.Context,
	req *pbcalendar.RSVPCalendarEntryRequest,
) (*pbcalendar.RSVPCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.GetEntry() == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	entry, err := s.store.GetEntry(
		ctx,
		userInfo,
		tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntry().GetEntryId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrFailedQuery
	}
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() !=
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.store.CheckIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		entry.GetId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	if entry.GetClosed() {
		return nil, errorscalendar.ErrEntryClosed
	}

	occurrenceKey := strings.TrimSpace(req.GetOccurrenceKey())
	if occurrenceKey == "" && req.GetEntry() != nil {
		occurrenceKey = strings.TrimSpace(req.GetEntry().GetOccurrenceKey())
	}

	if occurrenceKey != "" {
		if entry.GetRecurring() == nil {
			return nil, errorscalendar.ErrNoPerms
		}
		if err := s.store.ValidateRecurringOccurrenceKey(entry, occurrenceKey); err != nil {
			return nil, err
		}
	}

	if req.Remove != nil && req.GetRemove() && occurrenceKey == "" {
		req.Entry.Response = calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN
	}

	if err := s.store.SetCalendarEntryRSVP(
		ctx,
		req.GetEntry(),
		userInfo,
		occurrenceKey,
		req.GetRemove(),
	); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	rsvpEntry, err := s.store.GetRSVPCalendarEntry(
		ctx,
		req.GetEntry().GetEntryId(),
		userInfo.GetUserId(),
		occurrenceKey,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if rsvpEntry.GetUser() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, rsvpEntry.GetUser())
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbcalendar.RSVPCalendarEntryResponse{
		Entry: rsvpEntry,
	}, nil
}
