package calendar

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
)

const maxCalendarEntriesLimit = int64(125)

func (s *Server) ListCalendarEntries(
	ctx context.Context,
	req *pbcalendar.ListCalendarEntriesRequest,
) (*pbcalendar.ListCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	entries, err := s.store.ListCalendarEntries(ctx, userInfo, req)
	if err != nil {
		return nil, err
	}

	return &pbcalendar.ListCalendarEntriesResponse{
		Entries: s.finalizeCalendarEntries(
			entries,
			userInfo,
		),
	}, nil
}

func (s *Server) GetUpcomingEntries(
	ctx context.Context,
	req *pbcalendar.GetUpcomingEntriesRequest,
) (*pbcalendar.GetUpcomingEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcalendar.GetUpcomingEntriesResponse{
		Entries: []*calendarentries.CalendarEntry{},
	}
	entries, err := s.store.GetUpcomingEntries(ctx, userInfo, req)
	if err != nil {
		return nil, err
	}

	resp.Entries = s.store.FilterUpcomingCalendarEntries(
		s.finalizeCalendarEntries(
			entries,
			userInfo,
		),
		userInfo,
	)
	return resp, nil
}

func (s *Server) GetCalendarEntry(
	ctx context.Context,
	req *pbcalendar.GetCalendarEntryRequest,
) (*pbcalendar.GetCalendarEntryResponse, error) {
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
		return nil, errorscalendar.ErrNoPerms
	}
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() !=
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	// Check if user has access to existing calendar
	check, err := s.store.CheckIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		req.GetEntryId(),
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

	if entry.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, entry.GetCreator())
	}

	calAccess, err := s.store.ListTargetAccess(
		ctx,
		entry.GetCalendarId(),
		calendarSubjectAccessOptions,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	entry.Calendar.Access = calAccess

	return &pbcalendar.GetCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) CreateOrUpdateCalendarEntry(
	ctx context.Context,
	req *pbcalendar.CreateOrUpdateCalendarEntryRequest,
) (*pbcalendar.CreateOrUpdateCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.store.CheckIfUserHasAccessToCalendar(
		ctx,
		req.GetEntry().GetCalendarId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	calendar, err := s.getCalendar(
		ctx,
		userInfo,
		tCalendar.ID.EQ(mysql.Int64(req.GetEntry().GetCalendarId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if calendar == nil || calendar.GetClosed() {
		return nil, errorscalendar.ErrCalendarClosed
	}
	if calendar.GetSystemKind() != calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	defer tx.Rollback()

	tCalendarEntry := table.FivenetCalendarEntries
	if req.GetEntry().GetId() > 0 {
		oldEntry, err := s.store.GetEntry(
			ctx,
			userInfo,
			tCalendarEntry.AS("calendar_entry").ID.EQ(mysql.Int64(req.GetEntry().GetId())),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		if oldEntry == nil {
			return nil, errorscalendar.ErrNoPerms
		}
		lastID, err := s.store.UpsertCalendarEntry(ctx, tx, req.GetEntry(), oldEntry, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		req.Entry.Id = lastID

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	} else {
		req.GetEntry().CreatorId = &userInfo.UserId
		lastID, err := s.store.UpsertCalendarEntry(ctx, tx, req.GetEntry(), nil, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		req.Entry.Id = lastID

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	}

	newUsers := []int32{}
	if len(req.GetUserIds()) > 0 {
		newUsers, err = s.store.ShareCalendarEntry(
			ctx,
			tx,
			req.GetEntry().GetId(),
			req.GetUserIds(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	entry, err := s.store.GetEntry(
		ctx,
		userInfo,
		tCalendarEntry.AS("calendar_entry").ID.EQ(mysql.Int64(req.GetEntry().GetId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if len(newUsers) > 0 {
		if err := s.sendShareNotifications(ctx, userInfo.GetUserId(), entry, newUsers); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	if entry.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, entry.GetCreator())
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.CreateOrUpdateCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteCalendarEntry(
	ctx context.Context,
	req *pbcalendar.DeleteCalendarEntryRequest,
) (*pbcalendar.DeleteCalendarEntryResponse, error) {
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
		return nil, errorscalendar.ErrNoPerms
	}
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() !=
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.store.CheckIfUserHasAccessToCalendar(
		ctx,
		entry.GetCalendarId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	var deletedAtTime *timestamp.Timestamp
	if entry.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteCalendarEntry(
		ctx,
		s.db,
		req.GetEntryId(),
		entry.GetCalendarId(),
		deletedAtTime,
	); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.DeleteCalendarEntryResponse{}, nil
}
