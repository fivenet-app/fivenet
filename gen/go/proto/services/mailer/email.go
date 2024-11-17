package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tEmails      = table.FivenetMailerEmails.AS("email")
	tEmailsShort = tEmails.AS("email_short")

	tEmailsUserAccess = table.FivenetMailerEmailsUserAccess
	tEmailsJobAccess  = table.FivenetMailerEmailsJobAccess
)

func (s *Server) ListEmails(ctx context.Context, req *ListEmailsRequest) (*ListEmailsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.Bool(true)
	if !userInfo.SuperUser {
		condition = condition.AND(
			jet.OR(
				tEmailsShort.UserID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tEmailsUserAccess.Access.IS_NULL(),
					tEmailsJobAccess.Access.IS_NOT_NULL(),
					tEmailsJobAccess.Access.GT_EQ(jet.Int32(int32(mailer.AccessLevel_ACCESS_LEVEL_VIEW))),
				),
				jet.AND(
					tEmailsUserAccess.Access.IS_NOT_NULL(),
					tEmailsUserAccess.Access.GT_EQ(jet.Int32(int32(mailer.AccessLevel_ACCESS_LEVEL_VIEW))),
				),
			),
		)
	}

	stmt := tEmailsShort.
		SELECT(
			tEmailsShort.ID,
			tEmailsShort.CreatedAt,
			tEmailsShort.UpdatedAt,
			tEmailsShort.DeletedAt,
			tEmailsShort.Job,
			tEmailsShort.UserID,
			tEmailsShort.Email,
			tEmailsShort.Label,
			tEmailsShort.Internal,
		).
		FROM(
			tEmailsShort.
				LEFT_JOIN(tEmailsJobAccess,
					tEmailsJobAccess.EmailID.EQ(tEmailsShort.ID).
						AND(tEmailsJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tEmailsJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tEmailsUserAccess,
					tEmailsUserAccess.EmailID.EQ(tEmailsShort.ID).
						AND(tEmailsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				),
		).
		WHERE(condition).
		GROUP_BY(tEmailsShort.ID).
		ORDER_BY(tEmailsShort.Job.ASC(), tEmailsShort.Label.ASC())

	resp := &ListEmailsResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) getEmail(ctx context.Context, id uint64, withAccess bool) (*mailer.Email, error) {
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
			tEmails.Signature,
		).
		FROM(tEmails).
		WHERE(
			tEmails.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	dest := &mailer.Email{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	if withAccess {
		access, err := s.getEmailAccess(ctx, id)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		dest.Access = access
	}

	return dest, nil
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

	check, err := s.access.CanUserAccessTarget(ctx, req.Id, userInfo, mailer.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	resp := &GetEmailResponse{}
	resp.Email, err = s.getEmail(ctx, req.Id, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) getEmailAccess(ctx context.Context, id uint64) (*mailer.Access, error) {
	access := &mailer.Access{}

	jobsAccess, err := s.access.Jobs.List(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	access.Jobs = jobsAccess

	usersAccess, err := s.access.Users.List(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	access.Users = usersAccess

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

	// TODO validate email format

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if req.Email.Id <= 0 {
		if req.Email.UserId != nil {
			req.Email.UserId = &userInfo.UserId
			req.Email.Job = nil
		} else {
			req.Email.UserId = nil
			req.Email.Job = &userInfo.Job
		}

		lastId, err := s.createEmail(ctx, req.Email)
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

		stmt := tEmails.
			UPDATE(
				tEmails.Label,
				tEmails.Internal,
				tEmails.Signature,
			).
			SET(
				label,
				jet.Bool(req.Email.Internal),
				jet.String(*req.Email.Signature),
			).
			WHERE(jet.AND(
				tEmails.ID.EQ(jet.Uint64(req.Email.Id)),
				tEmails.Job.EQ(jet.String(userInfo.Job)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if _, err := s.access.HandleAccessChanges(ctx, tx, req.Email.Id, req.Email.Access.Jobs, req.Email.Access.Users); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	resp := &CreateOrUpdateEmailResponse{}
	resp.Email, err = s.getEmail(ctx, req.Email.Id, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return resp, nil
}

func (s *Server) createEmail(ctx context.Context, email *mailer.Email) (uint64, error) {
	stmt := tEmails.
		INSERT(
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.Label,
			tEmails.Internal,
			tEmails.Signature,
		).
		VALUES(
			email.Job,
			email.UserId,
			email.Email,
			email.Label,
			email.Internal,
			email.Signature,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
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

	stmt := tEmails.
		DELETE().
		WHERE(jet.AND(
			tEmails.ID.EQ(jet.Uint64(req.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteEmailResponse{}, nil
}
