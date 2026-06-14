package completor

import (
	context "context"
	"slices"

	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorscompletor "github.com/fivenet-app/fivenet/v2026/services/completor/errors"
	completorstore "github.com/fivenet-app/fivenet/v2026/stores/completor"
)

func (s *Server) CompleteCitizens(
	ctx context.Context,
	req *pbcompletor.CompleteCitizensRequest,
) (*pbcompletor.CompleteCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	currentJob := req.CurrentJob != nil && req.GetCurrentJob()
	dest, err := s.store.CompleteCitizens(ctx, completorstore.CitizensQuery{
		Search:      req.GetSearch(),
		CurrentJob:  currentJob,
		UserJob:     userInfo.GetJob(),
		UserIDs:     req.GetUserIds(),
		UserIDsOnly: req.GetUserIdsOnly(),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}

	if req.OnDuty != nil && req.GetOnDuty() {
		dest = slices.DeleteFunc(dest, func(us *usershort.UserShort) bool {
			return !s.tracker.IsUserOnDuty(us.GetUserId())
		})
	}

	if currentJob {
		jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
		for i := range dest {
			jobInfoFn(dest[i])
		}
	}

	return &pbcompletor.CompleteCitizensResponse{
		Users: dest,
	}, nil
}
