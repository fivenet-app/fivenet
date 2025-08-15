package tracker

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/tracker"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
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
	Subscribe(
		ctx context.Context,
	) (store.IKVWatcher[livemap.UserMarker, *livemap.UserMarker], error)
	GetFilteredUserMarkers(
		acl *permissions.JobGradeList,
		userInfo *userinfo.UserInfo,
	) []*livemap.UserMarker

	GetUserMapping(userId int32) (*tracker.UserMapping, error)
	SetUserMapping(ctx context.Context, mapping *tracker.UserMapping) error
	SetUserMappingForUser(ctx context.Context, userId int32, unitId *int64) error
	UnsetUnitIDForUser(ctx context.Context, userId int32) error
	ListUserMappings(ctx context.Context) (map[int32]*tracker.UserMapping, error)
}

type Tracker struct {
	logger *zap.Logger
	tracer trace.Tracer
	js     *events.JSWrapper

	jsCons jetstream.ConsumeContext

	userByIDStore     *store.Store[livemap.UserMarker, *livemap.UserMarker]
	userLocStore      *store.Store[livemap.UserMarker, *livemap.UserMarker]
	userMappingsStore *store.Store[tracker.UserMapping, *tracker.UserMapping]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	JS     *events.JSWrapper
	Cfg    *config.Config
}

func New(p Params) (ITracker, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.Named("tracker")
	t := &Tracker{
		logger: logger,
		tracer: p.TP.Tracer("tracker"),
		js:     p.JS,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		storeLogger := logger.WithOptions(
			zap.IncreaseLevel(
				p.Cfg.LogLevelOverrides.Get(config.LoggingComponentKVStore, p.Cfg.LogLevel),
			),
		)

		userMappingsStore, err := store.New[tracker.UserMapping, *tracker.UserMapping](
			ctxStartup, storeLogger, p.JS, BucketUserMappingsMap,
			store.WithLocks[tracker.UserMapping, *tracker.UserMapping](nil),
		)
		if err != nil {
			return err
		}
		if err := userMappingsStore.Start(ctxCancel, false); err != nil {
			return err
		}
		t.userMappingsStore = userMappingsStore

		userLocStore, err := store.New[livemap.UserMarker, *livemap.UserMarker](
			ctxStartup,
			storeLogger,
			p.JS,
			BucketUserLoc,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
		)
		if err != nil {
			return err
		}
		if err := userLocStore.Start(ctxCancel, false); err != nil {
			return err
		}
		t.userLocStore = userLocStore

		byID, err := store.New[livemap.UserMarker, *livemap.UserMarker](
			ctxStartup, storeLogger, p.JS, BucketUserLocByID,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
		)
		if err != nil {
			return err
		}
		if err := byID.Start(ctxCancel, false); err != nil {
			return err
		}
		t.userByIDStore = byID

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
		// key format = JOB.GRADE.USER_ID → cut at first dot
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
	marker, err := t.userByIDStore.Get(strconv.Itoa(int(id)))
	if err != nil {
		return nil, false
	}
	return marker, err == nil
}

func (t *Tracker) IsUserOnDuty(id int32) bool {
	um, err := t.userByIDStore.Get(strconv.Itoa(int(id)))
	return err == nil && um != nil && !um.GetHidden()
}

func (t *Tracker) Subscribe(
	ctx context.Context,
) (store.IKVWatcher[livemap.UserMarker, *livemap.UserMarker], error) {
	return t.userLocStore.WatchAll(ctx)
}

func (t *Tracker) GetFilteredUserMarkers(
	acl *permissions.JobGradeList,
	userInfo *userinfo.UserInfo,
) []*livemap.UserMarker {
	return t.userLocStore.ListFiltered("", func(key string, um *livemap.UserMarker) bool {
		if um == nil || um.GetHidden() {
			return false
		}

		jg := um.GetUser().GetJobGrade()
		if um.JobGrade != nil {
			jg = um.GetJobGrade()
		}

		if !userInfo.GetSuperuser() && !acl.HasJobGrade(um.GetJob(), jg) {
			return false
		}

		return true // keep this marker
	})
}

func (t *Tracker) GetUserMapping(userId int32) (*tracker.UserMapping, error) {
	mapping, err := t.userMappingsStore.Get(UserIdKey(userId))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve unit mapping for user %d. %w", userId, err)
	}
	return mapping, nil
}

func (t *Tracker) SetUserMapping(ctx context.Context, mapping *tracker.UserMapping) error {
	if mapping == nil {
		return errors.New("mapping cannot be nil")
	}

	if mapping.GetUserId() <= 0 {
		return fmt.Errorf("invalid user ID: %d", mapping.GetUserId())
	}

	if mapping.UnitId != nil && mapping.GetUnitId() == 0 {
		mapping.UnitId = nil // unset if zero
	}

	if mapping.GetCreatedAt() == nil {
		mapping.CreatedAt = timestamp.Now()
	}

	if err := t.userMappingsStore.Put(ctx, UserIdKey(mapping.GetUserId()), mapping); err != nil {
		return fmt.Errorf("failed to set unit mapping for user %d. %w", mapping.GetUserId(), err)
	}

	return nil
}

func (t *Tracker) SetUserMappingForUser(ctx context.Context, userId int32, unitId *int64) error {
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

func (t *Tracker) DeleteUserMapping(ctx context.Context, userId int32) error {
	if err := t.userMappingsStore.Delete(ctx, UserIdKey(userId)); err != nil {
		return fmt.Errorf("failed to delete unit mapping for user %d. %w", userId, err)
	}
	return nil
}

func (t *Tracker) ListUserMappings(ctx context.Context) (map[int32]*tracker.UserMapping, error) {
	mappings := t.userMappingsStore.List()
	ids := map[int32]*tracker.UserMapping{}
	for _, m := range mappings {
		ids[m.GetUserId()] = m
	}

	return ids, nil
}
