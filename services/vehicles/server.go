package vehicles

import (
	"database/sql"

	pbvehicles "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbvehicles.VehiclesServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.Enricher
	customDB config.CustomDB
}

type Params struct {
	fx.In

	DB       *sql.DB
	Ps       perms.Permissions
	Enricher *mstlystcdata.Enricher
	Config   *config.Config
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		ps:       p.Ps,
		enricher: p.Enricher,
		customDB: p.Config.Database.Custom,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbvehicles.RegisterVehiclesServiceServer(srv, s)
}
