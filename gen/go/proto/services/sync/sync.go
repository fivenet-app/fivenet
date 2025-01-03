package sync

import (
	"context"
	"database/sql"
	"slices"

	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	SyncServiceServer

	logger *zap.Logger
	db     *sql.DB
	auth   *auth.GRPCAuth

	tokens []string
}

type Params struct {
	fx.In

	Logger *zap.Logger
	DB     *sql.DB
	Auth   *auth.GRPCAuth

	Config *config.Config
}

func NewServer(p Params) *Server {
	return &Server{
		logger: p.Logger,
		db:     p.DB,
		auth:   p.Auth,

		tokens: p.Config.Auth.SyncAPITokens,
	}
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

func (s *Server) SyncData(ctx context.Context, req *SyncDataRequest) (*SyncDataResponse, error) {
	// TODO handle sync data request

	return nil, nil
}
