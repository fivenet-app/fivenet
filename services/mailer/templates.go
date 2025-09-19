package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tTemplates = table.FivenetMailerTemplates.AS("template")

func (s *Server) ListTemplates(
	ctx context.Context,
	req *pbmailer.ListTemplatesRequest,
) (*pbmailer.ListTemplatesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetEmailId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_READ,
	)
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
		WHERE(mysql.AND(
			tTemplates.EmailID.EQ(mysql.Int64(req.GetEmailId())),
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

func (s *Server) getTemplate(
	ctx context.Context,
	id int64,
	emailId *int64,
) (*mailer.Template, error) {
	condition := tTemplates.ID.EQ(mysql.Int64(id))

	if emailId == nil || *emailId <= 0 {
		condition = condition.AND(
			tTemplates.EmailID.EQ(mysql.IntExp(mysql.NULL)),
		)
	} else {
		condition = condition.AND(
			tTemplates.EmailID.EQ(mysql.Int64(*emailId)),
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

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) GetTemplate(
	ctx context.Context,
	req *pbmailer.GetTemplateRequest,
) (*pbmailer.GetTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "GetTemplate",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetEmailId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_READ,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	resp := &pbmailer.GetTemplateResponse{}
	resp.Template, err = s.getTemplate(ctx, req.GetTemplateId(), &req.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) CreateOrUpdateTemplate(
	ctx context.Context,
	req *pbmailer.CreateOrUpdateTemplateRequest,
) (*pbmailer.CreateOrUpdateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateTemplate",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetTemplate().GetEmailId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_MANAGE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrFailedQuery
	}

	tTemplates := table.FivenetMailerTemplates
	if req.GetTemplate().GetId() <= 0 {
		countStmt := tTemplates.
			SELECT(
				mysql.COUNT(tTemplates.ID).AS("data_count.total"),
			).
			FROM(tTemplates).
			WHERE(
				tTemplates.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
			)

		var count database.DataCount
		if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}

		// Max 5 templates per email
		if count.Total >= 5 {
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
				req.GetTemplate().GetEmailId(),
				req.GetTemplate().GetTitle(),
				req.GetTemplate().GetContent(),
				req.GetTemplate().GetCreatorJob(),
				userInfo.GetUserId(),
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		req.Template.Id = lastId

		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	} else {
		template, err := s.getTemplate(ctx, req.GetTemplate().GetId(), &req.Template.EmailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if template == nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		stmt := tTemplates.
			UPDATE().
			SET(
				tTemplates.Title.SET(mysql.String(req.GetTemplate().GetTitle())),
				tTemplates.Content.SET(mysql.String(req.GetTemplate().GetContent())),
			).
			WHERE(mysql.AND(
				tTemplates.ID.EQ(mysql.Int64(req.GetTemplate().GetId())),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
	}

	template, err := s.getTemplate(ctx, req.GetTemplate().GetId(), &req.Template.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.CreateOrUpdateTemplateResponse{
		Template: template,
	}, nil
}

func (s *Server) DeleteTemplate(
	ctx context.Context,
	req *pbmailer.DeleteTemplateRequest,
) (*pbmailer.DeleteTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteTemplate",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_MANAGE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrFailedQuery
	}

	stmt := tTemplates.
		UPDATE().
		SET(
			tTemplates.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tTemplates.ID.EQ(mysql.Int64(req.GetId())),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbmailer.DeleteTemplateResponse{}, nil
}
