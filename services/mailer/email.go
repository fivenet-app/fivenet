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
	"github.com/go-jet/jet/v2/mysql"
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

	condition := mysql.Bool(true)

	if !userInfo.GetSuperuser() || (userInfo.GetSuperuser() && req.All != nil && !req.GetAll()) {
		// Include deactivated e-mails
		condition = condition.AND(mysql.AND(
			tEmails.DeletedAt.IS_NULL(),
			mysql.OR(
				tEmails.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				mysql.OR(
					tEmailsAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					mysql.AND(
						tEmailsAccess.Job.EQ(mysql.String(userInfo.GetJob())),
						tEmailsAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					),
					mysql.AND(
						tEmailsAccess.QualificationID.IS_NOT_NULL(),
						tQualificationsResults.DeletedAt.IS_NULL(),
						tQualificationsResults.QualificationID.EQ(tEmailsAccess.QualificationID),
						tQualificationsResults.Status.EQ(
							mysql.Int32(
								int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL),
							),
						),
					),
				),
			),
		))
	}

	countStmt := tEmails.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tEmails.ID)).AS("data_count.total"),
		).
		FROM(
			tEmails.
				LEFT_JOIN(tEmailsAccess,
					tEmailsAccess.TargetID.EQ(tEmails.ID).
						AND(tEmailsAccess.Access.GT_EQ(mysql.Int32(int32(mailer.AccessLevel_ACCESS_LEVEL_READ)))),
				).
				LEFT_JOIN(tQualificationsResults,
					tQualificationsResults.QualificationID.EQ(tEmailsAccess.QualificationID).
						AND(tQualificationsResults.UserID.EQ(mysql.Int32(userInfo.GetUserId()))),
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

	return resp, nil
}

func ListUserEmails(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	pag *database.PaginationRequest,
	includeDisabled bool,
) ([]*mailer.Email, error) {
	condition := mysql.Bool(true)
	baseCondition := tEmails.DeletedAt.IS_NULL()
	if !includeDisabled {
		baseCondition = baseCondition.AND(tEmails.Deactivated.IS_FALSE())
	}

	if !userInfo.GetSuperuser() {
		access := int32(mailer.AccessLevel_ACCESS_LEVEL_READ)
		// Access predicates via EXISTS (no joins)
		userAccessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tEmailsAccess).
				WHERE(
					tEmailsAccess.TargetID.EQ(tEmails.ID).
						AND(tEmailsAccess.Access.GT_EQ(mysql.Int32(access))).
						AND(tEmailsAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId()))),
				),
		)

		jobAccessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tEmailsAccess).
				WHERE(
					tEmailsAccess.TargetID.EQ(tEmails.ID).
						AND(tEmailsAccess.Access.GT_EQ(mysql.Int32(access))).
						AND(tEmailsAccess.Job.EQ(mysql.String(userInfo.GetJob()))).
						AND(tEmailsAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade()))),
				),
		)

		// Qualification-based access: there exists an access row with a QualificationID
		// for this email AND the user has a successful (non-deleted) result for it.
		ea := tEmailsAccess.AS("ea")
		qr := tQualificationsResults.AS("qr")

		qualAccessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tEmailsAccess.AS("ea")).
				WHERE(mysql.AND(
					ea.TargetID.EQ(tEmails.ID),
					ea.Access.GT_EQ(mysql.Int32(access)),
					ea.QualificationID.IS_NOT_NULL(),
					mysql.EXISTS(
						mysql.
							SELECT(mysql.Int(1)).
							FROM(tQualificationsResults.AS("qr")).
							WHERE(
								mysql.AND(
									qr.QualificationID.EQ(ea.QualificationID),
									qr.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
									qr.DeletedAt.IS_NULL(),
									qr.Status.EQ(
										mysql.Int32(
											int32(
												qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL,
											),
										),
									),
								),
							),
					),
				),
				),
		)

		condition = condition.AND(
			baseCondition.AND(
				mysql.OR(
					// owner may always see their email
					tEmails.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					// access by explicit user, by job+grade, or by qualification
					userAccessExists,
					jobAccessExists,
					qualAccessExists,
				),
			),
		)
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
		FROM(tEmails).
		WHERE(condition).
		ORDER_BY(
			tEmails.Job.ASC(),
			tEmails.Label.ASC(),
		).
		LIMIT(listEmailsPageSize)

	if pag != nil {
		stmt = stmt.
			OFFSET(pag.GetOffset())
	}

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
	condition mysql.BoolExpression,
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
	emailId int64,
	withAccess bool,
	withSettings bool,
) (*mailer.Email, error) {
	email, err := s.getEmailByCondition(ctx, s.db, tEmails.ID.EQ(mysql.Int64(emailId)))
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

func (s *Server) getEmailAccess(ctx context.Context, emailId int64) (*mailer.Access, error) {
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
				tEmails.UserID.EQ(mysql.Int32(req.GetEmail().GetUserId())),
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

		label := mysql.NULL
		if req.Email.Label != nil && *req.Email.Label != "" {
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

		if userInfo.GetSuperuser() {
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
) (int64, error) {
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

	deletedAtTime := mysql.CURRENT_TIMESTAMP()
	if email != nil && email.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = mysql.TimestampExp(mysql.NULL)
	}

	tEmails := table.FivenetMailerEmails
	stmt := tEmails.
		UPDATE().
		SET(
			tEmails.DeletedAt.SET(deletedAtTime),
		).
		WHERE(mysql.AND(
			tEmails.ID.EQ(mysql.Int64(req.GetId())),
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
