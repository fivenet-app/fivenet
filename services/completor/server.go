package completor

import (
	"context"
	"database/sql"

	documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	completorstore "github.com/fivenet-app/fivenet/v2026/stores/completor"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type completorStore interface {
	CompleteCitizens(
		ctx context.Context,
		q completorstore.CitizensQuery,
	) ([]*usershort.UserShort, error)
	CompleteDocumentCategories(
		ctx context.Context,
		q completorstore.DocumentCategoriesQuery,
	) ([]*documentscategory.Category, error)
}

type Server struct {
	pbcompletor.CompletorServiceServer

	ps         perms.Permissions
	jobsSearch *mstlystcdata.JobsSearch
	laws       *mstlystcdata.Laws
	tracker    tracker.ITracker
	enricher   mstlystcdata.IUserAwareEnricher
	store      completorStore
}

type Params struct {
	fx.In

	DB         *sql.DB
	Perms      perms.Permissions
	JobsSearch *mstlystcdata.JobsSearch
	Laws       *mstlystcdata.Laws
	Tracker    tracker.ITracker
	Enricher   mstlystcdata.IUserAwareEnricher
	Config     *config.Config
	Store      completorStore `optional:"true"`
}

func NewServer(p Params) *Server {
	store := p.Store
	if store == nil {
		store = completorstore.New(p.DB, &p.Config.Database.Custom)
	}

	s := &Server{
		ps:         p.Perms,
		jobsSearch: p.JobsSearch,
		laws:       p.Laws,
		tracker:    p.Tracker,
		enricher:   p.Enricher,
		store:      store,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcompletor.RegisterCompletorServiceServer(srv, s)
}
