package tracker

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type ITracker interface {
	GetUsersByJob(job string) (*xsync.Map[int32, *livemap.UserMarker], bool)
	ListTrackedJobs() []string
	GetUserById(id int32) (*livemap.UserMarker, bool)
	IsUserOnDuty(userId int32) bool

	Subscribe(ctx context.Context) (chan *store.KeyValueEntry[livemap.UserMarker, *livemap.UserMarker], error)
}

type Tracker struct {
	ITracker

	logger *zap.Logger
	tracer trace.Tracer
	js     *events.JSWrapper

	jsCons jetstream.ConsumeContext

	userLocStore *store.Store[livemap.UserMarker, *livemap.UserMarker]
	unitMapStore *store.Store[centrum.UserUnitMapping, *centrum.UserUnitMapping]
	usersByJob   *xsync.Map[string, *xsync.Map[int32, *livemap.UserMarker]]
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

		usersByJob: xsync.NewMap[string, *xsync.Map[int32, *livemap.UserMarker]](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		userLocStore, err := store.New(ctxStartup, p.Logger, p.JS, BucketUserLoc,
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
			store.WithIgnoredKeys[livemap.UserMarker, *livemap.UserMarker]("_snapshot"),
			store.WithOnUpdateFn(func(s *store.Store[livemap.UserMarker, *livemap.UserMarker], um *livemap.UserMarker) (*livemap.UserMarker, error) {
				if um == nil {
					return um, nil
				}

				jobUsers, _ := t.usersByJob.LoadOrCompute(um.Job, func() (*xsync.Map[int32, *livemap.UserMarker], bool) {
					return xsync.NewMap[int32, *livemap.UserMarker](), false
				})
				// Maybe we can be smarter about updating the user marker here, but
				// without mutexes it will be problematic (data races and Co.)
				// Is `proto.Clone` really the solution to this?
				jobUsers.Store(um.UserId, proto.Clone(um).(*livemap.UserMarker))

				return um, nil
			}),

			store.WithOnDeleteFn(func(s *store.Store[livemap.UserMarker, *livemap.UserMarker], entry jetstream.KeyValueEntry, um *livemap.UserMarker) error {
				if um == nil {
					return nil
				}

				if jobUsers, ok := t.usersByJob.Load(um.Job); ok {
					jobUsers.Delete(um.UserId)
				}

				return nil
			}),
		)
		if err != nil {
			return err
		}
		if err := userLocStore.Start(ctxCancel, false); err != nil {
			return err
		}
		t.userLocStore = userLocStore

		unitStore, err := store.New[centrum.UserUnitMapping, *centrum.UserUnitMapping](
			ctxStartup, p.Logger, p.JS, BucketUnitMap,
			store.WithLocks[centrum.UserUnitMapping, *centrum.UserUnitMapping](nil),
		)
		if err != nil {
			return err
		}
		if err := unitStore.Start(ctxCancel, false); err != nil {
			return err
		}
		t.unitMapStore = unitStore

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

// Returns a `xsync.Map` with **copies** (proto cloned) of the `*livemap.UserMarker`
func (s *Tracker) GetUsersByJob(job string) (*xsync.Map[int32, *livemap.UserMarker], bool) {
	return s.usersByJob.Load(job)
}

func (s *Tracker) ListTrackedJobs() []string {
	var jobs []string
	s.usersByJob.Range(func(job string, _ *xsync.Map[int32, *livemap.UserMarker]) bool {
		jobs = append(jobs, job)

		return true
	})

	return jobs
}

func (t *Tracker) GetUserById(id int32) (*livemap.UserMarker, bool) {
	mapping, err := t.unitMapStore.Get(fmt.Sprint(id))
	if err != nil {
		return nil, false
	}
	key := fmt.Sprintf("%s.%d.%d", mapping.Job, mapping.JobGrade, id)
	marker, err := t.userLocStore.Get(key)
	return marker, err == nil
}

func (t *Tracker) IsUserOnDuty(id int32) bool {
	_, err := t.unitMapStore.Get(fmt.Sprint(id))
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
