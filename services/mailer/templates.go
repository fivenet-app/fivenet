package mailer

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tTemplates = table.FivenetMailerTemplates.AS("template")

func (s *Server) ListTemplates(ctx context.Context, req *pbmailer.ListTemplatesRequest) (*pbmailer.ListTemplatesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.EmailId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_READ)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	stmt := tTemplates.
		SELECT(
			tTemplates.ID,
			tTemplates.CreatedAt,
			tTemplates.UpdatedAt,
			tTemplates.DeletedAt,
			tTemplates.EmailID,
			tTemplates.Title,
			tTemplates.Content,
		).
		FROM(tTemplates).
		WHERE(jet.AND(
			tTemplates.EmailID.EQ(jet.Uint64(req.EmailId)),
		)).
		LIMIT(25)

	resp := &pbmailer.ListTemplatesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Templates); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) getTemplate(ctx context.Context, id uint64, emailId *uint64) (*mailer.Template, error) {
	condition := tTemplates.ID.EQ(jet.Uint64(id))

	if emailId == nil || *emailId <= 0 {
		condition = condition.AND(
			tTemplates.EmailID.EQ(jet.IntExp(jet.NULL)),
		)
	} else {
		condition = condition.AND(
			tTemplates.EmailID.EQ(jet.Uint64(*emailId)),
		)
	}

	stmt := tTemplates.
		SELECT(
			tTemplates.ID,
			tTemplates.CreatedAt,
			tTemplates.UpdatedAt,
			tTemplates.DeletedAt,
			tTemplates.EmailID,
			tTemplates.Title,
			tTemplates.Content,
			tTemplates.CreatorJob,
			tTemplates.CreatorID,
		).
		FROM(tTemplates).
		WHERE(condition).
		LIMIT(1)

	dest := &mailer.Template{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) GetTemplate(ctx context.Context, req *pbmailer.GetTemplateRequest) (*pbmailer.GetTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "GetTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.EmailId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_READ)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	resp := &pbmailer.GetTemplateResponse{}
	resp.Template, err = s.getTemplate(ctx, req.TemplateId, &req.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateOrUpdateTemplate(ctx context.Context, req *pbmailer.CreateOrUpdateTemplateRequest) (*pbmailer.CreateOrUpdateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Template.EmailId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrFailedQuery
	}

	tTemplates := table.FivenetMailerTemplates
	if req.Template.Id <= 0 {
		countStmt := tTemplates.
			SELECT(
				jet.COUNT(tTemplates.ID).AS("datacount.totalcount"),
			).
			FROM(tTemplates).
			WHERE(
				tTemplates.CreatorJob.EQ(jet.String(userInfo.Job)),
			)

		var count database.DataCount
		if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}

		// Max 5 templates per email
		if count.TotalCount >= 5 {
			return nil, errorsmailer.ErrTemplateLimitReached
		}

		if req.Template.CreatorJob != nil {
			req.Template.CreatorJob = &userInfo.Job
		}

		stmt := tTemplates.
			INSERT(
				tTemplates.EmailID,
				tTemplates.Title,
				tTemplates.Content,
				tTemplates.CreatorJob,
				tTemplates.CreatorID,
			).
			VALUES(
				req.Template.EmailId,
				req.Template.Title,
				req.Template.Content,
				req.Template.CreatorJob,
				userInfo.UserId,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		req.Template.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		template, err := s.getTemplate(ctx, req.Template.Id, &req.Template.EmailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if template == nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		stmt := tTemplates.
			UPDATE().
			SET(
				tTemplates.Title.SET(jet.String(req.Template.Title)),
				tTemplates.Content.SET(jet.String(req.Template.Content)),
			).
			WHERE(jet.AND(
				tTemplates.ID.EQ(jet.Uint64(req.Template.Id)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	template, err := s.getTemplate(ctx, req.Template.Id, &req.Template.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.CreateOrUpdateTemplateResponse{
		Template: template,
	}, nil
}

func (s *Server) DeleteTemplate(ctx context.Context, req *pbmailer.DeleteTemplateRequest) (*pbmailer.DeleteTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteTemplate",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Id, userInfo, mailer.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrFailedQuery
	}

	stmt := tTemplates.
		UPDATE().
		SET(
			tTemplates.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tTemplates.ID.EQ(jet.Uint64(req.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbmailer.DeleteTemplateResponse{}, nil
}
