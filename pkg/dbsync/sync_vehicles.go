package dbsync

import (
	"context"
	"errors"
	"fmt"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

type vehiclesSync struct {
	*syncer

	state *dbsyncconfig.TableSyncState
}

func newVehiclesSync(s *syncer, state *dbsyncconfig.TableSyncState) *vehiclesSync {
	return &vehiclesSync{
		syncer: s,
		state:  state,
	}
}

func (s *vehiclesSync) Sync(ctx context.Context) error {
	if !s.cfg.Tables.Vehicles.Enabled {
		return nil
	}

	limit := int64(500)
	var offset int64
	sOffset := s.state.GetOffset()
	if s.state != nil && sOffset > 0 {
		offset = sOffset
	}
	s.logger.Debug("vehiclesSync", zap.Int64("offset", offset))

	// Ensure to zero the last check time if the data hasn't fully synced yet
	if !s.state.GetSyncedUp() {
		s.state.SetLastCheck(nil)
	}

	q := s.cfg.Tables.Vehicles.GetQuery(s.state, offset, limit)
	s.logger.Debug("vehicles sync query", zap.String("query", q))

	vehicles := []*vehicles.Vehicle{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &vehicles); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	s.logger.Debug("vehiclesSync", zap.Int("len", len(vehicles)))

	if len(vehicles) == 0 {
		s.state.Set(0, nil)
		return nil
	}

	// Sync vehicles to FiveNet server
	if s.cli != nil {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Vehicles{
				Vehicles: &syncdata.DataVehicles{
					Vehicles: vehicles,
				},
			},
		}); err != nil {
			return fmt.Errorf("failed to send vehicles data to FiveNet server. %w", err)
		}
	}

	// If less vehicles than limit are returned, we probably have reached the "end" of the table
	// and need to reset the offset to 0
	if int64(len(vehicles)) < limit {
		offset = 0
		s.state.SetSyncedUp(true)
	}

	lastPlate := vehicles[len(vehicles)-1].GetPlate()
	s.state.Set(limit+offset, &lastPlate)

	return nil
}
