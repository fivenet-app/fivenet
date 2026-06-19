package vehiclesstore

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
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
	sorter   *database.SorterBuilder
}

func New(db *sql.DB, customDB *config.CustomDB) IStore {
	return &Store{
		db:       db,
		customDB: customDB,
		sorter: database.New(
			database.SpecMap{
				"model": database.Column{Col: table.FivenetOwnedVehicles.AS("vehicle").Model},
				"plate": database.Column{Col: table.FivenetOwnedVehicles.AS("vehicle").Plate},
			},
			[]mysql.OrderByClause{table.FivenetOwnedVehicles.AS("vehicle").Plate.ASC()},
			[]mysql.OrderByClause{table.FivenetOwnedVehicles.AS("vehicle").Type.ASC()},
			"plate",
			3,
		),
	}
}
