package housekeeper

import (
	"testing"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUnitPingKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		key string
		id  int64
		ok  bool
	}{
		{key: "ping.42", id: 42, ok: true},
		{key: "ping.0"},
		{key: "ping.-1"},
		{key: "ping.nope"},
		{key: "unit.42"},
		{key: "ping."},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			t.Parallel()

			id, err := parseUnitPingKey(tt.key)
			if tt.ok {
				require.NoError(t, err)
				assert.Equal(t, tt.id, id)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestUnitHasOnlyValidMembers(t *testing.T) {
	t.Parallel()

	unit := &centrumunits.Unit{
		Users: []*centrumunits.UnitAssignment{
			{UserId: 1},
		},
	}

	assert.True(t, unitHasOnlyValidMembers(unit, []int32{1}))
	assert.False(t, unitHasOnlyValidMembers(unit, nil))
	assert.False(t, unitHasOnlyValidMembers(&centrumunits.Unit{}, []int32{1}))
	assert.False(t, unitHasOnlyValidMembers(nil, []int32{1}))
}
