package tracker

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
)

type TestTracker struct {
	ITracker

	broker *broker.Broker[*livemap.UsersUpdateEvent]

	usersCache *xsync.Map[string, *xsync.Map[int32, *livemap.UserMarker]]
	usersIDs   *xsync.Map[int32, *livemap.UserMarker]
}

type TestParams struct {
	fx.In

	LC fx.Lifecycle
}

func NewForTests(p TestParams) ITracker {
	t := &TestTracker{
		usersCache: xsync.NewMap[string, *xsync.Map[int32, *livemap.UserMarker]](),
		usersIDs:   xsync.NewMap[int32, *livemap.UserMarker](),

		broker: broker.New[*livemap.UsersUpdateEvent](),
	}

	brokerCtx, brokerCancel := context.WithCancel(context.Background())
	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		go t.broker.Start(brokerCtx)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		brokerCancel()

		return nil
	}))

	return t
}

func (s *TestTracker) GetUsersByJob(job string) (*xsync.Map[int32, *livemap.UserMarker], bool) {
	return s.usersCache.Load(job)
}

func (s *TestTracker) GetUserByJobAndID(job string, userId int32) (*livemap.UserMarker, bool) {
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

func (s *TestTracker) IsUserOnDuty(userId int32) bool {
	if _, ok := s.usersIDs.Load(userId); !ok {
		return false
	}

	return true
}

func (s *TestTracker) GetUserById(id int32) (*livemap.UserMarker, bool) {
	info, ok := s.usersIDs.Load(id)
	if !ok {
		return nil, false
	}

	return s.GetUserByJobAndID(info.Job, id)
}

func (s *TestTracker) Subscribe() chan *livemap.UsersUpdateEvent {
	return s.broker.Subscribe()
}

func (s *TestTracker) Unsubscribe(c chan *livemap.UsersUpdateEvent) {
	s.broker.Unsubscribe(c)
}
