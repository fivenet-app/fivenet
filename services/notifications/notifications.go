package notifications

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
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
)

func (s *Server) GetNotifications(
	ctx context.Context,
	req *pbnotifications.GetNotificationsRequest,
) (*pbnotifications.GetNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	query := notificationsstore.ListQuery{
		UserID:     userInfo.GetUserId(),
		UnreadOnly: req.IncludeRead != nil && !req.GetIncludeRead(),
		Categories: req.GetCategories(),
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

	return resp, nil
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

	if updated > 0 {
		if req.GetUnread() {
			updated = -updated
		}

		s.js.PublishProto(
			ctx,
			fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.GetUserId()),
			&notificationsevents.UserEvent{
				Data: &notificationsevents.UserEvent_NotificationsReadCount{
					NotificationsReadCount: updated,
				},
			},
		)
	}

	return &pbnotifications.MarkNotificationsResponse{
		Updated: updated,
	}, nil
}
