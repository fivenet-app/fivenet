package settings

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
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

func (s *SettingsDB) GetJobAccessList(
	ctx context.Context,
	userJob string,
	_ int32,
) ([]string, *pbcentrum.JobAccess, error) {
	settings, err := s.Get(ctx, userJob)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get settings for job %s. %w", userJob, err)
	}

	jobsAccess := &pbcentrum.JobAccess{}
	jobs := []string{}

	// Add the user's own job to the access list
	jae := &pbcentrum.JobAccessEntry{
		Job:    userJob,
		Access: centrum.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_DISPATCH,
	}
	s.enricher.EnrichJobName(jae)

	jobsAccess.Dispatches = append(jobsAccess.Dispatches, jae)
	jobs = append(jobs, userJob)

	if settings == nil || settings.GetAccess() == nil {
		return jobs, jobsAccess, nil
	}

	for _, ja := range settings.GetAccess().GetJobs() {
		js, err := s.Get(ctx, ja.GetJob())
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to get settings from other job %s for job list %s. %w",
				ja.GetJob(),
				userJob,
				err,
			)
		}

		if js == nil || js.GetAccess() == nil || js.Access.Jobs == nil {
			continue
		}

		// Check if the job's share their access (equal or higher level)
		if !slices.ContainsFunc(js.GetAccess().GetJobs(), func(j *centrum.CentrumJobAccess) bool {
			return j.GetJob() == userJob &&
				j.GetAccess() > centrum.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_BLOCKED &&
				j.GetAccess() >= ja.GetAccess()
		}) {
			continue
		}

		jae := &pbcentrum.JobAccessEntry{
			Job:    ja.GetJob(),
			Access: ja.GetAccess(),
		}
		s.enricher.EnrichJobName(jae)

		jobsAccess.Dispatches = append(jobsAccess.Dispatches, jae)
		jobs = append(jobs, ja.GetJob())
	}

	jobs = utils.RemoveSliceDuplicates(jobs)

	return jobs, jobsAccess, nil
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

	if settings == nil || settings.GetAccess() == nil {
		return false, nil
	}

	// The source job can have a lower access level than the target job
	for _, ja := range settings.GetAccess().GetJobs() {
		// Find the target job in the access list and ensure user has access
		if ja.GetJob() != targetJob || ja.GetMinimumGrade() > userGrade {
			continue
		}

		// Retrieve target job settings
		js, err := s.Get(ctx, ja.GetJob())
		if err != nil {
			return false, fmt.Errorf(
				"failed to get settings from other job %s for job list %s. %w",
				ja.GetJob(),
				userJob,
				err,
			)
		}

		if js == nil || js.GetAccess() == nil || js.Access.Jobs == nil {
			continue
		}

		// Check if the target job's share their access on equal or higher level
		if slices.ContainsFunc(js.GetAccess().GetJobs(), func(j *centrum.CentrumJobAccess) bool {
			return j.GetJob() == userJob &&
				j.GetAccess() > centrum.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_BLOCKED &&
				j.GetAccess() >= level
		}) {
			return true, nil
		}
	}

	return false, nil
}
