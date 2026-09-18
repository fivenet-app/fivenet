package tracker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbtracker "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/tracker"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
)

type TestTracker struct {
	broker *broker.Broker[*store.KeyValueEntry[livemapmarkers.UserMarker, *livemapmarkers.UserMarker]]

	jobs         []string
	usersCache   *xsync.Map[string, *xsync.Map[int32, *livemapmarkers.UserMarker]]
	usersIDs     *xsync.Map[int32, *livemapmarkers.UserMarker]
	mappingsMu   sync.RWMutex
	userMappings map[int32]*pbtracker.UserMapping
}

type TestParams struct {
	fx.In

	LC fx.Lifecycle
}

func NewForTests(p TestParams) ITracker {
	t := &TestTracker{
		usersCache:   xsync.NewMap[string, *xsync.Map[int32, *livemapmarkers.UserMarker]](),
		usersIDs:     xsync.NewMap[int32, *livemapmarkers.UserMarker](),
		userMappings: map[int32]*pbtracker.UserMapping{},

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

func (s *TestTracker) SeedUserMarker(marker *livemapmarkers.UserMarker) {
	if marker == nil {
		return
	}

	s.usersIDs.Store(marker.GetUserId(), marker)

	users, _ := s.usersCache.LoadOrStore(
		marker.GetJob(),
		xsync.NewMap[int32, *livemapmarkers.UserMarker](),
	)
	users.Store(marker.GetUserId(), marker)
}

func (s *TestTracker) DeleteUserMarker(userId int32) {
	info, ok := s.usersIDs.Load(userId)
	if !ok {
		return
	}

	s.usersIDs.Delete(userId)

	if info == nil {
		return
	}

	if users, ok := s.usersCache.Load(info.GetJob()); ok {
		users.Delete(userId)
	}
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

func (s *TestTracker) GetUserMapping(userId int32) (*pbtracker.UserMapping, bool, error) {
	s.mappingsMu.RLock()
	defer s.mappingsMu.RUnlock()

	mapping, ok := s.userMappings[userId]
	return mapping, ok, nil
}

func (s *TestTracker) SetUserMapping(_ context.Context, mapping *pbtracker.UserMapping) error {
	if mapping == nil {
		return errors.New("mapping cannot be nil")
	}
	if mapping.GetUserId() <= 0 {
		return fmt.Errorf("invalid user ID: %d", mapping.GetUserId())
	}
	if mapping.UnitId != nil && mapping.GetUnitId() == 0 {
		mapping.UnitId = nil
	}
	if mapping.GetCreatedAt() == nil {
		mapping.CreatedAt = timestamp.Now()
	}

	s.mappingsMu.Lock()
	defer s.mappingsMu.Unlock()
	s.userMappings[mapping.GetUserId()] = mapping
	return nil
}

func (s *TestTracker) SetUserMappingForUser(
	ctx context.Context,
	userId int32,
	unitId *int64,
) error {
	return s.SetUserMapping(ctx, &pbtracker.UserMapping{
		UserId: userId,
		UnitId: unitId,
	})
}

func (s *TestTracker) UnsetUnitIDForUser(ctx context.Context, userId int32) error {
	return s.SetUserMappingForUser(ctx, userId, nil)
}

func (s *TestTracker) DeleteUserMapping(_ context.Context, userId int32) error {
	s.mappingsMu.Lock()
	defer s.mappingsMu.Unlock()
	delete(s.userMappings, userId)
	return nil
}

func (s *TestTracker) ListUserMappings(
	_ context.Context,
) (map[int32]*pbtracker.UserMapping, error) {
	s.mappingsMu.RLock()
	defer s.mappingsMu.RUnlock()

	mappings := make(map[int32]*pbtracker.UserMapping, len(s.userMappings))
	for userID, mapping := range s.userMappings {
		mappings[userID] = mapping
	}
	return mappings, nil
}

func (s *TestTracker) Subscribe(
	_ context.Context,
) (store.IKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker], error) {
	return &TestKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker]{
		broker: s.broker,
	}, nil
}

func (s *TestTracker) GetFilteredUserMarkers(
	acl *permissionsattributes.JobGradeList,
	userInfo *userinfo.UserInfo,
) []*livemapmarkers.UserMarker {
	markers := []*livemapmarkers.UserMarker{}
	s.usersIDs.Range(func(_ int32, marker *livemapmarkers.UserMarker) bool {
		if marker == nil || marker.GetHidden() {
			return true
		}

		grade := marker.GetUser().GetJobGrade()
		if marker.JobGrade != nil {
			grade = marker.GetJobGrade()
		}
		if userInfo != nil && !userInfo.GetJobAdmin() &&
			(acl == nil || !acl.HasJobGrade(marker.GetJob(), grade)) {
			return true
		}

		markers = append(markers, marker)
		return true
	})
	return markers
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
