package calendar

import (
	"context"
	"errors"
	"time"

	calendar "github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pbcalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/services/calendar/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jinzhu/now"
)

func (s *Server) ListCalendarEntries(ctx context.Context, req *pbcalendar.ListCalendarEntriesRequest) (*pbcalendar.ListCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	rsvpResponse := calendar.RsvpResponses_RSVP_RESPONSES_HIDDEN
	if req.ShowHidden != nil && *req.ShowHidden {
		rsvpResponse = calendar.RsvpResponses_RSVP_RESPONSES_UNSPECIFIED
	}

	condition := jet.AND(
		tCalendarEntry.DeletedAt.IS_NULL(),
		jet.OR(
			tCalendar.ID.IN(
				tCalendarSubs.
					SELECT(
						tCalendarSubs.CalendarID,
					).
					FROM(tCalendarSubs).
					WHERE(jet.AND(
						tCalendarSubs.UserID.EQ(jet.Int32(userInfo.UserId)),
					)),
			),
			tCalendarEntry.ID.IN(
				tCalendarRSVP.
					SELECT(
						tCalendarRSVP.EntryID,
					).
					FROM(tCalendarRSVP).
					WHERE(jet.AND(
						tCalendarRSVP.UserID.EQ(jet.Int32(userInfo.UserId)),
						tCalendarRSVP.Response.GT(jet.Int16(int16(rsvpResponse))),
					)),
			),
			tCalendarEntry.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			jet.OR(
				jet.AND(
					tCUserAccess.Access.IS_NOT_NULL(),
					tCUserAccess.Access.GT(jet.Int32(int32(calendar.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
				jet.AND(
					tCUserAccess.Access.IS_NULL(),
					tCJobAccess.Access.IS_NOT_NULL(),
					tCJobAccess.Access.GT(jet.Int32(int32(calendar.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		),
	)

	if req.After != nil {
		condition = condition.AND(tCalendar.UpdatedAt.GT_EQ(jet.TimestampT(req.After.AsTime())))
	}

	baseDate := now.New(time.Date(int(req.Year), time.Month(req.Month), 1, 0, 0, 0, 0, time.Local))
	startDate := baseDate.BeginningOfMonth()
	endDate := baseDate.EndOfMonth()

	condition = condition.AND(tCalendarEntry.StartTime.GT_EQ(jet.DateTimeT(startDate))).
		AND(tCalendarEntry.StartTime.LT(jet.DateTimeT(endDate)))

	resp := &pbcalendar.ListCalendarEntriesResponse{}

	if len(req.CalendarIds) > 0 {
		ids := []jet.Expression{}
		for i := 0; i < len(req.CalendarIds); i++ {
			if req.CalendarIds[i] == 0 {
				continue
			}
			ids = append(ids, jet.Uint64(req.CalendarIds[i]))
		}

		condition = condition.AND(tCalendarEntry.CalendarID.IN(ids...))
	}

	stmt := s.listCalendarEntriesQuery(condition, userInfo)

	if req.After != nil {
		stmt.ORDER_BY(tCalendar.UpdatedAt.GT_EQ(jet.TimestampT(req.After.AsTime())))
	}

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Entries); i++ {
		if resp.Entries[i].Creator != nil {
			jobInfoFn(resp.Entries[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) GetUpcomingEntries(ctx context.Context, req *pbcalendar.GetUpcomingEntriesRequest) (*pbcalendar.GetUpcomingEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcalendar.GetUpcomingEntriesResponse{
		Entries: []*calendar.CalendarEntry{},
	}

	condition := jet.AND(
		tCalendarEntry.DeletedAt.IS_NULL(),
		jet.OR(
			tCalendarEntry.ID.IN(
				tCalendarRSVP.
					SELECT(
						tCalendarRSVP.EntryID,
					).
					FROM(tCalendarRSVP).
					WHERE(jet.AND(
						tCalendarRSVP.UserID.EQ(jet.Int32(userInfo.UserId)),
						// RSVP responses: Maybe and Yes
						tCalendarRSVP.Response.GT(jet.Int16(int16(calendar.RsvpResponses_RSVP_RESPONSES_NO))),
					)),
			),
			tCalendarEntry.CreatorID.EQ(jet.Int32(userInfo.UserId)),
		),
		tCalendarEntry.StartTime.LT_EQ(
			// Now plus X seconds
			jet.CURRENT_TIMESTAMP().ADD(jet.INTERVALd(time.Duration(req.Seconds)*time.Second)),
		),
		tCalendarEntry.StartTime.GT_EQ(
			// Now minus 1 minute
			jet.CURRENT_TIMESTAMP().SUB(jet.INTERVALd(1*time.Minute)),
		),
	)

	stmt := s.listCalendarEntriesQuery(condition, userInfo)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) GetCalendarEntry(ctx context.Context, req *pbcalendar.GetCalendarEntryRequest) (*pbcalendar.GetCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	// Check if user has access to existing calendar
	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, entry.CalendarId, req.EntryId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	calAccess, err := s.getAccess(ctx, entry.CalendarId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	entry.Calendar.Access = calAccess

	return &pbcalendar.GetCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) CreateOrUpdateCalendarEntry(ctx context.Context, req *pbcalendar.CreateOrUpdateCalendarEntryRequest) (*pbcalendar.CreateOrUpdateCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateCalendarEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToCalendar(ctx, req.Entry.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_EDIT, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	calendar, err := s.getCalendar(ctx, userInfo, tCalendar.ID.EQ(jet.Uint64(req.Entry.CalendarId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if calendar == nil || calendar.Closed {
		return nil, errorscalendar.ErrCalendarClosed
	}

	req.Entry.CreatorId = &userInfo.UserId

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tCalendarEntry := table.FivenetCalendarEntries
	if req.Entry.Id > 0 {
		startTime := jet.TimestampExp(jet.NULL)
		if req.Entry.StartTime != nil {
			startTime = jet.TimestampT(req.Entry.StartTime.AsTime())
		}
		endTime := jet.TimestampExp(jet.NULL)
		if req.Entry.EndTime != nil {
			endTime = jet.TimestampT(req.Entry.EndTime.AsTime())
		}

		stmt := tCalendarEntry.
			UPDATE(
				tCalendarEntry.Title,
				tCalendarEntry.Content,
				tCalendarEntry.StartTime,
				tCalendarEntry.EndTime,
				tCalendarEntry.Closed,
				tCalendarEntry.RsvpOpen,
				tCalendarEntry.Recurring,
			).
			SET(
				req.Entry.Title,
				req.Entry.Content,
				startTime,
				endTime,
				req.Entry.Closed,
				req.Entry.RsvpOpen,
				req.Entry.Recurring,
			).
			WHERE(jet.AND(
				tCalendarEntry.ID.EQ(jet.Uint64(req.Entry.Id)),
				tCalendarEntry.CalendarID.EQ(jet.Uint64(req.Entry.CalendarId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	} else {
		stmt := tCalendarEntry.
			INSERT(
				tCalendarEntry.CalendarID,
				tCalendarEntry.Job,
				tCalendarEntry.StartTime,
				tCalendarEntry.EndTime,
				tCalendarEntry.Title,
				tCalendarEntry.Content,
				tCalendarEntry.Closed,
				tCalendarEntry.RsvpOpen,
				tCalendarEntry.Recurring,
				tCalendarEntry.CreatorID,
				tCalendarEntry.CreatorJob,
			).
			VALUES(
				req.Entry.CalendarId,
				userInfo.Job,
				req.Entry.StartTime,
				req.Entry.EndTime,
				req.Entry.Title,
				req.Entry.Content,
				req.Entry.Closed,
				req.Entry.RsvpOpen,
				req.Entry.Recurring,
				userInfo.UserId,
				userInfo.Job,
			)

		res, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		req.Entry.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	}

	newUsers := []int32{}
	if len(req.UserIds) > 0 {
		newUsers, err = s.shareCalendarEntry(ctx, tx, req.Entry.Id, req.UserIds)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.AS("calendar_entry").ID.EQ(jet.Uint64(req.Entry.Id)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if len(newUsers) > 0 {
		if err := s.sendShareNotifications(ctx, userInfo.UserId, entry, newUsers); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	return &pbcalendar.CreateOrUpdateCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteCalendarEntry(ctx context.Context, req *pbcalendar.DeleteCalendarEntryRequest) (*pbcalendar.DeleteCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "DeleteCalendarEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.checkIfUserHasAccessToCalendar(ctx, entry.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_MANAGE, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if entry.DeletedAt != nil && userInfo.SuperUser {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	stmt := tCalendarEntry.
		UPDATE(
			tCalendarEntry.DeletedAt,
		).
		SET(
			tCalendarEntry.DeletedAt.SET(deletedAtTime),
		).
		WHERE(jet.AND(
			tCalendarEntry.CalendarID.EQ(jet.Uint64(entry.CalendarId)),
			tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbcalendar.DeleteCalendarEntryResponse{}, nil
}

func (s *Server) getEntry(ctx context.Context, userInfo *userinfo.UserInfo, condition jet.BoolExpression) (*calendar.CalendarEntry, error) {
	tCreator := tables.Users().AS("creator")

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarEntry.CreatedAt,
			tCalendarEntry.UpdatedAt,
			tCalendarEntry.DeletedAt,
			tCalendarEntry.CalendarID,
			tCalendar.ID,
			tCalendar.Name,
			tCalendar.Color,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendarEntry.Job,
			tCalendarEntry.StartTime,
			tCalendarEntry.EndTime,
			tCalendarEntry.Title,
			tCalendarEntry.Content,
			tCalendarEntry.Closed,
			tCalendarEntry.RsvpOpen,
			tCalendarEntry.CreatorID,
			tCalendarEntry.CreatorJob,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
			tCalendarEntry.Recurring,
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		).
		FROM(tCalendarEntry.
			INNER_JOIN(tCalendar,
				tCalendar.ID.EQ(tCalendarEntry.CalendarID).
					AND(tCalendar.DeletedAt.IS_NULL()),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCalendarEntry.CreatorID),
			).
			LEFT_JOIN(tCalendarRSVP,
				tCalendarRSVP.UserID.EQ(jet.Int32(userInfo.UserId)).
					AND(tCalendarRSVP.EntryID.EQ(tCalendarEntry.ID)),
			),
		).
		GROUP_BY(tCalendarEntry.ID).
		WHERE(condition).
		LIMIT(1)

	dest := &calendar.CalendarEntry{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	if dest.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, dest.Creator)
	}

	return dest, nil
}
