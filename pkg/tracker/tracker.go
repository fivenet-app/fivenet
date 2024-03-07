package tracker

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/pkg/nats/store"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	tLocs     = table.FivenetUserLocations
	tUsers    = table.Users.AS("user")
	tJobProps = table.FivenetJobProps
)

type ITracker interface {
	GetUsersByJob(job string) (*xsync.MapOf[int32, *livemap.UserMarker], bool)
	GetUserById(id int32) (*livemap.UserMarker, bool)
	GetAllActiveUsers() ([]*livemap.UserMarker, error)
	IsUserOnDuty(userId int32) bool

	Subscribe() chan *livemap.UsersUpdateEvent
	Unsubscribe(c chan *livemap.UsersUpdateEvent)
}

type Tracker struct {
	ITracker

	logger *zap.Logger
	js     jetstream.JetStream

	jsCons jetstream.ConsumeContext

	userStore  *store.Store[livemap.UserMarker, *livemap.UserMarker]
	usersByJob *xsync.MapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]]

	broker *utils.Broker[*livemap.UsersUpdateEvent]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     jetstream.JetStream
}

func New(p Params) (ITracker, error) {
	ctx, cancel := context.WithCancel(context.Background())

	t := &Tracker{
		logger: p.Logger,
		js:     p.JS,

		usersByJob: xsync.NewMapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]](),

		broker: utils.NewBroker[*livemap.UsersUpdateEvent](),
	}

	p.LC.Append(fx.StartHook(func(c context.Context) error {
		if err := registerStreams(c, p.JS); err != nil {
			return err
		}

		go t.broker.Start(ctx)

		userIDs, err := store.NewWithLocks[livemap.UserMarker, *livemap.UserMarker](ctx, p.Logger, p.JS, "tracker", nil,
			func(s *store.Store[livemap.UserMarker, *livemap.UserMarker]) error {
				s.OnUpdate = func(um *livemap.UserMarker) (*livemap.UserMarker, error) {
					if um == nil || um.Info == nil {
						return um, nil
					}

					jobUsers, _ := t.usersByJob.LoadOrCompute(um.Info.Job, func() *xsync.MapOf[int32, *livemap.UserMarker] {
						return xsync.NewMapOf[int32, *livemap.UserMarker]()
					})
					if m, loaded := jobUsers.LoadOrStore(um.UserId, um); loaded {
						// Merge value if loaded from local data store
						m.Merge(um)
					}

					return um, nil
				}
				return nil
			},
			func(s *store.Store[livemap.UserMarker, *livemap.UserMarker]) error {
				s.OnDelete = func(kve jetstream.KeyValueEntry, um *livemap.UserMarker) error {
					if um == nil || um.Info == nil {
						return nil
					}

					if jobUsers, ok := t.usersByJob.Load(um.Info.Job); ok {
						jobUsers.Delete(um.UserId)
					}

					return nil
				}
				return nil
			})
		if err != nil {
			return err
		}

		if err := userIDs.Start(c); err != nil {
			return err
		}
		t.userStore = userIDs

		if err := t.registerSubscriptions(c); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		t.jsCons.Stop()

		return nil
	}))

	return t, nil
}

func (s *Tracker) watchForChanges(msg jetstream.Msg) {
	dest := &livemap.UsersUpdateEvent{}
	if err := proto.Unmarshal(msg.Data(), dest); err != nil {
		s.logger.Error("failed to unmarshal nats user update response", zap.Error(err))
		return
	}

	if s.broker.SubCount() <= 0 {
		return
	}

	s.broker.Publish(dest)
}

func (s *Tracker) GetUsersByJob(job string) (*xsync.MapOf[int32, *livemap.UserMarker], bool) {
	return s.usersByJob.Load(job)
}

func (s *Tracker) GetAllActiveUsers() ([]*livemap.UserMarker, error) {
	return s.userStore.List()
}

func (s *Tracker) IsUserOnDuty(userId int32) bool {
	_, ok := s.userStore.Get(userIdKey(userId))
	return ok
}

func (s *Tracker) GetUserById(id int32) (*livemap.UserMarker, bool) {
	return s.userStore.Get(userIdKey(id))
}

func (s *Tracker) Subscribe() chan *livemap.UsersUpdateEvent {
	return s.broker.Subscribe()
}

func (s *Tracker) Unsubscribe(c chan *livemap.UsersUpdateEvent) {
	s.broker.Unsubscribe(c)
}
