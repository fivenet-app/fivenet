package users

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
)

func CreateUserActivities(ctx context.Context, tx qrm.DB, activities ...*UserActivity) error {
	if len(activities) == 0 {
		return nil
	}

	tUserActivity := table.FivenetUserActivity

	stmt := tUserActivity.
		INSERT(
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
		)

	for _, activity := range activities {
		stmt = stmt.
			VALUES(
				activity.GetSourceUserId(),
				activity.GetTargetUserId(),
				activity.GetType(),
				activity.GetReason(),
				activity.GetData(),
			)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
