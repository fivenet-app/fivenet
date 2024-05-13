package calendar

import (
	"context"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorscalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

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

	check, err := s.checkIfUserHasAccessToCalendar(ctx, entry.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_SHARE, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	if entry.Closed {
		return nil, errorscalendar.ErrEntryClosed
	}

	req.UserIds = utils.RemoveSliceDuplicates(req.UserIds)

	resp := &ShareCalendarEntryResponse{}
	if len(req.UserIds) == 0 {
		return resp, nil
	}

	userIds := make([]jet.Expression, len(req.UserIds))
	for i := 0; i < len(req.UserIds); i++ {
		userIds[i] = jet.Int32(req.UserIds[i])
	}

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.UserID,
		).
		FROM(tCalendarRSVP).
		WHERE(jet.AND(
			tCalendarRSVP.EntryID.EQ(jet.Uint64(req.EntryId)),
			tCalendarRSVP.UserID.IN(userIds...),
		))

	var rsvps []*calendar.CalendarEntryRSVP
	if err := stmt.QueryContext(ctx, s.db, &rsvps); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	newUsers := []int32{}
	if len(rsvps) == 0 {
		newUsers = append(newUsers, req.UserIds...)
	}
	for _, rsvp := range rsvps {
		if !slices.Contains(req.UserIds, rsvp.UserId) {
			newUsers = append(newUsers, rsvp.UserId)
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tCalendarRSVP := table.FivenetCalendarRsvp
	insertStmt := tCalendarRSVP.
		INSERT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		)

	for i := 0; i < len(req.UserIds); i++ {
		insertStmt = insertStmt.VALUES(
			req.EntryId,
			req.UserIds[i],
			calendar.RsvpResponses_RSVP_RESPONSES_INVITED,
		)
	}

	if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
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

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return resp, nil
}

func (s *Server) sendShareNotifications(ctx context.Context, sourceUserId int32, entry *calendar.CalendarEntry, targetUsers []int32) error {
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

	sourceUser := users.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, &sourceUser); err != nil {
		return err
	}

	for _, newUser := range targetUsers {
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
			},
		}); err != nil {
			return errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	return nil
}
