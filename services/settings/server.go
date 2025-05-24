package settings

import (
	"database/sql"

	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
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
	pbsettings.FilestoreServiceServer
	pbsettings.LawsServiceServer
	pbsettings.AccountsServiceServer

	logger    *zap.Logger
	db        *sql.DB
	ps        perms.Permissions
	aud       audit.IAuditer
	enricher  *mstlystcdata.Enricher
	laws      *mstlystcdata.Laws
	st        storage.IStorage
	cfg       *config.Config
	appCfg    appconfig.IConfig
	js        *events.JSWrapper
	cronState *croner.Registry
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	PS        perms.Permissions
	Aud       audit.IAuditer
	Enricher  *mstlystcdata.Enricher
	Laws      *mstlystcdata.Laws
	Storage   storage.IStorage
	Config    *config.Config
	AppConfig appconfig.IConfig
	JS        *events.JSWrapper
	CronState *croner.Registry
}

func NewServer(p Params) *Server {
	return &Server{
		logger:    p.Logger,
		db:        p.DB,
		ps:        p.PS,
		aud:       p.Aud,
		enricher:  p.Enricher,
		laws:      p.Laws,
		st:        p.Storage,
		cfg:       p.Config,
		appCfg:    p.AppConfig,
		js:        p.JS,
		cronState: p.CronState,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbsettings.RegisterSettingsServiceServer(srv, s)
	pbsettings.RegisterConfigServiceServer(srv, s)
	pbsettings.RegisterCronServiceServer(srv, s)
	pbsettings.RegisterFilestoreServiceServer(srv, s)
	pbsettings.RegisterLawsServiceServer(srv, s)
	pbsettings.RegisterAccountsServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbsettings.PermsRemap
}
