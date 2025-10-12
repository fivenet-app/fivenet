package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	permscentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
)

func (s *Server) GetSettings(
	ctx context.Context,
	req *pbcentrum.GetSettingsRequest,
) (*pbcentrum.GetSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	settings, err := s.settings.Get(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if settings.GetAccess() != nil && settings.GetAccess().GetJobs() != nil {
		for _, ja := range settings.GetAccess().GetJobs() {
			// Lookup job info by using SourceJob field
			j := s.enricher.GetJobByName(ja.SourceJob)
			if j != nil {
				ja.JobLabel = &j.Label
			} else {
				ja.JobLabel = &ja.Job
			}
		}
	}

	if settings.GetOfferedAccess() != nil && settings.GetOfferedAccess().GetJobs() != nil {
		for _, ja := range settings.GetOfferedAccess().GetJobs() {
			// Lookup job info by using SourceJob field
			j := s.enricher.GetJobByName(ja.SourceJob)
			if j != nil {
				ja.JobLabel = &j.Label
			} else {
				ja.JobLabel = &ja.Job
			}
		}
	}

	if settings.GetEffectiveAccess() != nil {
		if settings.GetEffectiveAccess().GetDispatches() != nil {
			for _, ja := range settings.GetEffectiveAccess().GetDispatches().GetJobs() {
				s.enricher.EnrichJobName(ja)
			}
		}
	}

	return &pbcentrum.GetSettingsResponse{
		Settings: settings,
	}, nil
}

func (s *Server) UpdateSettings(
	ctx context.Context,
	req *pbcentrum.UpdateSettingsRequest,
) (*pbcentrum.UpdateSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	current, err := s.settings.Get(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	fields, err := s.ps.AttrStringList(
		userInfo,
		permscentrum.CentrumServicePerm,
		permscentrum.CentrumServiceUpdateSettingsPerm,
		permscentrum.CentrumServiceUpdateSettingsAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	// Reset access if the user doesn't have access to the shared dispatch center feature
	if !fields.Contains("Shared") {
		req.Settings.Access = current.GetAccess()
		req.Settings.OfferedAccess = current.GetOfferedAccess()
		req.Settings.EffectiveAccess = current.GetEffectiveAccess()
	}

	if !fields.Contains("Public") {
		req.Settings.Public = current.GetPublic()
	}

	settings, err := s.settings.Update(ctx, userInfo.GetJob(), req.GetSettings())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbcentrum.UpdateSettingsResponse{
		Settings: settings,
	}, nil
}
