package tracker

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/state"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/nats/store"
	"github.com/gin-gonic/gin"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Manager struct {
	logger   *zap.Logger
	tracer   trace.Tracer
	js       *events.JSWrapper
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	postals  postals.Postals
	state    *state.State
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
	State     *state.State
	AppConfig appconfig.IConfig
}

func NewManager(p ManagerParams) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		logger:   p.Logger,
		tracer:   p.TP.Tracer("tracker-manager"),
		js:       p.JS,
		db:       p.DB,
		enricher: p.Enricher,
		postals:  p.Postals,
		state:    p.State,
		appCfg:   p.AppConfig,
	}

	p.LC.Append(fx.StartHook(func(c context.Context) error {
		userStore, err := store.NewWithLocks[livemap.UserMarker, *livemap.UserMarker](c, p.Logger, p.JS, "tracker", nil)
		if err != nil {
			return err
		}

		if err := userStore.Start(ctx); err != nil {
			return err
		}
		m.userStore = userStore

		if err := registerStreams(c, m.js); err != nil {
			return err
		}

		go m.start(ctx)

		// Only run the tracker random user marker generator in debug mode
		if p.Config.Mode == gin.DebugMode {
			go m.randomizeUserMarkers(ctx)
		}

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

	if err := m.refreshUserLocations(ctx); err != nil {
		m.logger.Error("failed to refresh user tracker cache", zap.Error(err))
	}
}

func (m *Manager) cleanupUserIDs(ctx context.Context, found map[int32]interface{}) error {
	event := &livemap.UsersUpdateEvent{}

	keys, err := m.userStore.Keys(ctx, "")
	if err != nil && !errors.Is(err, jetstream.ErrNoKeysFound) {
		return err
	}

	for _, key := range keys {
		idKey, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			continue
		}

		if _, ok := found[int32(idKey)]; ok {
			continue
		}

		marker, ok := m.userStore.Get(key)
		if !ok {
			continue
		}

		// Marker has been updated in the latest 15 seconds, skip it
		if marker.Info.UpdatedAt != nil && time.Since(marker.Info.UpdatedAt.AsTime()) <= 15*time.Second {
			continue
		}

		if err := m.userStore.Delete(ctx, key); err != nil {
			return err
		}

		event.Removed = append(event.Removed, marker)
	}

	if len(event.Removed) > 0 {
		if err := m.sendUpdateEvent(ctx, UsersUpdate, event); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) refreshUserLocations(ctx context.Context) error {
	m.logger.Debug("refreshing user tracker cache")

	tLocs := tLocs.AS("markerInfo")
	stmt := tLocs.
		SELECT(
			tLocs.Identifier,
			tLocs.Job,
			tLocs.X,
			tLocs.Y,
			tLocs.UpdatedAt,
			tUsers.ID.AS("usermarker.userid"),
			tUsers.ID.AS("user.id"),
			tUsers.ID.AS("markerInfo.id"),
			tLocs.Job.AS("user.job"),
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.PhoneNumber,
			tJobProps.LivemapMarkerColor.AS("markerInfo.color"),
		).
		FROM(
			tLocs.
				INNER_JOIN(tUsers,
					tLocs.Identifier.EQ(tUsers.Identifier),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tUsers.Job),
				),
		).
		WHERE(jet.AND(
			tLocs.Hidden.IS_FALSE(),
			tLocs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(4, jet.HOUR))),
		))

	var dest []*livemap.UserMarker
	if err := stmt.QueryContext(ctx, m.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	foundUserIds := map[int32]interface{}{}

	errs := multierr.Combine()

	event := &livemap.UsersUpdateEvent{}
	for i := 0; i < len(dest); i++ {
		foundUserIds[dest[i].UserId] = nil

		m.enricher.EnrichJobInfo(dest[i].User)

		if dest[i].Info.Color == nil {
			defaultColor := users.DefaultLivemapMarkerColor
			dest[i].Info.Color = &defaultColor
		}

		postal := m.postals.Closest(dest[i].Info.X, dest[i].Info.Y)
		if postal != nil {
			dest[i].Info.Postal = postal.Code
		}

		userId := dest[i].User.UserId

		unitId, ok := m.state.GetUserUnitID(ctx, userId)
		if ok {
			dest[i].UnitId = &unitId
			job := dest[i].User.Job
			if unit, err := m.state.GetUnit(ctx, job, unitId); err == nil {
				dest[i].Unit = unit
			}
		}

		userMarker, ok := m.userStore.Get(userIdKey(userId))
		// No user marker in key value store nor locally
		if userMarker == nil || !ok {
			// User wasn't in the list, so they must be new so add the user to event for keeping track of users
			event.Added = append(event.Added, dest[i])

			if err := m.userStore.Put(ctx, userIdKey(userId), dest[i]); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		} else {
			// If not equal, update marker in store
			if !proto.Equal(userMarker, dest[i]) {
				userMarker.Merge(dest[i])

				if err := m.userStore.Put(ctx, userIdKey(userId), userMarker); err != nil {
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

	return nil
}

func userIdKey(id int32) string {
	return strconv.Itoa(int(id))
}
