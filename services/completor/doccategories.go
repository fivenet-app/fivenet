package completor

import (
	context "context"

	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorscompletor "github.com/fivenet-app/fivenet/v2026/services/completor/errors"
	completorstore "github.com/fivenet-app/fivenet/v2026/stores/completor"
)

func (s *Server) CompleteDocumentCategories(
	ctx context.Context,
	req *pbcompletor.CompleteDocumentCategoriesRequest,
) (*pbcompletor.CompleteDocumentCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobs, err := s.ps.AttrJobList(
		userInfo,
		permsdocuments.CategoriesService.ListCategories.Jobs,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}
	// Ensure user's job is always included
	if !jobs.Contains(userInfo.GetJob()) {
		jobs.Strings = append(jobs.Strings, userInfo.GetJob())
	}
	categories, err := s.store.CompleteDocumentCategories(
		ctx,
		completorstore.DocumentCategoriesQuery{
			Search:      req.GetSearch(),
			CategoryIDs: req.GetCategoryIds(),
			Jobs:        jobs.GetStrings(),
			CurrentJob:  userInfo.GetJob(),
		},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}

	return &pbcompletor.CompleteDocumentCategoriesResponse{Categories: categories}, nil
}
