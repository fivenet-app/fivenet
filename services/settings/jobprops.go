package settings

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var tJobProps = table.FivenetJobProps

func (s *Server) GetJobProps(
	ctx context.Context,
	req *pbsettings.GetJobPropsRequest,
) (*pbsettings.GetJobPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobProps, err := s.getJobProps(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	return &pbsettings.GetJobPropsResponse{
		JobProps: jobProps,
	}, nil
}

func (s *Server) getJobProps(ctx context.Context, job string) (*jobs.JobProps, error) {
	props, err := jobs.GetJobProps(ctx, s.db, job)
	if err != nil {
		return nil, err
	}

	s.enricher.EnrichJobName(props)

	return props, nil
}

func (s *Server) SetJobProps(
	ctx context.Context,
	req *pbsettings.SetJobPropsRequest,
) (*pbsettings.SetJobPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "SetJobProps",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	jobProps, err := s.getJobProps(ctx, userInfo.GetJob())
	if err != nil {
		return nil, err
	}

	// Ensure that the job is the user's job
	req.JobProps.Job = userInfo.GetJob()
	req.JobProps.LivemapMarkerColor = strings.ToLower(req.GetJobProps().GetLivemapMarkerColor())

	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordSyncSettings,
			tJobProps.Settings,
		).
		VALUES(
			req.GetJobProps().GetJob(),
			req.GetJobProps().GetLivemapMarkerColor(),
			req.GetJobProps().RadioFrequency,
			req.GetJobProps().GetQuickButtons(),
			req.GetJobProps().DiscordGuildId,
			req.GetJobProps().GetDiscordSyncSettings(),
			req.GetJobProps().GetSettings(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.LivemapMarkerColor.SET(jet.String(req.GetJobProps().GetLivemapMarkerColor())),
			tJobProps.RadioFrequency.SET(jet.StringExp(jet.Raw("VALUES(`radio_frequency`)"))),
			tJobProps.QuickButtons.SET(jet.StringExp(jet.Raw("VALUES(`quick_buttons`)"))),
			tJobProps.DiscordGuildID.SET(jet.StringExp(jet.Raw("VALUES(`discord_guild_id`)"))),
			tJobProps.DiscordSyncSettings.SET(jet.StringExp(jet.Raw("VALUES(`discord_sync_settings`)"))),
			tJobProps.Settings.SET(jet.StringExp(jet.Raw("VALUES(`settings`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	newJobProps, err := s.getJobProps(ctx, userInfo.GetJob())
	if err != nil {
		return nil, err
	}

	if !proto.Equal(req.GetJobProps(), jobProps) {
		if _, err := s.js.PublishAsyncProto(ctx,
			fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.GetJob()),
			&notifications.JobEvent{
				Data: &notifications.JobEvent_JobProps{
					JobProps: newJobProps,
				},
			}); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	return &pbsettings.SetJobPropsResponse{
		JobProps: newJobProps,
	}, nil
}

func (s *Server) UploadJobLogo(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "UploadJobLogo",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, nil)

	props, err := s.getJobProps(ctx, userInfo.GetJob())
	if err != nil {
		return errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if props.LogoFileId != nil && props.GetLogoFileId() > 0 {
		if err := s.jobPropsFileHandler.Delete(ctx, userInfo.GetJob(), props.GetLogoFileId()); err != nil {
			return errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	meta, err := s.jobPropsFileHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	name := filepath.Base(meta.GetOriginalName())
	ext := filepath.Ext(name)
	key := fmt.Sprintf("joblogos/%s%s", userInfo.GetJob(), ext)

	resp, err := s.jobPropsFileHandler.UploadFile(
		ctx,
		userInfo.GetJob(),
		key,
		meta.GetSize(),
		meta.GetContentType(),
		srv,
	)
	if err != nil {
		return err
	}

	if resp.GetId() != props.GetLogoFileId() {
		newJobProps, err := s.getJobProps(ctx, userInfo.GetJob())
		if err != nil {
			return errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		if _, err := s.js.PublishAsyncProto(ctx,
			fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.GetJob()),
			&notifications.JobEvent{
				Data: &notifications.JobEvent_JobProps{
					JobProps: newJobProps,
				},
			}); err != nil {
			return errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return nil
}

func (s *Server) DeleteJobLogo(
	ctx context.Context,
	req *pbsettings.DeleteJobLogoRequest,
) (*pbsettings.DeleteJobLogoResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "DeleteJobLogo",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, nil)

	props, err := s.getJobProps(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if props.LogoFileId == nil || props.GetLogoFileId() == 0 {
		return &pbsettings.DeleteJobLogoResponse{}, nil
	}

	if err := s.jobPropsFileHandler.Delete(ctx, userInfo.GetJob(), props.GetLogoFileId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	newJobProps, err := s.getJobProps(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if _, err := s.js.PublishAsyncProto(ctx,
		fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.GetJob()),
		&notifications.JobEvent{
			Data: &notifications.JobEvent_JobProps{
				JobProps: newJobProps,
			},
		}); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return &pbsettings.DeleteJobLogoResponse{}, nil
}
