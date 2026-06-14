package documents

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListCategories(
	ctx context.Context,
	req *pbdocuments.ListCategoriesRequest,
) (*pbdocuments.ListCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	categories, err := s.store.ListCategories(ctx, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	resp := &pbdocuments.ListCategoriesResponse{Categories: categories}

	return resp, nil
}

func (s *Server) getCategory(ctx context.Context, id int64) (*documentscategory.Category, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	return s.store.GetCategory(ctx, id, userInfo)
}

func (s *Server) CreateOrUpdateCategory(
	ctx context.Context,
	req *pbdocuments.CreateOrUpdateCategoryRequest,
) (*pbdocuments.CreateOrUpdateCategoryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.GetCategory().GetId() == 0 {
		lastId, err := s.store.CreateCategory(ctx, s.db, req.GetCategory(), userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		req.Category.Id = lastId

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		if err := s.store.UpdateCategory(ctx, s.db, req.GetCategory(), userInfo); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	category, err := s.getCategory(ctx, req.GetCategory().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.CreateOrUpdateCategoryResponse{
		Category: category,
	}, nil
}

func (s *Server) DeleteCategory(
	ctx context.Context,
	req *pbdocuments.DeleteCategoryRequest,
) (*pbdocuments.DeleteCategoryResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.category_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	category, err := s.getCategory(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if category.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteCategory(ctx, s.db, req.GetId(), userInfo, deletedAtTime); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteCategoryResponse{}, nil
}
