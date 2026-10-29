package vehiclesstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountAppliesListFilters(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db, &config.CustomDB{
		Columns: dbutils.CustomColumns{
			Vehicle: dbutils.VehicleColumns{Model: "model"},
		},
	}, nil)

	wanted := &properties.PropertyValue{}
	wanted.SetBoolValue(true)
	query := ListQuery{
		LicensePlate: "ABC",
		Model:        "adder",
		UserIDs:      []int32{3, 4},
		Job:          "police",
		PropertyFilters: []*properties.PropertyFilter{{
			Key:       vehicleWantedPropertyKey,
			Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS,
			Value:     wanted,
		}},
		CanFilterProperties: true,
	}

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_vehicle AS vehicle`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle.plate LIKE ?`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle.model LIKE ?`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle.user_id IN (?, ?)`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_vehicle_property_values AS property_filter_0`) +
		`(?s).*` + regexp.QuoteMeta("property_filter_0.plate = `vehicle`.plate") +
		`(?s).*` + regexp.QuoteMeta(`property_filter_0.value_type = ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(true, "%ABC%", "%adder%", int32(3), int32(4), vehicleWantedPropertyKey, propertyValueTypeBool, true, "vehicle", vehicleWantedPropertyKey).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(7)))

	total, err := store.Count(t.Context(), query)
	require.NoError(t, err)
	assert.Equal(t, int64(7), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListAppliesSortFallbackAndTieBreaker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db, &config.CustomDB{
		Columns: dbutils.CustomColumns{
			Vehicle: dbutils.VehicleColumns{Model: dbutils.DisableColumnName},
		},
	}, nil)

	query := ListQuery{
		Sort: &database.Sort{
			Columns: []*database.SortByColumn{
				{Id: "unknown", Desc: true},
			},
		},
		Limit: 20,
	}

	expectedQuery := regexp.QuoteMeta(
		`ORDER BY vehicle.plate DESC, vehicle.type ASC LIMIT ? OFFSET ?;`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs("_", " ", true, int64(20), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"vehicle.plate"}).AddRow("ABC DEF1"))

	vehicles, err := store.List(t.Context(), query)
	require.NoError(t, err)
	require.Len(t, vehicles, 1)
	assert.Equal(t, "ABC DEF1", vehicles[0].GetPlate())
	require.NoError(t, mock.ExpectationsWereMet())
}
