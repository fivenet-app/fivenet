package auth

import (
	"context"
	"database/sql"

	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	authstore "github.com/fivenet-app/fivenet/v2026/stores/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetAccounts,
		IDColumn:        table.FivenetAccounts.ID,
		DeletedAtColumn: table.FivenetAccounts.DeletedAt,

		MinDays: 21,
	})
}

type Server struct {
	pbauth.AuthServiceServer

	logger   *zap.Logger
	auth     *auth.GRPCAuth
	tm       *auth.TokenMgr
	ps       perms.Permissions
	enricher mstlystcdata.IEnricher
	ui       userinfo.UserInfoRetriever
	appCfg   appconfig.IConfig
	js       *events.JSWrapper
	store    authstore.IStore

	domain          string
	oauth2Providers []*config.OAuth2Provider
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
	Enricher  mstlystcdata.IEnricher
	UI        userinfo.UserInfoRetriever
	Config    *config.Config
	AppConfig appconfig.IConfig
	JS        *events.JSWrapper
	Store     authstore.IStore `optional:"true"`
}

func NewServer(p Params) *Server {
	return &Server{
		logger:   p.Logger.Named("grpc.auth"),
		auth:     p.Auth,
		tm:       p.TM,
		ps:       p.Perms,
		enricher: p.Enricher,
		ui:       p.UI,
		appCfg:   p.AppConfig,
		js:       p.JS,
		store:    p.Store,

		domain:          p.Config.HTTP.Sessions.Domain,
		oauth2Providers: p.Config.OAuth2.Providers,
		superuserGroups: p.Config.Auth.SuperuserGroups,
		superuserUsers:  p.Config.Auth.SuperuserUsers,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbauth.RegisterAuthServiceServer(srv, s)
}

// AuthFuncOverride is called instead of the original auth func.
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

func (s *Server) PermissionUnaryFuncOverride(
	ctx context.Context,
	info *grpc.UnaryServerInfo,
) (context.Context, error) {
	// Skip permission check for the auth services
	return ctx, nil
}
