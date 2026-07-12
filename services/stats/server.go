package stats

import (
	"context"

	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	statsstore "github.com/fivenet-app/fivenet/v2026/stores/stats"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pbstats.StatsServiceServer

	worker *worker
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	AppConfig appconfig.IConfig
	Store     statsstore.IStore
}

func NewServer(p Params) *Server {
	s := &Server{
		worker: newWorker(p.Logger, p.Store),
	}

	ctxCancel, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		if p.AppConfig.Get().Website.GetStatsPage() {
			s.worker.Start(ctxCancel)
		}

		go func() {
			configUpdateCh := p.AppConfig.Subscribe()

			for {
				select {
				case <-ctxCancel.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					if cfg == nil {
						continue
					}

					if cfg.Website.GetStatsPage() {
						s.worker.Start(ctxCancel)
					} else {
						s.worker.Stop()
					}
				}
			}
		}()

		return nil
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbstats.RegisterStatsServiceServer(srv, s)
}

func (s *Server) PermissionUnaryFuncOverride(
	ctx context.Context,
	info *grpc.UnaryServerInfo,
) (context.Context, error) {
	// Skip permission check for the stats services
	return ctx, nil
}
