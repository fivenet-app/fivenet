package completor

import (
	"context"

	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	completorstore "github.com/fivenet-app/fivenet/v2026/stores/completor"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbcompletor.CompletorServiceServer

	auth       *grpcauth.GRPCAuth
	ps         perms.Permissions
	jobsSearch mstlystcdata.IJobsSearch
	laws       mstlystcdata.ILaws
	tracker    tracker.ITracker
	enricher   mstlystcdata.IUserAwareEnricher
	store      completorstore.IStore
}

type Params struct {
	fx.In

	Auth       *grpcauth.GRPCAuth
	Perms      perms.Permissions
	JobsSearch mstlystcdata.IJobsSearch
	Laws       mstlystcdata.ILaws
	Tracker    tracker.ITracker
	Enricher   mstlystcdata.IUserAwareEnricher
	Store      completorstore.IStore
}

func NewServer(p Params) *Server {
	s := &Server{
		auth:       p.Auth,
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

// AuthFuncOverride lets CompleteJobs work with account-token-only access.
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	switch fullMethod {
	case pbcompletor.CompletorService_CompleteJobs_FullMethodName:
		if hasUserTokenInContext(ctx) {
			return s.auth.GRPCAuthFunc(ctx, fullMethod)
		}
		return s.auth.GRPCAuthFuncWithoutUserInfo(ctx, fullMethod)

	default:
		return s.auth.GRPCAuthFunc(ctx, fullMethod)
	}
}

func hasUserTokenInContext(ctx context.Context) bool {
	token, err := grpcauth.GetUserTokenFromGRPCContext(ctx)
	return err == nil && token != ""
}
