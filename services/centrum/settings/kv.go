package settings

import (
	"context"
	"errors"
	"fmt"

	centrumaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/access"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *SettingsDB) updateInKV(
	ctx context.Context,
	job string,
	in *centrumsettings.Settings,
) error {
	if err := s.store.Put(ctx, job, in); err != nil {
		return err
	}

	return nil
}

func (s *SettingsDB) Get(ctx context.Context, job string) (*centrumsettings.Settings, error) {
	settings, err := s.store.GetOrLoad(ctx, job)
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, err
		}

		settings = &centrumsettings.Settings{
			Job: job,
		}
	}
	settings.Default(job)

	return settings, nil
}

func (s *SettingsDB) List(_ context.Context) []*centrumsettings.Settings {
	return s.store.List()
}

func (s *SettingsDB) ListFunc(
	_ context.Context,
	fn func(key string, val *centrumsettings.Settings) bool,
) []*centrumsettings.Settings {
	return s.store.ListFiltered("", fn)
}

func (s *SettingsDB) GetAccessList(
	ctx context.Context,
	userJob string,
	_ int32,
) ([]string, *centrumsettings.EffectiveAccess, error) {
	settings, err := s.Get(ctx, userJob)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get settings for job %s. %w", userJob, err)
	}

	// Build the effective dispatch access from the job's accepted inbound offers.
	access := settings.GetEffectiveAccess()
	jobs := []string{}
	if access == nil {
		access = s.calculateEffectiveAccess(
			settings.GetJob(),
			settings.GetOfferedAccess(),
		)
		settings.SetEffectiveAccess(access)
	}
	if access.GetDispatches() == nil {
		access.SetDispatches(&centrumsettings.EffectiveDispatchAccess{
			Jobs: []*centrumsettings.JobAccessEntry{},
		})
	}

	// Always include the caller's own job.
	jae := &centrumsettings.JobAccessEntry{
		Job:    userJob,
		Access: centrumaccess.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_DISPATCH,
	}
	s.enricher.EnrichJobName(jae)

	access.Dispatches.Jobs = append(access.Dispatches.Jobs, jae)
	jobs = append(jobs, userJob)

	for _, ja := range access.GetDispatches().GetJobs() {
		if ja.GetJob() == userJob {
			continue // Skip the caller job, it was already added above.
		}

		s.enricher.EnrichJobName(ja)
		jobs = append(jobs, ja.GetJob())
	}

	jobs = utils.SliceDedup(jobs)

	return jobs, access, nil
}

func (s *SettingsDB) HasAccessToJob(
	ctx context.Context,
	userJob string,
	userGrade int32,
	targetJob string,
	level centrumaccess.CentrumAccessLevel,
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

			if ja.GetAccess() > centrumaccess.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_BLOCKED &&
				ja.GetAccess() >= level {
				return true, nil
			}
		}
	}

	return false, nil
}
