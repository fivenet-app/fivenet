package rector

import (
	"context"
	"errors"
	"strings"

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

	// Ensure that the job is the user's job
	req.JobProps.Job = userInfo.Job

	req.JobProps.LivemapMarkerColor = strings.ToLower(req.JobProps.LivemapMarkerColor)

	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.Theme,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordSyncSettings,
		).
		VALUES(
			req.JobProps.Job,
			req.JobProps.Theme,
			req.JobProps.LivemapMarkerColor,
			req.JobProps.RadioFrequency,
			req.JobProps.QuickButtons,
			req.JobProps.DiscordGuildId,
			req.JobProps.DiscordSyncSettings,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.Theme.SET(jet.String(req.JobProps.Theme)),
			tJobProps.LivemapMarkerColor.SET(jet.String(req.JobProps.LivemapMarkerColor)),
			tJobProps.RadioFrequency.SET(jet.StringExp(jet.Raw("VALUES(`radio_frequency`)"))),
			tJobProps.QuickButtons.SET(jet.StringExp(jet.Raw("VALUES(`quick_buttons`)"))),
			tJobProps.DiscordGuildID.SET(jet.IntExp(jet.Raw("VALUES(`discord_guild_id`)"))),
			tJobProps.DiscordSyncSettings.SET(jet.StringExp(jet.Raw("VALUES(`discord_sync_settings`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(errorsrector.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &SetJobPropsResponse{}, nil
}
