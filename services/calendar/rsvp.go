package calendar

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2025/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListCalendarEntryRSVP(
	ctx context.Context,
	req *pbcalendar.ListCalendarEntryRSVPRequest,
) (*pbcalendar.ListCalendarEntryRSVPResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrFailedQuery
	}

	check, err := s.checkIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		entry.GetId(),
		userInfo,
		calendar.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	tUser := tables.User().AS("user_short")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	condition := tCalendarRSVP.EntryID.EQ(mysql.Int64(entry.GetId())).
		AND(tCalendarRSVP.Response.GT(mysql.Int32(int32(calendar.RsvpResponses_RSVP_RESPONSES_HIDDEN))))

	countStmt := tCalendarRSVP.
		SELECT(
			mysql.COUNT(tCalendarRSVP.UserID).AS("data_count.total"),
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

	pag, limit := req.GetPagination().GetResponse(count.Total)
	resp := &pbcalendar.ListCalendarEntryRSVPResponse{
		Pagination: pag,
	}

	if count.Total <= 0 {
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
			tUserProps.AvatarFileID.AS("user_short.profile_picture_file_id"),
			tAvatar.FilePath.AS("user_short.profile_picture"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(condition).
		ORDER_BY(tCalendarRSVP.Response.DESC()).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetEntries() {
		if resp.GetEntries()[i].GetUser() != nil {
			jobInfoFn(resp.GetEntries()[i].GetUser())
		}
	}

	resp.GetPagination().Update(len(resp.GetEntries()))

	return resp, nil
}

func (s *Server) RSVPCalendarEntry(
	ctx context.Context,
	req *pbcalendar.RSVPCalendarEntryRequest,
) (*pbcalendar.RSVPCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "RSVPCalendarEntry",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getEntry(
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

	check, err := s.checkIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		entry.GetId(),
		userInfo,
		calendar.AccessLevel_ACCESS_LEVEL_VIEW,
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

	if req.Remove != nil && req.GetRemove() {
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
			req.GetEntry().GetEntryId(),
			userInfo.GetUserId(),
			req.GetEntry().GetResponse(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCalendarRSVP.Response.SET(mysql.Int32(int32(req.GetEntry().GetResponse()))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	rsvpEntry, err := s.getRSVPCalendarEntry(ctx, req.GetEntry().GetEntryId(), userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if rsvpEntry.GetUser() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, rsvpEntry.GetUser())
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcalendar.RSVPCalendarEntryResponse{
		Entry: rsvpEntry,
	}, nil
}

func (s *Server) getRSVPCalendarEntry(
	ctx context.Context,
	entryId int64,
	userId int32,
) (*calendar.CalendarEntryRSVP, error) {
	tUser := tables.User().AS("user_short")
	tAvatar := table.FivenetFiles.AS("profile_picture")

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
			tUserProps.AvatarFileID.AS("user_short.profile_picture_file_id"),
			tAvatar.FilePath.AS("user_short.profile_picture"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(mysql.AND(
			tCalendarRSVP.EntryID.EQ(mysql.Int64(entryId)),
			tCalendarRSVP.UserID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	var dest calendar.CalendarEntryRSVP
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetEntryId() == 0 {
		return nil, nil
	}

	return &dest, nil
}
