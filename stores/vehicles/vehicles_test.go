package vehiclesstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
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
	})

	wanted := true
	query := ListQuery{
		LicensePlate:    "ABC",
		Model:           "adder",
		UserIDs:         []int32{3, 4},
		Job:             "police",
		Wanted:          &wanted,
		CanFilterWanted: true,
	}

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_owned_vehicles AS vehicle`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_user AS user_short ON`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.id IN (?, ?)`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle.plate LIKE ?`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle.model LIKE ?`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.id IN (?, ?)`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle_props.wanted = ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int32(3), int32(4), true, "%ABC%", "%adder%", int32(3), int32(4), true).
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
	})

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
