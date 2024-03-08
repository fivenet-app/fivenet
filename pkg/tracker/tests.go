package tracker

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/pkg/utils/broker"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
)

type TestTracker struct {
	ITracker

	broker *broker.Broker[*livemap.UsersUpdateEvent]

	usersCache *xsync.MapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]]
	usersIDs   *xsync.MapOf[int32, *livemap.UserMarker]
}

type TestParams struct {
	fx.In

	LC fx.Lifecycle
}

func NewForTests(p TestParams) ITracker {
	broker := broker.New[*livemap.UsersUpdateEvent]()

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		go broker.Start(context.Background())

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		broker.Stop()

		return nil
	}))

	return &TestTracker{
		usersCache: xsync.NewMapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]](),
		usersIDs:   xsync.NewMapOf[int32, *livemap.UserMarker](),

		broker: broker,
	}
}

func (s *TestTracker) GetUsersByJob(job string) (*xsync.MapOf[int32, *livemap.UserMarker], bool) {
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

func (s *TestTracker) GetAllActiveUsers() ([]*livemap.UserMarker, error) {
	list := []*livemap.UserMarker{}

	s.usersIDs.Range(func(key int32, value *livemap.UserMarker) bool {
		list = append(list, value)
		return true
	})

	return list, nil
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

	return s.GetUserByJobAndID(info.Info.Job, id)
}

func (s *TestTracker) Subscribe() chan *livemap.UsersUpdateEvent {
	return s.broker.Subscribe()
}

func (s *TestTracker) Unsubscribe(c chan *livemap.UsersUpdateEvent) {
	s.broker.Unsubscribe(c)
}
