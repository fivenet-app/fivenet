package centrummanager

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Manager) UpdateSettingsInDB(ctx context.Context, job string, settings *centrum.Settings) (*centrum.Settings, error) {
	stmt := tCentrumSettings.
		INSERT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
			tCentrumSettings.PredefinedStatus,
			tCentrumSettings.Timings,
		).
		VALUES(
			job,
			settings.Enabled,
			settings.Mode,
			settings.FallbackMode,
			settings.PredefinedStatus,
			settings.Timings,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCentrumSettings.Job.SET(jet.String(job)),
			tCentrumSettings.Enabled.SET(jet.Bool(settings.Enabled)),
			tCentrumSettings.Mode.SET(jet.Int32(int32(settings.Mode))),
			tCentrumSettings.FallbackMode.SET(jet.Int32(int32(settings.FallbackMode))),
			tCentrumSettings.PredefinedStatus.SET(jet.StringExp(jet.Raw("VALUES(`predefined_status`)"))),
			tCentrumSettings.Timings.SET(jet.StringExp(jet.Raw("VALUES(`timings`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	// Load settings from database so they are updated in the "cache"
	if err := s.LoadSettingsFromDB(ctx, job); err != nil {
		return nil, err
	}

	set, err := s.GetSettings(ctx, job)
	if err != nil {
		return nil, err
	}

	return set, nil
}
