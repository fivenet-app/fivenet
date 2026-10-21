package citizensstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func buildUserActivityCondition(
	tUserActivity *table.FivenetUserActivityTable,
	opts UserActivityOptions,
) mysql.BoolExpression {
	condition := tUserActivity.TargetUserID.EQ(mysql.Int32(opts.UserID))
	if len(opts.Types) > 0 {
		types := make([]mysql.Expression, 0, len(opts.Types))
		for _, t := range opts.Types {
			types = append(types, mysql.Int32(int32(*t.Enum())))
		}
		condition = condition.AND(tUserActivity.Type.IN(types...))
	}
	if len(opts.PropertyKeys) > 0 {
		keys := make([]mysql.Expression, 0, len(opts.PropertyKeys))
		for _, key := range opts.PropertyKeys {
			keys = append(keys, mysql.String(key))
		}
		condition = condition.AND(
			tUserActivity.Type.EQ(mysql.Int32(int32(usersactivity.UserActivityType_USER_ACTIVITY_TYPE_PROPERTY))).
				AND(tUserActivity.Key.IN(keys...)),
		)
	}
	propertyActivityType := mysql.Int32(int32(
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_PROPERTY,
	))
	if !opts.CanViewAllProperties {
		if len(opts.VisiblePropertyKeys) == 0 {
			condition = condition.AND(tUserActivity.Type.NOT_EQ(propertyActivityType))
		} else {
			keys := make([]mysql.Expression, 0, len(opts.VisiblePropertyKeys))
			for _, key := range opts.VisiblePropertyKeys {
				keys = append(keys, mysql.String(key))
			}
			condition = condition.AND(mysql.OR(
				tUserActivity.Type.NOT_EQ(propertyActivityType),
				tUserActivity.Key.IN(keys...),
			))
		}
	}
	if opts.From != nil {
		condition = condition.AND(
			tUserActivity.CreatedAt.GT_EQ(dbutils.TimestampToMySQL(opts.From)),
		)
	}
	if opts.To != nil {
		condition = condition.AND(tUserActivity.CreatedAt.LT_EQ(dbutils.TimestampToMySQL(opts.To)))
	}

	return condition
}

func (s *Store) ListUserActivity(
	ctx context.Context,
	opts ListUserActivityOptions,
) ([]*usersactivity.UserActivity, error) {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")

	condition := buildUserActivityCondition(tUserActivity, opts.UserActivityOptions)

	orderBys := s.userActivitySorter.Build(opts.Sort)

	stmt := mysql.
		SELECT(
			tUserActivity.ID,
			tUserActivity.CreatedAt,
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Key,
			tUserActivity.Reason,
			tUserActivity.Data,
		).
		FROM(tUserActivity).
		WHERE(condition).
		OFFSET(opts.Offset).
		ORDER_BY(orderBys...).
		LIMIT(opts.Limit)

	var activities []*usersactivity.UserActivity
	if err := stmt.QueryContext(ctx, s.db, &activities); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return activities, nil
}

func (s *Store) CountUserActivity(
	ctx context.Context,
	opts CountUserActivityOptions,
) (int64, error) {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")

	condition := buildUserActivityCondition(tUserActivity, opts.UserActivityOptions)

	countStmt := tUserActivity.
		SELECT(
			mysql.COUNT(tUserActivity.ID).AS("data_count.total"),
		).
		FROM(tUserActivity).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}
