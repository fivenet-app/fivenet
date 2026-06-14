package vehicles

import (
	"context"
	"database/sql"

	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type vehicleStore interface {
	Count(ctx context.Context, q vehiclesstore.ListQuery) (int64, error)
	List(ctx context.Context, q vehiclesstore.ListQuery) ([]*resourcesvehicles.Vehicle, error)
	UpdateProps(
		ctx context.Context,
		in *vehiclesprops.VehicleProps,
	) (*vehiclesprops.VehicleProps, error)
}

type Server struct {
	pbvehicles.VehiclesServiceServer

	ps       perms.Permissions
	enricher mstlystcdata.IEnricher
	store    vehicleStore
}

type Params struct {
	fx.In

	DB       *sql.DB
	Ps       perms.Permissions
	Enricher mstlystcdata.IEnricher
	Config   *config.Config
	Store    vehicleStore
}

func NewServer(p Params) *Server {
	store := p.Store
	if store == nil {
		store = vehiclesstore.New(p.DB, &p.Config.Database.Custom)
	}

	return &Server{
		ps:       p.Ps,
		enricher: p.Enricher,
		store:    store,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbvehicles.RegisterVehiclesServiceServer(srv, s)
}
