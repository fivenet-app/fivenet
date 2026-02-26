package syncers

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/cache"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

const vehicleHashCacheTTL = 4 * time.Hour

type VehiclesSync struct {
	*Syncer

	state         *dbsyncconfig.TableSyncState
	saveLastCheck bool

	hashes *cache.LRUCache[string, uint64]
}

func NewVehiclesSync(
	s *Syncer,
	state *dbsyncconfig.TableSyncState,
	saveLastCheck bool,
) *VehiclesSync {
	return &VehiclesSync{
		Syncer: s,

		state:         state,
		saveLastCheck: saveLastCheck,

		hashes: cache.NewLRUCache[string, uint64](500),
	}
}

func (s *VehiclesSync) Sync(ctx context.Context) (int64, int64, string, error) {
	// Ensure last check is nil when we don't want to save it
	if !s.saveLastCheck {
		s.state.SetLastCheck(nil)
	}

	limit := s.cfg.Limits.Vehicles
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
			return 0, offset, "", err
		}
	}

	count := int64(len(vehicles))
	s.logger.Debug("vehiclesSync", zap.Int64("len", count))
	if len(vehicles) == 0 {
		s.state.Set(0, nil)
		return 0, offset, "", nil
	}

	// Sync vehicles to FiveNet server
	if err := s.sendData(ctx, &pbsync.SendDataRequest{
		Data: &pbsync.SendDataRequest_Vehicles{
			Vehicles: &syncdata.DataVehicles{
				Vehicles: vehicles,
			},
		},
	}); err != nil {
		return 0, offset, "", fmt.Errorf(
			"failed to send vehicles data to FiveNet server. %w",
			err,
		)
	}

	for i, v := range slices.Backward(vehicles) {
		// Get hash of user data to compare with existing hash and skip sending if data is the same (treat as not updated)
		_, hash, err := protoutils.JSONAndHash(v)
		if err != nil {
			s.logger.Warn(
				"failed to compute vehicle data hash, skipping hash check and treating as new/updated vehicle",
				zap.String("plate", v.GetPlate()),
				zap.Error(err),
			)
		}

		if existingHash, ok := s.hashes.Get(v.GetPlate()); ok {
			if existingHash == hash {
				s.logger.Debug(
					"vehicle data hash is the same as existing entry, skipping update for vehicle",
					zap.String("plate", v.GetPlate()),
				)
				// Remove "skipped" vehicle
				vehicles = slices.Delete(vehicles, i, i+1)
				continue
			}
		} else {
			s.hashes.Put(v.GetPlate(), hash, vehicleHashCacheTTL)
		}
	}

	// If less vehicles than limit are returned, we probably have reached the "end" of the table
	// and need to reset the offset to 0
	if count < limit {
		offset = 0
		s.state.SetSyncedUp(true)
	}

	count = int64(len(vehicles))
	lastPlate := vehicles[count-1].GetPlate()
	s.state.Set(limit+offset, &lastPlate)

	return count, offset, lastPlate, nil
}

func (s *VehiclesSync) Resync(ctx context.Context) error {
	return nil
}
