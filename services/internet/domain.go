package internet

import (
	"context"

	pbinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) CheckDomainAvailability(ctx context.Context, req *pbinternet.CheckDomainAvailabilityRequest) (*pbinternet.CheckDomainAvailabilityResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) ListDomains(ctx context.Context, req *pbinternet.ListDomainsRequest) (*pbinternet.ListDomainsResponse, error) {
	condition := jet.Bool(true)

	stmt := tDomains.
		SELECT(
			tDomains.ID,
		).
		FROM(tDomains).
		WHERE(condition)

	_ = stmt
	// TODO

	return nil, nil
}

func (s *Server) RegisterDomain(ctx context.Context, req *pbinternet.RegisterDomainRequest) (*pbinternet.RegisterDomainResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) UpdateDomain(ctx context.Context, req *pbinternet.UpdateDomainRequest) (*pbinternet.UpdateDomainResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) TransferDomain(ctx context.Context, req *pbinternet.TransferDomainRequest) (*pbinternet.TransferDomainResponse, error) {
	// TODO

	return nil, nil
}
