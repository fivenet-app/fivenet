package rector

import (
	"database/sql"

	pbrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbrector.RectorServiceServer
	pbrector.RectorConfigServiceServer
	pbrector.RectorCronServiceServer
	pbrector.RectorFilestoreServiceServer
	pbrector.RectorLawsServiceServer

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
	cronState *croner.State
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
	CronState *croner.State
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
	pbrector.RegisterRectorServiceServer(srv, s)
	pbrector.RegisterRectorConfigServiceServer(srv, s)
	pbrector.RegisterRectorCronServiceServer(srv, s)
	pbrector.RegisterRectorFilestoreServiceServer(srv, s)
	pbrector.RegisterRectorLawsServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbrector.PermsRemap
}
