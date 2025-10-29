package filestore

import (
	"database/sql"

	pbfilestore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pbfilestore.FilestoreServiceServer

	logger   *zap.Logger
	db       *sql.DB
	st       storage.IStorage
	fHandler *filestore.Handler[int64]
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	PS        perms.Permissions
	Enricher  *mstlystcdata.Enricher
	Laws      *mstlystcdata.Laws
	Storage   storage.IStorage
	Config    *config.Config
	AppConfig appconfig.IConfig
	CronState *croner.Registry
}

func NewServer(p Params) *Server {
	fHandler := filestore.NewHandler[int64](
		p.Storage,
		p.DB,
		nil,
		nil,
		nil,
		-1,
		nil,
		filestore.InsertJoinRow,
		false,
	)

	return &Server{
		logger:   p.Logger,
		db:       p.DB,
		st:       p.Storage,
		fHandler: fHandler,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbfilestore.RegisterFilestoreServiceServer(srv, s)
}
