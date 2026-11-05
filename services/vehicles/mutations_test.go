package vehicles

import (
	"testing"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestVehicleWantedStateClearReason(t *testing.T) {
	t.Parallel()
	assert.True(t, vehicleWantedStateClearReason([]*properties.PropertyMutation{{
		Key:         "vehicle.wanted",
		ClearReason: new(true),
	}}))
	assert.False(t, vehicleWantedStateClearReason([]*properties.PropertyMutation{{
		Key:   "vehicle.wanted",
		State: &properties.PropertyStatePatch{},
	}}))
}

func TestVehicleScalarClearKeepsMutationReasonInActivity(t *testing.T) {
	t.Parallel()
	reason := "seized"
	mutation := &properties.PropertyMutation{Key: "custom.note", Reason: &reason}
	mutation.SetClearValue(&emptypb.Empty{})
	activities := vehicleMutationActivities([]*properties.PropertyMutation{mutation})
	change := activities["custom.note"]
	if assert.NotNil(t, change) {
		assert.Equal(t, "seized", change.GetLifecycle().GetStateReason())
	}
}
