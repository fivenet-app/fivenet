package docstore

import (
	context "context"

	"github.com/galexrt/arpanet/pkg/auth"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	stmt := dTemplates.SELECT(
		dTemplates.ID,
		dTemplates.Job,
		dTemplates.JobGrade,
		dTemplates.Title,
		dTemplates.Description,
		dTemplates.CreatorID,
	).
		FROM(dTemplates).
		WHERE(
			jet.AND(
				dTemplates.Job.EQ(jet.String(job)),
				dTemplates.JobGrade.LT_EQ(jet.Int32(jobGrade)),
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

	stmt := dTemplates.SELECT(
		dTemplates.AllColumns,
	).
		FROM(dTemplates).
		WHERE(
			jet.AND(
				dTemplates.ID.EQ(jet.Uint64(req.TemplateId)),
				dTemplates.Job.EQ(jet.String(job)),
				dTemplates.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	resp := &GetTemplateResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Template); err != nil {
		return nil, err
	}

	return resp, nil
}
