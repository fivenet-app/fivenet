package vehiclesstore

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeWantedChangeUsesServerReasonAndNewWantedAt(t *testing.T) {
	t.Parallel()

	currentWanted := false
	current := &vehiclesprops.VehicleProps{Wanted: &currentWanted}
	wanted := true
	in := &vehiclesprops.VehicleProps{Wanted: &wanted}

	normalizeWantedChange(current, in, "stolen")

	assert.Equal(t, "stolen", in.GetWantedReason())
	assert.NotNil(t, in.GetWantedAt())
	assert.Nil(t, in.GetWantedTill())
}

func TestNormalizeWantedChangePreservesExistingWantedState(t *testing.T) {
	t.Parallel()

	wanted := true
	reason := "stolen"
	wantedAt := timestamp.New(time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC))
	wantedTill := timestamp.New(time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC))
	current := &vehiclesprops.VehicleProps{
		Wanted:       &wanted,
		WantedReason: &reason,
		WantedAt:     wantedAt,
		WantedTill:   wantedTill,
	}
	inWanted := true
	in := &vehiclesprops.VehicleProps{Wanted: &inWanted}

	normalizeWantedChange(current, in, "server reason")

	assert.Equal(t, wantedAt, in.GetWantedAt())
	assert.Equal(t, wantedTill, in.GetWantedTill())
	assert.Equal(t, "server reason", in.GetWantedReason())
}

func TestNormalizeWantedChangeClearsRevokedWantedState(t *testing.T) {
	t.Parallel()

	wanted := true
	reason := "stolen"
	current := &vehiclesprops.VehicleProps{Wanted: &wanted, WantedReason: &reason}
	inWanted := false
	in := &vehiclesprops.VehicleProps{Wanted: &inWanted}

	normalizeWantedChange(current, in, "ignored")

	require.False(t, in.GetWanted())
	assert.Nil(t, in.WantedReason)
	assert.Nil(t, in.WantedAt)
	assert.Nil(t, in.WantedTill)
}
