package jobs

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
)

func CreateColleagueActivity(
	ctx context.Context,
	tx qrm.DB,
	activities ...*ColleagueActivity,
) error {
	if len(activities) == 0 {
		return nil
	}

	tJobColleagueActivity := table.FivenetJobColleagueActivity

	stmt := tJobColleagueActivity.
		INSERT(
			tJobColleagueActivity.Job,
			tJobColleagueActivity.SourceUserID,
			tJobColleagueActivity.TargetUserID,
			tJobColleagueActivity.ActivityType,
			tJobColleagueActivity.Reason,
			tJobColleagueActivity.Data,
		)

	for _, activity := range activities {
		stmt = stmt.
			VALUES(
				activity.GetJob(),
				activity.GetSourceUserId(),
				activity.GetTargetUserId(),
				activity.GetActivityType(),
				activity.GetReason(),
				activity.GetData(),
			)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
