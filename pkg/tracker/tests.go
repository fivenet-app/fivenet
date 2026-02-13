package tracker

import (
	"context"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
)

type TestTracker struct {
	ITracker

	broker *broker.Broker[*store.KeyValueEntry[livemapmarkers.UserMarker, *livemapmarkers.UserMarker]]

	jobs       []string
	usersCache *xsync.Map[string, *xsync.Map[int32, *livemapmarkers.UserMarker]]
	usersIDs   *xsync.Map[int32, *livemapmarkers.UserMarker]
}

type TestParams struct {
	fx.In

	LC fx.Lifecycle
}

func NewForTests(p TestParams) ITracker {
	t := &TestTracker{
		usersCache: xsync.NewMap[string, *xsync.Map[int32, *livemapmarkers.UserMarker]](),
		usersIDs:   xsync.NewMap[int32, *livemapmarkers.UserMarker](),

		broker: broker.New[*store.KeyValueEntry[livemapmarkers.UserMarker, *livemapmarkers.UserMarker]](),
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

func (s *TestTracker) GetUserByJobAndID(
	job string,
	userId int32,
) (*livemapmarkers.UserMarker, bool) {
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

func (s *TestTracker) GetUserMarkerById(id int32) (*livemapmarkers.UserMarker, bool) {
	info, ok := s.usersIDs.Load(id)
	if !ok {
		return nil, false
	}

	return s.GetUserByJobAndID(info.GetJob(), id)
}

func (s *TestTracker) Subscribe(
	_ context.Context,
) (store.IKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker], error) {
	return &TestKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker]{
		broker: s.broker,
	}, nil
}

type TestKVWatcher[T any, U protoutils.ProtoMessageWithMerge[T]] struct {
	broker *broker.Broker[*store.KeyValueEntry[livemapmarkers.UserMarker, *livemapmarkers.UserMarker]]
}

func (w *TestKVWatcher[T, U]) Stop() error {
	return nil
}

func (w *TestKVWatcher[T, U]) Updates() <-chan *store.KeyValueEntry[livemapmarkers.UserMarker, *livemapmarkers.UserMarker] {
	return w.broker.Subscribe()
}

func (w *TestKVWatcher[T, U]) Unsubscribe() error {
	w.broker.Unsubscribe(w.broker.Subscribe())
	return nil
}
