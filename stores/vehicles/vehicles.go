package vehiclesstore

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

type ListQuery struct {
	UserInfo *userinfo.UserInfo

	LicensePlate    string
	Plates          []string
	Model           string
	UserIDs         []int32
	Job             string
	PropertyFilters []*properties.PropertyFilter

	CanFilterProperties bool
	IncludePropsUpdated bool
	IncludeWantedFields bool

	Sort   *database.Sort
	Offset int64
	Limit  int64
}

func (s *Store) Count(ctx context.Context, q ListQuery) (int64, error) {
	tVehicles := table.FivenetVehicle.AS("vehicle")
	condition := s.listConditions(q, tVehicles)

	stmt := tVehicles.
		SELECT(
			mysql.COUNT(tVehicles.Plate).AS("data_count.total"),
		).
		FROM(tVehicles).
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
	tVehicles := table.FivenetVehicle.AS("vehicle")
	condition := s.listConditions(q, tVehicles)
	orderBys := s.sorter.Build(q.Sort)

	columns := dbutils.Columns{
		s.customDB.Columns.Vehicle.GetModel(tVehicles.Alias()),
		mysql.REPLACE(tVehicles.Type, mysql.String("_"), mysql.String(" ")).AS("vehicle.type"),
		tVehicles.UserID.AS("vehicle.owner_id"),
		tVehicles.Job,
		tVehicles.Data,
	}

	stmt := tVehicles.
		SELECT(
			tVehicles.Plate,
			columns.Get()...,
		).
		FROM(tVehicles).
		WHERE(condition).
		OFFSET(q.Offset).
		ORDER_BY(orderBys...).
		LIMIT(q.Limit)

	var vehicles []*resourcesvehicles.Vehicle
	if err := stmt.QueryContext(ctx, s.db, &vehicles); err != nil {
		return nil, err
	}

	if q.IncludeWantedFields {
		plates := make([]string, 0, len(vehicles))
		for _, vehicle := range vehicles {
			plates = append(plates, vehicle.GetPlate())
		}
		propsByPlate, err := loadVehicleWantedProps(ctx, s.db, plates)
		if err != nil {
			return nil, err
		}
		for _, vehicle := range vehicles {
			if props, ok := propsByPlate[vehicle.GetPlate()]; ok {
				vehicle.Props = props
			}
		}
	}

	if s.hydrator == nil {
		return vehicles, nil
	}

	targets := make([]citizenshydrator.BasicTarget, 0, len(vehicles))
	for i, vehicle := range vehicles {
		if vehicle.GetOwnerId() <= 0 {
			continue
		}
		targets = append(targets, citizenshydrator.BasicTarget{
			UserID: vehicle.GetOwnerId(),
			Set: func(user *usershort.UserShort) {
				vehicles[i].Owner = user
			},
		})
	}
	if err := s.hydrator.HydrateBasicTargets(ctx, nil, q.UserInfo, targets); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (s *Store) IsVehicleOwner(ctx context.Context, plate string, userID int32) (bool, error) {
	tVehicles := table.FivenetVehicle

	var owner struct {
		Plate string `alias:"plate"`
	}
	stmt := tVehicles.
		SELECT(tVehicles.Plate.AS("plate")).
		FROM(tVehicles).
		WHERE(mysql.AND(
			tVehicles.Plate.EQ(mysql.String(plate)),
			tVehicles.UserID.EQ(mysql.Int32(userID)),
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

func (s *Store) GetProps(
	ctx context.Context,
	tx qrm.DB,
	plate string,
) (*vehiclesprops.VehicleProps, error) {
	propsByPlate, err := loadVehicleWantedProps(ctx, tx, []string{plate})
	if err != nil {
		return nil, err
	}
	if props, ok := propsByPlate[plate]; ok {
		return props, nil
	}
	return &vehiclesprops.VehicleProps{Plate: plate}, nil
}

func loadVehicleWantedProps(
	ctx context.Context,
	db qrm.DB,
	plates []string,
) (map[string]*vehiclesprops.VehicleProps, error) {
	result := make(map[string]*vehiclesprops.VehicleProps, len(plates))
	if len(plates) == 0 {
		return result, nil
	}

	placeholders := strings.TrimRight(strings.Repeat("?,", len(plates)), ",")
	args := make([]any, 0, len(plates)+1)
	args = append(args, vehicleWantedPropertyKey)
	for _, plate := range plates {
		args = append(args, plate)
	}

	rows, err := db.QueryContext(ctx, `
SELECT value.plate, value.updated_at, value.value_bool, state.reason, state.active_at, state.expires_at
FROM fivenet_vehicle_property_values value
LEFT JOIN fivenet_vehicle_property_state state
  ON state.plate = value.plate AND state.property_key = value.property_key
WHERE value.property_key = ? AND value.plate IN (`+placeholders+`)
ORDER BY value.plate ASC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			plate     string
			updatedAt time.Time
			wanted    *bool
			reason    *string
			activeAt  *time.Time
			expiresAt *time.Time
		)
		if err := rows.Scan(
			&plate,
			&updatedAt,
			&wanted,
			&reason,
			&activeAt,
			&expiresAt,
		); err != nil {
			return nil, err
		}

		props := &vehiclesprops.VehicleProps{
			Plate:     plate,
			UpdatedAt: timestamp.New(updatedAt),
		}
		if wanted != nil {
			props.Wanted = wanted
		}
		if reason != nil {
			props.WantedReason = reason
		}
		if activeAt != nil {
			props.WantedAt = timestamp.New(*activeAt)
		}
		if expiresAt != nil {
			props.WantedTill = timestamp.New(*expiresAt)
		}
		result[plate] = props
	}

	return result, rows.Err()
}

type PropsUpdate struct {
	Props           *vehiclesprops.VehicleProps
	CreatorID       *int32
	CreatorJob      string
	Reason          string
	PropertyChanges []*properties.PropertyChange
}

func (s *Store) UpdateProps(
	ctx context.Context,
	update PropsUpdate,
) (*vehiclesprops.VehicleProps, error) {
	in := update.Props
	props, err := s.GetProps(ctx, s.db, in.GetPlate())
	if err != nil {
		return nil, err
	}

	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}
	normalizeWantedChange(props, in, update.Reason)
	activity := buildWantedActivity(
		props,
		in,
		update.CreatorID,
		update.CreatorJob,
		update.Reason,
		false,
	)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := handlePropsChanges(ctx, tx, props, in); err != nil {
		return nil, err
	}
	if activity != nil {
		if _, err := s.addVehicleActivity(ctx, tx, activity); err != nil {
			return nil, err
		}
	}
	if err := s.applyPropertiesWithActivity(
		ctx,
		tx,
		in.GetPlate(),
		update.PropertyChanges,
		update.CreatorID,
		update.CreatorJob,
		update.Reason,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetProps(ctx, s.db, in.GetPlate())
}

func (s *Store) ListExpiredWanted(
	ctx context.Context,
	maxDays int64,
	limit int64,
) ([]string, error) {
	tValues := table.FivenetVehiclePropertyValues.AS("value")
	tState := table.FivenetVehiclePropertyState.AS("state")
	stmt := tValues.
		SELECT(tValues.Plate).
		FROM(tValues.INNER_JOIN(tState, mysql.AND(
			tState.Plate.EQ(tValues.Plate),
			tState.PropertyKey.EQ(tValues.PropertyKey),
		))).
		WHERE(mysql.AND(
			tValues.PropertyKey.EQ(mysql.String(vehicleWantedPropertyKey)),
			tValues.ValueType.EQ(mysql.Uint8(uint8(propertyValueTypeBool))),
			tValues.ValueBool.IS_TRUE(),
			mysql.OR(
				tState.ActiveAt.LT(
					mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(maxDays, mysql.DAY)),
				),
				tState.ExpiresAt.LT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		LIMIT(limit)

	plates := []string{}
	if err := stmt.QueryContext(ctx, s.db, &plates); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	return plates, nil
}

func (s *Store) ClearWanted(ctx context.Context, plate string) error {
	props, err := s.GetProps(ctx, s.db, plate)
	if err != nil {
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
	normalizeWantedChange(props, in, "wanted_expired")

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := handlePropsChanges(ctx, tx, props, in); err != nil {
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
	tVehicles *table.FivenetVehicleTable,
) mysql.BoolExpression {
	condition := mysql.Bool(true)

	if search := dbutils.PrepareForLikeSearch(q.LicensePlate); search != "" {
		condition = mysql.AND(condition, tVehicles.Plate.LIKE(mysql.String(search)))
	}
	if len(q.Plates) > 0 {
		plates := make([]mysql.Expression, 0, len(q.Plates))
		for _, plate := range q.Plates {
			plates = append(plates, mysql.String(plate))
		}
		condition = mysql.AND(condition, tVehicles.Plate.IN(plates...))
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
			tVehicles.UserID.IN(userIDs...),
		)
	} else if q.Job != "" {
		condition = mysql.AND(condition,
			tVehicles.Job.EQ(mysql.String(q.Job)),
		)
	}

	if q.CanFilterProperties {
		for i, filter := range q.PropertyFilters {
			propertyCondition, err := propertyFilterCondition(filter, tVehicles.Alias(), i)
			if err != nil {
				// ListQuery is constructed by the service after protobuf validation.
				// Preserve the invalid request as an empty result rather than ever
				// interpolating untrusted SQL.
				return mysql.Bool(false)
			}
			condition = mysql.AND(condition, propertyCondition)
		}
	}

	return condition
}
