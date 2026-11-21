package vehicles

import (
	"testing"

	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	"github.com/stretchr/testify/assert"
)

func TestVehicleActivityTypesForPropertyKeys(t *testing.T) {
	t.Parallel()

	activityTypes := []vehiclesactivity.VehicleActivityType{
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED,
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_FILE,
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_PROPERTY,
	}

	assert.Equal(t,
		[]vehiclesactivity.VehicleActivityType{
			vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED,
			vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_FILE,
		},
		vehicleActivityTypesForPropertyKeys(
			[]string{"vehicle.wanted", vehicleFilesPropertyKey, "vehicle.status"},
			activityTypes,
			[]string{"vehicle.wanted", vehicleFilesPropertyKey, "vehicle.status"},
		),
	)
	assert.Empty(t, vehicleActivityTypesForPropertyKeys(
		[]string{"vehicle.wanted"},
		activityTypes,
		[]string{vehicleFilesPropertyKey},
	))
}
