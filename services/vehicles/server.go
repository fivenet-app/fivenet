package vehicles

import (
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbvehicles.VehiclesServiceServer

	perms    perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	hydrator citizenshydrator.IHydrator
	store    vehiclesstore.IStore
}

type Params struct {
	fx.In

	Perms    perms.Permissions
	Enricher mstlystcdata.IUserAwareEnricher
	Hydrator citizenshydrator.IHydrator
	Store    vehiclesstore.IStore
}

func NewServer(p Params) *Server {
	return &Server{
		perms:    p.Perms,
		enricher: p.Enricher,
		hydrator: p.Hydrator,
		store:    p.Store,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbvehicles.RegisterVehiclesServiceServer(srv, s)
}
