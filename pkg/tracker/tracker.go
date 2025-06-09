package tracker

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
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

	userStore  *store.Store[livemap.UserMarker, *livemap.UserMarker]
	usersByJob *xsync.Map[string, *xsync.Map[int32, *livemap.UserMarker]]
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
		userIDs, err := store.New(ctxStartup, p.Logger, p.JS, "tracker",
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
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

		if err := userIDs.Start(ctxCancel, false); err != nil {
			return err
		}
		t.userStore = userIDs

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

func (s *Tracker) GetUserById(id int32) (*livemap.UserMarker, bool) {
	return s.userStore.Get(userIdKey(id))
}

func (s *Tracker) IsUserOnDuty(userId int32) bool {
	um, ok := s.userStore.Get(userIdKey(userId))
	if !ok {
		return false
	}

	return !um.Hidden
}

func (s *Tracker) Subscribe(ctx context.Context) (chan *store.KeyValueEntry[livemap.UserMarker, *livemap.UserMarker], error) {
	return s.userStore.WatchAll(ctx)
}
