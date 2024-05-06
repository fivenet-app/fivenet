package calendar

import (
	"context"
	"errors"
	"time"

	calendar "github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorscalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListCalendarEntries(ctx context.Context, req *ListCalendarEntriesRequest) (*ListCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(
		tCalendarEntry.DeletedAt.IS_NULL(),
		jet.OR(
			jet.OR(
				tCalendar.Public.IS_TRUE(),
				tCalendarEntry.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			),
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

	condition = condition.AND(tCalendarEntry.StartTime.GT_EQ(jet.DateTime(int(req.Year), time.Month(req.Month), 1, 0, 0, 0)))

	resp := &ListCalendarEntriesResponse{}

	if len(req.CalendarIds) == 0 {
		return resp, nil
	}

	ids := []jet.Expression{}
	for i := 0; i < len(req.CalendarIds); i++ {
		if req.CalendarIds[i] == 0 {
			continue
		}
		ids = append(ids, jet.Uint64(req.CalendarIds[i]))
	}

	condition = condition.AND(tCalendarEntry.CalendarID.IN(ids...))

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
			tCalendarEntry.RsvpOpen,
			tCalendarEntry.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
		).
		FROM(tCalendarEntry.
			INNER_JOIN(tCalendar,
				tCalendar.ID.EQ(tCalendarEntry.CalendarID).
					AND(tCalendar.DeletedAt.IS_NULL()),
			).
			LEFT_JOIN(tCUserAccess,
				jet.OR(
					tCUserAccess.CalendarID.EQ(tCalendarEntry.CalendarID).
						AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
					tCUserAccess.EntryID.EQ(tCalendarEntry.ID).
						AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				),
			).
			LEFT_JOIN(tCJobAccess,
				jet.OR(
					tCJobAccess.CalendarID.EQ(tCalendarEntry.CalendarID).
						AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
					tCJobAccess.EntryID.EQ(tCalendarEntry.ID).
						AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendarEntry.ID).
		WHERE(condition).
		LIMIT(100)

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

func (s *Server) GetCalendarEntry(ctx context.Context, req *GetCalendarEntryRequest) (*GetCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user has access to existing calendar
	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, req.CalendarId, req.EntryId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	access, err := s.getAccess(ctx, entry.CalendarId, &entry.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	entry.Access = access

	return &GetCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) CreateOrUpdateCalendarEntry(ctx context.Context, req *CreateOrUpdateCalendarEntryRequest) (*CreateOrUpdateCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CalendarService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateCalendarEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	tCalendarEntry := table.FivenetCalendarEntries
	if req.Entry.Id > 0 {
		check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, req.Entry.CalendarId, req.Entry.Id, userInfo, calendar.AccessLevel_ACCESS_LEVEL_EDIT, false)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		if !check {
			return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
		}

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
				tCalendarEntry.RsvpOpen,
			).
			SET(
				tCalendarEntry.Title.SET(jet.String(req.Entry.Title)),
				tCalendarEntry.Content.SET(jet.String(req.Entry.Content)),
				tCalendarEntry.StartTime.SET(startTime),
				tCalendarEntry.EndTime.SET(endTime),
				tCalendarEntry.RsvpOpen.SET(jet.Bool(*req.Entry.RsvpOpen)),
			).
			WHERE(jet.AND(
				tCalendarEntry.ID.EQ(jet.Uint64(req.Entry.Id)),
				tCalendarEntry.CalendarID.EQ(jet.Uint64(req.Entry.CalendarId)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	} else {
		check, err := s.checkIfUserHasAccessToCalendar(ctx, req.Entry.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_EDIT, false)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		if !check {
			return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
		}

		stmt := tCalendarEntry.
			INSERT(
				tCalendarEntry.CalendarID,
				tCalendarEntry.Job,
				tCalendarEntry.StartTime,
				tCalendarEntry.EndTime,
				tCalendarEntry.Title,
				tCalendarEntry.Content,
				tCalendarEntry.RsvpOpen,
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
				req.Entry.RsvpOpen,
				userInfo.UserId,
				userInfo.Job,
			)

		res, err := stmt.ExecContext(ctx, s.db)
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

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.AS("calendar_entry").ID.EQ(jet.Uint64(req.Entry.Id)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &CreateOrUpdateCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteCalendarEntry(ctx context.Context, req *DeleteCalendarEntryRequest) (*DeleteCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CalendarService_ServiceDesc.ServiceName,
		Method:  "DeleteCalendarEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, req.CalendarId, req.EntryId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_MANAGE, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	stmt := tCalendarEntry.
		UPDATE(
			tCalendarEntry.DeletedAt,
		).
		SET(
			tCalendarEntry.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tCalendarEntry.CalendarID.EQ(jet.Uint64(req.CalendarId)),
			tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteCalendarEntryResponse{}, nil
}

func (s *Server) ShareCalendarEntry(ctx context.Context, req *ShareCalendarEntryRequest) (*ShareCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CalendarService_ServiceDesc.ServiceName,
		Method:  "ShareCalendarEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, entry.CalendarId, req.EntryId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_SHARE, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	if err := s.handleCalendarAccessChanges(ctx, s.db, calendar.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED, entry.CalendarId, nil, req.Access); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	access, err := s.getAccess(ctx, entry.CalendarId, &entry.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &ShareCalendarEntryResponse{
		Access: access,
	}, nil
}

func (s *Server) getEntry(ctx context.Context, userInfo *userinfo.UserInfo, condition jet.BoolExpression) (*calendar.CalendarEntry, error) {
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
			tCalendarEntry.RsvpOpen,
			tCalendarEntry.CreatorID,
			tCalendarEntry.CreatorJob,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
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
