package dmv

import (
	"database/sql"

	pbdmv "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/dmv"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbdmv.DMVServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.Enricher
	aud      audit.IAuditer
	customDB config.CustomDB
}

type Params struct {
	fx.In

	DB       *sql.DB
	Ps       perms.Permissions
	Enricher *mstlystcdata.Enricher
	Aud      audit.IAuditer
	Config   *config.Config
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		ps:       p.Ps,
		enricher: p.Enricher,
		aud:      p.Aud,
		customDB: p.Config.Database.Custom,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbdmv.RegisterDMVServiceServer(srv, s)
}
