package mailer

import (
	"context"
	"errors"
	"slices"

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
	tSettings       = table.FivenetMailerSettings
	tSettingsBlocks = table.FivenetMailerSettingsBlocked
)

func (s *Server) GetEmailSettings(ctx context.Context, req *GetEmailSettingsRequest) (*GetEmailSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.EmailId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	settings, err := s.getEmailSettings(ctx, req.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &GetEmailSettingsResponse{
		Settings: settings,
	}, nil
}

func (s *Server) getEmailSettings(ctx context.Context, emailId uint64) (*mailer.EmailSettings, error) {
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
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Server) SetEmailSettings(ctx context.Context, req *SetEmailSettingsRequest) (*SetEmailSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
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

	signature := jet.StringExp(jet.NULL)
	if req.Settings.Signature != nil {
		signature = jet.String(*req.Settings.Signature)
	}

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

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	settings, err := s.getEmailSettings(ctx, req.Settings.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Handle blocked users changes
	if len(req.Settings.BlockedEmails) == 0 {
		if len(settings.BlockedEmails) > 0 {
			stmt := tSettingsBlocks.
				DELETE().
				WHERE(tSettingsBlocks.EmailID.EQ(jet.Int32(userInfo.UserId)))

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
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

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
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

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
			}
		}
	}

	settings, err = s.getEmailSettings(ctx, req.Settings.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_EmailSettingsUpdated{
			EmailSettingsUpdated: settings,
		},
	}, req.Settings.EmailId)

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &SetEmailSettingsResponse{
		Settings: settings,
	}, nil
}
