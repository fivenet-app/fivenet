package citizens

import (
	context "context"

	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListUserActivity(
	ctx context.Context,
	req *pbcitizens.ListUserActivityRequest,
) (*pbcitizens.ListUserActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{citizenIDLogFieldKey, req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcitizens.ListUserActivityResponse{
		Activity: []*usersactivity.UserActivity{},
	}

	// User can't see their own activities, unless they have "Own" perm attribute, or are a superuser
	fields, err := permscitizens.CitizensService.ListUserActivity.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if userInfo.GetUserId() == req.GetUserId() {
		// If isn't superuser or doesn't have 'Own' activity feed access
		if !userInfo.GetJobAdmin() &&
			!fields.Contains(permscitizens.CitizensServiceListUserActivityFieldsPermValueOwn) {
			return resp, nil
		}
	}
	queryOpts := citizensstore.CountUserActivityOptions{
		UserActivityOptions: citizensstore.UserActivityOptions{
			UserID: req.GetUserId(),
			Types:  req.GetTypes(),
		},
	}
	count, err := s.store.CountUserActivity(ctx, queryOpts)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, 20)
	resp.Pagination = pag
	if count <= 0 {
		return resp, nil
	}

	activities, err := s.store.ListUserActivity(ctx, citizensstore.ListUserActivityOptions{
		UserActivityOptions: citizensstore.UserActivityOptions{
			UserID: req.GetUserId(),
			Types:  req.GetTypes(),
		},
		Sort:   req.GetSort(),
		Offset: req.GetPagination().GetOffset(),
		Limit:  limit,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	resp.Activity = activities

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetActivity() {
		if resp.GetActivity()[i].GetSourceUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetSourceUser())
		}
		if resp.GetActivity()[i].GetTargetUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetTargetUser())
		}
	}

	return resp, nil
}
