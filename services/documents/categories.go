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
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tDCategory = table.FivenetDocumentsCategories.AS("category")

func (s *Server) ListCategories(
	ctx context.Context,
	req *pbdocuments.ListCategoriesRequest,
) (*pbdocuments.ListCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tDCategory.Job.EQ(jet.String(userInfo.GetJob()))
	if !userInfo.GetSuperuser() {
		condition = jet.AND(
			tDCategory.DeletedAt.IS_NULL(),
			tDCategory.Job.EQ(jet.String(userInfo.GetJob())),
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
			tDCategory.Job.EQ(jet.String(userInfo.GetJob())),
			tDCategory.ID.EQ(jet.Uint64(id)),
		)).
		LIMIT(1)

	var dest documents.Category
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Server) CreateOrUpdateCategory(
	ctx context.Context,
	req *pbdocuments.CreateOrUpdateCategoryRequest,
) (*pbdocuments.CreateOrUpdateCategoryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateCategory",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	tDCategory := table.FivenetDocumentsCategories

	if req.GetCategory().GetId() == 0 {
		stmt := tDCategory.
			INSERT(
				tDCategory.Name,
				tDCategory.Description,
				tDCategory.Job,
				tDCategory.Color,
				tDCategory.Icon,
			).
			VALUES(
				req.GetCategory().GetName(),
				req.GetCategory().Description,
				userInfo.GetJob(),
				req.GetCategory().Color,
				req.GetCategory().Icon,
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
				req.GetCategory().GetName(),
				req.GetCategory().Description,
				userInfo.GetJob(),
				req.GetCategory().Color,
				req.GetCategory().Icon,
			).
			WHERE(jet.AND(
				tDCategory.ID.EQ(jet.Uint64(req.GetCategory().GetId())),
				tDCategory.Job.EQ(jet.String(userInfo.GetJob())),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
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

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "DeleteCategory",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	category, err := s.getCategory(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if category.GetDeletedAt() != nil && userInfo.GetSuperuser() {
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
			tDCategory.Job.EQ(jet.String(userInfo.GetJob())),
			tDCategory.ID.EQ(jet.Uint64(req.GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbdocuments.DeleteCategoryResponse{}, nil
}
