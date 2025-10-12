package documents

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tDCategory = table.FivenetDocumentsCategories.AS("category")

func (s *Server) ListCategories(
	ctx context.Context,
	req *pbdocuments.ListCategoriesRequest,
) (*pbdocuments.ListCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tDCategory.Job.EQ(mysql.String(userInfo.GetJob()))
	if !userInfo.GetSuperuser() {
		condition = mysql.AND(
			tDCategory.DeletedAt.IS_NULL(),
			tDCategory.Job.EQ(mysql.String(userInfo.GetJob())),
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

func (s *Server) getCategory(ctx context.Context, id int64) (*documents.Category, error) {
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
		WHERE(mysql.AND(
			tDCategory.Job.EQ(mysql.String(userInfo.GetJob())),
			tDCategory.ID.EQ(mysql.Int64(id)),
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

		req.Category.Id = lastId

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
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
			WHERE(mysql.AND(
				tDCategory.ID.EQ(mysql.Int64(req.GetCategory().GetId())),
				tDCategory.Job.EQ(mysql.String(userInfo.GetJob())),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
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

	deletedAtTime := mysql.CURRENT_TIMESTAMP()
	if category.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = mysql.TimestampExp(mysql.NULL)
	}

	tDCategory := table.FivenetDocumentsCategories
	stmt := tDCategory.
		UPDATE(
			tDCategory.DeletedAt,
		).
		SET(
			tDCategory.DeletedAt.SET(deletedAtTime),
		).
		WHERE(mysql.AND(
			tDCategory.Job.EQ(mysql.String(userInfo.GetJob())),
			tDCategory.ID.EQ(mysql.Int64(req.GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbdocuments.DeleteCategoryResponse{}, nil
}
