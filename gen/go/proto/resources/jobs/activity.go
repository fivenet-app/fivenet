package jobs

import (
	"context"
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/encoding/protojson"
)

func CreateJobsUserActivities(ctx context.Context, tx qrm.DB, activities ...*JobsUserActivity) error {
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

// Scan implements driver.Valuer for protobuf JobsUserActivityData.
func (x *JobsUserActivityData) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *JobsUserActivityData) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
