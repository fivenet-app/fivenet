package mailer

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tSettings       = table.FivenetMailerSettings
	tSettingsBlocks = table.FivenetMailerSettingsBlocked
)

func (s *Server) GetEmailSettings(
	ctx context.Context,
	req *pbmailer.GetEmailSettingsRequest,
) (*pbmailer.GetEmailSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetEmailId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_MANAGE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	settings, err := s.getEmailSettings(ctx, s.db, req.GetEmailId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.GetEmailSettingsResponse{
		Settings: settings,
	}, nil
}

func (s *Server) getEmailSettings(
	ctx context.Context,
	tx qrm.DB,
	emailId int64,
) (*mailer.EmailSettings, error) {
	tSettings := tSettings.AS("email_settings")
	stmt := tSettings.
		SELECT(
			tSettings.EmailID,
			tSettings.Signature,
			tSettingsBlocks.TargetEmail.AS("email_settings.blocked_emails"),
		).
		FROM(
			tSettings.
				LEFT_JOIN(tSettingsBlocks,
					tSettingsBlocks.EmailID.EQ(tSettings.EmailID),
				),
		).
		WHERE(
			tSettings.EmailID.EQ(mysql.Int64(emailId)),
		).
		LIMIT(25)

	dest := &mailer.EmailSettings{
		EmailId: emailId,
	}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Server) SetEmailSettings(
	ctx context.Context,
	req *pbmailer.SetEmailSettingsRequest,
) (*pbmailer.SetEmailSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetSettings().GetEmailId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_MANAGE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, req.GetSettings().GetEmailId(), false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	signature := mysql.StringExp(mysql.NULL)
	if req.Settings.Signature != nil && *req.Settings.Signature != "" {
		signature = mysql.String(req.GetSettings().GetSignature())
	}

	// Make all emails lowercase, remove own email, and remove duplicates
	for idx := range req.GetSettings().GetBlockedEmails() {
		req.Settings.BlockedEmails[idx] = strings.ToLower(req.GetSettings().GetBlockedEmails()[idx])
	}
	req.Settings.BlockedEmails = slices.DeleteFunc(
		req.GetSettings().GetBlockedEmails(),
		func(e string) bool {
			return e == email.GetEmail()
		},
	)
	utils.RemoveSliceDuplicates(req.GetSettings().GetBlockedEmails())

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	stmt := tSettings.
		INSERT(
			tSettings.EmailID,
			tSettings.Signature,
		).
		VALUES(
			req.GetSettings().GetEmailId(),
			signature,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tSettings.Signature.SET(signature),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	settings, err := s.getEmailSettings(ctx, tx, req.GetSettings().GetEmailId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Handle blocked users changes
	if len(req.GetSettings().GetBlockedEmails()) == 0 {
		if len(settings.GetBlockedEmails()) > 0 {
			stmt := tSettingsBlocks.
				DELETE().
				WHERE(tSettingsBlocks.EmailID.EQ(mysql.Int32(userInfo.GetUserId())))

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}
	} else {
		toCreate := []string{}
		toDelete := []string{}

		for _, be := range req.GetSettings().GetBlockedEmails() {
			if !slices.ContainsFunc(settings.GetBlockedEmails(), func(a string) bool {
				return a == be
			}) {
				toCreate = append(toCreate, be)
			}
		}

		for _, be := range settings.GetBlockedEmails() {
			if !slices.ContainsFunc(req.GetSettings().GetBlockedEmails(), func(a string) bool {
				return a == be
			}) {
				toDelete = append(toDelete, be)
			}
		}

		if len(toCreate) > 0 {
			stmt := tSettingsBlocks.
				INSERT(
					tSettingsBlocks.EmailID,
					tSettingsBlocks.TargetEmail,
				)

			for _, be := range toCreate {
				stmt = stmt.
					VALUES(
						req.GetSettings().GetEmailId(),
						be,
					)
			}

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}

		if len(toDelete) > 0 {
			targets := []mysql.Expression{}

			stmt := tSettingsBlocks.
				DELETE().
				WHERE(mysql.AND(
					tSettingsBlocks.EmailID.EQ(mysql.Int64(req.GetSettings().GetEmailId())),
					tSettingsBlocks.TargetEmail.IN(targets...),
				))

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	settings, err = s.getEmailSettings(ctx, s.db, req.GetSettings().GetEmailId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_EmailSettingsUpdated{
			EmailSettingsUpdated: settings,
		},
	}, req.GetSettings().GetEmailId())

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbmailer.SetEmailSettingsResponse{
		Settings: settings,
	}, nil
}
