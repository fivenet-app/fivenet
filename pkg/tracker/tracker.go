package tracker

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ITracker interface {
	ListTrackedJobs() []string
	GetUserMarkerById(id int32) (*livemap.UserMarker, bool)
	IsUserOnDuty(userId int32) bool

	Subscribe(ctx context.Context) (chan *store.KeyValueEntry[livemap.UserMarker, *livemap.UserMarker], error)

	GetUserMapping(userId int32) (*tracker.UserMapping, error)
	SetUserMapping(ctx context.Context, mapping *tracker.UserMapping) error
	SetUserMappingForUser(ctx context.Context, userId int32, unitId *uint64) error
	UnsetUnitIDForUser(ctx context.Context, userId int32) error
	ListUserMappings(ctx context.Context) (map[int32]*tracker.UserMapping, error)
}

type Tracker struct {
	ITracker

	logger *zap.Logger
	tracer trace.Tracer
	js     *events.JSWrapper

	jsCons jetstream.ConsumeContext

	userLocStore      *store.Store[livemap.UserMarker, *livemap.UserMarker]
	byIDStore         *store.Store[livemap.UserMarker, *livemap.UserMarker]
	userMappingsStore *store.Store[tracker.UserMapping, *tracker.UserMapping]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	JS     *events.JSWrapper
}

func New(p Params) (ITracker, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	t := &Tracker{
		logger: p.Logger.Named("tracker"),
		tracer: p.TP.Tracer("tracker"),
		js:     p.JS,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		userMappingsStore, err := store.New[tracker.UserMapping, *tracker.UserMapping](
			ctxStartup, p.Logger, p.JS, BucketUserMappingsMap,
			store.WithLocks[tracker.UserMapping, *tracker.UserMapping](nil),
		)
		if err != nil {
			return err
		}
		if err := userMappingsStore.Start(ctxCancel, false); err != nil {
			return err
		}
		t.userMappingsStore = userMappingsStore

		byID, err := store.New[livemap.UserMarker, *livemap.UserMarker](
			ctxStartup, p.Logger, p.JS, BucketUserLocByID,
		)
		if err != nil {
			return err
		}
		if err := byID.Start(ctxCancel, false); err != nil {
			return err
		}
		t.byIDStore = byID

		userLocStore, err := store.New(ctxStartup, p.Logger, p.JS, BucketUserLoc,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
			store.WithIgnoredKeys[livemap.UserMarker, *livemap.UserMarker]("_snapshot"),
			store.WithOnUpdateFn(func(_ *store.Store[livemap.UserMarker, *livemap.UserMarker],
				um *livemap.UserMarker,
			) (*livemap.UserMarker, error) {
				if um == nil {
					return nil, nil
				}

				// Upsert mapping (unit_id may be 0 = no unit)
				if err := t.SetUserMappingForUser(ctxCancel, um.UserId, um.UnitId); err != nil {
					return nil, fmt.Errorf("failed to upsert user unit mapping. %w", err)
				}

				return um, nil
			}),
		)
		if err != nil {
			return err
		}
		if err := userLocStore.Start(ctxCancel, false); err != nil {
			return err
		}
		t.userLocStore = userLocStore

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		if t.jsCons != nil {
			t.jsCons.Stop()
			t.jsCons = nil
		}

		return nil
	}))

	return t, nil
}

// ListTrackedJobs returns the distinct job strings that currently have
// at least one live UserMarker in the KV cache.
//
// Cost: O(#unique-jobs) + one Range pass over userLocStore.
// Safe for concurrent callers.
func (t *Tracker) ListTrackedJobs() []string {
	seen := make(map[string]struct{})

	t.userLocStore.Range(func(key string, _ *livemap.UserMarker) bool {
		// key format = JOB.GRADE.USER_ID -> cut at first dot
		if i := strings.IndexByte(key, '.'); i > 0 {
			seen[key[:i]] = struct{}{}
		}
		return true // continue iteration
	})

	jobs := make([]string, 0, len(seen))
	for j := range seen {
		jobs = append(jobs, j)
	}
	slices.Sort(jobs)
	return jobs
}

func (t *Tracker) GetUserMarkerById(id int32) (*livemap.UserMarker, bool) {
	mapping, err := t.byIDStore.Get(fmt.Sprint(id))
	if err != nil {
		return nil, false
	}
	return mapping, err == nil
}

func (t *Tracker) IsUserOnDuty(id int32) bool {
	_, err := t.userMappingsStore.Get(fmt.Sprint(id))
	return err == nil
}

func (s *Tracker) Subscribe(ctx context.Context) (chan *store.KeyValueEntry[livemap.UserMarker, *livemap.UserMarker], error) {
	return s.userLocStore.WatchAll(ctx)
}

func FilterMarkers(markers []*livemap.UserMarker, acl *permissions.JobGradeList, userInfo *userinfo.UserInfo) []*livemap.UserMarker {
	if len(markers) == 0 || acl == nil {
		return markers
	}

	filtered := make([]*livemap.UserMarker, 0, len(markers))
	for _, um := range markers {
		jg := um.User.JobGrade
		if um.JobGrade != nil {
			jg = *um.JobGrade
		}

		if !userInfo.Superuser && !acl.HasJobGrade(um.Job, jg) {
			continue
		}

		filtered = append(filtered, um)
	}

	return filtered
}

func (t *Tracker) GetUserMapping(userId int32) (*tracker.UserMapping, error) {
	mapping, err := t.userMappingsStore.Get(userIdKey(userId))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve unit mapping for user %d. %w", userId, err)
	}
	return mapping, nil
}

func (t *Tracker) SetUserMapping(ctx context.Context, mapping *tracker.UserMapping) error {
	if mapping == nil {
		return fmt.Errorf("mapping cannot be nil")
	}

	if mapping.UserId <= 0 {
		return fmt.Errorf("invalid user ID: %d", mapping.UserId)
	}

	if mapping.UnitId != nil && *mapping.UnitId == 0 {
		mapping.UnitId = nil // unset if zero
	}

	if mapping.CreatedAt == nil {
		mapping.CreatedAt = timestamp.Now()
	}

	if err := t.userMappingsStore.Put(ctx, userIdKey(mapping.UserId), mapping); err != nil {
		return fmt.Errorf("failed to set unit mapping for user %d. %w", mapping.UserId, err)
	}

	return nil
}

func (t *Tracker) SetUserMappingForUser(ctx context.Context, userId int32, unitId *uint64) error {
	if err := t.SetUserMapping(ctx, &tracker.UserMapping{
		UserId: userId,
		UnitId: unitId,
	}); err != nil {
		return err
	}

	return nil
}

func (t *Tracker) UnsetUnitIDForUser(ctx context.Context, userId int32) error {
	return t.SetUserMappingForUser(ctx, userId, nil)
}

func (t *Tracker) ListUserMappings(ctx context.Context) (map[int32]*tracker.UserMapping, error) {
	mappings := t.userMappingsStore.List()
	ids := map[int32]*tracker.UserMapping{}
	for _, m := range mappings {
		ids[m.UserId] = m
	}

	return ids, nil
}
