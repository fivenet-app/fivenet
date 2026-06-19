package calendar

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ShareCalendarEntry(
	ctx context.Context,
	req *pbcalendar.ShareCalendarEntryRequest,
) (*pbcalendar.ShareCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.store.GetEntry(
		ctx,
		userInfo,
		tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())),
	)
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

	calendar, err := s.store.GetAccessibleCalendar(
		ctx,
		entry.GetCalendarId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_SHARE,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if calendar == nil {
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	defer tx.Rollback()

	newUsers, err := s.store.ShareCalendarEntry(ctx, tx, req.GetEntryId(), req.GetUserIds())
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if len(newUsers) > 0 {
		if err := s.sendShareNotifications(ctx, userInfo.GetUserId(), entry, newUsers); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return resp, nil
}

func (s *Server) sendShareNotifications(
	ctx context.Context,
	sourceUserId int32,
	entry *calendarentries.CalendarEntry,
	targetCitizens []int32,
) error {
	sourceUser, err := s.store.GetUserShortByID(ctx, sourceUserId)
	if err != nil {
		return err
	}

	for _, newUser := range targetCitizens {
		if err := s.notif.NotifyUser(ctx, &notifications.Notification{
			UserId: newUser,
			Title: &common.I18NItem{
				Key: "notifications.calendar.entry_shared_with_you.title",
				Parameters: map[string]string{
					"title": entry.GetTitle(),
					"name": fmt.Sprintf(
						"%s %s",
						sourceUser.GetFirstname(),
						sourceUser.GetLastname(),
					),
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
					To: fmt.Sprintf("/calendar?entryId=%d", entry.GetId()),
				},
				CausedBy: &usershort.UserShort{
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
