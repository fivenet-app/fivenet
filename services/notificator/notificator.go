package notificator

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbnotificator "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/notificator"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
)

var tNotifications = table.FivenetNotifications

var (
	ErrFailedRequest = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.NotificatorService.ErrFailedRequest"}, nil)
	ErrFailedStream  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.NotificatorService.ErrFailedStream"}, nil)
)

func (s *Server) GetNotifications(ctx context.Context, req *pbnotificator.GetNotificationsRequest) (*pbnotificator.GetNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tNotifications := tNotifications.AS("notification")
	condition := tNotifications.UserID.EQ(jet.Int32(userInfo.UserId))
	if req.IncludeRead != nil && !*req.IncludeRead {
		condition = condition.AND(tNotifications.ReadAt.IS_NULL())
	}

	if len(req.Categories) > 0 {
		categoryIds := make([]jet.Expression, len(req.Categories))
		for i := range req.Categories {
			categoryIds[i] = jet.Int16(int16(req.Categories[i]))
		}

		condition = condition.AND(tNotifications.Category.IN(categoryIds...))
	}

	countStmt := tNotifications.
		SELECT(
			jet.COUNT(tNotifications.ID).AS("datacount.totalcount"),
		).
		FROM(tNotifications).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &pbnotificator.GetNotificationsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
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
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tNotifications.ID.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Notifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	resp.Pagination.Update(len(resp.Notifications))

	return resp, nil
}

func (s *Server) MarkNotifications(ctx context.Context, req *pbnotificator.MarkNotificationsRequest) (*pbnotificator.MarkNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tNotifications.UserID.EQ(
		jet.Int32(userInfo.UserId)).AND(
		tNotifications.ReadAt.IS_NULL(),
	)

	// If not all
	if len(req.Ids) > 0 {
		ids := make([]jet.Expression, len(req.Ids))
		for i := range req.Ids {
			ids[i] = jet.Uint64(req.Ids[i])
		}
		condition = condition.AND(tNotifications.ID.IN(ids...))
	} else if req.All == nil || !*req.All {
		return &pbnotificator.MarkNotificationsResponse{}, nil
	}

	readAt := jet.CURRENT_TIMESTAMP()
	if req.Unread {
		// Allow users to mark notifications as unread
		readAt = jet.TimestampExp(jet.NULL)
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
		if req.Unread {
			affected = -affected
		}

		s.js.PublishProto(ctx, fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.UserId),
			&notifications.UserEvent{
				Data: &notifications.UserEvent_NotificationsReadCount{
					NotificationsReadCount: int32(affected),
				},
			},
		)
	}

	return &pbnotificator.MarkNotificationsResponse{
		Updated: uint64(affected),
	}, nil
}

func (s *Server) getNotificationCount(ctx context.Context, userId int32) (int32, error) {
	stmt := tNotifications.
		SELECT(
			jet.COUNT(tNotifications.ID).AS("count"),
		).
		FROM(tNotifications).
		WHERE(jet.AND(
			tNotifications.UserID.EQ(jet.Int32(userId)),
			tNotifications.ReadAt.IS_NULL(),
		)).
		ORDER_BY(tNotifications.ID.DESC())

	var dest struct {
		Count int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, errswrap.NewError(err, ErrFailedStream)
		}
	}

	return dest.Count, nil
}
