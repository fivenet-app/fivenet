package calendar

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pbcalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/services/calendar/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListCalendarEntryRSVP(ctx context.Context, req *pbcalendar.ListCalendarEntryRSVPRequest) (*pbcalendar.ListCalendarEntryRSVPResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrFailedQuery
	}

	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, entry.CalendarId, entry.Id, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	tUser := tables.Users().AS("user_short")

	condition := tCalendarRSVP.EntryID.EQ(jet.Uint64(entry.Id)).
		AND(tCalendarRSVP.Response.GT(jet.Int16(int16(calendar.RsvpResponses_RSVP_RESPONSES_HIDDEN))))

	countStmt := tCalendarRSVP.
		SELECT(
			jet.COUNT(tCalendarRSVP.UserID).AS("datacount.totalcount"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &pbcalendar.ListCalendarEntryRSVPResponse{
		Pagination: pag,
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tUserProps.Avatar.AS("user_short.avatar"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			),
		).
		WHERE(condition).
		ORDER_BY(tCalendarRSVP.Response.DESC()).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Entries); i++ {
		if resp.Entries[i].User != nil {
			jobInfoFn(resp.Entries[i].User)
		}
	}

	resp.Pagination.Update(len(resp.Entries))

	return resp, nil
}

func (s *Server) RSVPCalendarEntry(ctx context.Context, req *pbcalendar.RSVPCalendarEntryRequest) (*pbcalendar.RSVPCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "RSVPCalendarEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.Entry.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrFailedQuery
	}

	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, entry.CalendarId, entry.Id, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	if entry.Closed {
		return nil, errorscalendar.ErrEntryClosed
	}

	if req.Remove != nil && *req.Remove {
		req.Entry.Response = calendar.RsvpResponses_RSVP_RESPONSES_HIDDEN
	}

	tCalendarRSVP := table.FivenetCalendarRsvp
	stmt := tCalendarRSVP.
		INSERT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		).
		VALUES(
			req.Entry.EntryId,
			userInfo.UserId,
			req.Entry.Response,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCalendarRSVP.Response.SET(jet.Int16(int16(req.Entry.Response))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	rsvpEntry, err := s.getRSVPCalendarEntry(ctx, req.Entry.EntryId, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if rsvpEntry.User != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, rsvpEntry.User)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbcalendar.RSVPCalendarEntryResponse{
		Entry: rsvpEntry,
	}, nil
}

func (s *Server) getRSVPCalendarEntry(ctx context.Context, entryId uint64, userId int32) (*calendar.CalendarEntryRSVP, error) {
	tUser := tables.Users().AS("user_short")

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tUserProps.Avatar.AS("user_short.avatar"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			),
		).
		WHERE(jet.AND(
			tCalendarRSVP.EntryID.EQ(jet.Uint64(entryId)),
			tCalendarRSVP.UserID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	var dest calendar.CalendarEntryRSVP
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.EntryId == 0 {
		return nil, nil
	}

	return &dest, nil
}
