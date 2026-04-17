package syncers

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
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

func (s *VehiclesSync) Sync(ctx context.Context) (int64, int64, string, *time.Time, error) {
	limit := s.cfg.Limits.Vehicles
	windowEnd := time.Now()
	batchCap, lag := calculateDrainBatchCap(
		s.state.GetLastCheck(),
		windowEnd,
		maxDrainBatchesPerSync,
	)
	if lag >= lagWarnThreshold {
		s.logger.Warn(
			"vehicles cursor lag is high",
			zap.Duration("lag", lag),
			zap.Time("window_end", windowEnd),
			zap.Timep("cursor_time", s.state.GetLastCheck()),
			zap.Int("drain_batch_cap", batchCap),
		)
	}

	var totalFetched int64
	var totalSent int64
	lastID := ""
	var lastUpdatedAt *time.Time
	prevID := ""
	var prevUpdatedAt *time.Time
	noopStreak := 0

	for batches := 0; ; batches++ {
		fetched, sent, cursorID, cursorTime, err := s.syncOnce(ctx, &windowEnd)
		if err != nil {
			return totalFetched, totalSent, lastID, lastUpdatedAt, err
		}

		totalFetched += fetched
		totalSent += sent
		if cursorID != "" {
			lastID = cursorID
		}
		if cursorTime != nil {
			lastUpdatedAt = cursorTime
		}

		if fetched < limit {
			break
		}

		if batches+1 >= batchCap {
			s.logger.Info(
				"vehicles sync hit drain batch cap; remaining updates continue next interval",
				zap.Int64("fetched", fetched),
				zap.Int64("sent", sent),
				zap.String("cursor_id", cursorID),
				zap.Int("drain_batch_cap", batchCap),
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

		// Guard against repeated no-op batches where cursor keeps advancing but no rows are sent.
		if sent == 0 {
			noopStreak++
			if noopStreak >= maxNoopBatchesPerSync {
				s.logger.Info(
					"vehicles sync hit no-op batch cap; stopping drain loop",
					zap.Int("noop_streak", noopStreak),
					zap.Int64("fetched", fetched),
					zap.Int64("sent", sent),
					zap.String("cursor_id", cursorID),
					zap.Timep("cursor_time", cursorTime),
				)
				break
			}
		} else {
			noopStreak = 0
		}

		prevID = cursorID
		if cursorTime != nil {
			t := *cursorTime
			prevUpdatedAt = &t
		} else {
			prevUpdatedAt = nil
		}
	}

	return totalFetched, totalSent, lastID, lastUpdatedAt, nil
}

func (s *VehiclesSync) Resync(ctx context.Context) (int64, int64, string, *time.Time, error) {
	// Ensure last check is nil when we don't want to save it
	if !s.saveUpdatedAt {
		s.state.SetLastCheck(nil)
	}

	fetched, sent, cursorID, cursorTime, err := s.syncOnce(ctx, nil)
	return fetched, sent, cursorID, cursorTime, err
}

func (s *VehiclesSync) syncOnce(
	ctx context.Context,
	windowEnd *time.Time,
) (int64, int64, string, *time.Time, error) {
	limit := s.cfg.Limits.Vehicles
	sQuery := s.cfg.Tables.Vehicles
	q := sQuery.GetSyncQuery(
		s.state,
		limit,
		updatedAtUpperBoundCondition(sQuery.UpdatedTimeColumn, windowEnd),
	)
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
		skippedPlates := make([]string, 0, len(vehicles))
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
					skippedPlates = append(skippedPlates, vehicle.GetPlate())
					// Remove "skipped" vehicle
					vehicles = slices.Delete(vehicles, i, i+1)
					continue
				}
			}
			s.hashes.Put(vehicle.GetPlate(), hash, vehicleHashCacheTTL)
		}
		if len(skippedPlates) > 0 {
			s.logger.Debug(
				"vehicles skipped unchanged entries by hash",
				zap.Int("count", len(skippedPlates)),
				zap.Strings("plates", skippedPlates),
			)
		}
	}

	// Sync vehicles to FiveNet server (if there are any left after hash check)
	if len(vehicles) > 0 {
		for start := 0; start < len(vehicles); start += sync.MaxVehiclesPerRequest {
			end := min(start+sync.MaxVehiclesPerRequest, len(vehicles))
			req := &pbsync.SendVehiclesRequest{
				Vehicles: vehicles[start:end],
			}
			if err := s.send(
				ctx,
				req,
				func(ctx context.Context, cli pbsync.SyncServiceClient) error {
					_, err := cli.SendVehicles(ctx, req)
					return err
				},
			); err != nil {
				return 0, 0, "", nil, fmt.Errorf(
					"failed to send vehicles data to FiveNet server. %w",
					err,
				)
			}
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

	t := ts.GetTimestamp().AsTime().In(time.Local)
	t = t.Truncate(time.Millisecond)
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
