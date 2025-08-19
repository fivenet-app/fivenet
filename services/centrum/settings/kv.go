package settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *SettingsDB) updateInKV(ctx context.Context, job string, in *centrum.Settings) error {
	if err := s.store.Put(ctx, job, in); err != nil {
		return err
	}

	return nil
}

func (s *SettingsDB) Get(ctx context.Context, job string) (*centrum.Settings, error) {
	settings, err := s.store.GetOrLoad(ctx, job)
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, err
		}

		settings = &centrum.Settings{
			Job: job,
		}
	}
	settings.Default(job)

	return settings, nil
}

func (s *SettingsDB) List(_ context.Context) []*centrum.Settings {
	return s.store.List()
}

func (s *SettingsDB) ListFunc(
	_ context.Context,
	fn func(key string, val *centrum.Settings) bool,
) []*centrum.Settings {
	return s.store.ListFiltered("", fn)
}

func (s *SettingsDB) GetAccessList(
	ctx context.Context,
	userJob string,
	_ int32,
) ([]string, *centrum.EffectiveAccess, error) {
	settings, err := s.Get(ctx, userJob)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get settings for job %s. %w", userJob, err)
	}

	if settings.GetEffectiveAccess() == nil {
		settings.EffectiveAccess = &centrum.EffectiveAccess{
			Dispatches: &centrum.EffectiveDispatchAccess{
				Jobs: []*centrum.JobAccessEntry{},
			},
		}
	}
	if settings.GetEffectiveAccess().GetDispatches() == nil {
		settings.EffectiveAccess.Dispatches = &centrum.EffectiveDispatchAccess{
			Jobs: []*centrum.JobAccessEntry{},
		}
	}

	access := settings.GetEffectiveAccess()
	jobs := []string{}
	if access == nil {
		s.calculateEffectiveAccess(settings)
	}

	// Add the user's own job to the access list
	jae := &centrum.JobAccessEntry{
		Job:    userJob,
		Access: centrum.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_DISPATCH,
	}
	s.enricher.EnrichJobName(jae)

	access.Dispatches.Jobs = append(access.Dispatches.Jobs, jae)
	jobs = append(jobs, userJob)

	if settings == nil || settings.GetEffectiveAccess() == nil {
		return jobs, access, nil
	}

	if settings.GetEffectiveAccess().GetDispatches() != nil {
		for _, ja := range settings.GetEffectiveAccess().GetDispatches().GetJobs() {
			if ja.GetJob() == userJob {
				continue // Skip the user's own job, as it is already added
			}

			s.enricher.EnrichJobName(ja)
			jobs = append(jobs, ja.GetJob())
		}
	}

	jobs = utils.RemoveSliceDuplicates(jobs)

	return jobs, access, nil
}

func (s *SettingsDB) HasAccessToJob(
	ctx context.Context,
	userJob string,
	userGrade int32,
	targetJob string,
	level centrum.CentrumAccessLevel,
) (bool, error) {
	// Same job, no need to check access
	if userJob == targetJob {
		return true, nil
	}

	settings, err := s.Get(ctx, userJob)
	if err != nil {
		return false, fmt.Errorf("failed to get settings for job %s. %w", userJob, err)
	}

	if settings == nil || settings.GetEffectiveAccess() == nil {
		return false, nil
	}

	if settings.GetEffectiveAccess().GetDispatches() != nil {
		// The source job can have a lower access level than the target job
		for _, ja := range settings.GetEffectiveAccess().GetDispatches().GetJobs() {
			// Find the target job in the access list and ensure user has access
			if ja.GetJob() != targetJob {
				continue
			}

			if ja.GetAccess() > centrum.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_BLOCKED &&
				ja.GetAccess() >= level {
				return true, nil
			}
		}
	}

	return false, nil
}
