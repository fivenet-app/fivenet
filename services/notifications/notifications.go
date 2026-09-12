package notifications

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	notificationsstore "github.com/fivenet-app/fivenet/v2026/stores/notifications"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedRequest = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.notifications.NotificationsService.ErrFailedRequest"},
		nil,
	)
	ErrFailedStream = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.notifications.NotificationsService.ErrFailedStream"},
		nil,
	)
	ErrInvalidPreference = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.notifications.NotificationsService.ErrInvalidPreference"},
		nil,
	)
	ErrEmptyPreference = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.notifications.NotificationsService.ErrEmptyPreference"},
		nil,
	)
)

func (s *Server) GetNotifications(
	ctx context.Context,
	req *pbnotifications.GetNotificationsRequest,
) (*pbnotifications.GetNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	query := notificationsstore.ListQuery{
		UserID:       userInfo.GetUserId(),
		UnreadOnly:   req.IncludeRead != nil && !req.GetIncludeRead(),
		Categories:   req.GetCategories(),
		Kinds:        req.GetKinds(),
		Archived:     req.GetIncludeArchived(),
		ArchivedOnly: req.GetArchivedOnly(),
		Starred:      req.GetStarredOnly(),
	}

	count, err := s.store.Count(ctx, query)
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	pag, limit := req.GetPagination().GetResponse(count)
	resp := &pbnotifications.GetNotificationsResponse{
		Pagination: pag,
	}
	if count <= 0 {
		return resp, nil
	}

	query.Offset = req.GetPagination().GetOffset()
	query.Limit = limit
	resp.Notifications, err = s.store.List(ctx, query)
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}
	if err := s.hydrateNotifications(ctx, userInfo, resp.Notifications); err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	return resp, nil
}

// hydrateNotifications resolves the actor into the existing caused_by payload
// used by notification clients. actor_user_id is preferred; caused_by.user_id
// keeps notifications written before the metadata migration useful as well.
func (s *Server) hydrateNotifications(
	ctx context.Context,
	userInfo *pbuserinfo.UserInfo,
	notifications []*resourcesnotifications.Notification,
) error {
	targets := make([]citizenshydrator.BasicTarget, 0, len(notifications))
	for i, notification := range notifications {
		if notification == nil {
			continue
		}

		actorUserID := notification.GetActorUserId()
		if actorUserID <= 0 {
			actorUserID = notification.GetData().GetCausedBy().GetUserId()
		}
		if actorUserID <= 0 {
			continue
		}

		targets = append(targets, citizenshydrator.BasicTarget{
			UserID: actorUserID,
			Set: func(user *usershort.UserShort) {
				if notifications[i].Data == nil {
					notifications[i].Data = &resourcesnotifications.Data{}
				}
				notifications[i].Data.CausedBy = user
			},
		})
	}

	return s.hydrator.HydrateBasicTargetsSafeFunc(userInfo)(ctx, nil, targets)
}

func (s *Server) MarkNotifications(
	ctx context.Context,
	req *pbnotifications.MarkNotificationsRequest,
) (*pbnotifications.MarkNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	updated, err := s.store.MarkNotifications(ctx, notificationsstore.MarkQuery{
		UserID: userInfo.GetUserId(),
		IDs:    req.GetIds(),
		All:    req.GetAll(),
		Unread: req.GetUnread(),
	})
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	unreadCount, err := s.store.CountUnread(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	if updated > 0 {
		s.js.PublishProto(
			ctx,
			fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.GetUserId()),
			&notificationsevents.UserEvent{
				Data: &notificationsevents.UserEvent_NotificationsUnreadCount{
					NotificationsUnreadCount: unreadCount,
				},
			},
		)
	}

	return &pbnotifications.MarkNotificationsResponse{
		Updated:     updated,
		UnreadCount: unreadCount,
	}, nil
}

func (s *Server) UpdateNotificationState(
	ctx context.Context,
	req *pbnotifications.UpdateNotificationStateRequest,
) (*pbnotifications.UpdateNotificationStateResponse, error) {
	if !req.HasUnread() && !req.HasStarred() && !req.HasArchived() {
		return nil, ErrFailedRequest
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	updated, err := s.store.UpdateNotificationState(ctx, notificationsstore.StateQuery{
		UserID:   userInfo.GetUserId(),
		IDs:      req.GetIds(),
		Unread:   req.Unread,
		Starred:  req.Starred,
		Archived: req.Archived,
	})
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	unreadCount, err := s.store.CountUnread(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	if updated > 0 {
		s.js.PublishProto(
			ctx,
			fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.GetUserId()),
			&notificationsevents.UserEvent{
				Data: &notificationsevents.UserEvent_NotificationsUnreadCount{
					NotificationsUnreadCount: unreadCount,
				},
			},
		)
	}

	return &pbnotifications.UpdateNotificationStateResponse{
		Updated:     updated,
		UnreadCount: unreadCount,
	}, nil
}

func (s *Server) GetNotificationPreferences(
	ctx context.Context,
	_ *pbnotifications.GetNotificationPreferencesRequest,
) (*pbnotifications.GetNotificationPreferencesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	preferences, err := s.store.ListPreferences(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	return &pbnotifications.GetNotificationPreferencesResponse{Preferences: preferences}, nil
}

func (s *Server) UpdateNotificationPreference(
	ctx context.Context,
	req *pbnotifications.UpdateNotificationPreferenceRequest,
) (*pbnotifications.UpdateNotificationPreferenceResponse, error) {
	preference := req.GetPreference()
	if preference == nil || !validPreferenceScope(preference) {
		return nil, ErrInvalidPreference
	}
	if !req.GetReset() && !hasDeliveryPreferenceOverride(preference) {
		return nil, ErrEmptyPreference
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	var err error
	if req.GetReset() {
		err = s.store.DeletePreference(
			ctx,
			userInfo.GetUserId(),
			preference.GetCategory(),
			preference.GetKind(),
		)
	} else {
		err = s.store.UpsertPreference(ctx, userInfo.GetUserId(), preference)
	}
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	preferences, err := s.store.ListPreferences(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	return &pbnotifications.UpdateNotificationPreferenceResponse{Preferences: preferences}, nil
}

func validPreferenceScope(preference *resourcesnotifications.NotificationPreference) bool {
	if preference == nil {
		return false
	}

	switch preference.GetCategory() {
	case resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_UNSPECIFIED,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_QUALIFICATIONS,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_MAILER,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DISPATCH,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_SYSTEM:
	default:
		return false
	}

	switch preference.GetKind() {
	case resourcesnotifications.NotificationKind_NOTIFICATION_KIND_UNSPECIFIED:
		return true
	case resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_ADDED,
		resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_REMOVED,
		resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_MEMBER_ADDED,
		resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_MEMBER_REMOVED:
		return preference.GetCategory() == resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS
	case resourcesnotifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_APPROVAL_ASSIGNED,
		resourcesnotifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_REQUEST_CREATED,
		resourcesnotifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_REQUEST_DECIDED,
		resourcesnotifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_REQUEST_CANCELLED:
		return preference.GetCategory() == resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT
	default:
		return false
	}
}

func hasDeliveryPreferenceOverride(preference *resourcesnotifications.NotificationPreference) bool {
	return preference != nil &&
		(preference.InboxEnabled != nil || preference.ToastEnabled != nil || preference.SoundEnabled != nil)
}
