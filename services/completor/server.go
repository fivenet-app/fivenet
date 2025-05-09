package completor

import (
	"database/sql"

	pbcompletor "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/completor"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbcompletor.CompletorServiceServer

	db       *sql.DB
	p        perms.Permissions
	jobsS    *mstlystcdata.JobsSearch
	laws     *mstlystcdata.Laws
	tracker  tracker.ITracker
	enricher *mstlystcdata.UserAwareEnricher

	customDB config.CustomDB
}

type Params struct {
	fx.In

	DB         *sql.DB
	Perms      perms.Permissions
	JobsSearch *mstlystcdata.JobsSearch
	Laws       *mstlystcdata.Laws
	Tracker    tracker.ITracker
	Enricher   *mstlystcdata.UserAwareEnricher
	Config     *config.Config
}

func NewServer(p Params) *Server {
	s := &Server{
		db:       p.DB,
		p:        p.Perms,
		jobsS:    p.JobsSearch,
		laws:     p.Laws,
		tracker:  p.Tracker,
		enricher: p.Enricher,

		customDB: p.Config.Database.Custom,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcompletor.RegisterCompletorServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbcompletor.PermsRemap
}
