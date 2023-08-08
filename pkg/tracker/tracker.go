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
	"github.com/galexrt/fivenet/pkg/utils/maps"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/gin-gonic/gin"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/puzpuzpuz/xsync/v2"
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

type Event struct {
	Added   []int32
	Removed []int32
}

type userInfo struct {
	Time time.Time
	Job  string
}

type Tracker struct {
	ctx    context.Context
	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB
	c      *mstlystcdata.Enricher

	usersCache *xsync.MapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]]
	usersIDs   *xsync.MapOf[int32, userInfo]

	broker *utils.Broker[*Event]

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

	broker := utils.NewBroker[*Event](ctx)

	t := &Tracker{
		ctx:    ctx,
		logger: p.Logger,
		tracer: p.TP.Tracer("tracker-cache"),
		db:     p.DB,
		c:      p.Enricher,

		usersCache: xsync.NewTypedMapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]](maps.HashString),
		usersIDs:   xsync.NewTypedMapOf[int32, userInfo](maps.HashInt32),

		broker: broker,

		refreshTime: p.Config.Game.Livemap.RefreshTime,
		visibleJobs: p.Config.Game.Livemap.Jobs,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) {
		go broker.Start()

		// Only run the tracker random user marker generator in debug mode
		if p.Config.Mode == gin.DebugMode {
			go t.GenerateRandomUserMarker()
			go t.GenerateRandomDispatchMarker()
		}

		go t.start()
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) {
		cancel()
	}))

	return t
}

func (s *Tracker) start() {
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

	s.cleanupUserIDs()
}

func (s *Tracker) cleanupUserIDs() error {
	event := &Event{}

	now := time.Now()
	s.usersIDs.Range(func(key int32, info userInfo) bool {
		if now.After(info.Time.Add(3 * s.refreshTime)) {
			event.Removed = append(event.Removed, key)
			s.usersIDs.Delete(key)

		}

		return true
	})

	s.broker.Publish(event)

	return nil
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

	event := &Event{}
	now := time.Now()
	markers := map[string]*xsync.MapOf[int32, *livemap.UserMarker]{}
	for i := 0; i < len(dest); i++ {
		s.c.EnrichJobInfo(dest[i].User)

		job := dest[i].User.Job
		if _, ok := markers[job]; !ok {
			markers[job] = xsync.NewTypedMapOf[int32, *livemap.UserMarker](maps.HashInt32)
		}
		if dest[i].Marker.IconColor == "" {
			dest[i].Marker.IconColor = jobs.DefaultLivemapMarkerColor
		}

		userId := dest[i].User.UserId
		markers[job].Store(userId, dest[i])

		if _, ok := s.usersIDs.Load(userId); !ok {
			// User wasn't in the list, so they must be new so add to event for "announcement"
			if _, ok := s.GetUserByJobAndID(job, userId); !ok {
				event.Added = append(event.Added, userId)
			}
		}

		s.usersIDs.Store(userId, userInfo{
			Time: now,
			Job:  dest[i].User.Job,
		})
	}

	for job, v := range markers {
		s.usersCache.Store(job, v)
	}

	s.broker.Publish(event)

	return nil
}

func (s *Tracker) GetUsers(job string) (*xsync.MapOf[int32, *livemap.UserMarker], bool) {
	return s.usersCache.Load(job)
}

func (s *Tracker) GetUserByJobAndID(job string, userId int32) (*livemap.UserMarker, bool) {
	users, ok := s.usersCache.Load(job)
	if !ok {
		return nil, false
	}

	user, ok := users.Load(userId)
	if !ok {
		return nil, false
	}

	return user, true
}

func (s *Tracker) GetUserById(id int32) (*livemap.UserMarker, bool) {
	info, ok := s.usersIDs.Load(id)
	if !ok {
		return nil, false
	}

	return s.GetUserByJobAndID(info.Job, id)
}

func (s *Tracker) Subscribe() chan *Event {
	return s.broker.Subscribe()
}

func (s *Tracker) Unsubscribe(c chan *Event) {
	s.broker.Unsubscribe(c)
}
