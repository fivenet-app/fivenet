package vehicles

import (
	"database/sql"

	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbvehicles.VehiclesServiceServer

	ps       perms.Permissions
	enricher mstlystcdata.IEnricher
	store    vehiclesstore.IStore
}

type Params struct {
	fx.In

	DB       *sql.DB
	Ps       perms.Permissions
	Enricher mstlystcdata.IEnricher
	Config   *config.Config
	Store    vehiclesstore.IStore
}

func NewServer(p Params) *Server {
	return &Server{
		ps:       p.Ps,
		enricher: p.Enricher,
		store:    p.Store,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbvehicles.RegisterVehiclesServiceServer(srv, s)
}
