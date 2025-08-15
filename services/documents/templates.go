package documents

import (
	"bytes"
	context "context"
	"errors"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var (
	tDTemplates       = table.FivenetDocumentsTemplates.AS("template_short")
	tDTemplatesAccess = table.FivenetDocumentsTemplatesAccess.AS("template_job_access")
)

func (s *Server) ListTemplates(
	ctx context.Context,
	req *pbdocuments.ListTemplatesRequest,
) (*pbdocuments.ListTemplatesResponse, error) {
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
				tDTemplatesAccess.Job.EQ(jet.String(userInfo.GetJob())),
				tDTemplatesAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.GetJobGrade())),
			),
		)).
		ORDER_BY(
			tDTemplates.Weight.DESC(),
			tDTemplates.ID.ASC(),
		).
		GROUP_BY(tDTemplates.ID)

	resp := &pbdocuments.ListTemplatesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Templates); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	for i := range resp.GetTemplates() {
		s.enricher.EnrichJobName(resp.GetTemplates()[i])
	}

	return resp, nil
}

func (s *Server) GetTemplate(
	ctx context.Context,
	req *pbdocuments.GetTemplateRequest,
) (*pbdocuments.GetTemplateResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.template_id", req.GetTemplateId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetTemplateId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	resp := &pbdocuments.GetTemplateResponse{}
	resp.Template, err = s.getTemplate(ctx, req.GetTemplateId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if resp.GetTemplate() == nil {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	if req.Render == nil || !req.GetRender() {
		resp.Template.JobAccess, err = s.templateAccess.Jobs.List(ctx, s.db, req.GetTemplateId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else if req.Render != nil && req.GetRender() && req.GetData() != nil {
		resp.Template.ContentTitle, resp.Template.State, resp.Template.Content, err = s.renderTemplate(resp.GetTemplate(), req.GetData())
		if err != nil {
			if s.ps.Can(userInfo, permsdocuments.DocumentsServicePerm, permsdocuments.DocumentsServiceCreateTemplatePerm) {
				return nil, err
			} else {
				return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateFailed)
			}
		}

		resp.Rendered = true
	}

	s.enricher.EnrichJobName(resp.GetTemplate())

	return resp, nil
}

func (s *Server) getTemplate(ctx context.Context, templateId int64) (*documents.Template, error) {
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
			tDTemplates.ID.EQ(jet.Int64(templateId)),
		)).
		LIMIT(1)

	var dest documents.Template
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	if dest.GetId() <= 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Server) renderTemplate(
	docTmpl *documents.Template,
	data *documents.TemplateData,
) (string, string, string, error) {
	// Render Title template
	titleTpl, err := template.
		New("title").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.GetContentTitle())
	if err != nil {
		return "", "", "", err
	}
	buf := &bytes.Buffer{}
	err = titleTpl.Execute(buf, data)
	if err != nil {
		return "", "", "", err
	}
	outTile := buf.String()

	// Render State template
	stateTpl, err := template.
		New("state").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.GetState())
	if err != nil {
		return "", "", "", err
	}

	buf.Reset()
	err = stateTpl.Execute(buf, data)
	if err != nil {
		return "", "", "", err
	}
	outState := buf.String()

	// Render Content template
	contentTpl, err := template.
		New("content").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.GetContent())
	if err != nil {
		return "", "", "", err
	}

	buf.Reset()
	err = contentTpl.Execute(buf, data)
	if err != nil {
		return "", "", "", err
	}
	out := buf.String()

	return outTile, outState, out, err
}

func (s *Server) CreateTemplate(
	ctx context.Context,
	req *pbdocuments.CreateTemplateRequest,
) (*pbdocuments.CreateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "CreateTemplate",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if documents.TemplateAccessHasDuplicates(req.GetTemplate().GetJobAccess()) ||
		documents.DocumentAccessHasDuplicates(req.GetTemplate().GetContentAccess()) {
		return nil, errorsdocuments.ErrTemplateAccessDuplicate
	}

	categoryId := jet.NULL
	if req.GetTemplate().GetCategory() != nil {
		cat, err := s.getCategory(ctx, req.GetTemplate().GetCategory().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if cat != nil {
			categoryId = jet.Int64(cat.GetId())
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
			req.GetTemplate().GetWeight(),
			categoryId,
			req.GetTemplate().GetTitle(),
			req.GetTemplate().GetDescription(),
			req.GetTemplate().Color,
			req.GetTemplate().Icon,
			req.GetTemplate().GetContentTitle(),
			req.GetTemplate().GetContent(),
			req.GetTemplate().GetState(),
			req.GetTemplate().GetContentAccess(),
			req.GetTemplate().GetSchema(),
			req.GetTemplate().GetWorkflow(),
			userInfo.GetJob(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := s.templateAccess.HandleAccessChanges(ctx, tx, lastId, req.GetTemplate().GetJobAccess(), nil, nil); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbdocuments.CreateTemplateResponse{
		Id: lastId,
	}, nil
}

func (s *Server) UpdateTemplate(
	ctx context.Context,
	req *pbdocuments.UpdateTemplateRequest,
) (*pbdocuments.UpdateTemplateResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.documents.template_id", req.GetTemplate().GetId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "UpdateTemplate",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetTemplate().GetId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	if documents.TemplateAccessHasDuplicates(req.GetTemplate().GetJobAccess()) ||
		documents.DocumentAccessHasDuplicates(req.GetTemplate().GetContentAccess()) {
		return nil, errorsdocuments.ErrTemplateAccessDuplicate
	}

	categoryId := jet.NULL
	if req.GetTemplate().GetCategory() != nil {
		cat, err := s.getCategory(ctx, req.GetTemplate().GetCategory().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if cat != nil {
			categoryId = jet.Int64(cat.GetId())
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
			req.GetTemplate().GetWeight(),
			categoryId,
			req.GetTemplate().GetTitle(),
			req.GetTemplate().GetDescription(),
			req.GetTemplate().Color,
			req.GetTemplate().Icon,
			req.GetTemplate().GetContentTitle(),
			req.GetTemplate().GetContent(),
			req.GetTemplate().GetState(),
			req.GetTemplate().GetContentAccess(),
			req.GetTemplate().GetSchema(),
			req.GetTemplate().GetWorkflow(),
		).
		WHERE(
			tDTemplates.ID.EQ(jet.Int64(req.GetTemplate().GetId())),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := s.templateAccess.HandleAccessChanges(ctx, tx, req.GetTemplate().GetId(), req.GetTemplate().GetJobAccess(), nil, nil); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tmpl, err := s.getTemplate(ctx, req.GetTemplate().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tmpl.JobAccess, err = s.templateAccess.Jobs.List(ctx, s.db, req.GetTemplate().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbdocuments.UpdateTemplateResponse{
		Template: tmpl,
	}, nil
}

func (s *Server) DeleteTemplate(
	ctx context.Context,
	req *pbdocuments.DeleteTemplateRequest,
) (*pbdocuments.DeleteTemplateResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.template_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "DeleteTemplate",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	dTmpl, err := s.getTemplate(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if !check && !userInfo.GetSuperuser() {
		if dTmpl.GetCreatorJob() == "" {
			return nil, errorsdocuments.ErrTemplateNoPerms
		}

		// Make sure the highest job grade can delete the template
		grade := s.jobs.GetHighestJobGrade(userInfo.GetJob())
		if grade == nil ||
			(userInfo.GetJob() == dTmpl.GetCreatorJob() && grade.GetGrade() != userInfo.GetJobGrade()) {
			return nil, errorsdocuments.ErrTemplateNoPerms
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
			tDTemplates.CreatorJob.EQ(jet.String(userInfo.GetJob())),
			tDTemplates.ID.EQ(jet.Int64(req.GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbdocuments.DeleteTemplateResponse{}, nil
}
