package qualificationsstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	qualificationsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CreateQualificationActivity(
	ctx context.Context,
	tx qrm.DB,
	activity *qualificationsactivity.QualificationActivity,
) error {
	tActivity := table.FivenetQualificationsActivity
	_, err := tActivity.
		INSERT(
			tActivity.QualificationID,
			tActivity.ActivityType,
			tActivity.ActorUserID,
			tActivity.TargetUserID,
			tActivity.Data,
		).
		VALUES(
			activity.GetQualificationId(),
			int32(activity.GetType()), dbutils.Int32P(activity.GetActorUserId()),
			dbutils.Int32P(activity.GetTargetUserId()), activity.GetData(),
		).
		ExecContext(ctx, tx)
	return err
}

func qualificationActivityCondition(
	tActivity *table.FivenetQualificationsActivityTable,
	opts ListQualificationActivityOptions,
) mysql.BoolExpression {
	condition := tActivity.QualificationID.EQ(mysql.Int64(opts.QualificationID))
	if len(opts.Types) > 0 {
		types := make([]mysql.Expression, 0, len(opts.Types))
		for _, activityType := range opts.Types {
			if activityType != qualificationsactivity.QualificationActivityType_QUALIFICATION_ACTIVITY_TYPE_UNSPECIFIED {
				types = append(types, mysql.Int32(int32(activityType)))
			}
		}
		if len(types) == 0 {
			return mysql.Bool(false)
		}
		condition = condition.AND(tActivity.ActivityType.IN(types...))
	}
	if opts.UserID > 0 {
		condition = condition.AND(
			mysql.OR(
				tActivity.ActorUserID.EQ(mysql.Int32(opts.UserID)),
				tActivity.TargetUserID.EQ(mysql.Int32(opts.UserID)),
			),
		)
	}
	if opts.From != nil {
		condition = condition.AND(tActivity.CreatedAt.GT_EQ(dbutils.TimestampToMySQL(opts.From)))
	}
	if opts.To != nil {
		condition = condition.AND(tActivity.CreatedAt.LT_EQ(dbutils.TimestampToMySQL(opts.To)))
	}
	return condition
}

func (s *Store) CountQualificationActivity(
	ctx context.Context,
	opts ListQualificationActivityOptions,
) (int64, error) {
	tActivity := table.FivenetQualificationsActivity.AS("qualification_activity")
	var count database.DataCount

	err := tActivity.
		SELECT(mysql.COUNT(tActivity.ID).AS("data_count.total")).
		FROM(tActivity).
		WHERE(qualificationActivityCondition(tActivity, opts)).QueryContext(ctx, s.db, &count)

	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return 0, err
	}

	return count.Total, nil
}

func (s *Store) ListQualificationActivity(
	ctx context.Context,
	opts ListQualificationActivityOptions,
) ([]*qualificationsactivity.QualificationActivity, error) {
	tActivity := table.FivenetQualificationsActivity.AS("qualification_activity")
	activity := []*qualificationsactivity.QualificationActivity{}
	orderBys := []mysql.OrderByClause{tActivity.CreatedAt.DESC(), tActivity.ID.DESC()}
	if opts.Sort != nil && len(opts.Sort.GetColumns()) > 0 {
		orderBys = []mysql.OrderByClause{}
		for _, column := range opts.Sort.GetColumns() {
			if column.GetId() != "createdAt" && column.GetId() != "created_at" {
				continue
			}
			if column.GetDesc() {
				orderBys = append(orderBys, tActivity.CreatedAt.DESC(), tActivity.ID.DESC())
			} else {
				orderBys = append(orderBys, tActivity.CreatedAt.ASC(), tActivity.ID.ASC())
			}
		}
		if len(orderBys) == 0 {
			orderBys = []mysql.OrderByClause{tActivity.CreatedAt.DESC(), tActivity.ID.DESC()}
		}
	}
	err := tActivity.
		SELECT(
			tActivity.ID,
			tActivity.QualificationID,
			tActivity.ActivityType.AS("qualification_activity.type"),
			tActivity.ActorUserID,
			tActivity.TargetUserID,
			tActivity.Data,
			tActivity.CreatedAt,
		).
		FROM(tActivity).
		WHERE(qualificationActivityCondition(tActivity, opts)).
		ORDER_BY(orderBys...).
		OFFSET(opts.Offset).
		LIMIT(opts.Limit).
		QueryContext(ctx, s.db, &activity)
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}

	return activity, nil
}
