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
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ShareCalendarEntry(ctx context.Context, req *pbcalendar.ShareCalendarEntryRequest) (*pbcalendar.ShareCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcalendar.CalendarService_ServiceDesc.ServiceName,
		Method:  "ShareCalendarEntry",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.checkIfUserHasAccessToCalendar(ctx, entry.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_SHARE, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	if entry.Closed {
		return nil, errorscalendar.ErrEntryClosed
	}

	req.UserIds = utils.RemoveSliceDuplicates(req.UserIds)

	resp := &pbcalendar.ShareCalendarEntryResponse{}
	if len(req.UserIds) == 0 {
		return resp, nil
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	newUsers, err := s.shareCalendarEntry(ctx, tx, req.EntryId, req.UserIds)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if len(newUsers) > 0 {
		if err := s.sendShareNotifications(ctx, userInfo.UserId, entry, newUsers); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return resp, nil
}

func (s *Server) shareCalendarEntry(ctx context.Context, tx qrm.DB, entryId uint64, inUserIds []int32) ([]int32, error) {
	userIds := make([]jet.Expression, len(inUserIds))
	for i := range inUserIds {
		userIds[i] = jet.Int32(inUserIds[i])
	}

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.UserID,
		).
		FROM(tCalendarRSVP).
		WHERE(jet.AND(
			tCalendarRSVP.EntryID.EQ(jet.Uint64(entryId)),
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
			if !slices.Contains(inUserIds, rsvp.UserId) {
				newUsers = append(newUsers, rsvp.UserId)
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

func (s *Server) sendShareNotifications(ctx context.Context, sourceUserId int32, entry *calendar.CalendarEntry, targetCitizens []int32) error {
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
			tUsers.ID.EQ(jet.Int32(sourceUserId)),
		).
		LIMIT(1)

	sourceUser := &users.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, sourceUser); err != nil {
		return err
	}

	for _, newUser := range targetCitizens {
		if err := s.notif.NotifyUser(ctx, &notifications.Notification{
			UserId: newUser,
			Title: &common.TranslateItem{
				Key: "notifications.calendar.entry_shared_with_you.title",
				Parameters: map[string]string{
					"title": entry.Title,
					"name":  fmt.Sprintf("%s %s", sourceUser.Firstname, sourceUser.Lastname),
				},
			},
			Content: &common.TranslateItem{
				Key:        "notifications.calendar.entry_shared_with_you.content",
				Parameters: map[string]string{"title": entry.Title},
			},
			Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR,
			Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
			Data: &notifications.Data{
				Link: &notifications.Link{
					To: fmt.Sprintf("/calendar?entry_id=%d", entry.Id),
				},
				CausedBy: &users.UserShort{
					UserId:      sourceUserId,
					Firstname:   sourceUser.Firstname,
					Lastname:    sourceUser.Lastname,
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
