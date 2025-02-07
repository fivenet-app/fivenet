package users

import (
	"context"
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/encoding/protojson"
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
				activity.SourceUserId,
				activity.TargetUserId,
				activity.Type,
				activity.Reason,
				activity.Data,
			)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

// Scan implements driver.Valuer for protobuf UserActivityData.
func (x *UserActivityData) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *UserActivityData) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
