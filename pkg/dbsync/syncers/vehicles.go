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

// vehicleHashCacheTTL Cache vehicles hashes for 4 hours to avoid keeping stale hashes for vehicles that have been updated.
const vehicleHashCacheTTL = 4 * time.Hour

type VehiclesSync struct {
	*Syncer

	logger        *zap.Logger
	state         *dbsyncconfig.TableSyncState
	saveUpdatedAt bool

	hashes *cache.LRUCache[string, uint64]
}

func NewVehiclesSync(
	s *Syncer,
	state *dbsyncconfig.TableSyncState,
	saveUpdatedAt bool,
) *VehiclesSync {
	var hashes *cache.LRUCache[string, uint64]
	if saveUpdatedAt {
		// Cache up to 500 vehicle hashes to avoid memory bloat, as this is only used to compare against
		// the most recent hash for each vehicle and not all historical hashes
		hashes = cache.NewLRUCache[string, uint64](500)

		// Ensure a sane last check value is set for the "update" sync to work immediately
		if state.GetLastCheck() == nil {
			initialLastCheck := time.Now().Add(-15 * time.Minute)
			state.SetLastCheck(&initialLastCheck)
		}
	}

	logger := s.logger.With(
		zap.String("syncer", "vehicles"),
		zap.Bool("resync", !saveUpdatedAt),
	)

	return &VehiclesSync{
		Syncer: s,

		logger:        logger,
		state:         state,
		saveUpdatedAt: saveUpdatedAt,

		hashes: hashes,
	}
}

func (s *VehiclesSync) Sync(ctx context.Context) (int64, string, *time.Time, error) {
	limit := s.cfg.Limits.Vehicles

	var total int64
	lastID := ""
	var lastUpdatedAt *time.Time
	prevID := ""
	var prevUpdatedAt *time.Time

	for batches := 0; ; batches++ {
		fetched, sent, cursorID, cursorTime, err := s.syncOnce(ctx)
		if err != nil {
			return total, lastID, lastUpdatedAt, err
		}

		total += sent
		if cursorID != "" {
			lastID = cursorID
		}
		if cursorTime != nil {
			lastUpdatedAt = cursorTime
		}

		if fetched < limit {
			break
		}

		if batches+1 >= maxDrainBatchesPerSync {
			s.logger.Info(
				"vehicles sync hit drain batch cap; remaining updates continue next interval",
				zap.Int64("fetched", fetched),
				zap.Int64("sent", sent),
				zap.String("cursor_id", cursorID),
			)
			break
		}

		sameTime := (prevUpdatedAt == nil && cursorTime == nil) ||
			(prevUpdatedAt != nil && cursorTime != nil && prevUpdatedAt.Equal(*cursorTime))
		if cursorID != "" && cursorID == prevID && sameTime {
			s.logger.Warn(
				"vehicles sync cursor did not advance, stopping drain loop",
				zap.String("cursor_id", cursorID),
				zap.Timep("cursor_time", cursorTime),
				zap.Int64("fetched", fetched),
				zap.Int64("sent", sent),
			)
			break
		}

		prevID = cursorID
		if cursorTime != nil {
			t := *cursorTime
			prevUpdatedAt = &t
		} else {
			prevUpdatedAt = nil
		}
	}

	return total, lastID, lastUpdatedAt, nil
}

func (s *VehiclesSync) Resync(ctx context.Context) (int64, string, *time.Time, error) {
	// Ensure last check is nil when we don't want to save it
	if !s.saveUpdatedAt {
		s.state.SetLastCheck(nil)
	}

	_, sent, cursorID, cursorTime, err := s.syncOnce(ctx)
	return sent, cursorID, cursorTime, err
}

func (s *VehiclesSync) syncOnce(
	ctx context.Context,
) (fetched int64, sent int64, cursorID string, cursorTime *time.Time, err error) {
	limit := s.cfg.Limits.Vehicles
	sQuery := s.cfg.Tables.Vehicles
	q := sQuery.GetQuery(s.state, 0, limit)
	s.logger.Debug("vehicles sync query", zap.String("query", q))

	vehicles := []*vehicles.Vehicle{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &vehicles); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, 0, "", nil, err
		}
	}

	fetchedCount := int64(len(vehicles))
	s.logger.Debug("vehiclesSync", zap.Int64("len", fetchedCount))
	if len(vehicles) == 0 {
		if !s.saveUpdatedAt {
			s.state.ResetCursor()
		}
		return 0, 0, "", nil, nil
	}

	cursorTime, cursorLastPlate := s.cursorFromVehiclesResults(vehicles)
	if s.saveUpdatedAt && cursorTime == nil {
		return 0, 0, "", nil, errors.New(
			"vehicles result is missing updated_at, cannot persist cursor timestamp",
		)
	}

	if s.hashes != nil {
		for i, vehicle := range slices.Backward(vehicles) {
			vehicle.UpdatedAt = nil

			// Get hash of user data to compare with existing hash and skip sending if data is the same (treat as not updated)
			_, hash, err := protoutils.JSONAndHash(vehicle)
			if err != nil {
				s.logger.Warn(
					"failed to compute vehicle data hash, skipping hash check and treating as new/updated vehicle",
					zap.String("plate", vehicle.GetPlate()),
					zap.Error(err),
				)
			}

			if existingHash, ok := s.hashes.Get(vehicle.GetPlate()); ok {
				if existingHash == hash {
					s.logger.Debug(
						"vehicle data hash is the same as existing entry, skipping update for vehicle",
						zap.String("plate", vehicle.GetPlate()),
					)
					// Remove "skipped" vehicle
					vehicles = slices.Delete(vehicles, i, i+1)
					continue
				}
			} else {
				s.hashes.Put(vehicle.GetPlate(), hash, vehicleHashCacheTTL)
			}
		}
	}

	// Sync vehicles to FiveNet server (if there are any left after hash check)
	if len(vehicles) > 0 {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Vehicles{
				Vehicles: &syncdata.DataVehicles{
					Vehicles: vehicles,
				},
			},
		}); err != nil {
			return 0, 0, "", nil, fmt.Errorf(
				"failed to send vehicles data to FiveNet server. %w",
				err,
			)
		}
	}

	s.persistCursor(fetchedCount, limit, cursorTime, cursorLastPlate)

	lastPlate := ""
	if cursorLastPlate != nil {
		lastPlate = *cursorLastPlate
	}

	return fetchedCount, int64(len(vehicles)), lastPlate, cursorTime, nil
}

func (s *VehiclesSync) cursorFromVehiclesResults(
	vehicles []*vehicles.Vehicle,
) (*time.Time, *string) {
	if len(vehicles) == 0 {
		return nil, nil
	}

	last := vehicles[len(vehicles)-1]
	lastPlate := last.GetPlate()

	ts := last.GetUpdatedAt()
	if ts == nil || ts.GetTimestamp() == nil {
		return nil, &lastPlate
	}

	t := ts.GetTimestamp().AsTime()
	return &t, &lastPlate
}

func (s *VehiclesSync) persistCursor(
	fetchedCount int64,
	limit int64,
	cursorTime *time.Time,
	lastID *string,
) {
	if s.saveUpdatedAt {
		s.state.SetCursor(cursorTime, lastID)
		return
	}

	if fetchedCount < limit {
		s.state.ResetCursor()
		return
	}

	s.state.SetCursor(nil, lastID)
}
