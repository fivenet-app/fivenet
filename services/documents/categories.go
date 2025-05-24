package documents

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tDCategory = table.FivenetDocumentsCategories.AS("category")

func (s *Server) ListCategories(ctx context.Context, req *pbdocuments.ListCategoriesRequest) (*pbdocuments.ListCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tDCategory.Job.EQ(jet.String(userInfo.Job))
	if !userInfo.Superuser {
		condition = jet.AND(
			tDCategory.DeletedAt.IS_NULL(),
			tDCategory.Job.EQ(jet.String(userInfo.Job)),
		)
	}

	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.DeletedAt,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
		).
		FROM(
			tDCategory,
		).
		WHERE(condition).
		ORDER_BY(
			tDCategory.SortKey.ASC(),
		)

	resp := &pbdocuments.ListCategoriesResponse{
		Categories: []*documents.Category{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Categories); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) getCategory(ctx context.Context, id uint64) (*documents.Category, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.DeletedAt,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
		).
		FROM(
			tDCategory,
		).
		WHERE(jet.AND(
			tDCategory.Job.EQ(jet.String(userInfo.Job)),
			tDCategory.ID.EQ(jet.Uint64(id)),
		)).
		LIMIT(1)

	var dest documents.Category
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Server) CreateOrUpdateCategory(ctx context.Context, req *pbdocuments.CreateOrUpdateCategoryRequest) (*pbdocuments.CreateOrUpdateCategoryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateCategory",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	tDCategory := table.FivenetDocumentsCategories

	if req.Category.Id == 0 {
		stmt := tDCategory.
			INSERT(
				tDCategory.Name,
				tDCategory.Description,
				tDCategory.Job,
				tDCategory.Color,
				tDCategory.Icon,
			).
			VALUES(
				req.Category.Name,
				req.Category.Description,
				userInfo.Job,
				req.Category.Color,
				req.Category.Icon,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		req.Category.Id = uint64(lastId)

		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	} else {
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
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
	}

	category, err := s.getCategory(ctx, req.Category.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.CreateOrUpdateCategoryResponse{
		Category: category,
	}, nil
}

func (s *Server) DeleteCategory(ctx context.Context, req *pbdocuments.DeleteCategoryRequest) (*pbdocuments.DeleteCategoryResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.category_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "DeleteCategory",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	category, err := s.getCategory(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if category.DeletedAt != nil && userInfo.Superuser {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	tDCategory := table.FivenetDocumentsCategories
	stmt := tDCategory.
		UPDATE(
			tDCategory.DeletedAt,
		).
		SET(
			tDCategory.DeletedAt.SET(deletedAtTime),
		).
		WHERE(jet.AND(
			tDCategory.Job.EQ(jet.String(userInfo.Job)),
			tDCategory.ID.EQ(jet.Uint64(req.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbdocuments.DeleteCategoryResponse{}, nil
}
