package mailer

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	permsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tEmails = table.FivenetMailerEmails.AS("email")

	tEmailsJobAccess            = table.FivenetMailerEmailsJobAccess
	tEmailsUserAccess           = table.FivenetMailerEmailsUserAccess
	tEmailsQualificationsAccess = table.FivenetMailerEmailsQualificationsAccess

	tQualificationsResults = table.FivenetQualificationsResults
)

func (s *Server) ListEmails(ctx context.Context, req *ListEmailsRequest) (*ListEmailsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	emails, err := ListUserEmails(ctx, s.db, userInfo)
	if err != nil {
		return nil, err
	}

	resp := &ListEmailsResponse{
		Emails: emails,
	}

	return resp, nil
}

func ListUserEmails(ctx context.Context, tx qrm.DB, userInfo *userinfo.UserInfo) ([]*mailer.Email, error) {
	condition := jet.Bool(true)
	if !userInfo.SuperUser {
		access := int32(mailer.AccessLevel_ACCESS_LEVEL_READ)
		condition = condition.AND(jet.AND(
			tEmails.DeletedAt.IS_NULL(),
			tEmails.Disabled.IS_FALSE(),
			jet.OR(
				tEmails.UserID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tEmailsUserAccess.Access.IS_NULL(),
					tEmailsJobAccess.Access.IS_NOT_NULL(),
					tEmailsJobAccess.Access.GT_EQ(jet.Int32(access)),
				),
				jet.AND(
					tEmailsUserAccess.Access.IS_NOT_NULL(),
					tEmailsUserAccess.Access.GT_EQ(jet.Int32(access)),
				),
				jet.AND(
					tEmailsQualificationsAccess.Access.IS_NOT_NULL(),
					tEmailsQualificationsAccess.Access.GT_EQ(jet.Int32(access)),
					tQualificationsResults.DeletedAt.IS_NULL(),
					tQualificationsResults.QualificationID.EQ(tEmailsQualificationsAccess.QualificationID),
					tQualificationsResults.Status.EQ(jet.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL.Number()))),
				),
			),
		))
	}

	stmt := tEmails.
		SELECT(
			tEmails.ID,
			tEmails.CreatedAt,
			tEmails.UpdatedAt,
			tEmails.DeletedAt,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.Label,
			tEmails.Internal,
		).
		FROM(
			tEmails.
				LEFT_JOIN(tEmailsJobAccess,
					tEmailsJobAccess.EmailID.EQ(tEmails.ID).
						AND(tEmailsJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tEmailsJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tEmailsUserAccess,
					tEmailsUserAccess.EmailID.EQ(tEmails.ID).
						AND(tEmailsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tEmailsQualificationsAccess,
					tEmailsQualificationsAccess.EmailID.EQ(tEmails.ID),
				).
				LEFT_JOIN(tQualificationsResults,
					tQualificationsResults.QualificationID.EQ(tEmailsQualificationsAccess.QualificationID).
						AND(tQualificationsResults.UserID.EQ(jet.Int32(userInfo.UserId))),
				),
		).
		WHERE(condition).
		GROUP_BY(tEmails.ID).
		ORDER_BY(tEmails.Job.ASC(), tEmails.Label.ASC())

	resp := &ListEmailsResponse{}
	if err := stmt.QueryContext(ctx, tx, &resp.Emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return resp.Emails, nil
}

func (s *Server) getEmailByCondition(ctx context.Context, tx qrm.DB, condition jet.BoolExpression) (*mailer.Email, error) {
	stmt := tEmails.
		SELECT(
			tEmails.ID,
			tEmails.CreatedAt,
			tEmails.UpdatedAt,
			tEmails.DeletedAt,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.Label,
			tEmails.Internal,
		).
		FROM(tEmails).
		WHERE(condition).
		LIMIT(1)

	dest := &mailer.Email{}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) getEmail(ctx context.Context, emailId uint64, withAccess bool, withSettings bool) (*mailer.Email, error) {
	email, err := s.getEmailByCondition(ctx, s.db, tEmails.ID.EQ(jet.Uint64(emailId)))
	if err != nil {
		return nil, err
	}
	if email == nil {
		return nil, nil
	}

	if withAccess {
		access, err := s.getEmailAccess(ctx, emailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		email.Access = access
	}

	if withSettings {
		settings, err := s.getEmailSettings(ctx, emailId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		email.Settings = settings
	}

	return email, nil
}

func (s *Server) GetEmail(ctx context.Context, req *GetEmailRequest) (*GetEmailResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "GetEmail",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Id, userInfo, mailer.AccessLevel_ACCESS_LEVEL_READ)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, req.Id, true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if email.Disabled {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return &GetEmailResponse{
		Email: email,
	}, nil
}

func (s *Server) getEmailAccess(ctx context.Context, emailId uint64) (*mailer.Access, error) {
	access := &mailer.Access{}

	jobsAccess, err := s.access.Jobs.List(ctx, s.db, emailId)
	if err != nil {
		return nil, err
	}
	access.Jobs = jobsAccess

	usersAccess, err := s.access.Users.List(ctx, s.db, emailId)
	if err != nil {
		return nil, err
	}
	access.Users = usersAccess

	qualiAccess, err := s.access.Qualifications.List(ctx, s.db, emailId)
	if err != nil {
		return nil, err
	}
	access.Qualifications = qualiAccess

	return access, nil
}

func (s *Server) CreateOrUpdateEmail(ctx context.Context, req *CreateOrUpdateEmailRequest) (*CreateOrUpdateEmailResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateEmail",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	var err error
	req.Email.Email, err = s.validateEmailName(req.Email.Email)
	if err != nil {
		return nil, err
	}

	if req.Email.UserId != nil {
		req.Email.UserId = &userInfo.UserId
		req.Email.Job = nil
	} else {
		req.Email.UserId = nil
		req.Email.Job = &userInfo.Job

		// Field Permission Check
		fieldsAttr, err := s.ps.Attr(userInfo, permsmailer.MailerServicePerm, permsmailer.MailerServiceCreateOrUpdateEmailPerm, permsmailer.MailerServiceCreateOrUpdateEmailFieldsPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		if !slices.Contains(fields, "Job") {
			return nil, errswrap.NewError(err, errorsmailer.ErrEmailAccessRequired)
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if req.Email.Id <= 0 {
		// Check if user already has a personal email
		if req.Email.UserId != nil {
			email, err := s.getEmailByCondition(ctx, tx, tEmails.UserID.EQ(jet.Int32(*req.Email.UserId)))
			if err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}

			if email != nil {
				return nil, errorsmailer.ErrAddresseAlreadyTaken
			}
		}

		lastId, err := s.createEmail(ctx, tx, req.Email)
		if err != nil {
			return nil, err
		}

		req.Email.Id = uint64(lastId)
	} else {
		check, err := s.access.CanUserAccessTarget(ctx, req.Email.Id, userInfo, mailer.AccessLevel_ACCESS_LEVEL_MANAGE)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		if !check {
			return nil, errorsmailer.ErrNoPerms
		}

		label := jet.NULL
		if req.Email.Label != nil {
			label = jet.String(*req.Email.Label)
		}

		tEmails := table.FivenetMailerEmails
		condition := tEmails.Job.EQ(jet.String(userInfo.Job))
		if req.Email.UserId != nil {
			condition = tEmails.UserID.EQ(jet.Int32(userInfo.UserId))
		}

		stmt := tEmails.
			UPDATE(
				tEmails.Label,
				tEmails.Internal,
			).
			SET(
				label,
				jet.Bool(req.Email.Internal),
			).
			WHERE(jet.AND(
				tEmails.ID.EQ(jet.Uint64(req.Email.Id)),
				condition,
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	// Only handle access changes for job emails
	if req.Email.Job != nil && req.Email.Access != nil {
		if _, err := s.access.HandleAccessChanges(ctx, tx, req.Email.Id, req.Email.Access.Jobs, req.Email.Access.Users, req.Email.Access.Qualifications); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	resp := &CreateOrUpdateEmailResponse{}
	resp.Email, err = s.getEmail(ctx, req.Email.Id, true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_EmailUpdate{
			EmailUpdate: resp.Email,
		},
	},
		resp.Email.Id,
	)

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return resp, nil
}

func (s *Server) createEmail(ctx context.Context, tx qrm.DB, email *mailer.Email) (uint64, error) {
	tEmails := table.FivenetMailerEmails
	stmt := tEmails.
		INSERT(
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.Label,
			tEmails.Internal,
		).
		VALUES(
			email.Job,
			email.UserId,
			email.Email,
			email.Label,
			email.Internal,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return 0, errorsmailer.ErrAddresseAlreadyTaken
		}
		return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return uint64(lastId), nil
}

func (s *Server) DeleteEmail(ctx context.Context, req *DeleteEmailRequest) (*DeleteEmailResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteEmail",
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
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, req.Id, false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Make sure that the user is not deleting their own personal email
	if email.Job == nil && email.UserId != nil {
		return nil, errorsmailer.ErrCantDeleteOwnEmail
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_EmailDelete{
			EmailDelete: req.Id,
		},
	},
		req.Id,
	)

	tEmails := table.FivenetMailerEmails
	stmt := tEmails.
		UPDATE().
		SET(
			tEmails.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tEmails.ID.EQ(jet.Uint64(req.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteEmailResponse{}, nil
}

func (s *Server) validateEmailName(email string) (string, error) {
	if strings.Contains(email, "@") {
		before, _, found := strings.Cut(email, "@")
		if found {
			email = before
		}
	}
	email += "@fivenet.app"

	return email, nil
}

func (s *Server) GetEmailProposals(context.Context, *GetEmailProposalsRequest) (*GetEmailProposalsResponse, error) {
	// TODO

	return nil, nil
}
