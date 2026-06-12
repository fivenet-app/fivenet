package calendar

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jinzhu/now"
)

const maxCalendarEntriesLimit = int64(125)

func (s *Server) ListCalendarEntries(
	ctx context.Context,
	req *pbcalendar.ListCalendarEntriesRequest,
) (*pbcalendar.ListCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	rsvpResponse := calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN
	if req.ShowHidden != nil && req.GetShowHidden() {
		rsvpResponse = calendarentries.RsvpResponses_RSVP_RESPONSES_UNSPECIFIED
	}

	condition := mysql.AND(
		tCalendarEntry.DeletedAt.IS_NULL(),
		mysql.OR(
			// Allow access to user's calendar subscriptions
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
			// Allow access to invited entries
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
			// Allow access to invited recurring entries
			tCalendarEntry.ID.IN(
				tCalendarRSVPOccurrence.
					SELECT(tCalendarRSVPOccurrence.EntryID).
					FROM(tCalendarRSVPOccurrence).
					WHERE(mysql.AND(
						tCalendarRSVPOccurrence.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tCalendarRSVPOccurrence.Response.GT(mysql.Int32(int32(rsvpResponse))),
						tCalendarRSVPOccurrence.RecurrenceVersion.EQ(
							tCalendarEntry.RecurrenceVersion,
						),
					)),
			),
			// Allow birthday calendar entries access
			s.birthdayCalendarVisible(
				tCalendarEntry.CalendarID,
				calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
				userInfo,
			),
			// Allow entries from calendars the user can view directly
			calendarEntryVisibility(
				userInfo,
				calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
				rsvpResponse,
			),
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

	condition = condition.AND(tCalendarEntry.StartTime.LT_EQ(mysql.DateTimeT(endDate)))

	dateWindowCondition := mysql.OR(
		// Non-recurring entries overlapping the requested range.
		mysql.AND(
			tCalendarEntry.Recurring.IS_NULL(),
			mysql.OR(
				tCalendarEntry.EndTime.IS_NULL().
					AND(tCalendarEntry.StartTime.GT_EQ(mysql.DateTimeT(startDate))),
				tCalendarEntry.EndTime.GT_EQ(mysql.DateTimeT(startDate)),
			),
		),

		// Recurring entries that started before range end and have not ended before range start.
		mysql.AND(
			tCalendarEntry.Recurring.IS_NOT_NULL(),
			mysql.OR(
				// No recurring-until means it may still produce occurrences.
				tCalendarEntry.RecurringUntil.IS_NULL(),

				// Series still active in this range.
				tCalendarEntry.RecurringUntil.GT_EQ(mysql.DateTimeT(startDate)),
			),
		),
	)

	regularCondition := condition.
		AND(dateWindowCondition).
		AND(mysql.OR(
			tCalendar.SystemKind.IS_NULL(),
			tCalendar.SystemKind.NOT_EQ(
				mysql.Int32(
					int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
				),
			),
		))

	birthdayCondition := condition.AND(
		tCalendar.SystemKind.EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
	)

	if len(req.GetCalendarIds()) > 0 {
		ids := []mysql.Expression{}
		for i := range req.GetCalendarIds() {
			if req.GetCalendarIds()[i] == 0 {
				continue
			}
			ids = append(ids, mysql.Int64(req.GetCalendarIds()[i]))
		}

		regularCondition = regularCondition.AND(tCalendarEntry.CalendarID.IN(ids...))
		birthdayCondition = birthdayCondition.AND(tCalendarEntry.CalendarID.IN(ids...))
	}

	regularEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		regularCondition,
		calendarEntryVisibility(
			userInfo,
			calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
			rsvpResponse,
		),
		startDate,
		endDate,
		new(maxCalendarEntriesLimit),
	)
	if err != nil {
		return nil, err
	}

	birthdayEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		birthdayCondition,
		s.birthdayCalendarVisible(
			tCalendarEntry.CalendarID,
			calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
			userInfo,
		),
		startDate,
		endDate,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &pbcalendar.ListCalendarEntriesResponse{
		Entries: s.finalizeCalendarEntries(
			append(regularEntries, birthdayEntries...),
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

	rangeStart := time.Now().Add(-1 * time.Minute)
	rangeEnd := time.Now().Add(time.Duration(req.GetSeconds()) * time.Second)

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
							mysql.Int32(int32(calendarentries.RsvpResponses_RSVP_RESPONSES_NO)),
						),
					)),
			),
			tCalendarEntry.ID.IN(
				tCalendarRSVPOccurrence.
					SELECT(
						tCalendarRSVPOccurrence.EntryID,
					).
					FROM(tCalendarRSVPOccurrence).
					WHERE(mysql.AND(
						tCalendarRSVPOccurrence.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tCalendarRSVPOccurrence.Response.GT(
							mysql.Int32(int32(calendarentries.RsvpResponses_RSVP_RESPONSES_NO)),
						),
						tCalendarRSVPOccurrence.RecurrenceVersion.EQ(
							tCalendarEntry.RecurrenceVersion,
						),
					)),
			),
			tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
		),
		tCalendarEntry.StartTime.LT_EQ(mysql.TimestampT(rangeEnd)),
	)

	regularCondition := condition.AND(mysql.OR(
		tCalendar.SystemKind.IS_NULL(),
		tCalendar.SystemKind.NOT_EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
	))
	birthdayCondition := condition.AND(
		tCalendar.SystemKind.EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
	)

	regularEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		regularCondition,
		calendarEntryVisibility(
			userInfo,
			calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
			calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN,
		),
		rangeStart,
		rangeEnd,
		new(maxCalendarEntriesLimit),
	)
	if err != nil {
		return nil, err
	}

	birthdayEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		birthdayCondition,
		s.birthdayCalendarVisible(
			tCalendarEntry.CalendarID,
			calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
			userInfo,
		),
		rangeStart,
		rangeEnd,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp.Entries = filterUpcomingCalendarEntries(
		s.finalizeCalendarEntries(
			append(regularEntries, birthdayEntries...),
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

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())))
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
	check, err := s.checkIfUserHasAccessToCalendarEntry(
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

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tCalendarEntry := table.FivenetCalendarEntries
	if req.GetEntry().GetId() > 0 {
		oldEntry, err := s.getEntry(
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

		values := []interface{}{
			mysql.String(req.GetEntry().GetTitle()),
			req.GetEntry().GetContent(),
			dbutils.TimestampToMySQL(req.GetEntry().GetStartTime()),
			dbutils.TimestampToMySQL(req.GetEntry().GetEndTime()),
			mysql.Bool(req.GetEntry().GetClosed()),
			mysql.Bool(req.GetEntry().GetRsvpOpen()),
			req.GetEntry().GetRecurring(),
			dbutils.TimestampToMySQL(req.GetEntry().GetRecurring().GetUntil()),
		}

		if recurrenceShapeChanged(oldEntry, req.GetEntry()) {
			values = append(
				values,
				tCalendarEntry.RecurrenceVersion.SET(
					tCalendarEntry.RecurrenceVersion.ADD(mysql.Int32(1)),
				),
			)
		} else {
			values = append(values, oldEntry.RecurrenceVersion)
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
				tCalendarEntry.RecurringUntil,
				tCalendarEntry.RecurrenceVersion,
			).
			SET(
				values[0],
				values[1:]...,
			).
			WHERE(mysql.AND(
				tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntry().GetId())),
				tCalendarEntry.CalendarID.EQ(mysql.Int64(req.GetEntry().GetCalendarId())),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	} else {
		req.GetEntry().CreatorId = &userInfo.UserId

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
				tCalendarEntry.RecurringUntil,
				tCalendarEntry.RecurrenceVersion,
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
				req.GetEntry().GetRecurring().GetUntil(),
				1,
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
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() !=
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.checkIfUserHasAccessToCalendar(
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

	stmt := tCalendarEntry.
		UPDATE(
			tCalendarEntry.DeletedAt,
		).
		SET(
			tCalendarEntry.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAtTime)),
		).
		WHERE(mysql.AND(
			tCalendarEntry.CalendarID.EQ(mysql.Int64(entry.GetCalendarId())),
			tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &pbcalendar.DeleteCalendarEntryResponse{}, nil
}

func (s *Server) getEntry(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
) (*calendarentries.CalendarEntry, error) {
	tCalendarEntry := tCalendarEntry
	tCreator := table.FivenetUser.AS("creator")
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
			tCalendarEntry.Recurring,
			tCalendarEntry.RecurringUntil,
			tCalendarEntry.RecurrenceVersion,
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

	dest := &calendarentries.CalendarEntry{}
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
