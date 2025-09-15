package notifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbnotifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
)

var tNotifications = table.FivenetNotifications

var (
	ErrFailedRequest = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.NotificationsService.ErrFailedRequest"},
		nil,
	)
	ErrFailedStream = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.NotificationsService.ErrFailedStream"},
		nil,
	)
)

func (s *Server) GetNotifications(
	ctx context.Context,
	req *pbnotifications.GetNotificationsRequest,
) (*pbnotifications.GetNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tNotifications := tNotifications.AS("notification")
	condition := tNotifications.UserID.EQ(mysql.Int32(userInfo.GetUserId()))
	if req.IncludeRead != nil && !req.GetIncludeRead() {
		condition = condition.AND(tNotifications.ReadAt.IS_NULL())
	}

	if len(req.GetCategories()) > 0 {
		categoryIds := make([]mysql.Expression, len(req.GetCategories()))
		for i := range req.GetCategories() {
			categoryIds[i] = mysql.Int32(int32(req.GetCategories()[i]))
		}

		condition = condition.AND(tNotifications.Category.IN(categoryIds...))
	}

	countStmt := tNotifications.
		SELECT(
			mysql.COUNT(tNotifications.ID).AS("data_count.total"),
		).
		FROM(tNotifications).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	pag, limit := req.GetPagination().GetResponse(count.Total)
	resp := &pbnotifications.GetNotificationsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := tNotifications.
		SELECT(
			tNotifications.ID,
			tNotifications.CreatedAt,
			tNotifications.ReadAt,
			tNotifications.UserID,
			tNotifications.Title,
			tNotifications.Type,
			tNotifications.Content,
			tNotifications.Category,
			tNotifications.Data,
		).
		FROM(tNotifications).
		WHERE(
			condition,
		).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tNotifications.ID.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Notifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	resp.GetPagination().Update(len(resp.GetNotifications()))

	return resp, nil
}

func (s *Server) MarkNotifications(
	ctx context.Context,
	req *pbnotifications.MarkNotificationsRequest,
) (*pbnotifications.MarkNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tNotifications.UserID.EQ(
		mysql.Int32(userInfo.GetUserId())).AND(
		tNotifications.ReadAt.IS_NULL(),
	)

	// If not all
	if len(req.GetIds()) > 0 {
		ids := make([]mysql.Expression, len(req.GetIds()))
		for i := range req.GetIds() {
			ids[i] = mysql.Int64(req.GetIds()[i])
		}
		condition = condition.AND(tNotifications.ID.IN(ids...))
	} else if req.All == nil || !req.GetAll() {
		return &pbnotifications.MarkNotificationsResponse{}, nil
	}

	readAt := mysql.CURRENT_TIMESTAMP()
	if req.GetUnread() {
		// Allow users to mark notifications as unread
		readAt = mysql.TimestampExp(mysql.NULL)
	}

	stmt := tNotifications.
		UPDATE(
			tNotifications.ReadAt,
		).
		SET(
			tNotifications.ReadAt.SET(readAt),
		).
		WHERE(condition)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	if affected > 0 {
		if req.GetUnread() {
			affected = -affected
		}

		s.js.PublishProto(
			ctx,
			fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.GetUserId()),
			&notifications.UserEvent{
				Data: &notifications.UserEvent_NotificationsReadCount{
					NotificationsReadCount: affected,
				},
			},
		)
	}

	return &pbnotifications.MarkNotificationsResponse{
		Updated: affected,
	}, nil
}

func (s *Server) getNotificationCount(ctx context.Context, userId int32) (int64, error) {
	stmt := tNotifications.
		SELECT(
			mysql.COUNT(tNotifications.ID).AS("count"),
		).
		FROM(tNotifications).
		WHERE(mysql.AND(
			tNotifications.UserID.EQ(mysql.Int32(userId)),
			tNotifications.ReadAt.IS_NULL(),
		)).
		ORDER_BY(tNotifications.ID.DESC())

	var dest struct {
		Count int64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, errswrap.NewError(err, ErrFailedStream)
		}
	}

	return dest.Count, nil
}
