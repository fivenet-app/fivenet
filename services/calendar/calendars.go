package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar"
	permscalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2025/services/calendar/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListCalendars(ctx context.Context, req *pbcalendar.ListCalendarsRequest) (*pbcalendar.ListCalendarsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tCreator := tables.Users().AS("creator")

	subsCondition := tCalendar.ID.IN(tCalendarSubs.
		SELECT(
			tCalendarSubs.CalendarID,
		).
		FROM(tCalendarSubs).
		WHERE(jet.AND(
			tCalendarSubs.UserID.EQ(jet.Int32(userInfo.UserId)),
		)),
	)

	minAccessLevel := calendar.AccessLevel_ACCESS_LEVEL_VIEW
	if req.MinAccessLevel != nil {
		minAccessLevel = *req.MinAccessLevel
		subsCondition = jet.Bool(false)
	}

	condition := jet.AND(
		tCalendar.DeletedAt.IS_NULL(),
		jet.OR(
			jet.OR(
				subsCondition,
				tCalendar.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			),
			jet.OR(
				tCAccess.UserID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tCAccess.Job.EQ(jet.String(userInfo.Job)),
					tCAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade)),
				),
			),
		),
	)

	if req.OnlyPublic {
		condition = jet.AND(
			tCalendar.DeletedAt.IS_NULL(),
			tCalendar.Public.IS_TRUE(),
		)
	}

	if req.After != nil {
		condition = condition.AND(tCalendar.UpdatedAt.GT_EQ(jet.TimestampT(req.After.AsTime())))
	}

	countStmt := tCalendar.
		SELECT(
			jet.COUNT(jet.DISTINCT(tCalendar.ID)).AS("datacount.totalcount"),
		).
		FROM(tCalendar.
			LEFT_JOIN(tCAccess,
				tCAccess.TargetID.EQ(tCalendar.ID).
					AND(tCAccess.Access.GT_EQ(jet.Int32(int32(minAccessLevel)))),
			).
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &pbcalendar.ListCalendarsResponse{
		Pagination: pag,
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
			tCalendar.CreatedAt,
			tCalendar.UpdatedAt,
			tCalendar.DeletedAt,
			tCalendar.Job,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
			tCalendarSubs.CalendarID,
			tCalendarSubs.UserID,
			tCalendarSubs.CreatedAt,
			tCalendarSubs.Confirmed,
			tCalendarSubs.Muted,
			tCAccess.ID.AS("user_access.id"),
			tCAccess.TargetID.AS("user_access.target_id"),
			tCAccess.UserID.AS("user_access.user_id"),
			tCAccess.Access.AS("user_access.access"),
			tCAccess.ID.AS("job_access.id"),
			tCAccess.TargetID.AS("job_access.target_id"),
			tCAccess.Job.AS("job_access.job"),
			tCAccess.MinimumGrade.AS("job_access.minimum_grade"),
			tCAccess.Access.AS("job_access.access"),
		).
		FROM(tCalendar.
			LEFT_JOIN(tCAccess,
				tCAccess.TargetID.EQ(tCalendar.ID).
					AND(tCAccess.Access.GT_EQ(jet.Int32(int32(minAccessLevel)))),
			).
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tCalendarSubs,
				tCalendarSubs.CalendarID.EQ(tCalendar.ID).
					AND(tCalendarSubs.UserID.EQ(jet.Int32(userInfo.UserId))),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if req.After != nil {
		stmt.ORDER_BY(tCalendar.UpdatedAt.GT_EQ(jet.TimestampT(req.After.AsTime())))
	}

	if err := stmt.QueryContext(ctx, s.db, &resp.Calendars); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Calendars {
		if resp.Calendars[i].Creator != nil {
			jobInfoFn(resp.Calendars[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Calendars))

	return resp, nil
}

func (s *Server) GetCalendar(ctx context.Context, req *pbcalendar.GetCalendarRequest) (*pbcalendar.GetCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user has access to existing calendar
	check, err := s.checkIfUserHasAccessToCalendar(ctx, req.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	calendar, err := s.getCalendar(ctx, userInfo, tCalendar.ID.EQ(jet.Uint64(req.CalendarId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if calendar == nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	access, err := s.getAccess(ctx, calendar.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	for i := range access.Jobs {
		s.enricher.EnrichJobInfo(access.Jobs[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range access.Users {
		if access.Users[i].User != nil {
			jobInfoFn(access.Users[i].User)
		}
	}

	calendar.Access = access

	return &pbcalendar.GetCalendarResponse{
		Calendar: calendar,
	}, nil
}

func (s *Server) getAccess(ctx context.Context, calendarId uint64) (*calendar.CalendarAccess, error) {
	jobAccess, err := s.access.Jobs.List(ctx, s.db, calendarId)
	if err != nil {
		return nil, err
	}

	userAccess, err := s.access.Users.List(ctx, s.db, calendarId)
	if err != nil {
		return nil, err
	}

	return &calendar.CalendarAccess{
		Jobs:  jobAccess,
		Users: userAccess,
	}, nil
}

func (s *Server) CreateCalendar(ctx context.Context, req *pbcalendar.CreateCalendarRequest) (*pbcalendar.CreateCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "CreateCalendar",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	fields, err := s.p.AttrStringList(userInfo, permscalendar.CalendarServicePerm, permscalendar.CalendarServiceCreateCalendarPerm, permscalendar.CalendarServiceCreateCalendarFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if req.Calendar.Job != nil && !fields.Contains("Job") {
		return nil, errorscalendar.ErrFailedQuery
	}
	if req.Calendar.Color == "" {
		req.Calendar.Color = "blue"
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// Check if user has access to existing calendar
	if req.Calendar.Id > 0 {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	// Allow only one private calendar per user (job field will be null for private calendars)
	if req.Calendar.Job == nil {
		calendar, err := s.getCalendar(ctx, userInfo, jet.AND(
			tCalendar.DeletedAt.IS_NULL(),
			tCalendar.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			tCalendar.Job.IS_NULL(),
		))
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		if calendar != nil {
			return nil, errorscalendar.ErrOnePrivateCal
		}
	} else {
		req.Calendar.Job = &userInfo.Job
	}

	tCalendar := table.FivenetCalendar
	stmt := tCalendar.
		INSERT(
			tCalendar.Job,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.CreatorID,
			tCalendar.CreatorJob,
		).
		VALUES(
			req.Calendar.Job,
			req.Calendar.Name,
			req.Calendar.Description,
			req.Calendar.Public,
			req.Calendar.Closed,
			req.Calendar.Color,
			userInfo.UserId,
			userInfo.Job,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCalendar.Name.SET(jet.String(req.Calendar.Name)),
			tCalendar.Description.SET(jet.String("VALUES(`description`)")),
			tCalendar.Public.SET(jet.Bool(req.Calendar.Public)),
			tCalendar.Closed.SET(jet.Bool(req.Calendar.Closed)),
			tCalendar.Color.SET(jet.String(req.Calendar.Color)),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if req.Calendar.Id == 0 {
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		req.Calendar.Id = uint64(lastId)
	}

	if _, err := s.access.HandleAccessChanges(ctx, tx, req.Calendar.Id, req.Calendar.Access.Jobs, req.Calendar.Access.Users, nil); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	calendar, err := s.getCalendar(ctx, userInfo, tCalendar.AS("calendar").ID.EQ(jet.Uint64(req.Calendar.Id)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.CreateCalendarResponse{
		Calendar: calendar,
	}, nil
}

func (s *Server) UpdateCalendar(ctx context.Context, req *pbcalendar.UpdateCalendarRequest) (*pbcalendar.UpdateCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "UpdateCalendar",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	fields, err := s.p.AttrStringList(userInfo, permscalendar.CalendarServicePerm, permscalendar.CalendarServiceCreateCalendarPerm, permscalendar.CalendarServiceCreateCalendarFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if req.Calendar.Job != nil && !fields.Contains("Job") {
		return nil, errorscalendar.ErrFailedQuery
	}
	if req.Calendar.Color == "" {
		req.Calendar.Color = "blue"
	}

	// Check if user has access to existing calendar
	if req.Calendar.Id == 0 {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	check, err := s.checkIfUserHasAccessToCalendar(ctx, req.Calendar.Id, userInfo, calendar.AccessLevel_ACCESS_LEVEL_MANAGE, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	calendar, err := s.getCalendar(ctx, userInfo, tCalendar.ID.EQ(jet.Uint64(req.Calendar.Id)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	req.Calendar.Job = calendar.Job

	if req.Calendar.Description == nil {
		empty := ""
		req.Calendar.Description = &empty
	}

	if !fields.Contains("Public") && calendar.Public && req.Calendar.Public {
		req.Calendar.Public = false
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tCalendar := table.FivenetCalendar
	stmt := tCalendar.
		UPDATE(
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
		).
		SET(
			tCalendar.Name.SET(jet.String(req.Calendar.Name)),
			tCalendar.Description.SET(jet.String(*req.Calendar.Description)),
			tCalendar.Public.SET(jet.Bool(req.Calendar.Public)),
			tCalendar.Closed.SET(jet.Bool(req.Calendar.Closed)),
			tCalendar.Color.SET(jet.String(req.Calendar.Color)),
		).
		WHERE(jet.AND(
			tCalendar.ID.EQ(jet.Uint64(req.Calendar.Id)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	if _, err := s.access.HandleAccessChanges(ctx, tx, req.Calendar.Id, req.Calendar.Access.Jobs, req.Calendar.Access.Users, nil); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	calendar, err = s.getCalendar(ctx, userInfo, tCalendar.AS("calendar").ID.EQ(jet.Uint64(req.Calendar.Id)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.UpdateCalendarResponse{
		Calendar: calendar,
	}, nil
}

func (s *Server) DeleteCalendar(ctx context.Context, req *pbcalendar.DeleteCalendarRequest) (*pbcalendar.DeleteCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "DeleteCalendar",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToCalendar(ctx, req.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_MANAGE, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	calendar, err := s.getCalendar(ctx, userInfo, tCalendar.ID.EQ(jet.Uint64(req.CalendarId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if calendar == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if calendar.DeletedAt != nil && userInfo.SuperUser {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	stmt := tCalendar.
		UPDATE(
			tCalendar.DeletedAt,
		).
		SET(
			tCalendar.DeletedAt.SET(deletedAtTime),
		).
		WHERE(tCalendar.ID.EQ(jet.Uint64(req.CalendarId)))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbcalendar.DeleteCalendarResponse{}, nil
}

func (s *Server) getCalendar(ctx context.Context, userInfo *userinfo.UserInfo, condition jet.BoolExpression) (*calendar.Calendar, error) {
	tCreator := tables.Users().AS("creator")

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
			tCalendar.CreatedAt,
			tCalendar.UpdatedAt,
			tCalendar.DeletedAt,
			tCalendar.Job,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.CreatorID,
			tCalendar.CreatorJob,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
			tCalendarSubs.CalendarID,
			tCalendarSubs.UserID,
			tCalendarSubs.CreatedAt,
			tCalendarSubs.Confirmed,
			tCalendarSubs.Muted,
		).
		FROM(tCalendar.
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCalendar.CreatorID),
			).
			LEFT_JOIN(tCalendarSubs,
				tCalendarSubs.CalendarID.EQ(tCalendar.ID).
					AND(tCalendarSubs.UserID.EQ(jet.Int32(userInfo.UserId))),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(condition).
		LIMIT(1)

	dest := &calendar.Calendar{}
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
