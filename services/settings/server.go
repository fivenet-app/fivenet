package settings

import (
	"context"
	"database/sql"
	"errors"

	discordapi "github.com/diamondburned/arikawa/v3/api"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/crypt"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/reqs"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/pkg/updatecheck"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	syncservice "github.com/fivenet-app/fivenet/v2026/services/sync"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	settingsstore "github.com/fivenet-app/fivenet/v2026/stores/settings"
	"github.com/go-jet/jet/v2/mysql"
	grpcmetadata "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

const (
	roleIDLogFieldKey  = "fivenet.settings.role_id"
	jobNameLogFieldKey = "fivenet.settings.job"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetJobProps,
		JobColumn:       table.FivenetJobProps.Job,
		DeletedAtColumn: table.FivenetJobProps.DeletedAt,

		MinDays: 7,
	})
}

type Server struct {
	pbsettings.SettingsServiceServer
	pbsettings.ConfigServiceServer
	pbsettings.CronServiceServer
	pbsettings.LawsServiceServer
	pbsettings.AccountsServiceServer
	pbsettings.SystemServiceServer

	logger       *zap.Logger
	db           *sql.DB
	auth         *grpcauth.GRPCAuth
	ps           perms.Permissions
	enricher     mstlystcdata.IUserAwareEnricher
	laws         mstlystcdata.ILaws
	cfg          *config.Config
	appCfg       appconfig.IConfig
	js           *events.JSWrapper
	cronRegistry *croner.Registry
	croner       croner.IScheduler
	crypt        *crypt.Crypt
	notifi       notifi.INotifi

	jobPropsFileHandler *filestore.Handler[string]

	dc               *discordapi.Client
	dcOAuth2Provider *config.OAuth2Provider

	syncServer    *syncservice.Server
	dbReq         *reqs.DBReqs
	natsReq       *reqs.NatsReqs
	updateChecker *updatecheck.Checker
	store         settingsstore.IStore
	jobsStore     jobsstore.IStore
}

type Params struct {
	fx.In

	Logger       *zap.Logger
	DB           *sql.DB
	Auth         *grpcauth.GRPCAuth
	PS           perms.Permissions
	Enricher     mstlystcdata.IUserAwareEnricher
	Laws         mstlystcdata.ILaws
	Storage      storage.IStorage
	Config       *config.Config
	AppConfig    appconfig.IConfig
	JS           *events.JSWrapper
	Croner       croner.IScheduler
	CronRegistry *croner.Registry
	Crypt        *crypt.Crypt
	Notifi       notifi.INotifi

	SyncServer    *syncservice.Server
	DBReq         *reqs.DBReqs
	NatsReq       *reqs.NatsReqs
	UpdateChecker *updatecheck.Checker
	Store         settingsstore.IStore
	JobsStore     jobsstore.IStore
}

func NewServer(p Params) *Server {
	tJobProps := table.FivenetJobProps

	fHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tJobProps,
		tJobProps.Job,
		tJobProps.LogoFileID,
		2<<20,
		1,
		func(parentID string) mysql.BoolExpression {
			return tJobProps.Job.EQ(mysql.String(parentID))
		},
		filestore.UpdateJoinRow,
		true,
	).WithUploadFilter(filestore.NewImageUploadFilter())

	var dcOAuth2Provider *config.OAuth2Provider
	var dc *discordapi.Client
	if p.Config.Discord.Enabled {
		dc = discordapi.NewClient("Bot " + p.Config.Discord.Token)

		for _, provider := range p.Config.OAuth2.Providers {
			if provider.Name == "discord" {
				dcOAuth2Provider = provider
				break
			}
		}
	}

	s := &Server{
		logger:       p.Logger,
		db:           p.DB,
		auth:         p.Auth,
		ps:           p.PS,
		enricher:     p.Enricher,
		laws:         p.Laws,
		cfg:          p.Config,
		appCfg:       p.AppConfig,
		js:           p.JS,
		croner:       p.Croner,
		cronRegistry: p.CronRegistry,
		crypt:        p.Crypt,
		notifi:       p.Notifi,

		jobPropsFileHandler: fHandler,

		dc:               dc,
		dcOAuth2Provider: dcOAuth2Provider,

		syncServer:    p.SyncServer,
		dbReq:         p.DBReq,
		natsReq:       p.NatsReq,
		updateChecker: p.UpdateChecker,
		store:         p.Store,
		jobsStore:     p.JobsStore,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbsettings.RegisterSettingsServiceServer(srv, s)
	pbsettings.RegisterConfigServiceServer(srv, s)
	pbsettings.RegisterCronServiceServer(srv, s)
	pbsettings.RegisterLawsServiceServer(srv, s)
	pbsettings.RegisterAccountsServiceServer(srv, s)
	pbsettings.RegisterSystemServiceServer(srv, s)
}

// AuthFuncOverride is called instead of the global auth func for settings services.
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	switch fullMethod {
	case pbsettings.ConfigService_GetAppConfig_FullMethodName,
		pbsettings.ConfigService_UpdateAppConfig_FullMethodName:
		if s.auth == nil {
			return nil, errors.New("settings auth is not configured")
		}
		if hasUserTokenInContext(ctx) {
			return s.auth.GRPCAuthFunc(ctx, fullMethod)
		}
		return s.auth.GRPCAuthFuncWithoutUserInfo(ctx, fullMethod)

	default:
		if s.auth == nil {
			return nil, errors.New("settings auth is not configured")
		}
		return s.auth.GRPCAuthFunc(ctx, fullMethod)
	}
}

func hasUserTokenInContext(ctx context.Context) bool {
	md := grpcmetadata.ExtractIncoming(ctx)
	return len(md.Get("authorization")) > 0 && len(md.Get("cookie")) > 0
}
