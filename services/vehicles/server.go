package vehicles

import (
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbvehicles.VehiclesServiceServer

	ps       perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	store    vehiclesstore.IStore
}

type Params struct {
	fx.In

	Ps       perms.Permissions
	Enricher mstlystcdata.IUserAwareEnricher
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
