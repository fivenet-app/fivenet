package citizensstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
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

	return condition
}

func (s *Store) ListUserActivity(
	ctx context.Context,
	opts ListUserActivityOptions,
) ([]*usersactivity.UserActivity, error) {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")
	tUTarget := table.FivenetUser.AS("target_user")
	tUSource := tUTarget.AS("source_user")

	condition := buildUserActivityCondition(tUserActivity, opts.UserActivityOptions)

	orderBys := s.userActivitySorter.Build(opts.Sort)

	stmt := mysql.
		SELECT(
			tUserActivity.ID,
			tUserActivity.CreatedAt,
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
			tUTarget.ID,
			tUTarget.Job,
			tUTarget.JobGrade,
			tUTarget.Firstname,
			tUTarget.Lastname,
			tUSource.ID,
			tUSource.Job,
			tUSource.JobGrade,
			tUSource.Firstname,
			tUSource.Lastname,
		).
		FROM(
			tUserActivity.
				INNER_JOIN(tUTarget,
					tUTarget.ID.EQ(tUserActivity.TargetUserID),
				).
				LEFT_JOIN(tUSource,
					tUSource.ID.EQ(tUserActivity.SourceUserID),
				),
		).
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
