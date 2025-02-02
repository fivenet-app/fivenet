package rector

import (
	"context"
	"fmt"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/filestore"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorsrector "github.com/fivenet-app/fivenet/services/rector/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

var tJobProps = table.FivenetJobProps

func (s *Server) GetJobProps(ctx context.Context, req *pbrector.GetJobPropsRequest) (*pbrector.GetJobPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobProps, err := s.getJobProps(ctx, userInfo.Job)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}

	return &pbrector.GetJobPropsResponse{
		JobProps: jobProps,
	}, nil
}

func (s *Server) getJobProps(ctx context.Context, job string) (*users.JobProps, error) {
	props, err := users.GetJobProps(ctx, s.db, job)
	if err != nil {
		return nil, err
	}

	s.enricher.EnrichJobName(props)

	return props, nil
}

func (s *Server) SetJobProps(ctx context.Context, req *pbrector.SetJobPropsRequest) (*pbrector.SetJobPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "SetJobProps",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	jobProps, err := s.getJobProps(ctx, userInfo.Job)
	if err != nil {
		return nil, err
	}

	// Ensure that the job is the user's job
	req.JobProps.Job = userInfo.Job
	req.JobProps.LivemapMarkerColor = strings.ToLower(req.JobProps.LivemapMarkerColor)

	if req.JobProps.LogoUrl != nil {
		// Set "current" image's url so the system will delete it if still exists
		if jobProps != nil && jobProps.LogoUrl != nil {
			req.JobProps.LogoUrl.Url = jobProps.LogoUrl.Url
		}

		if len(req.JobProps.LogoUrl.Data) > 0 {
			if !req.JobProps.LogoUrl.IsImage() {
				return nil, errorsrector.ErrFailedQuery
			}

			if err := req.JobProps.LogoUrl.Optimize(ctx); err != nil {
				return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
			}

			if err := req.JobProps.LogoUrl.Upload(ctx, s.st, filestore.JobLogos, userInfo.Job); err != nil {
				return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
			}
		} else if req.JobProps.LogoUrl.Delete != nil && *req.JobProps.LogoUrl.Delete {
			// Delete avatar from store
			if jobProps.LogoUrl != nil && jobProps.LogoUrl.Url != nil {
				if err := s.st.Delete(ctx, filestore.StripURLPrefix(*jobProps.LogoUrl.Url)); err != nil {
					return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
				}
			}
		}
	} else {
		req.JobProps.LogoUrl = jobProps.LogoUrl
	}

	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.Theme,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordSyncSettings,
			tJobProps.LogoURL,
			tJobProps.Settings,
		).
		VALUES(
			req.JobProps.Job,
			req.JobProps.Theme,
			req.JobProps.LivemapMarkerColor,
			req.JobProps.RadioFrequency,
			req.JobProps.QuickButtons,
			req.JobProps.DiscordGuildId,
			req.JobProps.DiscordSyncSettings,
			req.JobProps.LogoUrl,
			req.JobProps.Settings,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.Theme.SET(jet.String(req.JobProps.Theme)),
			tJobProps.LivemapMarkerColor.SET(jet.String(req.JobProps.LivemapMarkerColor)),
			tJobProps.RadioFrequency.SET(jet.StringExp(jet.Raw("VALUES(`radio_frequency`)"))),
			tJobProps.QuickButtons.SET(jet.StringExp(jet.Raw("VALUES(`quick_buttons`)"))),
			tJobProps.DiscordGuildID.SET(jet.StringExp(jet.Raw("VALUES(`discord_guild_id`)"))),
			tJobProps.DiscordSyncSettings.SET(jet.StringExp(jet.Raw("VALUES(`discord_sync_settings`)"))),
			tJobProps.LogoURL.SET(jet.StringExp(jet.Raw("VALUES(`logo_url`)"))),
			tJobProps.LogoURL.SET(jet.StringExp(jet.Raw("VALUES(`settings`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	newJobProps, err := s.getJobProps(ctx, userInfo.Job)
	if err != nil {
		return nil, err
	}

	if !proto.Equal(req.JobProps, jobProps) {
		if _, err := s.js.PublishAsyncProto(ctx,
			fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.Job),
			&notifications.JobEvent{
				Data: &notifications.JobEvent_JobProps{
					JobProps: newJobProps,
				},
			}); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	return &pbrector.SetJobPropsResponse{
		JobProps: newJobProps,
	}, nil
}
