package citizens

import (
	"database/sql"

	pbcitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbcitizens.CitizensServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	st       storage.IStorage
	appCfg   appconfig.IConfig
	cfg      *config.Config
	customDB config.CustomDB

	avatarHandler  *filestore.Handler[int32]
	mugshotHandler *filestore.Handler[int32]
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
	tUserProps := table.FivenetUserProps

	avatarHandler := filestore.NewHandler(p.Storage, p.DB, tUserProps, tUserProps.UserID, tUserProps.AvatarFileID, 3<<20, func(parentId int32) jet.BoolExpression {
		return tUserProps.UserID.EQ(jet.Int32(parentId))
	}, filestore.UpdateJoinRow, true)
	mugshotHandler := filestore.NewHandler(p.Storage, p.DB, tUserProps, tUserProps.UserID, tUserProps.MugshotFileID, 3<<20, func(parentId int32) jet.BoolExpression {
		return tUserProps.UserID.EQ(jet.Int32(parentId))
	}, filestore.UpdateJoinRow, true)

	return &Server{
		db:       p.DB,
		ps:       p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		st:       p.Storage,
		appCfg:   p.AppConfig,
		cfg:      p.Config,
		customDB: p.Config.Database.Custom,

		avatarHandler:  avatarHandler,
		mugshotHandler: mugshotHandler,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcitizens.RegisterCitizensServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbcitizens.PermsRemap
}
