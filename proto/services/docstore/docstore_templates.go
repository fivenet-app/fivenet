package docstore

import (
	context "context"

	"github.com/galexrt/arpanet/pkg/auth"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	stmt := adt.SELECT(
		adt.ID,
		adt.Job,
		adt.JobGrade,
		adt.Title,
		adt.Description,
		adt.CreatorID,
	).
		FROM(adt).
		WHERE(
			jet.AND(
				adt.Job.EQ(jet.String(job)),
				adt.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	resp := &ListTemplatesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Templates); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error) {
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	stmt := adt.SELECT(
		adt.AllColumns,
	).
		FROM(adt).
		WHERE(
			jet.AND(
				adt.ID.EQ(jet.Uint64(req.TemplateId)),
				adt.Job.EQ(jet.String(job)),
				adt.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	resp := &GetTemplateResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Template); err != nil {
		return nil, err
	}

	return resp, nil
}
