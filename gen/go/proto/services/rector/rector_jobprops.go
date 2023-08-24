package rector

import (
	"context"
	"errors"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	jobProps = table.FivenetJobProps
)

func (s *Server) GetJobProps(ctx context.Context, req *GetJobPropsRequest) (*GetJobPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobProps := table.FivenetJobProps.AS("jobprops")
	stmt := jobProps.
		SELECT(
			jobProps.Job,
			jobProps.UpdatedAt,
			jobProps.Theme,
			jobProps.LivemapMarkerColor,
			jobProps.QuickButtons,
			jobProps.DiscordGuildID,
		).
		FROM(jobProps).
		WHERE(
			jobProps.Job.EQ(jet.String(userInfo.Job)),
		).
		LIMIT(1)

	resp := &GetJobPropsResponse{
		JobProps: &users.JobProps{},
	}
	if err := stmt.QueryContext(ctx, s.db, resp.JobProps); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrInvalidRequest
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
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Ensure that the job is the user's job
	req.JobProps.Job = userInfo.Job

	req.JobProps.LivemapMarkerColor = strings.ToLower(strings.ReplaceAll(req.JobProps.LivemapMarkerColor, "#", ""))

	if !s.validateJobPropsQuickButtons(req.JobProps.QuickButtons) {
		return nil, status.Error(codes.InvalidArgument, "Invalid quick access button found!")
	}

	stmt := jobProps.
		INSERT(
			jobProps.Job,
			jobProps.Theme,
			jobProps.LivemapMarkerColor,
			jobProps.QuickButtons,
			jobProps.DiscordGuildID,
		).
		VALUES(
			req.JobProps.Job,
			req.JobProps.Theme,
			req.JobProps.LivemapMarkerColor,
			req.JobProps.QuickButtons,
			req.JobProps.DiscordGuildId,
		).
		ON_DUPLICATE_KEY_UPDATE(
			jobProps.Theme.SET(jet.String(req.JobProps.Theme)),
			jobProps.LivemapMarkerColor.SET(jet.String(req.JobProps.LivemapMarkerColor)),
			jobProps.QuickButtons.SET(jet.String(req.JobProps.QuickButtons)),
			jobProps.DiscordGuildID.SET(jet.IntExp(jet.Raw("VALUES(`discord_guild_id`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &SetJobPropsResponse{}, nil
}

func (s *Server) validateJobPropsQuickButtons(in string) bool {
	if in == "" {
		return true
	}

	for _, comp := range strings.Split(in, ";") {
		if comp != "PenaltyCalculator" {
			return false
		}
	}

	return true
}
