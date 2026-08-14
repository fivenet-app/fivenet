package jobsstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CreateGroupActivity(
	ctx context.Context,
	db qrm.DB,
	activities ...*jobsgroups.GroupActivity,
) error {
	if len(activities) == 0 {
		return nil
	}

	tActivity := table.FivenetJobGroupActivity
	stmt := tActivity.
		INSERT(
			tActivity.Job,
			tActivity.GroupID,
			tActivity.ActivityType,
			tActivity.ActorUserID,
			tActivity.TargetUserID,
			tActivity.RuleID,
			tActivity.Reason,
			tActivity.Data,
		)

	for _, activity := range activities {
		stmt = stmt.VALUES(
			activity.GetJob(),
			activity.GetGroupId(),
			int32(activity.GetType()),
			activity.GetActorUserId(),
			dbutils.Int32P(activity.GetTargetUserId()),
			dbutils.Int64P(activity.GetRuleId()),
			dbutils.StringPP(activity.Reason),
			activity.GetData(),
		)
	}

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (s *Store) CountGroupActivity(ctx context.Context, db qrm.DB, q ListQuery) (int64, error) {
	tActivity := table.FivenetJobGroupActivity.AS("group_activity")
	condition := mysql.AND(tActivity.Job.EQ(mysql.String(q.Job)))
	if q.Where != nil {
		condition = condition.AND(q.Where)
	}

	countStmt := tActivity.
		SELECT(mysql.COUNT(tActivity.ID).AS("data_count.total")).
		FROM(tActivity).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListGroupActivity(
	ctx context.Context,
	db qrm.DB,
	q ListQuery,
) ([]*jobsgroups.GroupActivity, error) {
	tActivity := table.FivenetJobGroupActivity.AS("group_activity")
	condition := mysql.AND(tActivity.Job.EQ(mysql.String(q.Job)))
	if q.Where != nil {
		condition = condition.AND(q.Where)
	}

	orderBys := []mysql.OrderByClause{tActivity.CreatedAt.DESC(), tActivity.ID.DESC()}
	if q.Sort != nil && len(q.Sort.GetColumns()) > 0 {
		orderBys = []mysql.OrderByClause{}
		for _, sc := range q.Sort.GetColumns() {
			switch sc.GetId() {
			case "id":
				if sc.GetDesc() {
					orderBys = append(orderBys, tActivity.ID.DESC())
				} else {
					orderBys = append(orderBys, tActivity.ID.ASC())
				}
			case "type", "activity_type":
				if sc.GetDesc() {
					orderBys = append(orderBys, tActivity.ActivityType.DESC())
				} else {
					orderBys = append(orderBys, tActivity.ActivityType.ASC())
				}
			case "actor_user_id":
				if sc.GetDesc() {
					orderBys = append(orderBys, tActivity.ActorUserID.DESC())
				} else {
					orderBys = append(orderBys, tActivity.ActorUserID.ASC())
				}
			case "target_user_id":
				if sc.GetDesc() {
					orderBys = append(orderBys, tActivity.TargetUserID.DESC())
				} else {
					orderBys = append(orderBys, tActivity.TargetUserID.ASC())
				}
			case "rule_id":
				if sc.GetDesc() {
					orderBys = append(orderBys, tActivity.RuleID.DESC())
				} else {
					orderBys = append(orderBys, tActivity.RuleID.ASC())
				}
			case "created_at":
				if sc.GetDesc() {
					orderBys = append(orderBys, tActivity.CreatedAt.DESC())
				} else {
					orderBys = append(orderBys, tActivity.CreatedAt.ASC())
				}
			}
		}
		if len(orderBys) == 0 {
			orderBys = []mysql.OrderByClause{tActivity.CreatedAt.DESC(), tActivity.ID.DESC()}
		}
	}

	stmt := tActivity.
		SELECT(
			tActivity.ID,
			tActivity.Job,
			tActivity.GroupID,
			tActivity.ActivityType.AS("group_activity.type"),
			tActivity.ActorUserID,
			tActivity.TargetUserID,
			tActivity.RuleID,
			tActivity.Reason,
			tActivity.Data,
			tActivity.CreatedAt,
		).
		FROM(tActivity).
		WHERE(condition).
		OFFSET(q.Offset).
		ORDER_BY(orderBys...).
		LIMIT(q.Limit)

	activity := []*jobsgroups.GroupActivity{}
	if err := stmt.QueryContext(ctx, db, &activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return activity, nil
}
