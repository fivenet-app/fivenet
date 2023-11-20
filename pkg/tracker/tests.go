package tracker

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/puzpuzpuz/xsync/v3"
)

type TestTracker struct {
	broker *utils.Broker[*Event]

	usersCache *xsync.MapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]]
	usersIDs   *xsync.MapOf[int32, *UserInfo]
}

func NewForTests(ctx context.Context) *TestTracker {
	broker := utils.NewBroker[*Event](ctx)

	return &TestTracker{
		usersCache: xsync.NewMapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]](),
		usersIDs:   xsync.NewMapOf[int32, *UserInfo](),

		broker: broker,
	}
}

func (s *TestTracker) GetUsers(job string) (*xsync.MapOf[int32, *livemap.UserMarker], bool) {
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

func (s *TestTracker) IsUserOnDuty(job string, userId int32) bool {
	users, ok := s.usersCache.Load(job)
	if !ok {
		return false
	}

	if _, ok := users.Load(userId); !ok {
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

func (s *TestTracker) Subscribe() chan *Event {
	return s.broker.Subscribe()
}

func (s *TestTracker) Unsubscribe(c chan *Event) {
	s.broker.Unsubscribe(c)
}
