package notificationsstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var notificationsTable = table.FivenetNotifications

func (s *Store) Count(ctx context.Context, q ListQuery) (int64, error) {
	tNotifications := notificationsTable.AS("notification")
	condition := buildListCondition(q, tNotifications)

	stmt := tNotifications.
		SELECT(
			mysql.COUNT(tNotifications.ID).AS("data_count.total"),
		).
		FROM(tNotifications).
		WHERE(condition)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) List(
	ctx context.Context,
	q ListQuery,
) ([]*resourcesnotifications.Notification, error) {
	tNotifications := notificationsTable.AS("notification")
	condition := buildListCondition(q, tNotifications)

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
		WHERE(condition).
		OFFSET(q.Offset).
		ORDER_BY(tNotifications.ID.DESC()).
		LIMIT(q.Limit)

	var notifications []*resourcesnotifications.Notification
	if err := stmt.QueryContext(ctx, s.db, &notifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return notifications, nil
}

func (s *Store) MarkNotifications(ctx context.Context, q MarkQuery) (int64, error) {
	tNotifications := notificationsTable
	condition := tNotifications.UserID.EQ(mysql.Int32(q.UserID)).AND(
		tNotifications.ReadAt.IS_NULL(),
	)

	if len(q.IDs) > 0 {
		ids := make([]mysql.Expression, len(q.IDs))
		for i := range q.IDs {
			ids[i] = mysql.Int64(q.IDs[i])
		}

		condition = condition.AND(tNotifications.ID.IN(ids...))
	} else if !q.All {
		return 0, nil
	}

	readAt := mysql.CURRENT_TIMESTAMP()
	if q.Unread {
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
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func buildListCondition(
	q ListQuery,
	tNotifications *table.FivenetNotificationsTable,
) mysql.BoolExpression {
	condition := tNotifications.UserID.EQ(mysql.Int32(q.UserID))
	if q.UnreadOnly {
		condition = condition.AND(tNotifications.ReadAt.IS_NULL())
	}

	if len(q.Categories) > 0 {
		categoryIds := make([]mysql.Expression, len(q.Categories))
		for i := range q.Categories {
			categoryIds[i] = mysql.Int32(int32(q.Categories[i]))
		}

		condition = condition.AND(tNotifications.Category.IN(categoryIds...))
	}

	return condition
}
