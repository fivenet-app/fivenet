package docstore

import (
	"context"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	dCategory = table.FivenetDocumentsCategories.AS("category")
)

func (s *Server) CreateDocumentCategory(ctx context.Context, req *CreateDocumentCategoryRequest) (*CreateDocumentCategoryResponse, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	stmt := dCategory.
		INSERT(
			dCategory.Name,
			dCategory.Description,
			dCategory.Job,
		).
		VALUES(
			req.Category.Name,
			req.Category.Description,
			job,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &CreateDocumentCategoryResponse{
		Id: lastId,
	}, nil
}

func (s *Server) UpdateDocumentCategory(ctx context.Context, req *UpdateDocumentCategoryRequest) (*UpdateDocumentCategoryResponse, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	stmt := dCategory.
		UPDATE(
			dCategory.Name,
			dCategory.Description,
		).
		SET(
			req.Category.Name,
			req.Category.Description,
			job,
		).
		WHERE(
			dCategory.ID.EQ(jet.Uint64(req.Category.Id)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &UpdateDocumentCategoryResponse{}, nil
}

func (s *Server) DeleteDocumentCategory(ctx context.Context, req *DeleteDocumentCategoryRequest) (*DeleteDocumentCategoryResponse, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	ids := make([]jet.Expression, len(req.Ids))
	for i := 0; i < len(req.Ids); i++ {
		ids[i] = jet.Uint64(req.Ids[i])
	}

	stmt := dCategory.
		DELETE().
		WHERE(
			dCategory.ID.IN(ids...).
				AND(dCategory.Job.EQ(jet.String(job))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &DeleteDocumentCategoryResponse{}, nil
}
