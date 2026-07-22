package vehiclesstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

type ListQuery struct {
	LicensePlate string
	Model        string
	UserIDs      []int32
	Job          string
	Wanted       *bool

	CanFilterWanted     bool
	IncludePhoneNumber  bool
	IncludePropsUpdated bool
	IncludeWantedFields bool

	Sort   *database.Sort
	Offset int64
	Limit  int64
}

func (s *Store) Count(ctx context.Context, q ListQuery) (int64, error) {
	tVehicles := table.FivenetOwnedVehicles.AS("vehicle")
	tVehicleProps := table.FivenetVehiclesProps.AS("vehicle_props")
	tUsers := table.FivenetUser.AS("user_short")

	condition, userCondition := s.listConditions(q, tVehicles, tVehicleProps, tUsers)

	stmt := tVehicles.
		SELECT(
			mysql.COUNT(tVehicles.Plate).AS("data_count.total"),
		).
		FROM(
			tVehicles.
				LEFT_JOIN(tUsers,
					userCondition,
				).
				LEFT_JOIN(tVehicleProps,
					tVehicleProps.Plate.EQ(tVehicles.Plate),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) List(ctx context.Context, q ListQuery) ([]*resourcesvehicles.Vehicle, error) {
	tVehicles := table.FivenetOwnedVehicles.AS("vehicle")
	tVehicleProps := table.FivenetVehiclesProps.AS("vehicle_props")
	tUsers := table.FivenetUser.AS("user_short")

	condition, userCondition := s.listConditions(q, tVehicles, tVehicleProps, tUsers)
	orderBys := s.sorter.Build(q.Sort)

	columns := dbutils.Columns{
		s.customDB.Columns.Vehicle.GetModel(tVehicles.Alias()),
		mysql.REPLACE(tVehicles.Type, mysql.String("_"), mysql.String(" ")).AS("vehicle.type"),
		tUsers.ID.AS("vehicle.owner_id"),
		tUsers.ID,
		tUsers.Firstname,
		tUsers.Lastname,
		tUsers.Dateofbirth,
		tVehicleProps.Plate,
		tVehicles.Job,
		tVehicles.Data,
	}

	if q.IncludePhoneNumber {
		columns = append(columns, tUsers.PhoneNumber)
	}
	if q.IncludePropsUpdated {
		columns = append(columns, tVehicleProps.UpdatedAt)
	}
	if q.IncludeWantedFields {
		columns = append(columns,
			tVehicleProps.Wanted,
			tVehicleProps.WantedReason,
		)
	}

	stmt := tVehicles.
		SELECT(
			tVehicles.Plate,
			columns.Get()...,
		).
		FROM(
			tVehicles.
				LEFT_JOIN(tUsers,
					userCondition,
				).
				LEFT_JOIN(tVehicleProps,
					tVehicleProps.Plate.EQ(tVehicles.Plate),
				),
		).
		WHERE(condition).
		OFFSET(q.Offset).
		ORDER_BY(orderBys...).
		LIMIT(q.Limit)

	var vehicles []*resourcesvehicles.Vehicle
	if err := stmt.QueryContext(ctx, s.db, &vehicles); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (s *Store) IsVehicleOwner(ctx context.Context, plate string, userID int32) (bool, error) {
	tVehicles := table.FivenetOwnedVehicles.AS("vehicle")
	tUsers := table.FivenetUser.AS("user_short")

	var owner struct {
		Plate string `alias:"plate"`
	}
	stmt := tVehicles.
		SELECT(tVehicles.Plate.AS("plate")).
		FROM(tVehicles.
			INNER_JOIN(tUsers, tUsers.ID.EQ(tVehicles.UserID)),
		).
		WHERE(mysql.AND(
			tVehicles.Plate.EQ(mysql.String(plate)),
			tUsers.ID.EQ(mysql.Int32(userID)),
		)).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &owner); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return owner.Plate != "", nil
}

func (s *Store) GetProps(ctx context.Context, plate string) (*vehiclesprops.VehicleProps, error) {
	props := &vehiclesprops.VehicleProps{}
	if err := props.LoadFromDB(ctx, s.db, plate); err != nil {
		return nil, err
	}

	return props, nil
}

func (s *Store) UpdateProps(
	ctx context.Context,
	in *vehiclesprops.VehicleProps,
	creatorID *int32,
	creatorJob string,
	reason string,
) (*vehiclesprops.VehicleProps, error) {
	props, err := s.GetProps(ctx, in.GetPlate())
	if err != nil {
		return nil, err
	}

	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}
	props.NormalizeWantedChange(in, reason)
	activity := buildWantedActivity(props, in, creatorID, creatorJob, reason, false)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := props.HandleChanges(ctx, tx, in); err != nil {
		return nil, err
	}
	if activity != nil {
		if _, err := s.addVehicleActivity(ctx, tx, activity); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetProps(ctx, in.GetPlate())
}

func (s *Store) ListExpiredWanted(
	ctx context.Context,
	maxDays int64,
	limit int64,
) ([]string, error) {
	tVehicleProps := table.FivenetVehiclesProps
	stmt := tVehicleProps.
		SELECT(
			tVehicleProps.Plate,
		).
		FROM(tVehicleProps).
		WHERE(mysql.AND(
			tVehicleProps.Wanted.IS_TRUE(),
			mysql.OR(
				tVehicleProps.WantedAt.LT(
					mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(maxDays, "DAY")),
				),
				tVehicleProps.WantedTill.LT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		LIMIT(limit)

	var plates []string
	if err := stmt.QueryContext(ctx, s.db, &plates); err != nil {
		return nil, err
	}

	return plates, nil
}

func (s *Store) ClearWanted(ctx context.Context, plate string) error {
	props := &vehiclesprops.VehicleProps{}
	if err := props.LoadFromDB(ctx, s.db, plate); err != nil {
		return err
	}
	if props.Wanted == nil || !props.GetWanted() {
		return nil
	}

	in := proto.Clone(props).(*vehiclesprops.VehicleProps)
	wanted := false
	in.Wanted = &wanted
	in.WantedReason = nil
	in.WantedAt = nil
	in.WantedTill = nil
	props.NormalizeWantedChange(in, "wanted_expired")

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := props.HandleChanges(ctx, tx, in); err != nil {
		return err
	}
	if activity := buildWantedActivity(
		props,
		in,
		nil,
		"",
		"wanted_expired",
		true,
	); activity != nil {
		if _, err := s.addVehicleActivity(ctx, tx, activity); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func buildWantedActivity(
	props *vehiclesprops.VehicleProps,
	in *vehiclesprops.VehicleProps,
	creatorID *int32,
	creatorJob string,
	reason string,
	auto bool,
) *vehiclesactivity.VehicleActivity {
	if in.Wanted == nil {
		return nil
	}

	changed := props.GetWanted() != in.GetWanted() ||
		props.GetWantedReason() != in.GetWantedReason() ||
		!proto.Equal(props.GetWantedTill(), in.GetWantedTill())
	if !changed {
		return nil
	}

	previousWanted := props.GetWanted()
	activity := &vehiclesactivity.VehicleActivity{
		Plate:        in.GetPlate(),
		ActivityType: vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED,
		CreatorId:    creatorID,
		CreatorJob:   creatorJob,
		Data: &vehiclesactivity.VehicleActivityData{
			Data: &vehiclesactivity.VehicleActivityData_WantedChange{
				WantedChange: &vehiclesactivity.WantedChange{
					Wanted:         in.GetWanted(),
					PreviousWanted: &previousWanted,
					WantedReason:   in.WantedReason,
					WantedAt:       in.GetWantedAt(),
					WantedTill:     in.GetWantedTill(),
					Auto:           auto,
				},
			},
		},
	}
	if reason != "" {
		activity.Reason = &reason
	}

	return activity
}

func (s *Store) listConditions(
	q ListQuery,
	tVehicles *table.FivenetOwnedVehiclesTable,
	tVehicleProps *table.FivenetVehiclesPropsTable,
	tUsers *table.FivenetUserTable,
) (mysql.BoolExpression, mysql.BoolExpression) {
	condition := mysql.Bool(true)
	userCondition := tUsers.ID.EQ(tVehicles.UserID)

	if search := dbutils.PrepareForLikeSearch(q.LicensePlate); search != "" {
		condition = mysql.AND(condition, tVehicles.Plate.LIKE(mysql.String(search)))
	}

	modelColumn := s.customDB.Columns.Vehicle.GetModel(tVehicles.Alias())
	if modelColumn != nil {
		if search := dbutils.PrepareForLikeSearch(q.Model); search != "" {
			condition = mysql.AND(condition, tVehicles.Model.LIKE(mysql.String(search)))
		}
	}

	if len(q.UserIDs) > 0 {
		userIDs := make([]mysql.Expression, 0, len(q.UserIDs))
		for _, userID := range q.UserIDs {
			userIDs = append(userIDs, mysql.Int32(userID))
		}

		condition = mysql.AND(condition,
			tUsers.ID.EQ(tVehicles.UserID),
			tUsers.ID.IN(userIDs...),
		)
		userCondition = mysql.AND(userCondition, tUsers.ID.IN(userIDs...))
	} else if q.Job != "" {
		condition = mysql.AND(condition,
			tVehicles.Job.EQ(mysql.String(q.Job)),
		)
	}

	if q.CanFilterWanted && q.Wanted != nil && *q.Wanted {
		condition = mysql.AND(condition,
			tVehicleProps.Wanted.EQ(mysql.Bool(*q.Wanted)),
		)
	}

	return condition, userCondition
}
