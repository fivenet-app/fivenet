package mailer

import (
	"context"
	"strings"
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
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	emailLastChangedInterval = 14 * 24 * time.Hour
)

var tUserProps = table.FivenetUserProps

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

		e, err := s.getEmail(ctx, userInfo, resp.GetEmails()[idx].GetId(), true, true)
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
	userInfo *userinfo.UserInfo,
	emailId int64,
	withAccess bool,
	withSettings bool,
) (*maileremails.Email, error) {
	email, err := s.store.GetEmail(ctx, s.db, emailId, userInfo != nil && userInfo.GetJobAdmin())
	if err != nil {
		return nil, err
	}

	if withAccess {
		access, err := s.getEmailAccess(ctx, emailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		email.Access = access
	}

	if withSettings {
		settings, err := s.store.GetEmailSettings(ctx, s.db, emailId)
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

	email, err := s.getEmail(ctx, userInfo, req.GetId(), true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return &pbmailer.GetEmailResponse{
		Email: email,
	}, nil
}

func (s *Server) getEmailAccess(ctx context.Context, emailId int64) (*maileraccess.Access, error) {
	return s.access.ListTargetAccess(ctx, s.db, emailId, mailerSubjectAccessOptions)
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
		fields, err := permsmailer.MailerService.CreateOrUpdateEmail.FieldsTyped.Get(s.ps, userInfo)
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

	if req.GetEmail().GetId() <= 0 {
		// Check if user already has a personal email
		if req.Email.UserId != nil {
			email, err := s.store.GetEmailByUserID(
				ctx,
				tx,
				req.GetEmail().GetUserId(),
			)
			if err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}

			if email != nil {
				if email.GetDeactivated() {
					return nil, errorsmailer.ErrEmailDisabled
				}

				return nil, errorsmailer.ErrAddresseAlreadyTaken
			}
		}

		lastId, err := s.createEmail(ctx, tx, req.GetEmail(), userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		req.Email.SetId(lastId)
	} else {
		check, err := s.access.CanUserAccessTarget(
			ctx,
			req.GetEmail().GetId(),
			userInfo,
			int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		if !check {
			return nil, errorsmailer.ErrNoPerms
		}

		email, err := s.getEmail(ctx, userInfo, req.GetEmail().GetId(), false, false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if !userInfo.GetJobAdmin() && email.GetDeactivated() {
			return nil, errorsmailer.ErrEmailDisabled
		}

		tEmails := table.FivenetMailerEmails

		label := mysql.NULL
		if req.GetEmail().GetLabel() != "" {
			label = mysql.String(req.GetEmail().GetLabel())
		}

		sets := []any{
			tEmails.Label.SET(mysql.StringExp(label)),
			tEmails.CreatorID.SET(mysql.Int32(userInfo.GetUserId())),
		}

		// Update email only when necessary and allowed
		if strings.Compare(req.GetEmail().GetEmail(), email.GetEmail()) != 0 {
			if email.GetEmailChanged() != nil {
				// Check if last email change is at least 2 weeks ago
				since := time.Since(email.GetEmailChanged().AsTime())
				if since < emailLastChangedInterval {
					return nil, errorsmailer.ErrEmailChangeTooEarly
				}
			}

			sets = append(sets,
				tEmails.Email.SET(mysql.String(req.GetEmail().GetEmail())),
				tEmails.EmailChanged.SET(mysql.CURRENT_TIMESTAMP()),
			)
		}

		if userInfo.GetJobAdmin() {
			sets = append(sets,
				tEmails.Deactivated.SET(mysql.Bool(req.GetEmail().GetDeactivated())),
			)
		}

		condition := tEmails.ID.EQ(mysql.Int64(req.GetEmail().GetId()))
		if req.Email.Job != nil {
			condition = condition.AND(tEmails.Job.EQ(mysql.String(userInfo.GetJob())))
		} else {
			condition = condition.AND(tEmails.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
		}

		stmt := tEmails.
			UPDATE().
			SET(
				sets[0],
				sets[1:]...,
			).
			WHERE(mysql.AND(
				tEmails.ID.EQ(mysql.Int64(req.GetEmail().GetId())),
				condition,
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		// Update user email in the user props
		if req.Email.UserId != nil {
			upStmt := tUserProps.
				INSERT(
					tUserProps.UserID,
					tUserProps.Email,
				).
				VALUES(
					userInfo.GetUserId(),
					email.GetEmail(),
				).
				ON_DUPLICATE_KEY_UPDATE(
					tUserProps.Email.SET(mysql.String(email.GetEmail())),
				)

			if _, err := upStmt.ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}
	}

	// Only handle access changes for job emails
	if req.Email.Job != nil && req.GetEmail().GetAccess() != nil {
		if req.GetEmail().GetAccess().IsEmpty() {
			return nil, errorsmailer.ErrEmailAccessRequired
		}

		fallbackAccess := &maileraccess.Access{
			Jobs: []*maileraccess.JobAccess{{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
			}},
		}

		normalizedAccess, err := access.NormalizeAccess(
			req.GetEmail().GetAccess(),
			nil,
			fallbackAccess,
			15,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if _, err := s.access.ReplaceTargetAccess(
			ctx,
			tx,
			s.accessResolver,
			req.GetEmail().GetId(),
			normalizedAccess,
			mailerSubjectAccessOptions,
		); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	resp := &pbmailer.CreateOrUpdateEmailResponse{}
	resp.Email, err = s.getEmail(ctx, userInfo, req.GetEmail().GetId(), true, true)
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return resp, nil
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
		upStmt := tUserProps.
			INSERT(
				tUserProps.UserID,
				tUserProps.Email,
			).
			VALUES(
				userInfo.GetUserId(),
				email.GetEmail(),
			).
			ON_DUPLICATE_KEY_UPDATE(
				tUserProps.Email.SET(mysql.String(email.GetEmail())),
			)

		if _, err := upStmt.ExecContext(ctx, tx); err != nil {
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

	email, err := s.getEmail(ctx, userInfo, req.GetId(), false, false)
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
