package vehiclesstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type VehicleActivityOptions struct {
	Plate string
	Types []vehiclesactivity.VehicleActivityType
}

type CountVehicleActivityOptions struct {
	VehicleActivityOptions
}

type ListVehicleActivityOptions struct {
	VehicleActivityOptions

	Sort   *database.Sort
	Offset int64
	Limit  int64
}

func buildVehicleActivityCondition(
	tVehicleActivity *table.FivenetVehiclesActivityTable,
	opts VehicleActivityOptions,
) mysql.BoolExpression {
	condition := tVehicleActivity.Plate.EQ(mysql.String(opts.Plate))
	if len(opts.Types) > 0 {
		types := make([]mysql.Expression, 0, len(opts.Types))
		for _, t := range opts.Types {
			types = append(types, mysql.Int32(int32(*t.Enum())))
		}
		condition = condition.AND(tVehicleActivity.Type.IN(types...))
	}

	return condition
}

func (s *Store) CountVehicleActivity(
	ctx context.Context,
	opts CountVehicleActivityOptions,
) (int64, error) {
	tVehicleActivity := table.FivenetVehiclesActivity.AS("vehicle_activity")

	countStmt := tVehicleActivity.
		SELECT(
			mysql.COUNT(tVehicleActivity.ID).AS("data_count.total"),
		).
		FROM(tVehicleActivity).
		WHERE(buildVehicleActivityCondition(tVehicleActivity, opts.VehicleActivityOptions))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListVehicleActivity(
	ctx context.Context,
	opts ListVehicleActivityOptions,
) ([]*vehiclesactivity.VehicleActivity, error) {
	tVehicleActivity := table.FivenetVehiclesActivity.AS("vehicle_activity")
	tCreator := table.FivenetUser.AS("creator")

	orderBys := s.vehicleActivitySorter.Build(opts.Sort)

	stmt := tVehicleActivity.
		SELECT(
			tVehicleActivity.ID,
			tVehicleActivity.CreatedAt,
			tVehicleActivity.Plate,
			tVehicleActivity.Type,
			tVehicleActivity.CreatorID,
			tVehicleActivity.CreatorJob,
			tVehicleActivity.Reason,
			tVehicleActivity.Data,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
		).
		FROM(
			tVehicleActivity.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tVehicleActivity.CreatorID),
				),
		).
		WHERE(buildVehicleActivityCondition(tVehicleActivity, opts.VehicleActivityOptions)).
		OFFSET(opts.Offset).
		ORDER_BY(orderBys...).
		LIMIT(opts.Limit)

	activity := []*vehiclesactivity.VehicleActivity{}
	if err := stmt.QueryContext(ctx, s.db, &activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return activity, nil
}

func (s *Store) addVehicleActivity(
	ctx context.Context,
	tx qrm.DB,
	activity *vehiclesactivity.VehicleActivity,
) (int64, error) {
	tVehicleActivity := table.FivenetVehiclesActivity
	stmt := tVehicleActivity.
		INSERT(
			tVehicleActivity.CreatorID,
			tVehicleActivity.Plate,
			tVehicleActivity.Type,
			tVehicleActivity.CreatorJob,
			tVehicleActivity.Reason,
			tVehicleActivity.Data,
		).
		VALUES(
			activity.CreatorId,
			activity.GetPlate(),
			activity.GetActivityType(),
			activity.GetCreatorJob(),
			activity.Reason,
			activity.GetData(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}
