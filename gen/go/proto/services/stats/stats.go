package stats

import (
	"context"
	"database/sql"
	"sync/atomic"

	stats "github.com/fivenet-app/fivenet/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/pkg/events"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Stats = map[string]*stats.Stat

type Server struct {
	StatsServiceServer

	logger *zap.Logger
	db     *sql.DB
	js     *events.JSWrapper

	stats atomic.Pointer[Stats]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	JS     *events.JSWrapper
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("stats_worker"),
		db:     p.DB,
		js:     p.JS,

		stats: atomic.Pointer[map[string]*stats.Stat]{},
	}
	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go s.calculateStats(ctx)

		return nil
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterStatsServiceServer(srv, s)
}

func (s *Server) PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Skip permission check for the stats services
	return ctx, nil
}

func (s *Server) GetStats(ctx context.Context, req *GetStatsRequest) (*GetStatsResponse, error) {
	stats := s.stats.Load()
	return &GetStatsResponse{
		Stats: *stats,
	}, nil
}
