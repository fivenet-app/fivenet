package tracker

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/tracker"
	pblivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumstate"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/klauspost/compress/zstd"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	tLocs           = table.FivenetCentrumUserLocations
	tJobProps       = table.FivenetJobProps
	tColleagueProps = table.FivenetJobColleagueProps.AS("colleague_props")
)

type Manager struct {
	logger   *zap.Logger
	tracer   trace.Tracer
	js       *events.JSWrapper
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	postals  postals.Postals
	state    *centrumstate.State
	appCfg   appconfig.IConfig

	refreshTicker  *time.Ticker
	snapshotTicker *time.Ticker

	userLocStore      *store.Store[livemap.UserMarker, *livemap.UserMarker]
	byIDStore         *store.Store[livemap.UserMarker, *livemap.UserMarker]
	userMappingsStore *store.Store[tracker.UserMapping, *tracker.UserMapping]
}

type ManagerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	JS        *events.JSWrapper
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	Enricher  *mstlystcdata.Enricher
	Postals   postals.Postals
	Config    *config.Config
	State     *centrumstate.State
	AppConfig appconfig.IConfig
}

func NewManager(p ManagerParams) (*Manager, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	m := &Manager{
		logger:   p.Logger,
		tracer:   p.TP.Tracer("tracker.manager"),
		js:       p.JS,
		db:       p.DB,
		enricher: p.Enricher,
		postals:  p.Postals,
		state:    p.State,
		appCfg:   p.AppConfig,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		appCfg := p.AppConfig.Get()
		m.refreshTicker = time.NewTicker(appCfg.UserTracker.DbRefreshTime.AsDuration())

		go func() {
			configUpdateCh := p.AppConfig.Subscribe()
			for {
				select {
				case <-ctxCancel.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					if cfg == nil {
						continue
					}

					m.handleAppConfigUpdate(cfg)
				}
			}
		}()

		userMappingsStore, err := store.New[tracker.UserMapping, *tracker.UserMapping](
			ctxStartup, p.Logger, p.JS, BucketUserMappingsMap,
			store.WithLocks[tracker.UserMapping, *tracker.UserMapping](nil),
		)
		if err != nil {
			return fmt.Errorf("failed to create user mappings store. %w", err)
		}
		if err := userMappingsStore.Start(ctxCancel, false); err != nil {
			return fmt.Errorf("failed to start user mappings store. %w", err)
		}
		m.userMappingsStore = userMappingsStore

		byID, err := store.New[livemap.UserMarker, *livemap.UserMarker](
			ctxStartup, p.Logger, p.JS, BucketUserLocByID,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
		)
		if err != nil {
			return fmt.Errorf("failed to create user location by ID store. %w", err)
		}
		if err := byID.Start(ctxCancel, false); err != nil {
			return fmt.Errorf("failed to start user location by ID store. %w", err)
		}
		m.byIDStore = byID

		userLocStore, err := store.New[livemap.UserMarker, *livemap.UserMarker](
			ctxStartup, p.Logger, p.JS, BucketUserLoc,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
			store.WithIgnoredKeys[livemap.UserMarker, *livemap.UserMarker]("_snapshot"),
			store.WithOnUpdateFn(func(_ *store.Store[livemap.UserMarker, *livemap.UserMarker],
				um *livemap.UserMarker,
			) (*livemap.UserMarker, error) {
				if um == nil {
					return nil, nil
				}

				// Upsert mapping (unit_id may be 0 = no unit)
				if err := m.userMappingsStore.Put(ctxCancel, userIdKey(um.UserId), &tracker.UserMapping{
					UserId:    um.UserId,
					UnitId:    um.UnitId,
					CreatedAt: timestamp.Now(),
				}); err != nil {
					return nil, fmt.Errorf("failed to upsert user unit mapping. %w", err)
				}

				// Mirror latest marker by USER_ID (idempotent)
				if err := m.byIDStore.Put(ctxCancel, userIdKey(um.UserId), um); err != nil {
					return nil, fmt.Errorf("failed to upsert user location by ID. %w", err)
				}

				return um, nil
			}),
			store.WithOnDeleteFn(func(_ *store.Store[livemap.UserMarker, *livemap.UserMarker],
				entry jetstream.KeyValueEntry, um *livemap.UserMarker,
			) error {
				if um == nil {
					return nil
				}

				// Remove mapping
				if err := m.userMappingsStore.Delete(ctxCancel, userIdKey(um.UserId)); err != nil {
					m.logger.Error("failed to remove user unit mapping", zap.Error(err))
				}

				return nil
			}),
		)
		if err != nil {
			return fmt.Errorf("failed to create user location store. %w", err)
		}
		if err := userLocStore.Start(ctxCancel, false); err != nil {
			return fmt.Errorf("failed to start user location store. %w", err)
		}
		m.userLocStore = userLocStore

		// Periodic snapshot publisher
		m.snapshotTicker = time.NewTicker(defaultSnapEvery)
		go func() {
			if err := m.publishSnapshot(ctxCancel); err != nil {
				m.logger.Error("snapshot publish failed", zap.Error(err))
			}

			for {
				select {
				case <-ctxCancel.Done():
					return

				case <-m.snapshotTicker.C:
				}
			}
		}()

		go m.start(ctxCancel)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		m.refreshTicker.Stop()

		if m.snapshotTicker != nil {
			m.snapshotTicker.Stop()
		}

		return nil
	}))

	return m, nil
}

func (m *Manager) handleAppConfigUpdate(appCfg *appconfig.Cfg) {
	dbRefreshTime := appCfg.UserTracker.DbRefreshTime.AsDuration()
	m.refreshTicker.Reset(dbRefreshTime)
}

func (m *Manager) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-m.refreshTicker.C:
			m.refreshCache(ctx)
		}
	}
}

func (m *Manager) refreshCache(ctx context.Context) {
	ctx, span := m.tracer.Start(ctx, "tracker-refresh")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := m.refreshUserLocations(ctx); err != nil {
		m.logger.Error("failed to refresh user tracker cache", zap.Error(err))
	}
}

func (m *Manager) refreshUserLocations(ctx context.Context) error {
	m.logger.Debug("refreshing user tracker cache")

	tLocs := tLocs.AS("user_marker")
	tUsers := tables.User().AS("user")
	tFallbackJobProps := tJobProps.AS("fallback_job_props")
	tFallbackColleagueProps := tColleagueProps.AS("fallback_colleague_props")

	stmt := tLocs.
		SELECT(
			tLocs.Identifier,
			tLocs.Job,
			tLocs.JobGrade,
			tLocs.X,
			tLocs.Y,
			tLocs.UpdatedAt,
			tLocs.Hidden.AS("user_marker.hidden"),
			tUsers.ID.AS("user_marker.userid"),
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.PhoneNumber,
			jet.COALESCE(
				tColleagueProps.UserID,
				tFallbackColleagueProps.UserID,
			).AS("colleague_props.userid"),
			jet.COALESCE(
				tColleagueProps.Job,
				tFallbackColleagueProps.Job,
			).AS("colleague_props.job"),
			jet.COALESCE(
				tColleagueProps.NamePrefix,
				tFallbackColleagueProps.NamePrefix,
			).AS("colleague_props.name_prefix"),
			jet.COALESCE(
				tColleagueProps.NameSuffix,
				tFallbackColleagueProps.NameSuffix,
			).AS("colleague_props.name_suffix"),
			jet.COALESCE(
				tJobProps.LivemapMarkerColor,
				tFallbackJobProps.LivemapMarkerColor,
			).AS("user_marker.color"),
		).
		FROM(
			tLocs.
				INNER_JOIN(tUsers,
					tLocs.Identifier.EQ(tUsers.Identifier),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tLocs.Job),
				).
				LEFT_JOIN(tFallbackJobProps,
					tFallbackJobProps.Job.EQ(tUsers.Job),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tUsers.ID).
						AND(tColleagueProps.Job.EQ(tLocs.Job)),
				).
				LEFT_JOIN(tFallbackColleagueProps,
					tFallbackColleagueProps.UserID.EQ(tUsers.ID).
						AND(tFallbackColleagueProps.Job.EQ(tUsers.Job)),
				),
		).
		WHERE(jet.AND(
			tLocs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(4, jet.HOUR))),
		))

	var dest []*livemap.UserMarker
	if err := stmt.QueryContext(ctx, m.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	foundUserIds := map[int32]any{}
	added := 0

	errs := multierr.Combine()
	for i := range dest {
		if dest[i].User == nil {
			continue
		}

		foundUserIds[dest[i].UserId] = nil

		job := dest[i].User.Job
		if dest[i].Job != "" {
			dest[i].User.Job = dest[i].Job
			job = dest[i].Job // Use the job from the marker, not the user if set
		} else {
			dest[i].Job = job
		}

		jg := dest[i].User.JobGrade
		if dest[i].JobGrade != nil {
			jg = *dest[i].JobGrade
			dest[i].User.JobGrade = *dest[i].JobGrade
			dest[i].JobGrade = nil // Clear to avoid duplication, it is just used for overriding the user job grade
		}

		if dest[i].Color == nil {
			defaultColor := jobs.DefaultLivemapMarkerColor
			dest[i].Color = &defaultColor
		}

		postal, ok := m.postals.Closest(dest[i].X, dest[i].Y)
		if postal != nil && ok {
			dest[i].Postal = postal.Code
		}

		unitMapping, err := m.userMappingsStore.Get(userIdKey(dest[i].UserId))
		if err == nil && unitMapping.UnitId != nil && *unitMapping.UnitId > 0 {
			dest[i].UnitId = unitMapping.UnitId
			if unit, err := m.state.GetUnit(ctx, job, *unitMapping.UnitId); err == nil {
				dest[i].Unit = unit
			}
		} else {
			dest[i].UnitId = nil
			dest[i].Unit = nil
		}

		m.enricher.EnrichJobName(dest[i])
		m.enricher.EnrichJobInfo(dest[i].User)

		dest[i].User.JobGrade = jg

		um, err := m.userLocStore.Get(userMarkerKey(dest[i].UserId, job, jg))
		// No user marker in key value store nor locally
		if um == nil || err != nil {
			added++
			if err := m.userLocStore.Put(ctx, userMarkerKey(dest[i].UserId, job, jg), dest[i]); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		} else {
			// If not equal, update marker in store
			if proto.Equal(um, dest[i]) {
				continue
			}

			um.Merge(dest[i])

			oldKey := userMarkerKey(dest[i].UserId, job, jg) // job/jg are the *previous* ones
			ujg := jg
			if um.JobGrade != nil {
				ujg = *um.JobGrade // Use the job grade from the existing marker
			}
			newKey := userMarkerKey(dest[i].UserId, um.Job, ujg)
			if oldKey != newKey {
				if err := m.userLocStore.Delete(ctx, oldKey); err != nil {
					errs = multierr.Append(errs, err)
				}
			}

			if err := m.userLocStore.Put(ctx, userMarkerKey(dest[i].UserId, job, jg), um); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		}
	}

	removed, err := m.cleanupUserIDs(ctx, foundUserIds)
	if err != nil {
		return err
	}

	m.logger.Debug("completed user tracker cache refresh", zap.Int("added", added), zap.Int("removed", removed))

	return nil
}

func (m *Manager) cleanupUserIDs(ctx context.Context, foundUserIDs map[int32]any) (int, error) {
	m.logger.Debug("cleaning up user IDs", zap.Int32s("found_user_ids", utils.GetMapKeys(foundUserIDs)))

	var errs error
	now := time.Now()
	keys := m.userLocStore.Keys("")
	removed := []string{}
	for _, key := range keys {
		idKey, err := ExtractUserID(key)
		if err != nil {
			m.logger.Warn("failed to extract user ID from key", zap.String("key", key), zap.Error(err))
			continue
		}

		if _, ok := foundUserIDs[int32(idKey)]; ok {
			continue
		}

		marker, err := m.userLocStore.Get(key)
		if err != nil {
			continue
		}

		// Marker has been updated in the latest 30 seconds, skip it
		if marker.UpdatedAt != nil && now.Sub(marker.UpdatedAt.AsTime()) <= 30*time.Second {
			continue
		}

		if err := m.userLocStore.Delete(ctx, key); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		removed = append(removed, key)
	}

	m.logger.Debug("removed user ids from tracker cache", zap.Strings("user_ids", removed))

	return len(removed), errs
}

// Snapshot logic - one compressed roll-up every defaultSnapEvery
func (m *Manager) publishSnapshot(ctx context.Context) error {
	// build Snapshot proto
	snap := &pblivemap.Snapshot{}
	m.userLocStore.Range(func(_ string, um *livemap.UserMarker) bool {
		if um.Hidden {
			return true // Skip hidden markers
		}

		snap.Markers = append(snap.Markers, um)
		return true
	})

	raw, err := proto.Marshal(snap)
	if err != nil {
		return fmt.Errorf("marshal snapshot. %w", err)
	}

	// Compress (zstd keeps CPU low and ratio high)
	var dst bytes.Buffer
	enc, err := zstd.NewWriter(&dst)
	if err != nil {
		return fmt.Errorf("create zstd writer. %w", err)
	}
	if _, err := enc.Write(raw); err != nil {
		return fmt.Errorf("write snapshot to zstd writer. %w", err)
	}
	enc.Close()

	msg := &nats.Msg{
		Subject: SnapshotSubject,
		Data:    dst.Bytes(),
		Header: nats.Header{
			"Nats-Rollup":  []string{"all"}, // Atomic replace
			"KV-Operation": []string{"ROLLUP"},
		},
	}
	_, err = m.js.PublishMsg(ctx, msg)
	return err
}

func userMarkerKey(id int32, job string, grade int32) string {
	return fmt.Sprintf("%s.%d.%d", job, grade, id)
}

func decodeUserMarkerKey(key string) (int32, string, int32, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid user marker key: %s", key)
	}

	id, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid user marker id: %s", parts[2])
	}

	job := parts[1]
	grade, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid user marker grade: %s", parts[0])
	}

	return int32(id), job, int32(grade), nil
}
