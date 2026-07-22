package vehiclesprops

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeWantedChangeSetsReasonAsWantedReason(t *testing.T) {
	t.Parallel()

	currentWanted := false
	current := &VehicleProps{
		Plate:  "ABC DEF1",
		Wanted: &currentWanted,
	}

	wanted := true
	in := &VehicleProps{
		Plate:  "ABC DEF1",
		Wanted: &wanted,
	}

	current.NormalizeWantedChange(in, "stolen")

	require.NotNil(t, in.WantedReason)
	assert.Equal(t, "stolen", in.GetWantedReason())
	assert.NotNil(t, in.GetWantedAt())
	assert.Nil(t, in.GetWantedTill())
}

func TestNormalizeWantedChangeReasonOverridesClientWantedReason(t *testing.T) {
	t.Parallel()

	currentWanted := true
	currentReason := "old reason"
	wantedAt := timestamp.New(time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC))
	wantedTill := timestamp.New(time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC))
	current := &VehicleProps{
		Plate:        "ABC DEF1",
		Wanted:       &currentWanted,
		WantedReason: &currentReason,
		WantedAt:     wantedAt,
		WantedTill:   wantedTill,
	}

	clientReason := "client supplied reason"
	wanted := true
	in := &VehicleProps{
		Plate:        "ABC DEF1",
		Wanted:       &wanted,
		WantedReason: &clientReason,
	}

	current.NormalizeWantedChange(in, "server reason")

	require.NotNil(t, in.WantedReason)
	assert.Equal(t, "server reason", in.GetWantedReason())
	assert.Equal(t, wantedAt, in.GetWantedAt())
	assert.Equal(t, wantedTill, in.GetWantedTill())
}

func TestNormalizeWantedChangePreservesWantedAtWhenAlreadyWanted(t *testing.T) {
	t.Parallel()

	currentWanted := true
	currentReason := "stolen"
	wantedAt := timestamp.New(time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC))
	clientWantedAt := timestamp.New(time.Date(2027, 7, 22, 10, 0, 0, 0, time.UTC))
	current := &VehicleProps{
		Plate:        "ABC DEF1",
		Wanted:       &currentWanted,
		WantedReason: &currentReason,
		WantedAt:     wantedAt,
	}

	wanted := true
	in := &VehicleProps{
		Plate:    "ABC DEF1",
		Wanted:   &wanted,
		WantedAt: clientWantedAt,
	}

	current.NormalizeWantedChange(in, "server reason")

	assert.Equal(t, wantedAt, in.GetWantedAt())
}

func TestNormalizeWantedChangeRevokesWantedState(t *testing.T) {
	t.Parallel()

	currentWanted := true
	currentReason := "stolen"
	current := &VehicleProps{
		Plate:        "ABC DEF1",
		Wanted:       &currentWanted,
		WantedReason: &currentReason,
		WantedAt:     timestamp.New(time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC)),
		WantedTill:   timestamp.New(time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC)),
	}

	wanted := false
	in := &VehicleProps{
		Plate:  "ABC DEF1",
		Wanted: &wanted,
	}

	current.NormalizeWantedChange(in, "no longer wanted")

	assert.False(t, in.GetWanted())
	assert.Nil(t, in.WantedReason)
	assert.Nil(t, in.GetWantedAt())
	assert.Nil(t, in.GetWantedTill())
}

func TestNormalizeWantedChangeCopiesCurrentWhenWantedNotChanged(t *testing.T) {
	t.Parallel()

	currentWanted := true
	currentReason := "stolen"
	wantedAt := timestamp.New(time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC))
	wantedTill := timestamp.New(time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC))
	current := &VehicleProps{
		Plate:        "ABC DEF1",
		Wanted:       &currentWanted,
		WantedReason: &currentReason,
		WantedAt:     wantedAt,
		WantedTill:   wantedTill,
	}

	in := &VehicleProps{Plate: "ABC DEF1"}

	current.NormalizeWantedChange(in, "ignored reason")

	require.NotNil(t, in.Wanted)
	assert.True(t, in.GetWanted())
	assert.Equal(t, "stolen", in.GetWantedReason())
	assert.Equal(t, wantedAt, in.GetWantedAt())
	assert.Equal(t, wantedTill, in.GetWantedTill())
}
