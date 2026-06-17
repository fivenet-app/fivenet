package documents

import (
	"bytes"
	context "context"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var templateSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	},
}

func (s *Server) sanitizeTemplateAccess(
	tmpl *documentstemplates.Template,
	sanitizeJobAccess bool,
	sanitizeContentAccess bool,
) error {
	if tmpl == nil {
		return nil
	}

	if sanitizeJobAccess {
		jobAccess, err := access.SanitizeJobAccessEntries(s.jobs, tmpl.GetJobAccess())
		if err != nil {
			return err
		}
		tmpl.JobAccess = access.NormalizeRequiredJobAccessFloors(jobAccess)
	}

	if sanitizeContentAccess {
		contentAccess, err := access.SanitizeAccessJobs(s.jobs, tmpl.GetContentAccess())
		if err != nil {
			return err
		}
		tmpl.ContentAccess = access.NormalizeRequiredAccessFloors(contentAccess)
	}

	return nil
}

func templateJobAccess(jobs []*documentstemplates.TemplateJobAccess) *resourcesaccess.Access {
	return &resourcesaccess.Access{Jobs: jobs}
}

func normalizeTemplateJobAccess(
	userInfo *userinfo.UserInfo,
	jobs []*documentstemplates.TemplateJobAccess,
) (*resourcesaccess.Access, error) {
	return access.NormalizeAccess(
		templateJobAccess(jobs),
		nil,
		&resourcesaccess.Access{
			Jobs: []*resourcesaccess.JobAccess{{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			}},
		},
		15,
	)
}

func (s *Server) ListTemplates(
	ctx context.Context,
	req *pbdocuments.ListTemplatesRequest,
) (*pbdocuments.ListTemplatesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	templates, err := s.store.ListTemplates(ctx, false, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	resp := &pbdocuments.ListTemplatesResponse{
		Templates: templates,
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
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	resp := &pbdocuments.GetTemplateResponse{}
	resp.Template, err = s.store.GetTemplate(ctx, req.GetTemplateId(), userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if resp.GetTemplate() == nil {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	if req.Render == nil || !req.GetRender() {
		templateAccess, err := s.templateAccess.ListTargetAccess(
			ctx,
			s.db,
			req.GetTemplateId(),
			templateSubjectAccessOptions,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.Template.JobAccess = templateAccess.GetJobs()
		if err := s.sanitizeTemplateAccess(resp.Template, true, true); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else if req.Render != nil && req.GetRender() && req.GetData() != nil {
		resp.Template.ContentTitle, resp.Template.State, resp.Template.Content, err = s.renderTemplate(
			resp.GetTemplate(),
			req.GetData(),
		)
		if err != nil {
			if s.ps.Can(
				userInfo,
				permsdocuments.TemplatesService.CreateTemplate.Perm,
			) {
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

func (s *Server) renderTemplate(
	docTmpl *documentstemplates.Template,
	data *documentstemplates.TemplateData,
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

	if err := s.sanitizeTemplateAccess(req.GetTemplate(), true, true); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if documentstemplates.TemplateAccessHasDuplicates(req.GetTemplate().GetJobAccess()) ||
		documentsaccess.DocumentAccessHasDuplicates(req.GetTemplate().GetContentAccess()) {
		return nil, errorsdocuments.ErrTemplateAccessDuplicate
	}

	var categoryId *int64
	if req.GetTemplate().GetCategory() != nil {
		cat, err := s.getCategory(ctx, req.GetTemplate().GetCategory().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if cat != nil {
			id := cat.GetId()
			categoryId = &id
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	lastId, err := s.store.CreateTemplate(ctx, tx, req.GetTemplate(), userInfo.GetJob(), categoryId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	normalizedAccess, err := normalizeTemplateJobAccess(userInfo, req.GetTemplate().GetJobAccess())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := s.templateAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		lastId,
		normalizedAccess,
		templateSubjectAccessOptions,
	); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

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

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetTemplate().GetId(),
		userInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	if err := s.sanitizeTemplateAccess(req.GetTemplate(), true, true); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if documentstemplates.TemplateAccessHasDuplicates(req.GetTemplate().GetJobAccess()) ||
		documentsaccess.DocumentAccessHasDuplicates(req.GetTemplate().GetContentAccess()) {
		return nil, errorsdocuments.ErrTemplateAccessDuplicate
	}

	var categoryId *int64
	if req.GetTemplate().GetCategory() != nil {
		cat, err := s.getCategory(ctx, req.GetTemplate().GetCategory().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if cat != nil {
			id := cat.GetId()
			categoryId = &id
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.UpdateTemplate(ctx, tx, req.GetTemplate(), categoryId); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	normalizedAccess, err := normalizeTemplateJobAccess(userInfo, req.GetTemplate().GetJobAccess())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := s.templateAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		req.GetTemplate().GetId(),
		normalizedAccess,
		templateSubjectAccessOptions,
	); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tmpl, err := s.store.GetTemplate(ctx, req.GetTemplate().GetId(), false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	templateAccess, err := s.templateAccess.ListTargetAccess(
		ctx,
		s.db,
		req.GetTemplate().GetId(),
		templateSubjectAccessOptions,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	tmpl.JobAccess = templateAccess.GetJobs()
	if err := s.sanitizeTemplateAccess(tmpl, true, true); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

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

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	dTmpl, err := s.store.GetTemplate(ctx, req.GetId(), true)
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

	var deletedAtTime *timestamp.Timestamp
	// Check if page has any un-deleted child pages
	if dTmpl.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteTemplate(
		ctx,
		s.db,
		req.GetId(),
		userInfo.GetJob(),
		deletedAtTime,
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteTemplateResponse{}, nil
}
