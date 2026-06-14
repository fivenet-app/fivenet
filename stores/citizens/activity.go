package citizensstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	pb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func newUserActivitySorter() *database.Builder {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")

	return database.New(
		database.SpecMap{
			"createdAt": database.Column{Col: tUserActivity.CreatedAt, NullsLast: true},
		},
		[]mysql.OrderByClause{
			tUserActivity.CreatedAt.DESC().NULLS_LAST(),
		},
		[]mysql.OrderByClause{
			tUserActivity.ID.DESC(),
		},
		3,
	)
}

func (s *Store) ListUserActivity(
	ctx context.Context,
	req *pb.ListUserActivityRequest,
	limit int64,
) ([]*usersactivity.UserActivity, error) {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")
	tUTarget := table.FivenetUser.AS("target_user")
	tUSource := tUTarget.AS("source_user")

	condition := tUserActivity.TargetUserID.EQ(mysql.Int32(req.GetUserId()))
	if len(req.GetTypes()) > 0 {
		types := make([]mysql.Expression, 0, len(req.GetTypes()))
		for _, t := range req.GetTypes() {
			types = append(types, mysql.Int32(int32(*t.Enum())))
		}
		condition = condition.AND(tUserActivity.Type.IN(types...))
	}

	orderBys := newUserActivitySorter().Build(req.GetSort())

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
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit)

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
	req *pb.ListUserActivityRequest,
) (int64, error) {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")

	condition := tUserActivity.TargetUserID.EQ(mysql.Int32(req.GetUserId()))
	if len(req.GetTypes()) > 0 {
		types := make([]mysql.Expression, 0, len(req.GetTypes()))
		for _, t := range req.GetTypes() {
			types = append(types, mysql.Int32(int32(*t.Enum())))
		}
		condition = condition.AND(tUserActivity.Type.IN(types...))
	}

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
