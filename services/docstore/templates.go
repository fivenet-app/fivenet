package docstore

import (
	"bytes"
	context "context"
	"errors"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pbdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore"
	permsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/perms"
	"github.com/fivenet-app/fivenet/pkg/dbutils"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorsdocstore "github.com/fivenet-app/fivenet/services/docstore/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tDTemplates       = table.FivenetDocumentsTemplates.AS("templateshort")
	tDTemplatesAccess = table.FivenetDocumentsTemplatesAccess.AS("templatejobaccess")
)

func (s *Server) ListTemplates(ctx context.Context, req *pbdocstore.ListTemplatesRequest) (*pbdocstore.ListTemplatesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tDTemplates.
		SELECT(
			tDTemplates.ID,
			tDTemplates.Weight,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.Schema,
			tDTemplates.Workflow,
			tDTemplates.CreatorJob,
		).
		FROM(
			tDTemplates.
				LEFT_JOIN(tDTemplatesAccess,
					tDTemplatesAccess.TargetID.EQ(tDTemplates.ID).
						AND(tDTemplatesAccess.Access.GT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)))),
				).
				LEFT_JOIN(tDCategory,
					tDCategory.ID.EQ(tDTemplates.CategoryID),
				),
		).
		WHERE(jet.AND(
			tDTemplates.DeletedAt.IS_NULL(),
			jet.AND(
				tDTemplatesAccess.Job.EQ(jet.String(userInfo.Job)),
				tDTemplatesAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade)),
			),
		)).
		ORDER_BY(
			tDTemplates.Weight.DESC(),
			tDTemplates.ID.ASC(),
		).
		GROUP_BY(tDTemplates.ID)

	resp := &pbdocstore.ListTemplatesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Templates); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	for i := range resp.Templates {
		s.enricher.EnrichJobName(resp.Templates[i])
	}

	return resp, nil
}

func (s *Server) GetTemplate(ctx context.Context, req *pbdocstore.GetTemplateRequest) (*pbdocstore.GetTemplateResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.template_id", int64(req.TemplateId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.templateAccess.CanUserAccessTarget(ctx, req.TemplateId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrTemplateNoPerms
	}

	resp := &pbdocstore.GetTemplateResponse{}
	resp.Template, err = s.getTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if resp.Template == nil {
		return nil, errorsdocstore.ErrTemplateNoPerms
	}

	if req.Render == nil || !*req.Render {
		resp.Template.JobAccess, err = s.templateAccess.Jobs.List(ctx, s.db, req.TemplateId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	} else if req.Render != nil && *req.Render && req.Data != nil {
		resp.Template.ContentTitle, resp.Template.State, resp.Template.Content, err = s.renderTemplate(resp.Template, req.Data)
		if err != nil {
			if s.ps.Can(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceCreateTemplatePerm) {
				return nil, err
			} else {
				return nil, errswrap.NewError(err, errorsdocstore.ErrTemplateFailed)
			}
		}

		resp.Rendered = true
	}

	s.enricher.EnrichJobName(resp.Template)

	return resp, nil
}

func (s *Server) getTemplate(ctx context.Context, templateId uint64) (*documents.Template, error) {
	tDTemplates := tDTemplates.AS("template")
	stmt := tDTemplates.
		SELECT(
			tDTemplates.ID,
			tDTemplates.Weight,
			tDTemplates.CreatedAt,
			tDTemplates.UpdatedAt,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.Workflow,
			tDTemplates.CreatorJob,
		).
		FROM(
			tDTemplates.
				LEFT_JOIN(tDCategory,
					tDCategory.ID.EQ(tDTemplates.CategoryID),
				),
		).
		WHERE(jet.AND(
			tDTemplates.DeletedAt.IS_NULL(),
			tDTemplates.ID.EQ(jet.Uint64(templateId)),
		)).
		LIMIT(1)

	var dest documents.Template
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	if dest.Id <= 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Server) renderTemplate(docTmpl *documents.Template, data *documents.TemplateData) (outTile string, outState string, out string, err error) {
	// Render Title template
	titleTpl, err := template.
		New("title").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.ContentTitle)
	if err != nil {
		return
	}
	buf := &bytes.Buffer{}
	err = titleTpl.Execute(buf, data)
	if err != nil {
		return
	}
	outTile = buf.String()

	// Render State template
	stateTpl, err := template.
		New("state").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.State)
	if err != nil {
		return
	}

	buf.Reset()
	err = stateTpl.Execute(buf, data)
	if err != nil {
		return
	}
	outState = buf.String()

	// Render Content template
	contentTpl, err := template.
		New("content").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.Content)
	if err != nil {
		return
	}

	buf.Reset()
	err = contentTpl.Execute(buf, data)
	if err != nil {
		return
	}
	out = buf.String()

	return
}

func (s *Server) CreateTemplate(ctx context.Context, req *pbdocstore.CreateTemplateRequest) (*pbdocstore.CreateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "CreateTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if documents.TemplateAccessHasDuplicates(req.Template.JobAccess) || documents.DocumentAccessHasDuplicates(req.Template.ContentAccess) {
		return nil, errorsdocstore.ErrTemplateAccessDuplicate
	}

	categoryId := jet.NULL
	if req.Template.Category != nil {
		cat, err := s.getCategory(ctx, req.Template.Category.Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
		if cat != nil {
			categoryId = jet.Uint64(cat.Id)
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		INSERT(
			tDTemplates.Weight,
			tDTemplates.CategoryID,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.Workflow,
			tDTemplates.CreatorJob,
		).
		VALUES(
			req.Template.Weight,
			categoryId,
			req.Template.Title,
			req.Template.Description,
			req.Template.Color,
			req.Template.Icon,
			req.Template.ContentTitle,
			req.Template.Content,
			req.Template.State,
			req.Template.ContentAccess,
			req.Template.Schema,
			req.Template.Workflow,
			userInfo.Job,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := s.templateAccess.HandleAccessChanges(ctx, tx, uint64(lastId), req.Template.JobAccess, nil, nil); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &pbdocstore.CreateTemplateResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) UpdateTemplate(ctx context.Context, req *pbdocstore.UpdateTemplateRequest) (*pbdocstore.UpdateTemplateResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.template_id", int64(req.Template.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "UpdateTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.templateAccess.CanUserAccessTarget(ctx, req.Template.Id, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrTemplateNoPerms
	}

	if documents.TemplateAccessHasDuplicates(req.Template.JobAccess) || documents.DocumentAccessHasDuplicates(req.Template.ContentAccess) {
		return nil, errorsdocstore.ErrTemplateAccessDuplicate
	}

	categoryId := jet.NULL
	if req.Template.Category != nil {
		cat, err := s.getCategory(ctx, req.Template.Category.Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
		if cat != nil {
			categoryId = jet.Uint64(cat.Id)
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		UPDATE(
			tDTemplates.Weight,
			tDTemplates.CategoryID,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.Workflow,
		).
		SET(
			req.Template.Weight,
			categoryId,
			req.Template.Title,
			req.Template.Description,
			req.Template.Color,
			req.Template.Icon,
			req.Template.ContentTitle,
			req.Template.Content,
			req.Template.State,
			req.Template.ContentAccess,
			req.Template.Schema,
			req.Template.Workflow,
		).
		WHERE(
			tDTemplates.ID.EQ(jet.Uint64(req.Template.Id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := s.templateAccess.HandleAccessChanges(ctx, tx, req.Template.Id, req.Template.JobAccess, nil, nil); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	tmpl, err := s.getTemplate(ctx, req.Template.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	tmpl.JobAccess, err = s.templateAccess.Jobs.List(ctx, s.db, req.Template.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbdocstore.UpdateTemplateResponse{
		Template: tmpl,
	}, nil
}

func (s *Server) DeleteTemplate(ctx context.Context, req *pbdocstore.DeleteTemplateRequest) (*pbdocstore.DeleteTemplateResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.template_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.templateAccess.CanUserAccessTarget(ctx, req.Id, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	dTmpl, err := s.getTemplate(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if !check && !userInfo.SuperUser {
		if dTmpl.CreatorJob == "" {
			return nil, errorsdocstore.ErrTemplateNoPerms
		}

		// Make sure the highest job grade can delete the template
		grade := s.jobs.GetHighestJobGrade(userInfo.Job)
		if grade == nil || (userInfo.Job == dTmpl.CreatorJob && grade.Grade != userInfo.JobGrade) {
			return nil, errorsdocstore.ErrTemplateNoPerms
		}
	}

	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		UPDATE(
			tDTemplates.DeletedAt,
		).
		SET(
			tDTemplates.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tDTemplates.CreatorJob.EQ(jet.String(userInfo.Job)),
			tDTemplates.ID.EQ(jet.Uint64(req.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbdocstore.DeleteTemplateResponse{}, nil
}
