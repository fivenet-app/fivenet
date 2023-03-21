package docstore

import (
	context "context"
	"encoding/json"
	"html/template"
	"os"

	"github.com/Masterminds/sprig/v3"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	dTemplates = table.ArpanetDocumentsTemplates.AS("documenttemplateshort")
)

func (s *Server) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	stmt := dTemplates.
		SELECT(
			dTemplates.ID,
			dTemplates.Job,
			dTemplates.JobGrade,
			dCategory.ID,
			dCategory.Name,
			dCategory.Description,
			dCategory.Job,
			dTemplates.Title,
			dTemplates.Description,
			dTemplates.CreatorID,
		).
		FROM(
			dTemplates.
				LEFT_JOIN(dCategory,
					dCategory.ID.EQ(dTemplates.CategoryID),
				),
		).
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

	dTemplates := dTemplates.AS("documenttemplate")
	stmt := dTemplates.
		SELECT(
			dTemplates.ID,
			dTemplates.CreatedAt,
			dTemplates.UpdatedAt,
			dTemplates.Job,
			dTemplates.JobGrade,
			dCategory.ID,
			dCategory.Name,
			dCategory.Description,
			dCategory.Job,
			dTemplates.Title,
			dTemplates.Description,
			dTemplates.ContentTitle,
			dTemplates.Content,
			dTemplates.AdditionalData,
			dTemplates.CreatorID,
		).
		FROM(
			dTemplates.
				LEFT_JOIN(dCategory,
					dCategory.ID.EQ(dTemplates.CategoryID),
				),
		).
		WHERE(
			jet.AND(
				dTemplates.ID.EQ(jet.Uint64(req.TemplateId)),
				dTemplates.Job.EQ(jet.String(job)),
				dTemplates.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	resp := &GetTemplateResponse{}
	if err := stmt.QueryContext(ctx, s.db, resp); err != nil {
		return nil, err
	}

	if req.Process {
		// Parse data as json for the templating process
		var data map[string]interface{}
		err := json.Unmarshal([]byte(req.Data), &data)
		if err != nil {
			return nil, err
		}

		resp.Template.Content, resp.Template.ContentTitle, err = s.renderDocumentTemplate(resp.Template, data)
		if err != nil {
			return nil, err
		}

		resp.Processed = true
	}

	return resp, nil
}

func (s *Server) renderDocumentTemplate(docTmpl *documents.DocumentTemplate, data map[string]interface{}) (out string, outTile string, err error) {
	tmplTitle, err := template.
		New("title").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.ContentTitle)
	if err != nil {
		return
	}
	err = tmplTitle.Execute(os.Stdout, data)
	if err != nil {
		return
	}

	tmpl, err := template.
		New("content").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.Content)
	if err != nil {
		return
	}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return
	}

	return
}
