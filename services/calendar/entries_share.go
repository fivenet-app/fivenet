package calendar

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2025/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ShareCalendarEntry(
	ctx context.Context,
	req *pbcalendar.ShareCalendarEntryRequest,
) (*pbcalendar.ShareCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "ShareCalendarEntry",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

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
		calendar.AccessLevel_ACCESS_LEVEL_SHARE,
		false,
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

	req.UserIds = utils.RemoveSliceDuplicates(req.GetUserIds())

	resp := &pbcalendar.ShareCalendarEntryResponse{}
	if len(req.GetUserIds()) == 0 {
		return resp, nil
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	newUsers, err := s.shareCalendarEntry(ctx, tx, req.GetEntryId(), req.GetUserIds())
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if len(newUsers) > 0 {
		if err := s.sendShareNotifications(ctx, userInfo.GetUserId(), entry, newUsers); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return resp, nil
}

func (s *Server) shareCalendarEntry(
	ctx context.Context,
	tx qrm.DB,
	entryId int64,
	inUserIds []int32,
) ([]int32, error) {
	userIds := make([]mysql.Expression, len(inUserIds))
	for i := range inUserIds {
		userIds[i] = mysql.Int32(inUserIds[i])
	}

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.UserID,
		).
		FROM(tCalendarRSVP).
		WHERE(mysql.AND(
			tCalendarRSVP.EntryID.EQ(mysql.Int64(entryId)),
			tCalendarRSVP.UserID.IN(userIds...),
		))

	var currentRSVPs []*calendar.CalendarEntryRSVP
	if err := stmt.QueryContext(ctx, tx, &currentRSVPs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	newUsers := []int32{}
	if len(currentRSVPs) == 0 {
		newUsers = append(newUsers, inUserIds...)
	} else {
		for _, rsvp := range currentRSVPs {
			if !slices.Contains(inUserIds, rsvp.GetUserId()) {
				newUsers = append(newUsers, rsvp.GetUserId())
			}
		}
	}

	tCalendarRSVP := table.FivenetCalendarRsvp
	insertStmt := tCalendarRSVP.
		INSERT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		)
	for i := range userIds {
		insertStmt = insertStmt.VALUES(
			entryId,
			userIds[i],
			calendar.RsvpResponses_RSVP_RESPONSES_INVITED,
		)
	}

	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	return newUsers, nil
}

func (s *Server) sendShareNotifications(
	ctx context.Context,
	sourceUserId int32,
	entry *calendar.CalendarEntry,
	targetCitizens []int32,
) error {
	tUsers := tables.User().AS("user_short")

	stmt := tUsers.
		SELECT(
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.PhoneNumber,
		).
		FROM(
			tUsers,
		).
		WHERE(
			tUsers.ID.EQ(mysql.Int32(sourceUserId)),
		).
		LIMIT(1)

	sourceUser := &users.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, sourceUser); err != nil {
		return err
	}

	for _, newUser := range targetCitizens {
		if err := s.notif.NotifyUser(ctx, &notifications.Notification{
			UserId: newUser,
			Title: &common.I18NItem{
				Key: "notifications.calendar.entry_shared_with_you.title",
				Parameters: map[string]string{
					"title": entry.GetTitle(),
					"name":  fmt.Sprintf("%s %s", sourceUser.GetFirstname(), sourceUser.GetLastname()),
				},
			},
			Content: &common.I18NItem{
				Key:        "notifications.calendar.entry_shared_with_you.content",
				Parameters: map[string]string{"title": entry.GetTitle()},
			},
			Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR,
			Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
			Data: &notifications.Data{
				Link: &notifications.Link{
					To: fmt.Sprintf("/calendar?entry_id=%d", entry.GetId()),
				},
				CausedBy: &users.UserShort{
					UserId:      sourceUserId,
					Firstname:   sourceUser.GetFirstname(),
					Lastname:    sourceUser.GetLastname(),
					PhoneNumber: sourceUser.PhoneNumber,
				},
				Calendar: &notifications.CalendarData{
					CalendarEntryId: &entry.Id,
				},
			},
		}); err != nil {
			return errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	return nil
}
