package rector

import (
	"database/sql"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

type Server struct {
	RectorServiceServer
	RectorConfigServiceServer

	logger   *zap.Logger
	db       *sql.DB
	ps       perms.Permissions
	aud      audit.IAuditer
	enricher *mstlystcdata.Enricher
	cache    *mstlystcdata.Cache
	st       storage.IStorage
	cfg      *config.BaseConfig
}

type Params struct {
	fx.In

	Logger   *zap.Logger
	DB       *sql.DB
	PS       perms.Permissions
	Aud      audit.IAuditer
	Enricher *mstlystcdata.Enricher
	Cache    *mstlystcdata.Cache
	Storage  storage.IStorage
	Config   *config.BaseConfig
}

func NewServer(p Params) *Server {
	return &Server{
		logger:   p.Logger,
		db:       p.DB,
		ps:       p.PS,
		aud:      p.Aud,
		enricher: p.Enricher,
		cache:    p.Cache,
		st:       p.Storage,
		cfg:      p.Config,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterRectorServiceServer(srv, s)
	RegisterRectorConfigServiceServer(srv, s)
}
