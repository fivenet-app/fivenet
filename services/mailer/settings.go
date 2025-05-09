package mailer

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tSettings       = table.FivenetMailerSettings
	tSettingsBlocks = table.FivenetMailerSettingsBlocked
)

func (s *Server) GetEmailSettings(ctx context.Context, req *pbmailer.GetEmailSettingsRequest) (*pbmailer.GetEmailSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.EmailId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	settings, err := s.getEmailSettings(ctx, s.db, req.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.GetEmailSettingsResponse{
		Settings: settings,
	}, nil
}

func (s *Server) getEmailSettings(ctx context.Context, tx qrm.DB, emailId uint64) (*mailer.EmailSettings, error) {
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
			tSettings.EmailID.EQ(jet.Uint64(emailId)),
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

func (s *Server) SetEmailSettings(ctx context.Context, req *pbmailer.SetEmailSettingsRequest) (*pbmailer.SetEmailSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "SetEmailSettings",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Settings.EmailId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	email, err := s.getEmail(ctx, req.Settings.EmailId, false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	signature := jet.StringExp(jet.NULL)
	if req.Settings.Signature != nil {
		signature = jet.String(*req.Settings.Signature)
	}

	// Make all emails lowercase, remove own email, and remove duplicates
	for idx := range req.Settings.BlockedEmails {
		req.Settings.BlockedEmails[idx] = strings.ToLower(req.Settings.BlockedEmails[idx])
	}
	req.Settings.BlockedEmails = slices.DeleteFunc(req.Settings.BlockedEmails, func(e string) bool {
		return e == email.Email
	})
	utils.RemoveSliceDuplicates(req.Settings.BlockedEmails)

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
			req.Settings.EmailId,
			req.Settings.Signature,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tSettings.Signature.SET(signature),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	settings, err := s.getEmailSettings(ctx, tx, req.Settings.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Handle blocked users changes
	if len(req.Settings.BlockedEmails) == 0 {
		if len(settings.BlockedEmails) > 0 {
			stmt := tSettingsBlocks.
				DELETE().
				WHERE(tSettingsBlocks.EmailID.EQ(jet.Int32(userInfo.UserId)))

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}
	} else {
		toCreate := []string{}
		toDelete := []string{}

		for _, be := range req.Settings.BlockedEmails {
			if !slices.ContainsFunc(settings.BlockedEmails, func(a string) bool {
				return a == be
			}) {
				toCreate = append(toCreate, be)
			}
		}

		for _, be := range settings.BlockedEmails {
			if !slices.ContainsFunc(req.Settings.BlockedEmails, func(a string) bool {
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
						req.Settings.EmailId,
						be,
					)
			}

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}

		if len(toDelete) > 0 {
			targets := []jet.Expression{}

			stmt := tSettingsBlocks.
				DELETE().
				WHERE(jet.AND(
					tSettingsBlocks.EmailID.EQ(jet.Uint64(req.Settings.EmailId)),
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

	settings, err = s.getEmailSettings(ctx, s.db, req.Settings.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_EmailSettingsUpdated{
			EmailSettingsUpdated: settings,
		},
	}, req.Settings.EmailId)

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbmailer.SetEmailSettingsResponse{
		Settings: settings,
	}, nil
}
