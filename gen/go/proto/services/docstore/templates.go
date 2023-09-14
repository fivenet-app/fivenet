package docstore

import (
	"bytes"
	context "context"
	"errors"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tDTemplates          = table.FivenetDocumentsTemplates.AS("templateshort")
	tDTemplatesJobAccess = table.FivenetDocumentsTemplatesJobAccess.AS("templatejobaccess")
)

var (
	ErrTemplateFailed = status.Error(codes.InvalidArgument, "errors.DocStoreService.ErrTemplateFailed")
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *Server) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tDTemplates.
		SELECT(
			tDTemplates.ID,
			tDTemplates.Weight,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Schema,
			tDTemplates.CreatorID,
			tDTemplates.CreatorJob,
		).
		FROM(
			tDTemplates.
				INNER_JOIN(tDTemplatesJobAccess,
					tDTemplatesJobAccess.TemplateID.EQ(tDTemplates.ID).
						AND(tDTemplatesJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tDTemplatesJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tDCategory,
					tDCategory.ID.EQ(tDTemplates.CategoryID),
				),
		).
		WHERE(
			tDTemplates.DeletedAt.IS_NULL(),
		).
		ORDER_BY(
			tDTemplates.Weight.DESC(),
			tDTemplates.ID.ASC(),
		).
		GROUP_BY(tDTemplates.ID)

	resp := &ListTemplatesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Templates); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.checkIfUserHasAccessToTemplate(ctx, req.TemplateId, userInfo, false, documents.ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, ErrFailedQuery
	}
	if !check && !userInfo.SuperUser {
		return nil, ErrTemplateNoPerms
	}

	resp := &GetTemplateResponse{}
	resp.Template, err = s.getTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, ErrFailedQuery
	}

	if req.Render == nil || !*req.Render {
		resp.Template.JobAccess, err = s.getTemplateJobAccess(ctx, req.TemplateId)
		if err != nil {
			return nil, ErrFailedQuery
		}
	} else if req.Render != nil && *req.Render && req.Data != nil {
		// Parse data as json for the templating process
		var data map[string]interface{}
		if err := json.UnmarshalFromString(*req.Data, &data); err != nil {
			return nil, ErrTemplateFailed
		}

		resp.Template.ContentTitle, resp.Template.State, resp.Template.Content, err = s.renderTemplate(resp.Template, data)
		if err != nil {
			return nil, ErrTemplateFailed
		}

		resp.Rendered = true
	}

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
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.CreatorID,
			tDTemplates.CreatorJob,
		).
		FROM(
			tDTemplates.
				LEFT_JOIN(tDCategory,
					tDCategory.ID.EQ(tDTemplates.CategoryID),
				),
		).
		WHERE(
			tDTemplates.ID.EQ(jet.Uint64(templateId)),
		).
		LIMIT(1)

	var dest documents.Template
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) renderTemplate(docTmpl *documents.Template, data map[string]interface{}) (outTile string, outState string, out string, err error) {
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

func (s *Server) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "CreateTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	categoryId := jet.NULL
	if req.Template.Category != nil {
		cat, err := s.getCategory(ctx, req.Template.Category.Id)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			categoryId = jet.Uint64(cat.Id)
		}
	}

	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		INSERT(
			tDTemplates.Weight,
			tDTemplates.CategoryID,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.CreatorID,
			tDTemplates.CreatorJob,
		).
		VALUES(
			req.Template.Weight,
			categoryId,
			req.Template.Title,
			req.Template.Description,
			req.Template.ContentTitle,
			req.Template.Content,
			req.Template.State,
			req.Template.ContentAccess,
			req.Template.Schema,
			userInfo.UserId,
			userInfo.Job,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.handleTemplateAccessChanges(ctx, tx, uint64(lastId), req.Template.JobAccess); err != nil {
		return nil, ErrFailedQuery
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return &CreateTemplateResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) UpdateTemplate(ctx context.Context, req *UpdateTemplateRequest) (*UpdateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "UpdateTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToTemplate(ctx, req.Template.Id, userInfo, false, documents.ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, ErrFailedQuery
	}
	if !check && !userInfo.SuperUser {
		return nil, ErrTemplateNoPerms
	}

	categoryId := jet.NULL
	if req.Template.Category != nil {
		cat, err := s.getCategory(ctx, req.Template.Category.Id)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			categoryId = jet.Uint64(cat.Id)
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
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
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
		).
		SET(
			req.Template.Weight,
			categoryId,
			req.Template.Title,
			req.Template.Description,
			req.Template.ContentTitle,
			req.Template.Content,
			req.Template.State,
			req.Template.ContentAccess,
			req.Template.Schema,
		).
		WHERE(
			tDTemplates.ID.EQ(jet.Uint64(req.Template.Id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, ErrFailedQuery
	}

	if err := s.handleTemplateAccessChanges(ctx, tx, uint64(req.Template.Id), req.Template.JobAccess); err != nil {
		return nil, ErrFailedQuery
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &UpdateTemplateResponse{
		Id: req.Template.Id,
	}, nil
}

func (s *Server) DeleteTemplate(ctx context.Context, req *DeleteTemplateRequest) (*DeleteTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToTemplate(ctx, req.Id, userInfo, false, documents.ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, ErrFailedQuery
	}

	dTmpl, err := s.getTemplate(ctx, req.Id)
	if err != nil {
		return nil, ErrFailedQuery
	}

	if !check && !userInfo.SuperUser {
		if dTmpl.CreatorJob == nil {
			return nil, ErrTemplateNoPerms
		}

		// Make sure the highest job grade can delete the template
		grade := s.cache.GetHighestJobGrade(userInfo.Job)
		if grade == nil || (userInfo.Job == *dTmpl.CreatorJob && grade.Grade != userInfo.JobGrade) {
			return nil, ErrTemplateNoPerms
		}
	}

	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		DELETE().
		WHERE(
			tDTemplates.CreatorJob.EQ(jet.String(userInfo.Job)).AND(
				tDTemplates.ID.EQ(jet.Uint64(req.Id))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &DeleteTemplateResponse{}, nil
}
