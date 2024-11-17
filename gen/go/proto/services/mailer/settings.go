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
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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
	stmt := tSettingsBlocks.
		SELECT(
			tSettingsBlocks.EmailID,
			tSettingsBlocks.TargetEmail,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
		).
		FROM(
			tSettingsBlocks.
				INNER_JOIN(tUsers,
					tUsers.ID.EQ(tSettingsBlocks.EmailID),
				),
		).
		WHERE(
			tSettingsBlocks.EmailID.EQ(jet.Uint64(emailId)),
		).
		LIMIT(25)

	dest := &mailer.EmailSettings{}
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
		toUpdate := []string{}

		for _, be := range req.Settings.BlockedEmails {
			if slices.ContainsFunc(settings.BlockedEmails, func(a string) bool {
				return a == be
			}) {
				toUpdate = append(toUpdate, be)
			} else {
				toCreate = append(toCreate, be)
			}
		}
	}

	// Handle blocked emails changes
	// TODO

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return nil, nil
}
