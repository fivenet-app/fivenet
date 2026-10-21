package vehiclesstore

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	Count(ctx context.Context, q ListQuery) (int64, error)
	List(ctx context.Context, q ListQuery) ([]*resourcesvehicles.Vehicle, error)
	IsVehicleOwner(ctx context.Context, plate string, userID int32) (bool, error)
	GetProps(ctx context.Context, tx qrm.DB, plate string) (*vehiclesprops.VehicleProps, error)
	GetProperties(ctx context.Context, plate string) ([]*properties.Property, error)
	ListProperties(ctx context.Context, plates []string) (map[string][]*properties.Property, error)
	ListCustomPropertyKeys(ctx context.Context) ([]string, error)
	ResolvePropertyDefinitions(
		ctx context.Context,
		keys []string,
	) (map[string]*properties.PropertyDefinition, error)
	ResolvePropertyGroups(
		ctx context.Context,
		definitions map[string]*properties.PropertyDefinition,
	) ([]*properties.PropertyGroup, error)
	ResolvePropertyCatalog(
		ctx context.Context,
		keys []string,
	) (*propertyregistry.ResolvedCatalog, error)
	IsSystemPropertyEnabled(ctx context.Context, key string) (bool, error)
	UpdateProps(
		ctx context.Context,
		update PropsUpdate,
	) (*vehiclesprops.VehicleProps, error)
	ListVehicleActivity(
		ctx context.Context,
		opts ListVehicleActivityOptions,
	) ([]*vehiclesactivity.VehicleActivity, error)
	CountVehicleActivity(ctx context.Context, opts CountVehicleActivityOptions) (int64, error)
	ListExpiredWanted(ctx context.Context, maxDays int64, limit int64) ([]string, error)
	ClearWanted(ctx context.Context, plate string) error
}

type Store struct {
	db                    *sql.DB
	customDB              *config.CustomDB
	hydrator              citizenshydrator.IHydrator
	sorter                *database.SorterBuilder
	vehicleActivitySorter *database.SorterBuilder
	propertyCache         *propertyregistry.DefinitionCache
}

func New(
	db *sql.DB,
	customDB *config.CustomDB,
	hydrator citizenshydrator.IHydrator,
) IStore {
	return NewWithCache(db, customDB, hydrator, propertyregistry.NewDefinitionCache())
}

func NewWithCache(
	db *sql.DB,
	customDB *config.CustomDB,
	hydrator citizenshydrator.IHydrator,
	propertyCache *propertyregistry.DefinitionCache,
) IStore {
	return &Store{
		db:            db,
		customDB:      customDB,
		hydrator:      hydrator,
		propertyCache: propertyCache,
		sorter: database.New(
			database.SpecMap{
				"model": database.Column{Col: table.FivenetVehicle.AS("vehicle").Model},
				"plate": database.Column{Col: table.FivenetVehicle.AS("vehicle").Plate},
			},
			[]mysql.OrderByClause{table.FivenetVehicle.AS("vehicle").Plate.ASC()},
			[]mysql.OrderByClause{table.FivenetVehicle.AS("vehicle").Type.ASC()},
			"plate",
			3,
		),
		vehicleActivitySorter: database.New(
			database.SpecMap{
				"createdAt": database.Column{
					Col:       table.FivenetVehicleActivity.AS("vehicle_activity").CreatedAt,
					NullsLast: true,
				},
			},
			[]mysql.OrderByClause{
				table.FivenetVehicleActivity.AS("vehicle_activity").CreatedAt.DESC().NULLS_LAST(),
			},
			[]mysql.OrderByClause{table.FivenetVehicleActivity.AS("vehicle_activity").ID.DESC()},
			"createdAt",
			3,
		),
	}
}
