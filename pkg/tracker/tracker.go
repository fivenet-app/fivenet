package tracker

import (
	"context"
	"database/sql"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/tracker/postals"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/gin-gonic/gin"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/puzpuzpuz/xsync/v3"
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

type UserInfo struct {
	Job    string
	UserID int32
	Time   time.Time
}

type Event struct {
	Added   []*UserInfo
	Removed []*UserInfo
	Current []*UserInfo
}

type Tracker struct {
	ctx      context.Context
	logger   *zap.Logger
	tracer   trace.Tracer
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	postals  *postals.Postals
	state    *state.State

	usersCache *xsync.MapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]]
	usersIDs   *xsync.MapOf[int32, *UserInfo]

	broker *utils.Broker[*Event]

	refreshTime time.Duration
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	TP       *tracesdk.TracerProvider
	DB       *sql.DB
	Enricher *mstlystcdata.Enricher
	Postals  *postals.Postals
	Config   *config.Config
	State    *state.State
}

func New(p Params) *Tracker {
	ctx, cancel := context.WithCancel(context.Background())

	broker := utils.NewBroker[*Event](ctx)

	t := &Tracker{
		ctx:      ctx,
		logger:   p.Logger,
		tracer:   p.TP.Tracer("tracker-cache"),
		db:       p.DB,
		enricher: p.Enricher,
		postals:  p.Postals,
		state:    p.State,

		usersCache: xsync.NewMapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]](),
		usersIDs:   xsync.NewMapOf[int32, *UserInfo](),

		broker: broker,

		refreshTime: p.Config.Game.Livemap.RefreshTime,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go broker.Start()

		// Only run the tracker random user marker generator in debug mode
		if p.Config.Mode == gin.DebugMode {
			go t.GenerateRandomUserMarker()
			go t.GenerateRandomDispatchMarker()
		}

		go t.start()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		return nil
	}))

	return t
}

func (s *Tracker) start() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-ticker.C:
			s.refreshCache(true)

		case <-time.After(s.refreshTime):
			// Only refresh cache if broker has a subscriber
			if s.broker.SubCount() <= 0 {
				break
			}

			s.refreshCache(false)
		}
	}
}

func (s *Tracker) refreshCache(force bool) {
	ctx, span := s.tracer.Start(s.ctx, "livemap-refresh-cache")
	defer span.End()

	if err := s.refreshUserLocations(ctx, force); err != nil {
		s.logger.Error("failed to refresh livemap users cache", zap.Error(err))
	}

	s.cleanupUserIDs()
}

func (s *Tracker) cleanupUserIDs() error {
	event := &Event{}

	now := time.Now()
	s.usersIDs.Range(func(key int32, info *UserInfo) bool {
		if now.After(info.Time) {
			event.Removed = append(event.Removed, info)
			s.usersIDs.Delete(key)

			jobUsers, ok := s.usersCache.Load(info.Job)
			if ok {
				jobUsers.Delete(key)
			}
		}

		return true
	})

	s.broker.Publish(event)

	return nil
}

func (s *Tracker) refreshUserLocations(ctx context.Context, force bool) error {
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
			tUsers.Identifier,
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
			jet.OR(
				tLocs.UpdatedAt.IS_NULL(),
				tLocs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(4, jet.HOUR))),
			),
		))

	var dest []*livemap.UserMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	event := &Event{}
	expiration := time.Now().Add(7 * s.refreshTime)
	markers := map[string]*xsync.MapOf[int32, *livemap.UserMarker]{}
	for i := 0; i < len(dest); i++ {
		s.enricher.EnrichJobInfo(dest[i].User)

		job := dest[i].User.Job
		if _, ok := markers[job]; !ok {
			markers[job] = xsync.NewMapOf[int32, *livemap.UserMarker]()
		}

		if dest[i].Info.Color == nil {
			defaultColor := users.DefaultLivemapMarkerColor
			dest[i].Info.Color = &defaultColor
		}

		postal := s.postals.Closest(dest[i].Info.X, dest[i].Info.Y)
		if postal != nil {
			dest[i].Info.Postal = postal.Code
		}

		userId := dest[i].User.UserId

		unitId, ok := s.state.GetUserUnitID(userId)
		if ok {
			dest[i].UnitId = &unitId
			if unit, err := s.state.GetUnit(job, unitId); err == nil {
				dest[i].Unit = unit
			}
		}

		markers[job].Store(userId, dest[i])

		userInfo := &UserInfo{
			Job:    dest[i].User.Job,
			UserID: userId,
			Time:   expiration,
		}
		if ui, ok := s.usersIDs.LoadOrStore(userId, userInfo); ok {
			ui.Job = userInfo.Job
			ui.UserID = userInfo.UserID
			ui.Time = userInfo.Time
		} else {
			// User wasn't in the list, so they must be new so add the user to event for keeping track of users
			if _, ok := s.GetUserByJobAndID(job, userId); !ok {
				event.Added = append(event.Added, userInfo)
			}
		}

		if force {
			event.Current = append(event.Current, userInfo)
		}
	}

	for job, marker := range markers {
		s.usersCache.Store(job, marker)
	}

	if len(event.Added) > 0 || len(event.Removed) > 0 || len(event.Current) > 0 {
		s.broker.Publish(event)
	}

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

func (s *Tracker) IsUserOnDuty(job string, userId int32) bool {
	users, ok := s.usersCache.Load(job)
	if !ok {
		return false
	}

	if _, ok := users.Load(userId); !ok {
		return false
	}

	return true
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
