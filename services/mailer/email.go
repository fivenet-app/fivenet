package mailer

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	permsmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	emailLastChangedInterval = 14 * 24 * time.Hour
	listEmailsPageSize       = 20
)

var (
	tEmails = table.FivenetMailerEmails.AS("email")

	tEmailsAccess = table.FivenetMailerEmailsAccess

	tQualificationsResults = table.FivenetQualificationsResults

	tUserProps = table.FivenetUserProps
)

func (s *Server) ListEmails(
	ctx context.Context,
	req *pbmailer.ListEmailsRequest,
) (*pbmailer.ListEmailsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.Bool(true)

	if !userInfo.GetSuperuser() || (userInfo.GetSuperuser() && req.All != nil && !req.GetAll()) {
		// Include deactivated e-mails
		condition = condition.AND(jet.AND(
			tEmails.DeletedAt.IS_NULL(),
			jet.OR(
				tEmails.UserID.EQ(jet.Int32(userInfo.GetUserId())),
				jet.OR(
					tEmailsAccess.UserID.EQ(jet.Int32(userInfo.GetUserId())),
					jet.AND(
						tEmailsAccess.Job.EQ(jet.String(userInfo.GetJob())),
						tEmailsAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.GetJobGrade())),
					),
					jet.AND(
						tEmailsAccess.QualificationID.IS_NOT_NULL(),
						tQualificationsResults.DeletedAt.IS_NULL(),
						tQualificationsResults.QualificationID.EQ(tEmailsAccess.QualificationID),
						tQualificationsResults.Status.EQ(
							jet.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
						),
					),
				),
			),
		))
	}

	countStmt := tEmails.
		SELECT(
			jet.COUNT(jet.DISTINCT(tEmails.ID)).AS("data_count.total"),
		).
		FROM(
			tEmails.
				LEFT_JOIN(tEmailsAccess,
					tEmailsAccess.TargetID.EQ(tEmails.ID).
						AND(tEmailsAccess.Access.GT_EQ(jet.Int32(int32(mailer.AccessLevel_ACCESS_LEVEL_READ)))),
				).
				LEFT_JOIN(tQualificationsResults,
					tQualificationsResults.QualificationID.EQ(tEmailsAccess.QualificationID).
						AND(tQualificationsResults.UserID.EQ(jet.Int32(userInfo.GetUserId()))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	pag, _ := req.GetPagination().GetResponseWithPageSize(count.Total, listEmailsPageSize)
	resp := &pbmailer.ListEmailsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	emails, err := ListUserEmails(ctx, s.db, userInfo, req.GetPagination(), true)
	if err != nil {
		return nil, err
	}
	resp.Emails = emails

	// Retrieve user's private email with access and settings
	for idx := range resp.GetEmails() {
		if resp.GetEmails()[idx] == nil || resp.Emails[idx].UserId == nil {
			continue
		}

		if resp.GetEmails()[idx].GetUserId() != userInfo.GetUserId() {
			continue
		}

		e, err := s.getEmail(ctx, resp.GetEmails()[idx].GetId(), true, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		resp.Emails[idx] = e
		break
	}

	resp.GetPagination().Update(len(resp.GetEmails()))

	return resp, nil
}

func ListUserEmails(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	pag *database.PaginationRequest,
	includeDisabled bool,
) ([]*mailer.Email, error) {
	condition := jet.Bool(true)
	baseCondition := tEmails.DeletedAt.IS_NULL()
	if !includeDisabled {
		baseCondition = baseCondition.AND(tEmails.Deactivated.IS_FALSE())
	}

	if !userInfo.GetSuperuser() {
		condition = condition.AND(jet.AND(
			baseCondition,
			jet.OR(
				tEmails.UserID.EQ(jet.Int32(userInfo.GetUserId())),
				jet.OR(
					tEmailsAccess.UserID.EQ(jet.Int32(userInfo.GetUserId())),
					jet.AND(
						tEmailsAccess.Job.EQ(jet.String(userInfo.GetJob())),
						tEmailsAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.GetJobGrade())),
					),
					jet.AND(
						tEmailsAccess.QualificationID.IS_NOT_NULL(),
						tQualificationsResults.DeletedAt.IS_NULL(),
						tQualificationsResults.QualificationID.EQ(tEmailsAccess.QualificationID),
						tQualificationsResults.Status.EQ(
							jet.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
						),
					),
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
			tEmails.Deactivated,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.EmailChanged,
			tEmails.Label,
		).
		FROM(
			tEmails.
				LEFT_JOIN(tEmailsAccess,
					tEmailsAccess.TargetID.EQ(tEmails.ID).
						AND(tEmailsAccess.Access.GT_EQ(jet.Int32(int32(mailer.AccessLevel_ACCESS_LEVEL_READ)))),
				).
				LEFT_JOIN(tQualificationsResults,
					tQualificationsResults.QualificationID.EQ(tEmailsAccess.QualificationID).
						AND(tQualificationsResults.UserID.EQ(jet.Int32(userInfo.GetUserId()))),
				),
		).
		WHERE(condition)

	if pag != nil {
		stmt = stmt.
			OFFSET(pag.GetOffset())
	}

	stmt = stmt.
		GROUP_BY(tEmails.ID).
		ORDER_BY(tEmails.Job.ASC(), tEmails.Label.ASC()).
		LIMIT(listEmailsPageSize)

	resp := &pbmailer.ListEmailsResponse{}
	if err := stmt.QueryContext(ctx, tx, &resp.Emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return resp.GetEmails(), nil
}

func (s *Server) getEmailByCondition(
	ctx context.Context,
	tx qrm.DB,
	condition jet.BoolExpression,
) (*mailer.Email, error) {
	stmt := tEmails.
		SELECT(
			tEmails.ID,
			tEmails.CreatedAt,
			tEmails.UpdatedAt,
			tEmails.DeletedAt,
			tEmails.Deactivated,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.EmailChanged,
			tEmails.Label,
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

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) getEmail(
	ctx context.Context,
	emailId uint64,
	withAccess bool,
	withSettings bool,
) (*mailer.Email, error) {
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
		settings, err := s.getEmailSettings(ctx, s.db, emailId)
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

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "GetEmail",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_READ,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, req.GetId(), true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return &pbmailer.GetEmailResponse{
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

func (s *Server) CreateOrUpdateEmail(
	ctx context.Context,
	req *pbmailer.CreateOrUpdateEmailRequest,
) (*pbmailer.CreateOrUpdateEmailResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateEmail",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if req.Email.UserId != nil {
		req.Email.UserId = &userInfo.UserId
		req.Email.Job = nil
	} else {
		req.Email.UserId = nil
		req.Email.Job = &userInfo.Job

		// Field Permission Check
		fields, err := s.ps.AttrStringList(userInfo, permsmailer.MailerServicePerm, permsmailer.MailerServiceCreateOrUpdateEmailPerm, permsmailer.MailerServiceCreateOrUpdateEmailFieldsPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if !fields.Contains("Job") {
			return nil, errswrap.NewError(err, errorsmailer.ErrEmailAccessDenied)
		}
	}

	if err := s.validateEmail(ctx, userInfo, req.GetEmail().GetEmail(), req.Email.Job != nil); err != nil {
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
			email, err := s.getEmailByCondition(
				ctx,
				tx,
				tEmails.UserID.EQ(jet.Int32(req.GetEmail().GetUserId())),
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
			return nil, err
		}

		req.Email.Id = lastId
	} else {
		check, err := s.access.CanUserAccessTarget(ctx, req.GetEmail().GetId(), userInfo, mailer.AccessLevel_ACCESS_LEVEL_MANAGE)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		if !check {
			return nil, errorsmailer.ErrNoPerms
		}

		email, err := s.getEmail(ctx, req.GetEmail().GetId(), false, false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		if !userInfo.GetSuperuser() && email.GetDeactivated() {
			return nil, errorsmailer.ErrEmailDisabled
		}

		tEmails := table.FivenetMailerEmails

		label := jet.NULL
		if req.Email.Label != nil && *req.Email.Label != "" {
			label = jet.String(req.GetEmail().GetLabel())
		}

		sets := []any{
			tEmails.Label.SET(jet.StringExp(label)),
			tEmails.CreatorID.SET(jet.Int32(userInfo.GetUserId())),
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
				tEmails.Email.SET(jet.String(req.GetEmail().GetEmail())),
				tEmails.EmailChanged.SET(jet.CURRENT_TIMESTAMP()),
			)
		}

		if userInfo.GetSuperuser() {
			sets = append(sets,
				tEmails.Deactivated.SET(jet.Bool(req.GetEmail().GetDeactivated())),
			)
		}

		condition := tEmails.ID.EQ(jet.Uint64(req.GetEmail().GetId()))
		if req.Email.Job != nil {
			condition = condition.AND(tEmails.Job.EQ(jet.String(userInfo.GetJob())))
		} else {
			condition = condition.AND(tEmails.UserID.EQ(jet.Int32(userInfo.GetUserId())))
		}

		stmt := tEmails.
			UPDATE().
			SET(
				sets[0],
				sets[1:]...,
			).
			WHERE(jet.AND(
				tEmails.ID.EQ(jet.Uint64(req.GetEmail().GetId())),
				condition,
			))

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
					tUserProps.Email.SET(jet.String(email.GetEmail())),
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

		if _, err := s.access.HandleAccessChanges(ctx, tx, req.GetEmail().GetId(), req.GetEmail().GetAccess().GetJobs(), req.GetEmail().GetAccess().GetUsers(), req.GetEmail().GetAccess().GetQualifications()); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	resp := &pbmailer.CreateOrUpdateEmailResponse{}
	resp.Email, err = s.getEmail(ctx, req.GetEmail().GetId(), true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_EmailUpdate{
			EmailUpdate: resp.GetEmail(),
		},
	},
		resp.GetEmail().GetId(),
	)

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return resp, nil
}

func (s *Server) createEmail(
	ctx context.Context,
	tx qrm.DB,
	email *mailer.Email,
	userInfo *userinfo.UserInfo,
) (uint64, error) {
	tEmails := table.FivenetMailerEmails
	stmt := tEmails.
		INSERT(
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.Label,
			tEmails.CreatorID,
		).
		VALUES(
			email.GetJob(),
			email.GetUserId(),
			email.GetEmail(),
			email.Label,
			userInfo.GetUserId(),
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

	// Update user email in the user props
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
				tUserProps.Email.SET(jet.String(email.GetEmail())),
			)

		if _, err := upStmt.ExecContext(ctx, tx); err != nil {
			return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return uint64(lastId), nil
}

func (s *Server) DeleteEmail(
	ctx context.Context,
	req *pbmailer.DeleteEmailRequest,
) (*pbmailer.DeleteEmailResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteEmail",
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
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, req.GetId(), false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Make sure that the user is not deleting their own personal email
	if email.Job == nil && email.UserId != nil {
		return nil, errorsmailer.ErrCantDeleteOwnEmail
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_EmailDelete{
			EmailDelete: req.GetId(),
		},
	},
		req.GetId(),
	)

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if email != nil && email.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	tEmails := table.FivenetMailerEmails
	stmt := tEmails.
		UPDATE().
		SET(
			tEmails.DeletedAt.SET(deletedAtTime),
		).
		WHERE(jet.AND(
			tEmails.ID.EQ(jet.Uint64(req.GetId())),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbmailer.DeleteEmailResponse{}, nil
}

func (s *Server) GetEmailProposals(
	ctx context.Context,
	req *pbmailer.GetEmailProposalsRequest,
) (*pbmailer.GetEmailProposalsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.UserId != nil && userInfo.GetSuperuser() {
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
