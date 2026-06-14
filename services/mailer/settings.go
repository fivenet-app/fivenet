package mailer

import (
	"context"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	mailerevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/events"
	mailersettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/settings"
	pbmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
	"github.com/go-jet/jet/v2/qrm"
)

const SignatureMaxLength = 1024

func (s *Server) GetEmailSettings(
	ctx context.Context,
	req *pbmailer.GetEmailSettingsRequest,
) (*pbmailer.GetEmailSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetEmailId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
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
) (*mailersettings.EmailSettings, error) {
	settings, err := s.store.GetEmailSettings(ctx, tx, emailId)
	if err != nil {
		return nil, err
	}

	return settings, nil
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
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
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

	var signature *content.Content
	if req.GetSettings().GetSignature() != nil {
		// Check comment length
		extracted := req.GetSettings().GetSignature().Extract()
		if len(extracted.GetText()) > SignatureMaxLength {
			return nil, errorsmailer.ErrSignatureTooLong
		}

		signature = req.GetSettings().GetSignature()
	} else if email.GetSettings() != nil {
		signature = email.GetSettings().GetSignature()
	}

	// Make all emails lowercase, remove own email, and remove duplicates
	blockedEmails := make([]string, 0, len(req.GetSettings().GetBlockedEmails()))
	for _, be := range req.GetSettings().GetBlockedEmails() {
		blockedEmails = append(blockedEmails, strings.ToLower(be))
	}
	blockedEmails = slices.DeleteFunc(
		blockedEmails,
		func(e string) bool {
			return e == email.GetEmail()
		},
	)
	blockedEmails = utils.RemoveSliceDuplicates(blockedEmails)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.UpsertEmailSettingsSignature(
		ctx,
		tx,
		req.GetSettings().GetEmailId(),
		signature,
	); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	settings, err := s.getEmailSettings(ctx, tx, req.GetSettings().GetEmailId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Handle blocked users changes
	if len(blockedEmails) == 0 {
		if len(settings.GetBlockedEmails()) > 0 {
			if err := s.store.DeleteBlockedEmails(
				ctx,
				tx,
				req.GetSettings().GetEmailId(),
				settings.GetBlockedEmails(),
			); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}
	} else {
		toCreate, toDelete := utils.SlicesDifference(blockedEmails, settings.GetBlockedEmails())

		if len(toCreate) > 0 {
			if err := s.store.AddBlockedEmails(
				ctx,
				tx,
				req.GetSettings().GetEmailId(),
				toCreate,
			); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}

		if len(toDelete) > 0 {
			if err := s.store.DeleteBlockedEmails(
				ctx,
				tx,
				req.GetSettings().GetEmailId(),
				toDelete,
			); err != nil {
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

	s.sendUpdate(ctx, &mailerevents.MailerEvent{
		Data: &mailerevents.MailerEvent_EmailSettingsUpdated{
			EmailSettingsUpdated: settings,
		},
	}, req.GetSettings().GetEmailId())

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbmailer.SetEmailSettingsResponse{
		Settings: settings,
	}, nil
}
