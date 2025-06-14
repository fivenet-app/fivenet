package tracker

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
)

type TestTracker struct {
	ITracker

	broker *broker.Broker[*store.KeyValueEntry[livemap.UserMarker, *livemap.UserMarker]]

	jobs       []string
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

		broker: broker.New[*store.KeyValueEntry[livemap.UserMarker, *livemap.UserMarker]](),
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

func (s *TestTracker) ListTrackedJobs() []string {
	return s.jobs
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

func (s *TestTracker) GetUserMarkerById(id int32) (*livemap.UserMarker, bool) {
	info, ok := s.usersIDs.Load(id)
	if !ok {
		return nil, false
	}

	return s.GetUserByJobAndID(info.Job, id)
}

func (s *TestTracker) Subscribe(ctx context.Context) (chan *store.KeyValueEntry[livemap.UserMarker, *livemap.UserMarker], error) {
	return s.broker.Subscribe(), nil
}
