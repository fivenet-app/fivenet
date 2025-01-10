package sync

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	SyncServiceServer

	logger *zap.Logger

	db   *sql.DB
	js   *events.JSWrapper
	auth *auth.GRPCAuth
	cfg  *config.Config

	esxCompat bool
	tokens    []string
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	JS     *events.JSWrapper
	Auth   *auth.GRPCAuth
	Config *config.Config
}

func NewServer(p Params) *Server {
	if !p.Config.Sync.Enabled {
		return nil
	}

	s := &Server{
		logger: p.Logger,
		db:     p.DB,
		js:     p.JS,
		auth:   p.Auth,
		cfg:    p.Config,

		esxCompat: p.Config.Database.ESXCompat,
		tokens:    p.Config.Sync.APITokens,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if _, err := s.registerStream(ctxStartup, s.js); err != nil {
			return fmt.Errorf("failed to register stream: %w", err)
		}

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterSyncServiceServer(srv, s)
}

// AuthFuncOverride is called instead of the original auth func
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	if fullMethod == "/services.sync.SyncService/GetStatus" {
		c, _ := s.auth.GRPCAuthFunc(ctx, fullMethod)
		if c != nil {
			return c, nil
		}
	}

	t, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(s.tokens, t) {
		return nil, auth.ErrInvalidToken
	}

	return ctx, nil
}

func (s *Server) PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Skip permission check for the sync service
	return ctx, nil
}

func (s *Server) StreamServerInterceptor(ctx context.Context, srv interface{}, info *grpc.StreamServerInfo) (context.Context, error) {
	// Skip permission check for the sync service
	return ctx, nil
}
