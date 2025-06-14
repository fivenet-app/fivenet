package tracker

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
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

	userStore *store.Store[livemap.UserMarker, *livemap.UserMarker]
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
		userStore, err := store.New[livemap.UserMarker, *livemap.UserMarker](ctxStartup, p.Logger, p.JS, "tracker",
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
		)
		if err != nil {
			return err
		}

		if err := userStore.Start(ctxCancel, false); err != nil {
			return err
		}
		m.userStore = userStore

		if err := registerStreams(ctxStartup, m.js); err != nil {
			return err
		}

		go m.start(ctxCancel)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return m, nil
}

func (m *Manager) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(m.appCfg.Get().UserTracker.DbRefreshTime.AsDuration()):
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

func (m *Manager) cleanupUserIDs(ctx context.Context, foundUserIDs map[int32]any) error {
	event := &livemap.UsersUpdateEvent{}

	var errs error

	now := time.Now()
	m.logger.Debug("cleaning up user IDs", zap.Int32s("found_user_ids", utils.GetMapKeys(foundUserIDs)))
	keys := m.userStore.Keys(ctx, "")
	removed := []string{}
	for _, key := range keys {
		idKey, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			continue
		}

		if _, ok := foundUserIDs[int32(idKey)]; ok {
			continue
		}

		marker, ok := m.userStore.Get(key)
		if !ok {
			continue
		}

		// Marker has been updated in the latest 30 seconds, skip it
		if marker.UpdatedAt != nil && now.Sub(marker.UpdatedAt.AsTime()) <= 30*time.Second {
			continue
		}

		if err := m.userStore.Delete(ctx, key); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		event.Removed = append(event.Removed, marker)
		removed = append(removed, key)
	}

	m.logger.Debug("removed user ids from tracker cache", zap.Strings("user_ids", removed))

	if len(event.Removed) > 0 {
		if err := m.sendUpdateEvent(ctx, UsersUpdate, event); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (m *Manager) refreshUserLocations(ctx context.Context) error {
	m.logger.Debug("refreshing user tracker cache")

	tLocs := tLocs.AS("user_marker")
	tUsers := tables.User().AS("user")

	stmt := tLocs.
		SELECT(
			tLocs.Identifier,
			tLocs.Job,
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
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tJobProps.LivemapMarkerColor.AS("user_marker.color"),
		).
		FROM(
			tLocs.
				INNER_JOIN(tUsers,
					tLocs.Identifier.EQ(tUsers.Identifier),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tLocs.Job).OR(tJobProps.Job.EQ(tUsers.Job)),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tUsers.ID).
						AND(tColleagueProps.Job.EQ(tLocs.Job).OR(tColleagueProps.Job.EQ(tUsers.Job))),
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

	errs := multierr.Combine()

	event := &livemap.UsersUpdateEvent{}
	for i := range dest {
		if dest[i].User == nil {
			continue
		}

		foundUserIds[dest[i].UserId] = nil

		m.enricher.EnrichJobName(dest[i])
		m.enricher.EnrichJobInfo(dest[i].User)

		if dest[i].Color == nil {
			defaultColor := jobs.DefaultLivemapMarkerColor
			dest[i].Color = &defaultColor
		}

		postal, ok := m.postals.Closest(dest[i].X, dest[i].Y)
		if postal != nil && ok {
			dest[i].Postal = postal.Code
		}

		unitId, ok := m.state.GetUserUnitID(ctx, dest[i].UserId)
		if ok {
			dest[i].UnitId = &unitId
			job := dest[i].User.Job
			if unit, err := m.state.GetUnit(ctx, job, unitId); err == nil {
				dest[i].Unit = unit
			}
		}

		userMarker, ok := m.userStore.Get(userIdKey(dest[i].UserId))
		// No user marker in key value store nor locally
		if userMarker == nil || !ok {
			// User wasn't in the list, so they must be new so add the user to event for keeping track of users
			event.Added = append(event.Added, dest[i])

			if err := m.userStore.Put(ctx, userIdKey(dest[i].UserId), dest[i]); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		} else {
			// If not equal, update marker in store
			if !proto.Equal(userMarker, dest[i]) {
				userMarker.Merge(dest[i])

				if err := m.userStore.Put(ctx, userIdKey(dest[i].UserId), userMarker); err != nil {
					errs = multierr.Append(errs, err)
					continue
				}
			}
		}
	}

	if len(event.Added) > 0 {
		if err := m.sendUpdateEvent(ctx, UsersUpdate, event); err != nil {
			return err
		}
	}

	if err := m.cleanupUserIDs(ctx, foundUserIds); err != nil {
		return err
	}

	m.logger.Debug("completed user tracker cache refresh", zap.Int("added", len(event.Added)), zap.Int("removed", len(event.Removed)))

	return nil
}

func userIdKey(id int32) string {
	return strconv.FormatInt(int64(id), 10)
}
