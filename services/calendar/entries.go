package calendar

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	calendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2025/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jinzhu/now"
)

func (s *Server) ListCalendarEntries(
	ctx context.Context,
	req *pbcalendar.ListCalendarEntriesRequest,
) (*pbcalendar.ListCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	rsvpResponse := calendar.RsvpResponses_RSVP_RESPONSES_HIDDEN
	if req.ShowHidden != nil && req.GetShowHidden() {
		rsvpResponse = calendar.RsvpResponses_RSVP_RESPONSES_UNSPECIFIED
	}

	condition := mysql.AND(
		tCalendarEntry.DeletedAt.IS_NULL(),
		mysql.OR(
			tCalendar.ID.IN(
				tCalendarSubs.
					SELECT(
						tCalendarSubs.CalendarID,
					).
					FROM(tCalendarSubs).
					WHERE(mysql.AND(
						tCalendarSubs.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					)),
			),
			tCalendarEntry.ID.IN(
				tCalendarRSVP.
					SELECT(
						tCalendarRSVP.EntryID,
					).
					FROM(tCalendarRSVP).
					WHERE(mysql.AND(
						tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tCalendarRSVP.Response.GT(mysql.Int32(int32(rsvpResponse))),
					)),
			),
			tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
		),
	)

	if req.GetAfter() != nil {
		condition = condition.AND(
			tCalendar.UpdatedAt.GT_EQ(mysql.TimestampT(req.GetAfter().AsTime())),
		)
	}

	baseDate := now.New(
		time.Date(int(req.GetYear()), time.Month(req.GetMonth()), 1, 0, 0, 0, 0, time.UTC),
	)
	startDate := baseDate.BeginningOfMonth()
	endDate := baseDate.EndOfMonth()

	condition = condition.
		AND(mysql.AND(
			tCalendarEntry.StartTime.GT_EQ(mysql.DateTimeT(startDate)),
			tCalendarEntry.StartTime.LT(mysql.DateTimeT(endDate)),
		))

	resp := &pbcalendar.ListCalendarEntriesResponse{}

	if len(req.GetCalendarIds()) > 0 {
		ids := []mysql.Expression{}
		for i := range req.GetCalendarIds() {
			if req.GetCalendarIds()[i] == 0 {
				continue
			}
			ids = append(ids, mysql.Int64(req.GetCalendarIds()[i]))
		}

		condition = condition.AND(tCalendarEntry.CalendarID.IN(ids...))
	}

	stmt := s.listCalendarEntriesQuery(condition, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW)

	if req.GetAfter() != nil {
		stmt.ORDER_BY(tCalendar.UpdatedAt.GT_EQ(mysql.TimestampT(req.GetAfter().AsTime())))
	}

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetEntries() {
		if resp.GetEntries()[i].GetCreator() != nil {
			jobInfoFn(resp.GetEntries()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) GetUpcomingEntries(
	ctx context.Context,
	req *pbcalendar.GetUpcomingEntriesRequest,
) (*pbcalendar.GetUpcomingEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcalendar.GetUpcomingEntriesResponse{
		Entries: []*calendar.CalendarEntry{},
	}

	condition := mysql.AND(
		tCalendarEntry.DeletedAt.IS_NULL(),
		mysql.OR(
			tCalendarEntry.ID.IN(
				tCalendarRSVP.
					SELECT(
						tCalendarRSVP.EntryID,
					).
					FROM(tCalendarRSVP).
					WHERE(mysql.AND(
						tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						// RSVP responses: Maybe and Yes
						tCalendarRSVP.Response.GT(
							mysql.Int32(int32(calendar.RsvpResponses_RSVP_RESPONSES_NO)),
						),
					)),
			),
			tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
		),
		tCalendarEntry.StartTime.LT_EQ(
			// Now plus X seconds
			mysql.CURRENT_TIMESTAMP().
				ADD(mysql.INTERVALd(time.Duration(req.GetSeconds())*time.Second)),
		),
		tCalendarEntry.StartTime.GT_EQ(
			// Now minus 1 minute
			mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVALd(1*time.Minute)),
		),
	)

	stmt := s.listCalendarEntriesQuery(condition, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) GetCalendarEntry(
	ctx context.Context,
	req *pbcalendar.GetCalendarEntryRequest,
) (*pbcalendar.GetCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	// Check if user has access to existing calendar
	check, err := s.checkIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		req.GetEntryId(),
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

	calAccess, err := s.getAccess(ctx, entry.GetCalendarId())
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

	check, err := s.checkIfUserHasAccessToCalendar(
		ctx,
		req.GetEntry().GetCalendarId(),
		userInfo,
		calendar.AccessLevel_ACCESS_LEVEL_EDIT,
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

	req.Entry.CreatorId = &userInfo.UserId

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tCalendarEntry := table.FivenetCalendarEntries
	if req.GetEntry().GetId() > 0 {
		startTime := mysql.TimestampExp(mysql.NULL)
		if req.GetEntry().GetStartTime() != nil {
			startTime = mysql.TimestampT(req.GetEntry().GetStartTime().AsTime())
		}
		endTime := mysql.TimestampExp(mysql.NULL)
		if req.GetEntry().GetEndTime() != nil {
			endTime = mysql.TimestampT(req.GetEntry().GetEndTime().AsTime())
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
				req.GetEntry().GetTitle(),
				req.GetEntry().GetContent(),
				startTime,
				endTime,
				req.GetEntry().GetClosed(),
				req.GetEntry().GetRsvpOpen(),
				req.GetEntry().GetRecurring(),
			).
			WHERE(mysql.AND(
				tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntry().GetId())),
				tCalendarEntry.CalendarID.EQ(mysql.Int64(req.GetEntry().GetCalendarId())),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
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
				req.GetEntry().GetCalendarId(),
				userInfo.GetJob(),
				req.GetEntry().GetStartTime(),
				req.GetEntry().GetEndTime(),
				req.GetEntry().GetTitle(),
				req.GetEntry().GetContent(),
				req.GetEntry().GetClosed(),
				req.GetEntry().GetRsvpOpen(),
				req.GetEntry().GetRecurring(),
				userInfo.GetUserId(),
				userInfo.GetJob(),
			)

		res, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		req.Entry.Id = lastId

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	}

	newUsers := []int32{}
	if len(req.GetUserIds()) > 0 {
		newUsers, err = s.shareCalendarEntry(ctx, tx, req.GetEntry().GetId(), req.GetUserIds())
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	entry, err := s.getEntry(
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

	return &pbcalendar.CreateOrUpdateCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteCalendarEntry(
	ctx context.Context,
	req *pbcalendar.DeleteCalendarEntryRequest,
) (*pbcalendar.DeleteCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.checkIfUserHasAccessToCalendar(
		ctx,
		entry.GetCalendarId(),
		userInfo,
		calendar.AccessLevel_ACCESS_LEVEL_MANAGE,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	deletedAtTime := mysql.CURRENT_TIMESTAMP()
	if entry.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = mysql.TimestampExp(mysql.NULL)
	}

	stmt := tCalendarEntry.
		UPDATE(
			tCalendarEntry.DeletedAt,
		).
		SET(
			tCalendarEntry.DeletedAt.SET(deletedAtTime),
		).
		WHERE(mysql.AND(
			tCalendarEntry.CalendarID.EQ(mysql.Int64(entry.GetCalendarId())),
			tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbcalendar.DeleteCalendarEntryResponse{}, nil
}

func (s *Server) getEntry(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
) (*calendar.CalendarEntry, error) {
	tCreator := tables.User().AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")

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
			tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
			tAvatar.FilePath.AS("creator.profile_picture"),
			tCalendarEntry.Recurring,
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		).
		FROM(tCalendarEntry.
			INNER_JOIN(tCalendar,
				mysql.AND(
					tCalendar.ID.EQ(tCalendarEntry.CalendarID),
					tCalendar.DeletedAt.IS_NULL(),
				),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCalendarEntry.CreatorID),
			).
			LEFT_JOIN(tCalendarRSVP,
				mysql.AND(
					tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarRSVP.EntryID.EQ(tCalendarEntry.ID),
				),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(condition).
		LIMIT(1)

	dest := &calendar.CalendarEntry{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	if dest.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, dest.GetCreator())
	}

	return dest, nil
}
