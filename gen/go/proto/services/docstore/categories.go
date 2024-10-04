package docstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tDCategory = table.FivenetDocumentsCategories.AS("category")

func (s *Server) ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
		).
		FROM(
			tDCategory,
		).
		WHERE(
			tDCategory.Job.EQ(jet.String(userInfo.Job)),
		).
		ORDER_BY(
			tDCategory.Name.ASC(),
		)

	resp := &ListCategoriesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Category); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) getCategory(ctx context.Context, id uint64) (*documents.Category, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
		).
		FROM(
			tDCategory,
		).
		WHERE(
			jet.AND(
				tDCategory.Job.EQ(jet.String(userInfo.Job)),
				tDCategory.ID.EQ(jet.Uint64(id)),
			),
		).
		LIMIT(1)

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
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	tDCategory := table.FivenetDocumentsCategories
	stmt := tDCategory.
		INSERT(
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
		).
		VALUES(
			req.Category.Name,
			req.Category.Description,
			userInfo.Job,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &CreateCategoryResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.category_id", int64(req.Category.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "UpdateCategory",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	tDCategory := table.FivenetDocumentsCategories
	stmt := tDCategory.
		UPDATE(
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
		).
		SET(
			req.Category.Name,
			req.Category.Description,
			userInfo.Job,
			req.Category.Color,
			req.Category.Icon,
		).
		WHERE(jet.AND(
			tDCategory.ID.EQ(jet.Uint64(req.Category.Id)),
			tDCategory.Job.EQ(jet.String(userInfo.Job)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateCategoryResponse{}, nil
}

func (s *Server) DeleteCategory(ctx context.Context, req *DeleteCategoryRequest) (*DeleteCategoryResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64Slice("fivenet.docstore.category_ids", utils.SliceUint64ToInt64(req.Ids)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteCategory",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	ids := make([]jet.Expression, len(req.Ids))
	for i := 0; i < len(req.Ids); i++ {
		ids[i] = jet.Uint64(req.Ids[i])
	}

	tDCategory := table.FivenetDocumentsCategories
	stmt := tDCategory.
		DELETE().
		WHERE(
			jet.AND(
				tDCategory.Job.EQ(jet.String(userInfo.Job)),
				tDCategory.ID.IN(ids...),
			),
		).
		LIMIT(int64(len(req.Ids)))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteCategoryResponse{}, nil
}
