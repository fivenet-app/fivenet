package vehiclesstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

type IStore interface {
	Count(ctx context.Context, q ListQuery) (int64, error)
	List(ctx context.Context, q ListQuery) ([]*resourcesvehicles.Vehicle, error)
	GetProps(ctx context.Context, plate string) (*vehiclesprops.VehicleProps, error)
	UpdateProps(
		ctx context.Context,
		in *vehiclesprops.VehicleProps,
	) (*vehiclesprops.VehicleProps, error)
	ListExpiredWanted(ctx context.Context, maxDays int64, limit int64) ([]string, error)
	ClearWanted(ctx context.Context, plate string) error
}

type Store struct {
	db       *sql.DB
	customDB *config.CustomDB
}

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

func New(db *sql.DB, customDB *config.CustomDB) IStore {
	return &Store{
		db:       db,
		customDB: customDB,
	}
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
	orderBys := buildListOrder(q.Sort, tVehicles)

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
) (*vehiclesprops.VehicleProps, error) {
	props, err := s.GetProps(ctx, in.GetPlate())
	if err != nil {
		return nil, err
	}

	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := props.HandleChanges(ctx, tx, in); err != nil {
		return nil, err
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

	in := proto.Clone(props).(*vehiclesprops.VehicleProps)
	wanted := false
	in.Wanted = &wanted
	in.WantedReason = nil
	in.WantedAt = nil

	return props.HandleChanges(ctx, s.db, in)
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

func buildListOrder(
	sort *database.Sort,
	tVehicles *table.FivenetOwnedVehiclesTable,
) []mysql.OrderByClause {
	orderBys := []mysql.OrderByClause{}
	if sort != nil && len(sort.GetColumns()) > 0 {
		for _, sc := range sort.GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "model":
				column = tVehicles.Model
			case "plate":
				fallthrough
			default:
				column = tVehicles.Plate
			}

			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}
	} else {
		orderBys = append(orderBys, tVehicles.Plate.ASC())
	}

	return append(orderBys, tVehicles.Type.ASC())
}
