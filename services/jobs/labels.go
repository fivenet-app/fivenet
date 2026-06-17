package jobs

import (
	context "context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
)

func (s *Server) GetColleagueLabels(
	ctx context.Context,
	req *pbjobs.GetColleagueLabelsRequest,
) (*pbjobs.GetColleagueLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbjobs.GetColleagueLabelsResponse{
		Labels: []*jobslabels.Label{},
	}

	// Fields Permission Check
	fields, err := permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetSuperuser() {
		fields.Set(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels)
	}
	if !fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) {
		// Fallback to checking if user has manage colleague labels permission
		if !s.ps.Can(
			userInfo,
			permsjobs.ColleaguesService.ManageLabels.Perm,
		) {
			return nil, errorsjobs.ErrLabelsNoPerms
		}
	}

	labels, err := s.store.GetColleagueLabels(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetSearch(),
		userInfo.GetSuperuser(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	resp.Labels = labels

	return resp, nil
}

func (s *Server) ManageLabels(
	ctx context.Context,
	req *pbjobs.ManageLabelsRequest,
) (*pbjobs.ManageLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	labels, err := s.store.ManageLabels(ctx, s.db, userInfo.GetJob(), req.GetLabels())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	resp := &pbjobs.ManageLabelsResponse{Labels: labels}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return resp, nil
}

func (s *Server) GetColleagueLabelsStats(
	ctx context.Context,
	req *pbjobs.GetColleagueLabelsStatsRequest,
) (*pbjobs.GetColleagueLabelsStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Types Permission Check
	fields, err := permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetSuperuser() {
		fields.Set(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels)
	}
	if !fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) {
		return &pbjobs.GetColleagueLabelsStatsResponse{}, nil
	}

	dest, err := s.store.GetColleagueLabelsStats(ctx, s.db, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbjobs.GetColleagueLabelsStatsResponse{
		Count: dest,
	}, nil
}
