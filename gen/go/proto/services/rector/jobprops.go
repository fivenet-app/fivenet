package rector

import (
	"context"
	"errors"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/filestore"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	errorsrector "github.com/galexrt/fivenet/gen/go/proto/services/rector/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tJobProps = table.FivenetJobProps
)

func (s *Server) GetJobProps(ctx context.Context, req *GetJobPropsRequest) (*GetJobPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tJobProps := table.FivenetJobProps.AS("jobprops")
	stmt := tJobProps.
		SELECT(
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.Theme,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordLastSync,
			tJobProps.DiscordSyncSettings,
			tJobProps.LogoURL,
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(jet.String(userInfo.Job)),
		).
		LIMIT(1)

	resp := &GetJobPropsResponse{
		JobProps: &users.JobProps{},
	}
	if err := stmt.QueryContext(ctx, s.db, resp.JobProps); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsrector.ErrInvalidRequest, err)
		}
	}

	resp.JobProps.Default(userInfo.Job)

	if resp.JobProps != nil {
		s.enricher.EnrichJobName(resp.JobProps)
	}

	return resp, nil
}
func (s *Server) SetJobProps(ctx context.Context, req *SetJobPropsRequest) (*SetJobPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "SetJobProps",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	jobProps, err := s.GetJobProps(ctx, &GetJobPropsRequest{})
	if err != nil {
		return nil, err
	}

	// Ensure that the job is the user's job
	req.JobProps.Job = userInfo.Job
	req.JobProps.LivemapMarkerColor = strings.ToLower(req.JobProps.LivemapMarkerColor)

	if req.JobProps.LogoUrl != nil && len(req.JobProps.LogoUrl.Data) > 0 {
		if !req.JobProps.LogoUrl.IsImage() {
			return nil, errorsrector.ErrFailedQuery
		}

		// Set "current" image's url so the system will delete it if still exists
		if jobProps.JobProps != nil && jobProps.JobProps.LogoUrl != nil {
			req.JobProps.LogoUrl.Url = jobProps.JobProps.LogoUrl.Url
		}
		if err := req.JobProps.LogoUrl.Upload(ctx, s.st, filestore.JobLogos, userInfo.Job); err != nil {
			return nil, errswrap.NewError(errorsrector.ErrFailedQuery, err)
		}
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
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.Theme.SET(jet.String(req.JobProps.Theme)),
			tJobProps.LivemapMarkerColor.SET(jet.String(req.JobProps.LivemapMarkerColor)),
			tJobProps.RadioFrequency.SET(jet.StringExp(jet.Raw("VALUES(`radio_frequency`)"))),
			tJobProps.QuickButtons.SET(jet.StringExp(jet.Raw("VALUES(`quick_buttons`)"))),
			tJobProps.DiscordGuildID.SET(jet.IntExp(jet.Raw("VALUES(`discord_guild_id`)"))),
			tJobProps.DiscordSyncSettings.SET(jet.StringExp(jet.Raw("VALUES(`discord_sync_settings`)"))),
			tJobProps.LogoURL.SET(jet.StringExp(jet.Raw("VALUES(`logo_url`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(errorsrector.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	jobProps, err = s.GetJobProps(ctx, &GetJobPropsRequest{})
	if err != nil {
		return nil, err
	}

	return &SetJobPropsResponse{
		JobProps: jobProps.JobProps,
	}, nil
}
