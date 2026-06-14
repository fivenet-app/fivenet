package jobs

import (
	"context"

	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
)

func (s *Server) GetMOTD(
	ctx context.Context,
	req *pbjobs.GetMOTDRequest,
) (*pbjobs.GetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	motd, err := s.store.GetMOTD(ctx, s.db, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.GetMOTDResponse{Motd: motd}, nil
}

func (s *Server) SetMOTD(
	ctx context.Context,
	req *pbjobs.SetMOTDRequest,
) (*pbjobs.SetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	if err := s.store.SetMOTD(ctx, s.db, userInfo.GetJob(), req.GetMotd()); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.SetMOTDResponse{
		Motd: req.GetMotd(),
	}, nil
}
