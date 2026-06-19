package notificationsstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CountUnread(ctx context.Context, userID int32) (int64, error) {
	tNotifications := table.FivenetNotifications
	stmt := tNotifications.
		SELECT(
			mysql.COUNT(tNotifications.ID).AS("count"),
		).
		FROM(tNotifications).
		WHERE(mysql.AND(
			tNotifications.UserID.EQ(mysql.Int32(userID)),
			tNotifications.ReadAt.IS_NULL(),
		)).
		ORDER_BY(tNotifications.ID.DESC())

	var dest struct {
		Count int64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return dest.Count, nil
}
