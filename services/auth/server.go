package auth

import (
	"context"
	"database/sql"

	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbauth.AuthServiceServer

	logger   *zap.Logger
	db       *sql.DB
	auth     *auth.GRPCAuth
	tm       *auth.TokenMgr
	ps       perms.Permissions
	enricher *mstlystcdata.Enricher
	aud      audit.IAuditer
	ui       userinfo.UserInfoRetriever
	appCfg   appconfig.IConfig
	js       *events.JSWrapper

	domain          string
	oauth2Providers []*config.OAuth2Provider
	customDB        config.CustomDB
	superuserGroups []string
	superuserUsers  []string
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	Auth      *auth.GRPCAuth
	TM        *auth.TokenMgr
	Perms     perms.Permissions
	Enricher  *mstlystcdata.Enricher
	Aud       audit.IAuditer
	UI        userinfo.UserInfoRetriever
	Config    *config.Config
	AppConfig appconfig.IConfig
	JS        *events.JSWrapper
}

func NewServer(p Params) *Server {
	return &Server{
		logger:   p.Logger.Named("grpc.auth"),
		db:       p.DB,
		auth:     p.Auth,
		tm:       p.TM,
		ps:       p.Perms,
		enricher: p.Enricher,
		aud:      p.Aud,
		ui:       p.UI,
		appCfg:   p.AppConfig,
		js:       p.JS,

		domain:          p.Config.HTTP.Sessions.Domain,
		oauth2Providers: p.Config.OAuth2.Providers,
		customDB:        p.Config.Database.Custom,
		superuserGroups: p.Config.Auth.SuperuserGroups,
		superuserUsers:  p.Config.Auth.SuperuserUsers,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbauth.RegisterAuthServiceServer(srv, s)
}

// AuthFuncOverride is called instead of the original auth func
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	// Skip authentication for the anon accessible endpoints
	if fullMethod == "/services.auth.AuthService/CreateAccount" ||
		fullMethod == "/services.auth.AuthService/Login" ||
		fullMethod == "/services.auth.AuthService/ForgotPassword" {
		return ctx, nil
	}

	if fullMethod == "/services.auth.AuthService/Logout" {
		c, _ := s.auth.GRPCAuthFunc(ctx, fullMethod)
		if c != nil {
			return c, nil
		}
		return ctx, nil
	}

	if fullMethod == "/services.auth.AuthService/SetSuperuserMode" {
		return s.auth.GRPCAuthFunc(ctx, fullMethod)
	}

	return s.auth.GRPCAuthFuncWithoutUserInfo(ctx, fullMethod)
}

func (s *Server) PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Skip permission check for the auth services
	return ctx, nil
}
