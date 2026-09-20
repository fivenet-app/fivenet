package mailer

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	maileremails "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/emails"
	mailerevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/events"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer"
	permsmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
	mailerstore "github.com/fivenet-app/fivenet/v2026/stores/mailer"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	emailLastChangedInterval = 14 * 24 * time.Hour
)

var mailerSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(maileraccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_WRITE),
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
	},
}

func (s *Server) ListEmails(
	ctx context.Context,
	req *pbmailer.ListEmailsRequest,
) (*pbmailer.ListEmailsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	pag, emails, err := s.store.ListEmails(
		ctx,
		s.db,
		userInfo,
		req.GetPagination(),
		req.GetAll(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	resp := &pbmailer.ListEmailsResponse{
		Pagination: pag,
		Emails:     emails,
	}

	// Retrieve user's private email with access and settings
	for idx := range resp.GetEmails() {
		if resp.GetEmails()[idx] == nil || resp.Emails[idx].UserId == nil {
			continue
		}

		if resp.GetEmails()[idx].GetUserId() != userInfo.GetUserId() {
			continue
		}

		e, err := s.getEmail(ctx, s.db, userInfo, resp.GetEmails()[idx].GetId(), true, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		resp.Emails[idx] = e
		break
	}

	return resp, nil
}

func (s *Server) getEmail(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	emailId int64,
	withAccess bool,
	withSettings bool,
) (*maileremails.Email, error) {
	email, err := s.store.GetEmail(ctx, db, emailId, userInfo != nil && userInfo.GetJobAdmin())
	if err != nil {
		return nil, err
	}

	if withAccess {
		access, err := s.getEmailAccess(ctx, db, emailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		email.Access = access
	}

	if withSettings {
		settings, err := s.store.GetEmailSettings(ctx, db, emailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		email.Settings = settings
	}

	return email, nil
}

func (s *Server) GetEmail(
	ctx context.Context,
	req *pbmailer.GetEmailRequest,
) (*pbmailer.GetEmailResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, s.db, userInfo, req.GetId(), true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return &pbmailer.GetEmailResponse{
		Email: email,
	}, nil
}

func (s *Server) getEmailAccess(ctx context.Context, db qrm.DB, emailId int64) (*maileraccess.Access, error) {
	return s.access.ListTargetAccess(ctx, db, emailId, mailerSubjectAccessOptions)
}

func (s *Server) CreateOrUpdateEmail(
	ctx context.Context,
	req *pbmailer.CreateOrUpdateEmailRequest,
) (*pbmailer.CreateOrUpdateEmailResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.Email.UserId != nil {
		req.Email.UserId = &userInfo.UserId
		req.Email.Job = nil
	} else {
		req.Email.UserId = nil
		req.Email.Job = &userInfo.Job

		// Field Permission Check
		fields, err := permsmailer.MailerService.CreateOrUpdateEmail.FieldsTyped.Get(
			s.perms,
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if !fields.Contains(permsmailer.MailerServiceCreateOrUpdateEmailFieldsPermValueJob) {
			return nil, errswrap.NewError(err, errorsmailer.ErrEmailAccessDenied)
		}
	}

	if err := s.validateEmail(
		ctx,
		userInfo,
		req.GetEmail().GetEmail(),
		req.Email.Job != nil,
	); err != nil {
		return nil, err
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()
	auditAction := audit.EventAction_EVENT_ACTION_CREATED

	if req.GetEmail().GetId() <= 0 {
		if req.Email.UserId != nil {
			emailID, action, err := s.createOrRestorePersonalEmail(ctx, tx, req.GetEmail(), userInfo)
			if err != nil {
				return nil, err
			}
			auditAction = action
			req.Email.SetId(emailID)
		} else {
			lastID, err := s.createEmail(ctx, tx, req.GetEmail(), userInfo)
			if err != nil {
				return nil, err
			}
			req.Email.SetId(lastID)
		}
	} else {
		auditAction = audit.EventAction_EVENT_ACTION_UPDATED
		if err := s.updateExistingEmail(ctx, tx, req.GetEmail(), userInfo); err != nil {
			return nil, err
		}
	}

	if err := s.applyJobEmailAccess(ctx, tx, req.GetEmail(), userInfo); err != nil {
		return nil, err
	}

	// Keep the calculated visibility map in sync for private emails as well.
	// Job emails refresh their visibility through applyJobEmailAccess, while
	// private emails have no ACL entries to trigger that refresh.
	if req.GetEmail().GetUserId() > 0 {
		if err := s.access.RefreshTargetVisibility(ctx, tx, req.GetEmail().GetId()); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	resp := &pbmailer.CreateOrUpdateEmailResponse{}
	resp.Email, err = s.getEmail(ctx, s.db, userInfo, req.GetEmail().GetId(), true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailerevents.MailerEvent{
		Data: &mailerevents.MailerEvent_EmailUpdate{
			EmailUpdate: resp.GetEmail(),
		},
	},
		resp.GetEmail().GetId(),
	)

	grpc_audit.SetAction(ctx, auditAction)

	return resp, nil
}

func (s *Server) createOrRestorePersonalEmail(
	ctx context.Context,
	tx qrm.DB,
	email *maileremails.Email,
	userInfo *userinfo.UserInfo,
) (int64, audit.EventAction, error) {
	existing, err := s.store.GetEmailByUserID(ctx, tx, email.GetUserId())
	if err != nil {
		return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if existing == nil {
		id, err := s.createEmail(ctx, tx, email, userInfo)
		return id, audit.EventAction_EVENT_ACTION_CREATED, err
	}

	if existing.GetDeactivated() {
		return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errorsmailer.ErrEmailDisabled
	}

	if existing.GetDeletedAt() == nil {
		// A repeated create request for the same personal address is idempotent.
		if existing.GetEmail() == email.GetEmail() {
			if err := s.store.UpdateEmail(ctx, tx, mailerstore.EmailUpdate{
				ID:        existing.GetId(),
				Email:     email.GetEmail(),
				Label:     email.Label,
				UserID:    &userInfo.UserId,
				CreatorID: userInfo.GetUserId(),
			}); err != nil {
				return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
			if err := s.store.UpdateUserEmailProperty(ctx, tx, userInfo.GetUserId(), email.GetEmail()); err != nil {
				return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
			return existing.GetId(), audit.EventAction_EVENT_ACTION_UPDATED, nil
		}
		return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errorsmailer.ErrAddresseAlreadyTaken
	}

	// Deleted personal emails are hidden from normal listings, but their unique
	// user_id/email indexes still reserve the row. Reuse the row instead.
	if err := s.store.RestoreEmail(
		ctx,
		tx,
		existing.GetId(),
		email.GetEmail(),
		email.Label,
		userInfo.GetUserId(),
	); err != nil {
		if dbutils.IsDuplicateError(err) {
			return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errorsmailer.ErrAddresseAlreadyTaken
		}
		return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if err := s.store.UpdateUserEmailProperty(ctx, tx, userInfo.GetUserId(), email.GetEmail()); err != nil {
		return 0, audit.EventAction_EVENT_ACTION_UNSPECIFIED, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return existing.GetId(), audit.EventAction_EVENT_ACTION_RESTORED, nil
}

func (s *Server) updateExistingEmail(
	ctx context.Context,
	tx qrm.DB,
	email *maileremails.Email,
	userInfo *userinfo.UserInfo,
) error {
	check, err := s.access.CanUserAccessTarget(
		ctx,
		email.GetId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
	)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return errorsmailer.ErrNoPerms
	}

	existing, err := s.getEmail(ctx, tx, userInfo, email.GetId(), false, false)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if !userInfo.GetJobAdmin() && existing.GetDeactivated() {
		return errorsmailer.ErrEmailDisabled
	}

	emailChanged := email.GetEmail() != existing.GetEmail()
	if emailChanged {
		if existing.GetEmailChanged() != nil {
			since := time.Since(existing.GetEmailChanged().AsTime())
			if since < emailLastChangedInterval {
				return errorsmailer.ErrEmailChangeTooEarly
			}
		}
	}

	update := mailerstore.EmailUpdate{
		ID:           email.GetId(),
		Email:        email.GetEmail(),
		EmailChanged: emailChanged,
		Label:        email.Label,
		CreatorID:    userInfo.GetUserId(),
	}
	var deactivated *bool
	if userInfo.GetJobAdmin() {
		deactivatedValue := email.GetDeactivated()
		deactivated = &deactivatedValue
	}
	update.Deactivated = deactivated

	if email.Job != nil {
		update.Job = &userInfo.Job
	} else {
		update.UserID = &userInfo.UserId
	}

	if err := s.store.UpdateEmail(ctx, tx, update); err != nil {
		if dbutils.IsDuplicateError(err) {
			return errorsmailer.ErrAddresseAlreadyTaken
		}
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if email.UserId != nil {
		if err := s.store.UpdateUserEmailProperty(ctx, tx, userInfo.GetUserId(), email.GetEmail()); err != nil {
			return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return nil
}

func (s *Server) applyJobEmailAccess(
	ctx context.Context,
	tx qrm.DB,
	email *maileremails.Email,
	userInfo *userinfo.UserInfo,
) error {
	if email.Job == nil {
		return nil
	}

	fallbackAccess := &maileraccess.Access{
		Jobs: []*maileraccess.JobAccess{{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
		}},
	}
	if highestGrade, ok := s.enricher.GetHighestJobGrade(userInfo.GetJob()); ok {
		fallbackAccess.GetJobs()[0].MinimumGrade = highestGrade
	}

	currentAccess := email.GetAccess()
	if currentAccess == nil {
		currentAccess = &maileraccess.Access{}
	}
	normalizedAccess, err := access.NormalizeAccess(currentAccess, nil, fallbackAccess, 15)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if _, err := s.access.ReplaceTargetAccess(
		ctx,
		tx,
		s.accessResolver,
		email.GetId(),
		normalizedAccess,
		mailerSubjectAccessOptions,
	); err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return nil
}

func (s *Server) createEmail(
	ctx context.Context,
	tx qrm.DB,
	email *maileremails.Email,
	userInfo *userinfo.UserInfo,
) (int64, error) {
	lastId, err := s.store.CreateEmail(ctx, tx, email, userInfo.GetUserId())
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return 0, errorsmailer.ErrAddresseAlreadyTaken
		}
		return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Update user email in the user props if it is a "private" email
	if email.UserId != nil {
		if err := s.store.UpdateUserEmailProperty(ctx, tx, userInfo.GetUserId(), email.GetEmail()); err != nil {
			return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return lastId, nil
}

func (s *Server) DeleteEmail(
	ctx context.Context,
	req *pbmailer.DeleteEmailRequest,
) (*pbmailer.DeleteEmailResponse, error) {
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
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, s.db, userInfo, req.GetId(), false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Make sure that the user is not deleting their own personal email
	if email.Job == nil && email.UserId != nil {
		return nil, errorsmailer.ErrCantDeleteOwnEmail
	}

	s.sendUpdate(ctx, &mailerevents.MailerEvent{
		Data: &mailerevents.MailerEvent_EmailDelete{
			EmailDelete: req.GetId(),
		},
	},
		req.GetId(),
	)

	var deletedAtTime *timestamp.Timestamp
	if email == nil || email.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteEmail(ctx, s.db, req.GetId(), deletedAtTime); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.DeleteEmailResponse{}, nil
}

func (s *Server) GetEmailProposals(
	ctx context.Context,
	req *pbmailer.GetEmailProposalsRequest,
) (*pbmailer.GetEmailProposalsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.UserId != nil && userInfo.GetJobAdmin() {
		userInfo.UserId = req.GetUserId()
	}

	forJob := req.Job != nil && req.GetJob()
	emails, domains, err := s.generateEmailProposals(ctx, userInfo, forJob)
	if err != nil {
		return nil, err
	}

	return &pbmailer.GetEmailProposalsResponse{
		Emails:  emails,
		Domains: domains,
	}, nil
}
