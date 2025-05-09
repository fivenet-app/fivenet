package citizenstore

import (
	"database/sql"

	pbcitizenstore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizenstore"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbcitizenstore.CitizenStoreServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	st       storage.IStorage
	appCfg   appconfig.IConfig
	cfg      *config.Config

	customDB config.CustomDB
}

type Params struct {
	fx.In

	DB        *sql.DB
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	Aud       audit.IAuditer
	Config    *config.Config
	Storage   storage.IStorage
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		ps:       p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		st:       p.Storage,
		appCfg:   p.AppConfig,
		cfg:      p.Config,

		customDB: p.Config.Database.Custom,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcitizenstore.RegisterCitizenStoreServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbcitizenstore.PermsRemap
}
