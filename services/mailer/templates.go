package mailer

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	mailertemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/templates"
	pbmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
)

func (s *Server) ListTemplates(
	ctx context.Context,
	req *pbmailer.ListTemplatesRequest,
) (*pbmailer.ListTemplatesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetEmailId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	templates, err := s.store.ListTemplates(ctx, s.db, req.GetEmailId(), 25)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.ListTemplatesResponse{Templates: templates}, nil
}

func (s *Server) getTemplate(
	ctx context.Context,
	id int64,
	emailId *int64,
) (*mailertemplates.Template, error) {
	template, err := s.store.GetTemplate(ctx, s.db, id, emailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return template, nil
}

func (s *Server) GetTemplate(
	ctx context.Context,
	req *pbmailer.GetTemplateRequest,
) (*pbmailer.GetTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetEmailId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return resp, nil
}

func (s *Server) CreateOrUpdateTemplate(
	ctx context.Context,
	req *pbmailer.CreateOrUpdateTemplateRequest,
) (*pbmailer.CreateOrUpdateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetTemplate().GetEmailId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrFailedQuery
	}

	if req.GetTemplate().GetId() <= 0 {
		count, err := s.store.CountTemplatesByCreatorJob(ctx, s.db, userInfo.GetJob())
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		// Max 5 templates per email
		if count >= 5 {
			return nil, errorsmailer.ErrTemplateLimitReached
		}

		if req.Template.CreatorJob != nil {
			req.Template.CreatorJob = &userInfo.Job
		}

		lastID, err := s.store.CreateTemplate(ctx, s.db, req.GetTemplate(), userInfo.GetUserId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		req.Template.Id = lastID

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		template, err := s.getTemplate(ctx, req.GetTemplate().GetId(), &req.Template.EmailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if template == nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if err := s.store.UpdateTemplate(ctx, s.db, req.GetTemplate()); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	template, err := s.getTemplate(ctx, req.GetTemplate().GetId(), &req.Template.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.CreateOrUpdateTemplateResponse{Template: template}, nil
}

func (s *Server) DeleteTemplate(
	ctx context.Context,
	req *pbmailer.DeleteTemplateRequest,
) (*pbmailer.DeleteTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrFailedQuery
	}

	if err := s.store.DeleteTemplate(ctx, s.db, req.GetId()); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbmailer.DeleteTemplateResponse{}, nil
}
