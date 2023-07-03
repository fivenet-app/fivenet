package docstore

import (
	"context"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tDCategory = table.FivenetDocumentsCategories.AS("category")
)

func (s *Server) ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tDCategory := table.FivenetDocumentsCategories.AS("category")
	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
		).
		FROM(
			tDCategory,
		).
		WHERE(
			tDCategory.Job.EQ(jet.String(userInfo.Job)),
		).
		ORDER_BY(
			tDCategory.Name,
		)

	resp := &ListCategoriesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Category); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) getCategory(ctx context.Context, id uint64) (*documents.Category, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tDCategory := table.FivenetDocumentsCategories.AS("category")
	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
		).
		FROM(
			tDCategory,
		).
		WHERE(
			jet.AND(
				tDCategory.Job.EQ(jet.String(userInfo.Job)),
				tDCategory.ID.EQ(jet.Uint64(id)),
			),
		)

	var dest documents.Category
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CreateCategoryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "CreateCategory",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	dCategory := table.FivenetDocumentsCategories
	stmt := dCategory.
		INSERT(
			dCategory.Name,
			dCategory.Description,
			dCategory.Job,
		).
		VALUES(
			req.Category.Name,
			req.Category.Description,
			userInfo.Job,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return &CreateCategoryResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "UpdateCategory",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	dCategory := table.FivenetDocumentsCategories
	stmt := dCategory.
		UPDATE(
			dCategory.Name,
			dCategory.Description,
			dCategory.Job,
		).
		SET(
			req.Category.Name,
			req.Category.Description,
			userInfo.Job,
		).
		WHERE(jet.AND(
			dCategory.ID.EQ(jet.Uint64(req.Category.Id)),
			dCategory.Job.EQ(jet.String(userInfo.Job)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &UpdateCategoryResponse{}, nil
}

func (s *Server) DeleteCategory(ctx context.Context, req *DeleteCategoryRequest) (*DeleteCategoryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteCategory",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	ids := make([]jet.Expression, len(req.Ids))
	for i := 0; i < len(req.Ids); i++ {
		ids[i] = jet.Uint64(req.Ids[i])
	}

	dCategory := table.FivenetDocumentsCategories
	stmt := dCategory.
		DELETE().
		WHERE(
			jet.AND(
				dCategory.Job.EQ(jet.String(userInfo.Job)),
				dCategory.ID.IN(ids...),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &DeleteCategoryResponse{}, nil
}
