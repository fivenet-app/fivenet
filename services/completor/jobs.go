package completor

import (
	context "context"

	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorscompletor "github.com/fivenet-app/fivenet/v2026/services/completor/errors"
)

func (s *Server) CompleteJobs(
	ctx context.Context,
	req *pbcompletor.CompleteJobsRequest,
) (*pbcompletor.CompleteJobsResponse, error) {
	var search string
	if req.Search != nil && req.GetSearch() != "" {
		search = req.GetSearch()
	}
	if req.CurrentJob != nil && req.GetCurrentJob() {
		userInfo, ok := auth.GetUserInfoFromContext(ctx)
		if !ok {
			return nil, errorsgrpcauth.ErrNoUserInfo
		}
		search = userInfo.GetJob()
	}
	exactMatch := false
	if req.ExactMatch != nil {
		exactMatch = req.GetExactMatch()
	}

	resp := &pbcompletor.CompleteJobsResponse{}
	if search != "" {
		var err error
		resp.Jobs, err = s.jobsSearch.Search(ctx, search, exactMatch)
		if err != nil {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	} else {
		resp.Jobs = s.jobsSearch.List()
	}

	return resp, nil
}
