package users

import (
	"context"

	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
)

func CreateUserActivities(ctx context.Context, tx qrm.DB, activities ...*UserActivity) error {
	tUserActivity := table.FivenetUserActivity

	stmt := tUserActivity.
		INSERT(
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Key,
			tUserActivity.OldValue,
			tUserActivity.NewValue,
			tUserActivity.Reason,
		)

	for _, activity := range activities {
		stmt = stmt.
			VALUES(
				activity.SourceUserId,
				activity.TargetUserId,
				activity.Type,
				activity.Key,
				activity.OldValue,
				activity.NewValue,
				activity.Reason,
			)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
