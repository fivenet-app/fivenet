package jobs

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
)

func CreateJobsUserActivities(ctx context.Context, tx qrm.DB, activities ...*JobsUserActivity) error {
	if len(activities) == 0 {
		return nil
	}

	tJobsUserActivity := table.FivenetJobsUserActivity

	stmt := tJobsUserActivity.
		INSERT(
			tJobsUserActivity.Job,
			tJobsUserActivity.SourceUserID,
			tJobsUserActivity.TargetUserID,
			tJobsUserActivity.ActivityType,
			tJobsUserActivity.Reason,
			tJobsUserActivity.Data,
		)

	for _, activity := range activities {
		stmt = stmt.
			VALUES(
				activity.Job,
				activity.SourceUserId,
				activity.TargetUserId,
				activity.ActivityType,
				activity.Reason,
				activity.Data,
			)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
