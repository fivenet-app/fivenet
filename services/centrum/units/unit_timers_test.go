package units

import (
	"errors"
	"testing"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	testnats "github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitNeedsPing(t *testing.T) {
	t.Parallel()

	status := func(status centrumunits.StatusUnit) *centrumunits.UnitStatus {
		return &centrumunits.UnitStatus{Status: status}
	}

	tests := []struct {
		name string
		unit *centrumunits.Unit
		want bool
	}{
		{name: "nil", want: false},
		{
			name: "empty unavailable",
			unit: &centrumunits.Unit{Status: status(centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE)},
			want: false,
		},
		{
			name: "empty available",
			unit: &centrumunits.Unit{Status: status(centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE)},
			want: true,
		},
		{
			name: "empty without status",
			unit: &centrumunits.Unit{},
			want: true,
		},
		{
			name: "occupied unavailable",
			unit: &centrumunits.Unit{
				Status: &centrumunits.UnitStatus{Status: centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE},
				Users:  []*centrumunits.UnitAssignment{{UserId: 1}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, unitNeedsPing(tt.unit))
		})
	}
}

func TestSyncUnitPingCreatesOnlyNeededTimers(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	js := testnats.NewServer(t, testnats.ServerOptions{InProcess: true}).GetJS()
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:         "test_unit_ping",
		Storage:        jetstream.MemoryStorage,
		History:        1,
		LimitMarkerTTL: 2 * PingTTL,
	})
	require.NoError(t, err)

	db := &UnitDB{KVPing: kv}
	unit := &centrumunits.Unit{
		Id: 42,
		Status: &centrumunits.UnitStatus{
			Status: centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE,
		},
	}

	require.NoError(t, db.SyncUnitPing(ctx, unit))
	_, err = kv.Get(ctx, "ping.42")
	assert.ErrorIs(t, err, jetstream.ErrKeyNotFound)

	unit.Status.Status = centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE
	require.NoError(t, db.SyncUnitPing(ctx, unit))
	_, err = kv.Get(ctx, "ping.42")
	require.NoError(t, err)

	unit.Status.Status = centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE
	unit.Users = []*centrumunits.UnitAssignment{{UserId: 7}}
	require.NoError(t, db.SyncUnitPing(ctx, unit))
	_, err = kv.Get(ctx, "ping.42")
	require.NoError(t, err)

	unit.Users = nil
	require.NoError(t, db.SyncUnitPing(ctx, unit))
	_, err = kv.Get(ctx, "ping.42")
	assert.True(t, errors.Is(err, jetstream.ErrKeyNotFound))
}
