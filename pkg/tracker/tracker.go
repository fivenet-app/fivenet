package tracker

import (
	"context"
	"database/sql"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/gin-gonic/gin"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	MaxDispatchMarkerLimit = 120
)

var (
	tLocs     = table.FivenetUserLocations
	tUsers    = table.Users.AS("user")
	tJobProps = table.FivenetJobProps
)

// TODO keep log of on duty (non hidden) player locations

// TODO keep track of all markers in database (in the future also "Sperrzonen" and Panicbuttons)

type Tracker struct {
	ctx    context.Context
	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB
	c      *mstlystcdata.Enricher

	usersCache syncx.Map[string, []*livemap.UserMarker]

	broker *utils.Broker[interface{}]

	refreshTime time.Duration
	visibleJobs []string
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	TP       *tracesdk.TracerProvider
	DB       *sql.DB
	Enricher *mstlystcdata.Enricher
	Config   *config.Config
}

func New(p Params) *Tracker {
	ctx, cancel := context.WithCancel(context.Background())

	broker := utils.NewBroker[interface{}](ctx)

	t := &Tracker{
		ctx:    ctx,
		logger: p.Logger,
		tracer: p.TP.Tracer("tracker-cache"),
		db:     p.DB,
		c:      p.Enricher,

		usersCache: syncx.Map[string, []*livemap.UserMarker]{},

		broker: broker,

		refreshTime: p.Config.Cache.RefreshTime,
		visibleJobs: p.Config.Game.Livemap.Jobs,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) {
		go broker.Start()

		// Only run the tracker random user marker generator in debug mode
		if p.Config.Mode == gin.DebugMode {
			go t.GenerateRandomUserMarker()
			go t.GenerateRandomDispatchMarker()
		}

		go t.Start()
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) {
		cancel()
	}))

	return t
}

func (s *Tracker) Start() {
	for {
		s.refreshCache()

		select {
		case <-s.ctx.Done():
			s.broker.Stop()
			return
		case <-time.After(s.refreshTime):
		}
	}
}

func (s *Tracker) refreshCache() {
	ctx, span := s.tracer.Start(s.ctx, "livemap-refresh-cache")
	defer span.End()

	if err := s.refreshUserLocations(ctx); err != nil {
		s.logger.Error("failed to refresh livemap users cache", zap.Error(err))
	}

	s.broker.Publish(nil)
}

func (s *Tracker) refreshUserLocations(ctx context.Context) error {
	tLocs := tLocs.AS("genericmarker")
	stmt := tLocs.
		SELECT(
			tLocs.Identifier,
			tLocs.Job,
			tLocs.X,
			tLocs.Y,
			tLocs.UpdatedAt,
			tUsers.ID.AS("user.id"),
			tUsers.ID.AS("genericmarker.id"),
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tJobProps.LivemapMarkerColor.AS("genericmarker.icon_color"),
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
			tLocs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(60, jet.MINUTE))),
		))

	var dest []*livemap.UserMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	markers := map[string][]*livemap.UserMarker{}
	for i := 0; i < len(dest); i++ {
		s.c.EnrichJobInfo(dest[i].User)

		job := dest[i].User.Job
		if _, ok := markers[job]; !ok {
			markers[job] = []*livemap.UserMarker{}
		}
		if dest[i].Marker.IconColor == "" {
			dest[i].Marker.IconColor = jobs.DefaultLivemapMarkerColor
		}

		markers[job] = append(markers[job], dest[i])
	}

	// TODO compare the markers with the current markers before updating them and use these as events for players going on/off duty

	for job, v := range markers {
		s.usersCache.Store(job, v)
	}

	return nil
}

func (s *Tracker) GetPlayers(job string) ([]*livemap.UserMarker, bool) {
	return s.usersCache.Load(job)
}

func (s *Tracker) GetPlayerFromJob(job string, userId int32) (float64, float64, bool) {
	users, ok := s.usersCache.Load(job)
	if !ok {
		return 0, 0, ok
	}

	// TODO
	_ = users

	return 0, 0, true
}
