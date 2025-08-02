package sync

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"sync/atomic"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	pkggrpc "github.com/fivenet-app/fivenet/v2025/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pbsync.SyncServiceServer

	logger *zap.Logger

	db   *sql.DB
	js   *events.JSWrapper
	auth *auth.GRPCAuth
	cfg  *config.Config

	dispatches *dispatches.DispatchDB

	esxCompat bool
	tokens    []string

	lastSyncedData     atomic.Int64
	lastSyncedActivity atomic.Int64
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger     *zap.Logger
	DB         *sql.DB
	JS         *events.JSWrapper
	Auth       *auth.GRPCAuth
	Config     *config.Config
	DispatchDB *dispatches.DispatchDB
}

type Result struct {
	fx.Out

	Server  *Server
	Service pkggrpc.Service `group:"grpcservices"`
}

func NewServer(p Params) (Result, error) {
	s := &Server{
		logger: p.Logger.Named("sync"),
		db:     p.DB,
		js:     p.JS,
		auth:   p.Auth,
		cfg:    p.Config,

		dispatches: p.DispatchDB,

		esxCompat: p.Config.Database.ESXCompat,
		tokens:    p.Config.Sync.APITokens,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if _, err := s.registerStream(ctxStartup, s.js); err != nil {
			return fmt.Errorf("failed to register stream. %w", err)
		}

		return nil
	}))

	return Result{
		Server:  s,
		Service: s,
	}, nil
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbsync.RegisterSyncServiceServer(srv, s)
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
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	return ctx, nil
}

func (s *Server) PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Skip permission check for the sync service
	return ctx, nil
}

func (s *Server) PermissionStreamFuncOverride(ctx context.Context, srv any, info *grpc.StreamServerInfo) (context.Context, error) {
	// Skip permission check for the sync service
	return ctx, nil
}

func (s *Server) GetSyncTimes() *pbsettings.DBSyncStatus {
	st := &pbsettings.DBSyncStatus{
		Enabled: s.cfg.Sync.Enabled,
	}

	lastSyncedData := s.lastSyncedData.Load()
	if lastSyncedData > 0 {
		st.LastSyncedData = timestamp.New(time.Unix(lastSyncedData, 0))
	}

	lastSyncedActivity := s.lastSyncedActivity.Load()
	if lastSyncedActivity > 0 {
		st.LastSyncedActivity = timestamp.New(time.Unix(lastSyncedActivity, 0))
	}

	return st
}
