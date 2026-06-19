package calendar

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	calendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	permscalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar/perms"
	permssettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	calendarstore "github.com/fivenet-app/fivenet/v2026/stores/calendar"
	"github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) canEditCalendarDiscordSettings(userInfo *userinfo.UserInfo) bool {
	return s.ps.Can(userInfo, permssettings.SettingsService.SetJobProps.Perm)
}

func (s *Server) ListCalendars(
	ctx context.Context,
	req *pbcalendar.ListCalendarsRequest,
) (*pbcalendar.ListCalendarsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	query := calendarstore.ListQuery{
		UserInfo:    userInfo,
		OnlyPublic:  req.GetOnlyPublic(),
		After:       req.GetAfter(),
		CalendarIDs: req.GetCalendarIds(),
	}
	minAccessLevel := calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW
	if req.MinAccessLevel != nil {
		minAccessLevel = req.GetMinAccessLevel()
	}
	query.MinAccessLevel = &minAccessLevel

	total, err := s.store.CountCalendars(ctx, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponse(total)
	resp := &pbcalendar.ListCalendarsResponse{
		Pagination: pag,
	}

	if total <= 0 {
		return resp, nil
	}

	resp.Calendars, err = s.store.ListCalendars(ctx, query, req.GetPagination().GetOffset(), limit)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetCalendars() {
		if resp.GetCalendars()[i].GetCreator() != nil {
			jobInfoFn(resp.GetCalendars()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) GetCalendar(
	ctx context.Context,
	req *pbcalendar.GetCalendarRequest,
) (*pbcalendar.GetCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	cal, err := s.store.GetAccessibleCalendar(
		ctx,
		req.GetCalendarId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if cal == nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	access, err := s.store.ListTargetAccess(ctx, cal.GetId(), calendarSubjectAccessOptions)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	for i := range access.GetJobs() {
		s.enricher.EnrichJobInfo(access.GetJobs()[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range access.GetUsers() {
		if access.GetUsers()[i].GetUser() != nil {
			jobInfoFn(access.GetUsers()[i].GetUser())
		}
	}

	cal.Access = access

	return &pbcalendar.GetCalendarResponse{
		Calendar: cal,
	}, nil
}

func (s *Server) CreateCalendar(
	ctx context.Context,
	req *pbcalendar.CreateCalendarRequest,
) (*pbcalendar.CreateCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	fields, err := permscalendar.CalendarService.CreateCalendar.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if req.Calendar.Job != nil &&
		!fields.Contains(permscalendar.CalendarServiceCreateCalendarFieldsPermValueJob) {
		return nil, errorscalendar.ErrFailedQuery
	}
	if req.GetCalendar().GetColor() == "" {
		req.Calendar.Color = "blue"
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	defer tx.Rollback()

	if req.GetCalendar().GetId() > 0 {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if req.GetCalendar().
		GetSystemKind() !=
		calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	if req.Calendar.Job == nil {
		calendar, err := s.store.GetCalendar(ctx, userInfo, mysql.AND(
			tCalendar.DeletedAt.IS_NULL(),
			tCalendar.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
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

	discordSettings, discordSettingsJSON, err := s.prepareCalendarDiscordSettings(
		ctx,
		req.GetCalendar(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if discordSettings != nil && !s.canEditCalendarDiscordSettings(userInfo) {
		return nil, status.Error(
			codes.PermissionDenied,
			"missing permission to configure calendar discord reminders",
		)
	}
	req.Calendar.DiscordSettings = discordSettings

	lastID, err := s.store.CreateCalendar(ctx, tx, req.GetCalendar(), userInfo, discordSettingsJSON)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	req.Calendar.Id = lastID

	calendarAccess := req.GetCalendar().GetAccess()
	if calendarAccess == nil || len(calendarAccess.GetJobs()) == 0 {
		calendarAccess = &calendaraccess.CalendarAccess{
			Jobs: []*calendaraccess.CalendarJobAccess{
				{
					TargetId:     req.GetCalendar().GetId(),
					Job:          userInfo.GetJob(),
					MinimumGrade: userInfo.GetJobGrade(),
					Access:       int32(calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
				},
			},
		}
	}
	fallbackAccess := &calendaraccess.CalendarAccess{
		Jobs: []*calendaraccess.CalendarJobAccess{{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
		}},
	}
	normalizedAccess, err := access.NormalizeAccess(calendarAccess, nil, fallbackAccess, 15)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if _, err := s.access.ReplaceTargetAccess(
		ctx,
		tx,
		s.accessResolver,
		req.GetCalendar().GetId(),
		normalizedAccess,
		calendarSubjectAccessOptions,
	); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	calendar, err := s.getCalendar(
		ctx,
		userInfo,
		tCalendar.AS("calendar").ID.EQ(mysql.Int64(req.GetCalendar().GetId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.CreateCalendarResponse{
		Calendar: calendar,
	}, nil
}

func (s *Server) UpdateCalendar(
	ctx context.Context,
	req *pbcalendar.UpdateCalendarRequest,
) (*pbcalendar.UpdateCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.GetCalendar().GetId() == 0 {
		return nil, errorscalendar.ErrFailedQuery
	}

	currentCalendar, err := s.store.GetAccessibleCalendar(
		ctx,
		req.GetCalendar().GetId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if currentCalendar == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	isBirthdayCalendar := currentCalendar.GetSystemKind() == calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS
	if !isBirthdayCalendar &&
		currentCalendar.GetSystemKind() != calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	if req.GetCalendar().GetColor() == "" {
		if isBirthdayCalendar && currentCalendar.GetColor() != "" {
			req.Calendar.Color = currentCalendar.GetColor()
		} else {
			req.Calendar.Color = "blue"
		}
	}

	fields, err := permscalendar.CalendarService.CreateCalendar.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if !isBirthdayCalendar && req.Calendar.Job != nil &&
		!fields.Contains(permscalendar.CalendarServiceCreateCalendarFieldsPermValueJob) {
		return nil, errorscalendar.ErrFailedQuery
	}

	if isBirthdayCalendar {
		req.Calendar.Job = currentCalendar.Job
		req.Calendar.Name = currentCalendar.Name
		req.Calendar.Description = currentCalendar.Description
		req.Calendar.Public = currentCalendar.Public
		req.Calendar.Closed = currentCalendar.Closed
		if currentCalendar.Color != "" && req.GetCalendar().GetColor() == "" {
			req.Calendar.Color = currentCalendar.Color
		}
	}

	if !isBirthdayCalendar &&
		!fields.Contains(permscalendar.CalendarServiceCreateCalendarFieldsPermValuePublic) &&
		currentCalendar.GetPublic() &&
		req.GetCalendar().GetPublic() {
		req.Calendar.Public = false
	}

	if !s.canEditCalendarDiscordSettings(userInfo) {
		req.Calendar.DiscordSettings = currentCalendar.GetDiscordSettings()
	}

	discordSettings, discordSettingsJSON, err := s.prepareCalendarDiscordSettings(
		ctx,
		req.GetCalendar(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	req.Calendar.DiscordSettings = discordSettings

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.UpdateCalendar(ctx, tx, req.GetCalendar(), discordSettingsJSON); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	if !isBirthdayCalendar {
		calendarAccess := req.GetCalendar().GetAccess()
		if calendarAccess == nil || len(calendarAccess.GetJobs()) == 0 {
			calendarAccess = &calendaraccess.CalendarAccess{
				Jobs: []*calendaraccess.CalendarJobAccess{
					{
						TargetId:     req.GetCalendar().GetId(),
						Job:          userInfo.GetJob(),
						MinimumGrade: userInfo.GetJobGrade(),
						Access:       int32(calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
					},
				},
			}
		}
		fallbackAccess := &calendaraccess.CalendarAccess{
			Jobs: []*calendaraccess.CalendarJobAccess{{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
			}},
		}
		normalizedAccess, err := access.NormalizeAccess(calendarAccess, nil, fallbackAccess, 15)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		if _, err := s.access.ReplaceTargetAccess(
			ctx,
			tx,
			s.accessResolver,
			req.GetCalendar().GetId(),
			normalizedAccess,
			calendarSubjectAccessOptions,
		); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	calendar, err := s.getCalendar(
		ctx,
		userInfo,
		tCalendar.AS("calendar").ID.EQ(mysql.Int64(req.GetCalendar().GetId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.UpdateCalendarResponse{
		Calendar: calendar,
	}, nil
}

func (s *Server) DeleteCalendar(
	ctx context.Context,
	req *pbcalendar.DeleteCalendarRequest,
) (*pbcalendar.DeleteCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	cal, err := s.store.GetAccessibleCalendar(
		ctx,
		req.GetCalendarId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if cal == nil {
		return nil, errorscalendar.ErrNoPerms
	}
	if cal.GetSystemKind() != calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	var deletedAtTime *timestamp.Timestamp
	if cal.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteCalendar(ctx, s.db, req.GetCalendarId(), deletedAtTime); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.DeleteCalendarResponse{}, nil
}

func (s *Server) getCalendar(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
) (*calendar.Calendar, error) {
	dest, err := s.store.GetCalendar(ctx, userInfo, condition)
	if err != nil {
		return nil, err
	}

	if dest != nil && dest.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, dest.GetCreator())
	}

	return dest, nil
}
