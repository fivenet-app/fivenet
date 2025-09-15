package manager

import (
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
	pbtracker "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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
	appCfg   appconfig.IConfig

	units       *units.UnitDB
	helpers     *helpers.Helpers
	dispatchers *dispatchers.DispatchersDB

	refreshTicker *time.Ticker

	userByIDStore     *store.Store[livemap.UserMarker, *livemap.UserMarker]
	userLocStore      *store.Store[livemap.UserMarker, *livemap.UserMarker]
	userMappingsStore *store.Store[pbtracker.UserMapping, *pbtracker.UserMapping]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	JS        *events.JSWrapper
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	Enricher  *mstlystcdata.Enricher
	Postals   postals.Postals
	Cfg       *config.Config
	AppConfig appconfig.IConfig

	Units       *units.UnitDB
	Helpers     *helpers.Helpers
	Dispatchers *dispatchers.DispatchersDB
}

func New(p Params) (*Manager, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.Named("tracker.manager")

	m := &Manager{
		logger:   logger,
		tracer:   p.TP.Tracer("tracker.manager"),
		js:       p.JS,
		db:       p.DB,
		enricher: p.Enricher,
		postals:  p.Postals,
		appCfg:   p.AppConfig,

		units:       p.Units,
		helpers:     p.Helpers,
		dispatchers: p.Dispatchers,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		appCfg := p.AppConfig.Get()
		m.refreshTicker = time.NewTicker(appCfg.UserTracker.GetDbRefreshTime().AsDuration())

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

		storeLogger := logger.WithOptions(
			zap.IncreaseLevel(
				p.Cfg.Log.LevelOverrides.Get(config.LoggingComponentKVStore, p.Cfg.LogLevel),
			),
		)

		userMappingsStore, err := store.New[pbtracker.UserMapping, *pbtracker.UserMapping](
			ctxStartup,
			storeLogger,
			p.JS,
			tracker.BucketUserMappingsMap,
			store.WithLocks[pbtracker.UserMapping, *pbtracker.UserMapping](nil),
			store.WithOnDeleteFn(
				func(ctx context.Context, key string, um *pbtracker.UserMapping) error {
					if um == nil {
						return nil
					}

					// Remove user from unit if it has a unit_id
					if um.UnitId != nil && um.GetUnitId() > 0 {
						if err := m.units.UpdateUnitAssignments(ctx, "", &um.UserId, um.GetUnitId(), nil, []int32{um.GetUserId()}); err != nil {
							m.logger.Error("failed to remove user from unit", zap.Error(err))
						}
					}

					return nil
				},
			),
		)
		if err != nil {
			return fmt.Errorf("failed to create user mappings store. %w", err)
		}
		if err := userMappingsStore.Start(ctxCancel, false); err != nil {
			return fmt.Errorf("failed to start user mappings store. %w", err)
		}
		m.userMappingsStore = userMappingsStore

		userLocStore, err := store.New[livemap.UserMarker, *livemap.UserMarker](
			ctxStartup,
			storeLogger,
			p.JS,
			tracker.BucketUserLoc,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
		)
		if err != nil {
			return fmt.Errorf("failed to create user location store. %w", err)
		}
		if err := userLocStore.Start(ctxCancel, false); err != nil {
			return fmt.Errorf("failed to start user location store. %w", err)
		}
		m.userLocStore = userLocStore

		byID, err := store.New[livemap.UserMarker, *livemap.UserMarker](
			ctxStartup,
			storeLogger,
			p.JS,
			tracker.BucketUserLocByID,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
			store.WithOnUpdateFn(
				func(ctx context.Context, _ *livemap.UserMarker, newValue *livemap.UserMarker) (*livemap.UserMarker, error) {
					if newValue == nil {
						return nil, nil
					}

					if !m.userMappingsStore.Has(tracker.UserIdKey(newValue.GetUserId())) {
						// Upsert mapping (unit_id may be nil/0 = no unit)
						if err := m.userMappingsStore.Put(ctx, tracker.UserIdKey(newValue.GetUserId()), &pbtracker.UserMapping{
							UserId:    newValue.GetUserId(),
							UnitId:    newValue.UnitId,
							Hidden:    newValue.GetHidden(),
							CreatedAt: timestamp.Now(),
						}); err != nil {
							return nil, fmt.Errorf("failed to upsert user unit mapping. %w", err)
						}
					}

					if newValue.JobGrade != nil {
						if err := m.userLocStore.Put(ctx, userMarkerKey(newValue.GetUserId(), newValue.GetJob(), newValue.GetJobGrade()), newValue); err != nil {
							return nil, fmt.Errorf("failed to upsert user marker in store. %w", err)
						}
					}

					return newValue, nil
				},
			),
			store.WithOnDeleteFn(func(ctx context.Context,
				key string, um *livemap.UserMarker,
			) error {
				if um == nil {
					return nil
				}

				// Remove user marker if we have the info we need
				if um.JobGrade != nil {
					if err := m.userLocStore.Delete(ctx, userMarkerKey(um.GetUserId(), um.GetJob(), um.GetJobGrade())); err != nil {
						m.logger.Error(
							"failed to remove user marker from store",
							zap.Error(err),
							zap.Int32("user_id", um.GetUserId()),
							zap.String("job", um.GetJob()),
						)
					}
				}

				// Remove user mapping
				if err := m.userMappingsStore.Delete(ctx, key); err != nil {
					m.logger.Error(
						"failed to remove user unit mapping",
						zap.Error(err),
						zap.Int32("user_id", um.GetUserId()),
						zap.String("job", um.GetJob()),
					)
				}

				// Sign-off user from dispatchers
				if m.helpers.CheckIfUserIsDispatcher(ctx, um.GetJob(), um.GetUserId()) {
					if err := m.dispatchers.SetUserState(ctx, um.GetJob(), um.GetUserId(), false); err != nil {
						m.logger.Error(
							"failed to remove user from dispatchers",
							zap.Error(err),
							zap.Int32("user_id", um.GetUserId()),
							zap.String("job", um.GetJob()),
						)
					}
				}

				return nil
			}),
		)
		if err != nil {
			return fmt.Errorf("failed to create user location by ID store. %w", err)
		}
		if err := byID.Start(ctxCancel, false); err != nil {
			return fmt.Errorf("failed to start user location by ID store. %w", err)
		}
		m.userByIDStore = byID

		go m.start(ctxCancel)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		m.refreshTicker.Stop()

		return nil
	}))

	return m, nil
}

func (m *Manager) handleAppConfigUpdate(appCfg *appconfig.Cfg) {
	dbRefreshTime := appCfg.UserTracker.GetDbRefreshTime().AsDuration()
	m.refreshTicker.Reset(dbRefreshTime)
}

func (m *Manager) start(ctx context.Context) {
	m.refreshCache(ctx, true)

	for {
		select {
		case <-ctx.Done():
			return

		case <-m.refreshTicker.C:
			m.refreshCache(ctx, false)
		}
	}
}

func (m *Manager) refreshCache(ctx context.Context, initial bool) {
	ctx, span := m.tracer.Start(ctx, "tracker-refresh")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := m.refreshUserLocations(ctx, initial); err != nil {
		m.logger.Error("failed to refresh user tracker cache", zap.Error(err))
	}
}

func (m *Manager) refreshUserLocations(ctx context.Context, initial bool) error {
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
			mysql.COALESCE(
				tColleagueProps.UserID,
				tFallbackColleagueProps.UserID,
			).AS("colleague_props.userid"),
			mysql.COALESCE(
				tColleagueProps.Job,
				tFallbackColleagueProps.Job,
			).AS("colleague_props.job"),
			mysql.COALESCE(
				tColleagueProps.NamePrefix,
				tFallbackColleagueProps.NamePrefix,
			).AS("colleague_props.name_prefix"),
			mysql.COALESCE(
				tColleagueProps.NameSuffix,
				tFallbackColleagueProps.NameSuffix,
			).AS("colleague_props.name_suffix"),
			mysql.COALESCE(
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
		WHERE(mysql.AND(
			tLocs.UpdatedAt.GT_EQ(mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(4, mysql.HOUR))),
		))

	var dest []*livemap.UserMarker
	if err := stmt.QueryContext(ctx, m.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	foundUserIds := map[int32]any{}
	added := 0
	updated := 0

	errs := multierr.Combine()
	for i := range dest {
		if dest[i].GetUser() == nil {
			continue
		}

		foundUserIds[dest[i].GetUserId()] = nil

		// Use (override) job and job grade if set
		job := dest[i].GetUser().GetJob()
		if dest[i].GetJob() != "" {
			job = dest[i].GetJob() // Use the job from the marker, not the user if set
			dest[i].User.Job = dest[i].GetJob()
		} else {
			dest[i].Job = job
		}
		jg := dest[i].GetUser().GetJobGrade()
		if dest[i].JobGrade != nil {
			jg = dest[i].GetJobGrade()
			dest[i].User.JobGrade = dest[i].GetJobGrade()
		} else {
			dest[i].JobGrade = &jg // Ensure JobGrade is set, even if it is 0
		}

		if dest[i].Color == nil {
			defaultColor := jobs.DefaultLivemapMarkerColor
			dest[i].Color = &defaultColor
		}

		postal, ok := m.postals.Closest(dest[i].GetX(), dest[i].GetY())
		if postal != nil && ok {
			dest[i].Postal = postal.Code
		}

		unitMapping, err := m.userMappingsStore.Get(tracker.UserIdKey(dest[i].GetUserId()))
		if err == nil && unitMapping.UnitId != nil && unitMapping.GetUnitId() > 0 {
			dest[i].UnitId = unitMapping.UnitId
			if unit, err := m.units.Get(ctx, unitMapping.GetUnitId()); err == nil {
				dest[i].Unit = unit
			}
		} else {
			dest[i].UnitId = nil
			dest[i].Unit = nil
		}

		m.enricher.EnrichJobName(dest[i])
		m.enricher.EnrichJobInfo(dest[i].GetUser())

		um, err := m.userByIDStore.Get(tracker.UserIdKey(dest[i].GetUserId()))
		// No user marker in key value store nor locally
		if um == nil || err != nil {
			added++
			if err := m.userByIDStore.Put(ctx, tracker.UserIdKey(dest[i].GetUserId()), dest[i]); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		} else {
			// If not equal, update marker in store
			if !initial && proto.Equal(um, dest[i]) {
				continue
			}
			updated++

			uj := um.GetUser().GetJob()
			if um.GetJob() != "" {
				uj = um.GetJob()
			}
			ujg := um.GetUser().GetJobGrade()
			if um.JobGrade != nil {
				ujg = um.GetJobGrade() // Use the job grade from the existing marker
			}
			oldKey := userMarkerKey(dest[i].GetUserId(), uj, ujg) // uj/jg are the *previous* ones
			newKey := userMarkerKey(dest[i].GetUserId(), job, jg)
			if oldKey != newKey {
				if err := m.userLocStore.Delete(ctx, oldKey); err != nil {
					errs = multierr.Append(errs, err)
				}
			}

			um.Merge(dest[i])

			if err := m.userByIDStore.Put(ctx, tracker.UserIdKey(dest[i].GetUserId()), um); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		}
	}

	removed, err := m.cleanupUserIDs(ctx, foundUserIds)
	if err != nil {
		return err
	}

	m.logger.Debug(
		"completed user tracker cache refresh",
		zap.Int("added", added),
		zap.Int("updated", updated),
		zap.Int("removed", removed),
	)

	return nil
}

func (m *Manager) cleanupUserIDs(ctx context.Context, foundUserIds map[int32]any) (int, error) {
	var errs error
	keys := m.userLocStore.Keys("")
	removed := []string{}
	for _, key := range keys {
		userIdKey, err := extractUserID(key)
		if err != nil {
			m.logger.Warn(
				"failed to extract user ID from key",
				zap.String("key", key),
				zap.Error(err),
			)
			continue
		}

		// If the user ID is not in the foundUserIds map, we can remove it
		if _, ok := foundUserIds[userIdKey]; !ok {
			if err := m.userByIDStore.Delete(ctx, strconv.FormatInt(int64(userIdKey), 10)); err != nil {
				errs = multierr.Append(errs, err)
			}
			continue
		}

		// Lookup user by id
		marker, err := m.userByIDStore.Get(strconv.FormatInt(int64(userIdKey), 10))
		if err != nil {
			if !errors.Is(err, jetstream.ErrKeyNotFound) {
				continue
			}
		}
		// Short path if marker by id is not nil, we can remove the user location by key
		if marker != nil {
			jg := int32(0)
			if marker.JobGrade != nil {
				jg = marker.GetJobGrade()
			}
			oldKey := userMarkerKey(marker.GetUserId(), marker.GetJob(), jg)
			if key == oldKey {
				continue
			}

			if err := m.userLocStore.Delete(ctx, oldKey); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}

			continue
		}

		if err := m.userByIDStore.Delete(ctx, key); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		removed = append(removed, key)
	}

	m.logger.Debug("removed user ids from tracker cache", zap.Strings("user_ids", removed))

	return len(removed), errs
}

func userMarkerKey(id int32, job string, grade int32) string {
	return fmt.Sprintf("%s.%d.%d", job, grade, id)
}

// extractUserID takes a key like "police.3.123"  ➜  123.
func extractUserID(key string) (int32, error) {
	idx := strings.LastIndexByte(key, '.')
	if idx < 0 || idx+1 >= len(key) {
		return 0, fmt.Errorf("key %q does not contain a numeric suffix", key)
	}

	id, err := strconv.ParseInt(key[idx+1:], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("key %q does not contain a valid numeric suffix. %w", key, err)
	}
	return int32(id), nil
}
