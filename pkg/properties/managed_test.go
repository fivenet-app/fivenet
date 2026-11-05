package properties

import (
	"testing"
	"time"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestApplyWantedMutations(t *testing.T) {
	t.Parallel()
	expiresAt := timestamp.New(time.Now().UTC().Add(time.Hour))
	value := &propertypb.PropertyValue{}
	value.SetBoolValue(true)
	mutation := &propertypb.PropertyMutation{
		Key:   "user.wanted",
		State: &propertypb.PropertyStatePatch{},
	}
	mutation.SetSetValue(value)
	mutation.State.SetExpiresAt(expiresAt)

	user := &usersprops.UserProps{}
	require.NoError(t, ApplyWantedMutation(user, mutation))
	assert.True(t, user.GetWanted())
	assert.Equal(t, expiresAt, user.GetWantedTill())

	clearPropMutation := &propertypb.PropertyMutation{Key: "vehicle.wanted"}
	clearPropMutation.SetClearValue(&emptypb.Empty{})
	vehicle := &vehiclesprops.VehicleProps{}
	require.NoError(t, ApplyWantedMutation(vehicle, clearPropMutation))
	assert.False(t, vehicle.GetWanted())
}

func TestLifecyclePatchForMutation(t *testing.T) {
	t.Parallel()
	state := &propertypb.PropertyStatePatch{}
	state.SetClearExpiry(&emptypb.Empty{})
	mutation := &propertypb.PropertyMutation{Key: "user.note", State: state}

	expiresAt, cleared, hasState := LifecyclePatchForMutation(mutation)
	assert.Nil(t, expiresAt)
	assert.True(t, cleared)
	assert.True(t, hasState)
}

func TestClampWantedExpiry(t *testing.T) {
	t.Parallel()

	maxExpiry := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	later := timestamp.New(maxExpiry.Add(time.Hour))
	mutation := &propertypb.PropertyMutation{
		Key:   "user.wanted",
		State: &propertypb.PropertyStatePatch{},
	}
	mutation.State.SetExpiresAt(later)

	ClampWantedExpiry(
		propertypb.PropertyScope_PROPERTY_SCOPE_USER,
		[]*propertypb.PropertyMutation{mutation},
		maxExpiry,
	)
	assert.Equal(t, maxExpiry, mutation.GetState().GetExpiresAt().AsTime())
}

func TestClampWantedExpiryKeepsClear(t *testing.T) {
	t.Parallel()

	mutation := &propertypb.PropertyMutation{
		Key:   "vehicle.wanted",
		State: &propertypb.PropertyStatePatch{},
	}
	mutation.State.SetClearExpiry(&emptypb.Empty{})

	ClampWantedExpiry(
		propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		[]*propertypb.PropertyMutation{mutation},
		time.Now().UTC(),
	)
	assert.NotNil(t, mutation.GetState().GetClearExpiry())
	assert.Nil(t, mutation.GetState().GetExpiresAt())
}

func TestManagedMutationRegistry(t *testing.T) {
	t.Parallel()

	user, ok := ManagedMutation(propertypb.PropertyScope_PROPERTY_SCOPE_USER, "user.wanted")
	require.True(t, ok)
	assert.Equal(t, "wanted", user.HandlerName)

	vehicle, ok := ManagedMutation(
		propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		"vehicle.wanted",
	)
	require.True(t, ok)
	assert.Equal(t, "wanted", vehicle.HandlerName)

	_, ok = ManagedMutation(propertypb.PropertyScope_PROPERTY_SCOPE_USER, "user.note")
	assert.False(t, ok)
}
