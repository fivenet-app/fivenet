package docstore

import (
	"bytes"
	context "context"
	"encoding/json"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	dTemplates = table.FivenetDocumentsTemplates.AS("documenttemplateshort")
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
			dTemplates.Schema,
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
			dTemplates.Schema,
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

	if req.Render != nil && *req.Render {
		// Parse data as json for the templating process
		var data map[string]interface{}
		err := json.Unmarshal([]byte(req.Data), &data)
		if err != nil {
			return nil, err
		}

		resp.Template.ContentTitle, resp.Template.Content, err = s.renderDocumentTemplate(resp.Template, data)
		if err != nil {
			return nil, err
		}

		resp.Rendered = true
	}

	return resp, nil
}

func (s *Server) renderDocumentTemplate(docTmpl *documents.DocumentTemplate, data map[string]interface{}) (outTile string, out string, err error) {
	// Render Title template
	tmplTitle, err := template.
		New("title").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.ContentTitle)
	if err != nil {
		return
	}
	buf := &bytes.Buffer{}
	err = tmplTitle.Execute(buf, data)
	if err != nil {
		return
	}
	outTile = buf.String()

	// Render Content Template
	tmpl, err := template.
		New("content").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.Content)
	if err != nil {
		return
	}

	buf.Reset()
	err = tmpl.Execute(buf, data)
	if err != nil {
		return
	}
	out = buf.String()

	return
}

func (s *Server) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "CreateTemplate",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	categoryId := jet.NULL
	if req.Template.Category != nil {
		cat, err := s.getDocumentCategory(ctx, req.Template.Category.Id)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			categoryId = jet.Uint64(cat.Id)
		}
	}

	dTemplates := table.FivenetDocumentsTemplates
	stmt := dTemplates.
		INSERT(
			dTemplates.Job,
			dTemplates.JobGrade,
			dTemplates.CategoryID,
			dTemplates.Title,
			dTemplates.Description,
			dTemplates.ContentTitle,
			dTemplates.Content,
			dTemplates.Schema,
			dTemplates.CreatorID,
		).
		VALUES(
			job,
			jobGrade,
			categoryId,
			req.Template.Title,
			req.Template.Description,
			req.Template.ContentTitle,
			req.Template.Content,
			req.Template.Schema,
			userId,
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

	return &CreateTemplateResponse{
		Id: lastId,
	}, nil
}

func (s *Server) UpdateTemplate(ctx context.Context, req *UpdateTemplateRequest) (*UpdateTemplateResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "UpdateTemplate",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	categoryId := jet.NULL
	if req.Template.Category != nil {
		cat, err := s.getDocumentCategory(ctx, req.Template.Category.Id)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			categoryId = jet.Uint64(cat.Id)
		}
	}

	dTemplates := table.FivenetDocumentsTemplates
	stmt := dTemplates.
		UPDATE(
			dTemplates.Job,
			dTemplates.JobGrade,
			dTemplates.CategoryID,
			dTemplates.Title,
			dTemplates.Description,
			dTemplates.ContentTitle,
			dTemplates.Content,
			dTemplates.Schema,
		).
		SET(
			job,
			jobGrade,
			categoryId,
			req.Template.Title,
			req.Template.Description,
			req.Template.ContentTitle,
			req.Template.Content,
			req.Template.Schema,
		).
		WHERE(
			jet.AND(
				dTemplates.ID.EQ(jet.Uint64(req.Template.Id)),
				dTemplates.Job.EQ(jet.String(job)),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &UpdateTemplateResponse{}, nil
}

func (s *Server) DeleteTemplate(ctx context.Context, req *DeleteTemplateRequest) (*DeleteTemplateResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteTemplate",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	dTemplates := table.FivenetDocumentsTemplates
	stmt := dTemplates.
		DELETE().
		WHERE(
			jet.AND(
				dTemplates.Job.EQ(jet.String(job)),
				dTemplates.ID.EQ(jet.Uint64(req.Id)),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &DeleteTemplateResponse{}, nil
}
