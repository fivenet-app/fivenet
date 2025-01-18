package dbsync

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/vehicles"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/go-jet/jet/v2/qrm"
)

type vehiclesSync struct {
	*syncer

	state *TableSyncState
}

func newVehiclesSync(s *syncer, state *TableSyncState) *vehiclesSync {
	return &vehiclesSync{
		syncer: s,
		state:  state,
	}
}

func (s *vehiclesSync) Sync(ctx context.Context) error {
	if !s.cfg.Tables.Vehicles.Enabled {
		return nil
	}

	limit := 1000
	var offset uint64
	if s.state != nil && s.state.Offset > 0 {
		offset = s.state.Offset
	}

	// Ensure to zero the last check time if the data hasn't fully synced yet
	if !s.state.SyncedUp {
		s.state.LastCheck = nil
	}

	sQuery := s.cfg.Tables.Vehicles
	query := prepareStringQuery(sQuery, s.state, offset, limit)

	vehicles := []*vehicles.Vehicle{}
	if _, err := qrm.Query(ctx, s.db, query, []interface{}{}, &vehicles); err != nil {
		return err
	}

	if len(vehicles) == 0 {
		s.state.Set(0, nil)
		return nil
	}

	// Sync vehicles to FiveNet server
	if s.cli != nil {
		if _, err := s.cli.SendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Vehicles{
				Vehicles: &sync.DataVehicles{
					Vehicles: vehicles,
				},
			},
		}); err != nil {
			return err
		}
	}

	// If less vehicles than limit are returned, we probably have reached the "end" of the table
	// and need to reset the offset to 0
	if len(vehicles) < limit {
		offset = 0
		s.state.SyncedUp = true
	}

	lastPlate := vehicles[len(vehicles)-1].Plate
	s.state.Set(uint64(limit)+offset, &lastPlate)

	return nil
}
