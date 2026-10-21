package vehiclesstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type VehicleActivityOptions struct {
	Plate               string
	Types               []vehiclesactivity.VehicleActivityType
	PropertyKeys        []string
	VisiblePropertyKeys []string
	From                *timestamp.Timestamp
	To                  *timestamp.Timestamp
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
	tVehicleActivity *table.FivenetVehicleActivityTable,
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
	if len(opts.PropertyKeys) > 0 {
		keys := make([]mysql.Expression, 0, len(opts.PropertyKeys))
		for _, key := range opts.PropertyKeys {
			keys = append(keys, mysql.String(key))
		}
		condition = condition.AND(
			tVehicleActivity.Type.EQ(mysql.Int32(int32(vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_PROPERTY))).
				AND(tVehicleActivity.Key.IN(keys...)),
		)
	}
	propertyActivityType := mysql.Int32(int32(
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_PROPERTY,
	))
	if len(opts.VisiblePropertyKeys) == 0 {
		condition = condition.AND(tVehicleActivity.Type.NOT_EQ(propertyActivityType))
	} else {
		keys := make([]mysql.Expression, 0, len(opts.VisiblePropertyKeys))
		for _, key := range opts.VisiblePropertyKeys {
			keys = append(keys, mysql.String(key))
		}
		condition = condition.AND(mysql.OR(
			tVehicleActivity.Type.NOT_EQ(propertyActivityType),
			tVehicleActivity.Key.IN(keys...),
		))
	}
	activeDefinitionExists := mysql.EXISTS(
		mysql.SELECT(mysql.Int(1)).
			FROM(table.FivenetPropertyDefinitions).
			WHERE(mysql.AND(
				table.FivenetPropertyDefinitions.Scope.EQ(mysql.String("vehicle")),
				table.FivenetPropertyDefinitions.PropertyKey.EQ(tVehicleActivity.Key),
				table.FivenetPropertyDefinitions.DeletedAt.IS_NULL(),
			)),
	)
	customPropertyActivity := tVehicleActivity.Type.EQ(mysql.Int32(int32(
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_PROPERTY,
	))).AND(tVehicleActivity.Key.IS_NOT_NULL()).AND(mysql.NOT(activeDefinitionExists))
	for _, definition := range propertyregistry.SystemPropertyDefinitionsForScope(
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
	) {
		customPropertyActivity = customPropertyActivity.AND(
			tVehicleActivity.Key.NOT_EQ(mysql.String(definition.GetKey())),
		)
	}
	condition = condition.AND(mysql.NOT(customPropertyActivity))
	if opts.From != nil {
		condition = condition.AND(
			tVehicleActivity.CreatedAt.GT_EQ(dbutils.TimestampToMySQL(opts.From)),
		)
	}
	if opts.To != nil {
		condition = condition.AND(
			tVehicleActivity.CreatedAt.LT_EQ(dbutils.TimestampToMySQL(opts.To)),
		)
	}

	return condition
}

func (s *Store) CountVehicleActivity(
	ctx context.Context,
	opts CountVehicleActivityOptions,
) (int64, error) {
	tVehicleActivity := table.FivenetVehicleActivity.AS("vehicle_activity")

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
	tVehicleActivity := table.FivenetVehicleActivity.AS("vehicle_activity")

	orderBys := s.vehicleActivitySorter.Build(opts.Sort)

	stmt := tVehicleActivity.
		SELECT(
			tVehicleActivity.ID,
			tVehicleActivity.CreatedAt,
			tVehicleActivity.Plate,
			tVehicleActivity.Type,
			tVehicleActivity.Key,
			tVehicleActivity.CreatorID,
			tVehicleActivity.CreatorJob,
			tVehicleActivity.Reason,
			tVehicleActivity.Data,
		).
		FROM(tVehicleActivity).
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
	tVehicleActivity := table.FivenetVehicleActivity
	stmt := tVehicleActivity.INSERT(
		tVehicleActivity.CreatorID, tVehicleActivity.Plate, tVehicleActivity.Type,
		tVehicleActivity.CreatorJob, tVehicleActivity.Reason, tVehicleActivity.Data,
	)
	if activity.Key != nil {
		stmt = tVehicleActivity.INSERT(
			tVehicleActivity.CreatorID,
			tVehicleActivity.Plate,
			tVehicleActivity.Type,
			tVehicleActivity.Key,
			tVehicleActivity.CreatorJob,
			tVehicleActivity.Reason,
			tVehicleActivity.Data,
		).VALUES(
			activity.CreatorId, activity.GetPlate(), activity.GetActivityType(), activity.Key,
			activity.GetCreatorJob(), activity.Reason, activity.GetData(),
		)
	} else {
		stmt = stmt.VALUES(
			activity.CreatorId, activity.GetPlate(), activity.GetActivityType(),
			activity.GetCreatorJob(), activity.Reason, activity.GetData(),
		)
	}

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return 0, nil
		}
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}
