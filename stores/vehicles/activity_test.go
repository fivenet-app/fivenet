package vehiclesstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountVehicleActivityAppliesPlateAndTypeFilter(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_vehicles_activity AS vehicle_activity`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle_activity.plate = ?`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle_activity.type IN (?)`)
	mock.ExpectQuery(expectedQuery).
		WithArgs("ABC DEF1", int32(vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(3)))

	total, err := store.CountVehicleActivity(t.Context(), CountVehicleActivityOptions{
		VehicleActivityOptions: VehicleActivityOptions{
			Plate: "ABC DEF1",
			Types: []vehiclesactivity.VehicleActivityType{
				vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED,
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListVehicleActivityAppliesSortAndCreatorJoin(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_vehicles_activity AS vehicle_activity`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_user AS creator ON`) +
		`(?s).*` + regexp.QuoteMeta(`vehicle_activity.plate = ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY vehicle_activity.created_at DESC, vehicle_activity.id DESC LIMIT ? OFFSET ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs("ABC DEF1", int64(20), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"vehicle_activity.id",
			"vehicle_activity.created_at",
			"vehicle_activity.plate",
			"vehicle_activity.type",
			"vehicle_activity.creator_id",
			"vehicle_activity.creator_job",
			"vehicle_activity.reason",
			"vehicle_activity.data",
			"creator.id",
			"creator.job",
			"creator.job_grade",
			"creator.firstname",
			"creator.lastname",
		}).AddRow(
			int64(1),
			time.Now(),
			"ABC DEF1",
			int32(vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED),
			int32(7),
			"police",
			"updated",
			[]byte(`{"wantedChange":{"wanted":true}}`),
			int32(7),
			"police",
			int32(2),
			"Jane",
			"Doe",
		))

	activity, err := store.ListVehicleActivity(t.Context(), ListVehicleActivityOptions{
		VehicleActivityOptions: VehicleActivityOptions{Plate: "ABC DEF1"},
		Offset:                 0,
		Limit:                  20,
	})
	require.NoError(t, err)
	require.Len(t, activity, 1)
	assert.Equal(t, "ABC DEF1", activity[0].GetPlate())
	assert.Equal(t, int32(7), activity[0].GetCreator().GetUserId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdatePropsWritesWantedActivity(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*plate.*LIMIT \?`).
		WithArgs("ABC DEF1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"vehicle_props.plate", "vehicle_props.wanted", "vehicle_props.wanted_reason", "vehicle_props.wanted_at", "vehicle_props.wanted_till"}).
			AddRow("ABC DEF1", false, nil, nil, nil))
	mock.ExpectBegin()
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_vehicles_props.*ON DUPLICATE KEY UPDATE.*wanted = \?.*wanted_till = .*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT INTO fivenet_vehicles_activity.*`).
		WithArgs(int32(7), "ABC DEF1", int32(vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED), "police", "stolen", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*plate.*LIMIT \?`).
		WithArgs("ABC DEF1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"vehicle_props.plate", "vehicle_props.wanted", "vehicle_props.wanted_reason", "vehicle_props.wanted_at", "vehicle_props.wanted_till"}).
			AddRow("ABC DEF1", true, "stolen", time.Now(), nil))

	creatorID := int32(7)
	wanted := true
	reason := "stolen"
	props, err := store.UpdateProps(t.Context(), &vehiclesprops.VehicleProps{
		Plate:        "ABC DEF1",
		Wanted:       &wanted,
		WantedReason: &reason,
	}, &creatorID, "police", "stolen")
	require.NoError(t, err)
	assert.True(t, props.GetWanted())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdatePropsNoOpDoesNotWriteWantedActivity(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	wantedAt := time.Now().UTC()
	wantedTill := time.Now().UTC().Add(time.Hour)
	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*plate.*LIMIT \?`).
		WithArgs("ABC DEF1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"vehicle_props.plate", "vehicle_props.wanted", "vehicle_props.wanted_reason", "vehicle_props.wanted_at", "vehicle_props.wanted_till"}).
			AddRow("ABC DEF1", true, "stolen", wantedAt, wantedTill))
	mock.ExpectBegin()
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_vehicles_props.*ON DUPLICATE KEY UPDATE.*wanted = \?.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*plate.*LIMIT \?`).
		WithArgs("ABC DEF1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"vehicle_props.plate", "vehicle_props.wanted", "vehicle_props.wanted_reason", "vehicle_props.wanted_at", "vehicle_props.wanted_till"}).
			AddRow("ABC DEF1", true, "stolen", wantedAt, wantedTill))

	wanted := true
	reason := "stolen"
	creatorID := int32(7)
	_, err = store.UpdateProps(t.Context(), &vehiclesprops.VehicleProps{
		Plate:        "ABC DEF1",
		Wanted:       &wanted,
		WantedReason: &reason,
		WantedTill:   timestamp.New(wantedTill),
	}, &creatorID, "police", "stolen")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
