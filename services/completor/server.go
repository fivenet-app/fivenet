package completor

import (
	"database/sql"

	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	completorstore "github.com/fivenet-app/fivenet/v2026/stores/completor"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbcompletor.CompletorServiceServer

	ps         perms.Permissions
	jobsSearch mstlystcdata.IJobsSearch
	laws       mstlystcdata.ILaws
	tracker    tracker.ITracker
	enricher   mstlystcdata.IUserAwareEnricher
	store      completorstore.IStore
}

type Params struct {
	fx.In

	DB         *sql.DB
	Perms      perms.Permissions
	JobsSearch mstlystcdata.IJobsSearch
	Laws       mstlystcdata.ILaws
	Tracker    tracker.ITracker
	Enricher   mstlystcdata.IUserAwareEnricher
	Config     *config.Config
	Store      completorstore.IStore `optional:"true"`
}

func NewServer(p Params) *Server {
	s := &Server{
		ps:         p.Perms,
		jobsSearch: p.JobsSearch,
		laws:       p.Laws,
		tracker:    p.Tracker,
		enricher:   p.Enricher,
		store:      p.Store,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcompletor.RegisterCompletorServiceServer(srv, s)
}
