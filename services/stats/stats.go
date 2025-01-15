package stats

import (
	"context"
	"database/sql"

	stats "github.com/fivenet-app/fivenet/gen/go/proto/resources/stats"
	pbstats "github.com/fivenet-app/fivenet/gen/go/proto/services/stats"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Stats = map[string]*stats.Stat

type Server struct {
	pbstats.StatsServiceServer

	logger *zap.Logger
	js     *events.JSWrapper

	worker *worker
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	DB        *sql.DB
	JS        *events.JSWrapper
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("stats_worker"),
		js:     p.JS,

		worker: newWorker(p.Logger, p.DB),
	}

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		if p.AppConfig.Get().Website.StatsPage {
			s.worker.Start()
		}

		go func() {
			configUpdateCh := p.AppConfig.Subscribe()

			for {
				select {
				case <-ctx.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					if cfg == nil {
						continue
					}

					if cfg.Website.StatsPage {
						s.worker.Start()
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

func (s *Server) PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Skip permission check for the stats services
	return ctx, nil
}

func (s *Server) GetStats(ctx context.Context, req *pbstats.GetStatsRequest) (*pbstats.GetStatsResponse, error) {
	stats := s.worker.GetStats()
	if stats == nil {
		return &pbstats.GetStatsResponse{}, nil
	}

	return &pbstats.GetStatsResponse{
		Stats: *stats,
	}, nil
}
