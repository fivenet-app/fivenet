package croner

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var StateModule = fx.Module("cron_scheduler",
	fx.Provide(
		NewState,
	),
)

type jobWrapper struct {
	ctx    context.Context
	cancel context.CancelFunc

	schedule string
}

type StateParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
}

type State struct {
	logger *zap.Logger

	ctx   context.Context
	js    *events.JSWrapper
	store *store.Store[cron.Cronjob, *cron.Cronjob]

	cronjobs *xsync.Map[string, *jobWrapper]
}

func NewState(p StateParams) *State {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &State{
		logger: p.Logger.Named("cron_scheduler"),
		ctx:    ctxCancel,
		js:     p.JS,

		cronjobs: xsync.NewMap[string, *jobWrapper](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxStartup, s.js); err != nil {
			return err
		}

		st, err := store.New(ctxStartup, p.Logger, p.JS, "cron",
			store.WithOnUpdateFn(func(_ *store.Store[cron.Cronjob, *cron.Cronjob], cj *cron.Cronjob) (*cron.Cronjob, error) {
				if cj == nil {
					return cj, nil
				}

				jw, ok := s.cronjobs.Load(cj.Name)
				if !ok {
					ctx, cancel := context.WithCancel(ctxCancel)
					s.cronjobs.Store(cj.Name, &jobWrapper{
						ctx:    ctx,
						cancel: cancel,

						schedule: cj.Schedule,
					})
				} else {
					if cj.Schedule != jw.schedule {
						jw.schedule = cj.Schedule
					}
				}

				return cj, nil
			}),

			store.WithOnDeleteFn(func(_ *store.Store[cron.Cronjob, *cron.Cronjob], entry jetstream.KeyValueEntry, cj *cron.Cronjob) error {
				jw, ok := s.cronjobs.LoadAndDelete(entry.Key())
				if !ok {
					return nil
				}

				jw.cancel()
				jw = nil

				return nil
			}),
		)
		if err != nil {
			return err
		}
		s.store = st

		if err := st.Start(ctxCancel, true); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return nil
	}))

	return s
}

func (s *State) ListCronjobs(ctx context.Context) []*cron.Cronjob {
	cj := []*cron.Cronjob{}

	s.store.Range(ctx, func(_ string, entry *cron.Cronjob) bool {
		cj = append(cj, entry)

		return true
	})

	return cj
}
