package tracker

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/nats/store"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	tLocs     = table.FivenetUserLocations
	tUsers    = table.Users.AS("user")
	tJobProps = table.FivenetJobProps
)

type ITracker interface {
	GetUsersByJob(job string) (*xsync.MapOf[int32, *livemap.UserMarker], bool)
	GetUserById(id int32) (*livemap.UserMarker, bool)
	GetAllActiveUsers() ([]*livemap.UserMarker, error)
	IsUserOnDuty(userId int32) bool

	Subscribe() chan *livemap.UsersUpdateEvent
	Unsubscribe(c chan *livemap.UsersUpdateEvent)
}

type Tracker struct {
	ITracker

	logger *zap.Logger
	tracer trace.Tracer
	js     *events.JSWrapper

	jsCons jetstream.ConsumeContext

	userStore  *store.Store[livemap.UserMarker, *livemap.UserMarker]
	usersByJob *xsync.MapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]]

	broker *broker.Broker[*livemap.UsersUpdateEvent]
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

		usersByJob: xsync.NewMapOf[string, *xsync.MapOf[int32, *livemap.UserMarker]](),

		broker: broker.New[*livemap.UsersUpdateEvent](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerStreams(ctxStartup, p.JS); err != nil {
			return err
		}

		go t.broker.Start(ctxCancel)

		userIDs, err := store.New(ctxStartup, p.Logger, p.JS, "tracker",
			store.WithLocks[livemap.UserMarker, *livemap.UserMarker](nil),
			store.WithOnUpdateFn(func(s *store.Store[livemap.UserMarker, *livemap.UserMarker], um *livemap.UserMarker) (*livemap.UserMarker, error) {
				if um == nil || um.Info == nil {
					return um, nil
				}

				jobUsers, _ := t.usersByJob.LoadOrCompute(um.Info.Job, func() *xsync.MapOf[int32, *livemap.UserMarker] {
					return xsync.NewMapOf[int32, *livemap.UserMarker]()
				})
				// Maybe we can be smarter about updating the user marker here, but
				// without mutexes it will be problematic (data races and Co.)
				jobUsers.Store(um.UserId, um)

				return um, nil
			}),

			store.WithOnDeleteFn(func(s *store.Store[livemap.UserMarker, *livemap.UserMarker], entry jetstream.KeyValueEntry, um *livemap.UserMarker) error {
				if um == nil || um.Info == nil {
					return nil
				}

				if jobUsers, ok := t.usersByJob.Load(um.Info.Job); ok {
					jobUsers.Delete(um.UserId)
				}

				return nil
			}),
		)
		if err != nil {
			return err
		}

		if err := userIDs.Start(ctxCancel); err != nil {
			return err
		}
		t.userStore = userIDs

		if err := t.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to register tracker nats subscriptions. %w", err)
		}

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

func (s *Tracker) watchForChanges(msg jetstream.Msg) {
	remoteCtx, _ := events.GetJetstreamMsgContext(msg)
	_, span := s.tracer.Start(trace.ContextWithRemoteSpanContext(context.Background(), remoteCtx), msg.Subject())
	defer span.End()

	if err := msg.Ack(); err != nil {
		s.logger.Error("failed to ack message", zap.Error(err))
	}

	dest := &livemap.UsersUpdateEvent{}
	if err := proto.Unmarshal(msg.Data(), dest); err != nil {
		s.logger.Error("failed to unmarshal nats user update response", zap.Error(err))
		return
	}

	if s.broker.SubCount() <= 0 {
		return
	}

	s.broker.Publish(dest)
}

func (s *Tracker) GetUsersByJob(job string) (*xsync.MapOf[int32, *livemap.UserMarker], bool) {
	return s.usersByJob.Load(job)
}

func (s *Tracker) GetAllActiveUsers() ([]*livemap.UserMarker, error) {
	return s.userStore.List()
}

func (s *Tracker) IsUserOnDuty(userId int32) bool {
	_, ok := s.userStore.Get(userIdKey(userId))
	return ok
}

func (s *Tracker) GetUserById(id int32) (*livemap.UserMarker, bool) {
	return s.userStore.Get(userIdKey(id))
}

func (s *Tracker) Subscribe() chan *livemap.UsersUpdateEvent {
	return s.broker.Subscribe()
}

func (s *Tracker) Unsubscribe(c chan *livemap.UsersUpdateEvent) {
	s.broker.Unsubscribe(c)
}
