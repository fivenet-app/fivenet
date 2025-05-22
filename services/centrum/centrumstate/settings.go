package centrumstate

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *State) GetSettings(ctx context.Context, job string) (*centrum.Settings, error) {
	settings, err := s.settings.GetOrLoad(ctx, job)
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, fmt.Errorf("failed to get settings for job %s. %w", job, err)
		}
	}

	if settings != nil {
		return settings, nil
	}

	settings = &centrum.Settings{}
	settings.Default(job)

	return settings, nil
}

func (s *State) UpdateSettings(ctx context.Context, job string, in *centrum.Settings) error {
	current, err := s.GetSettings(ctx, job)
	if err != nil {
		if !errors.Is(err, nats.ErrKeyNotFound) {
			return fmt.Errorf("failed to get settings for job %s. %w", job, err)
		}
	}
	current.Merge(in)

	return s.settings.Put(ctx, job, current)
}

func (s *State) ListSettings(ctx context.Context) []*centrum.Settings {
	list := []*centrum.Settings{}

	s.settings.Range(ctx, func(_ string, settings *centrum.Settings) bool {
		list = append(list, settings)
		return true
	})

	return list
}

func (s *State) GetJobList(ctx context.Context, userJob string, userGrade int32) ([]string, error) {
	settings, err := s.GetSettings(ctx, userJob)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings for job %s. %w", userJob, err)
	}

	if settings == nil || settings.Access == nil {
		return nil, nil
	}

	jobs := []string{}
	for _, ja := range settings.Access.Jobs {
		if ja.MinimumGrade > userGrade {
			continue
		}

		js, err := s.GetSettings(ctx, ja.Job)
		if err != nil {
			return nil, fmt.Errorf("failed to get settings from other job %s for job list %s. %w", ja.Job, userJob, err)
		}

		if js == nil || js.Access == nil || js.Access.Jobs == nil {
			continue
		}

		// Check if the job's share their access (equal or higher level)
		if !slices.ContainsFunc(js.Access.Jobs, func(j *centrum.JobAccess) bool {
			return j.Job == userJob && j.MinimumGrade <= userGrade &&
				j.Access > centrum.AccessLevel_ACCESS_LEVEL_BLOCKED && j.Access >= ja.Access
		}) {
			continue
		}

		jobs = append(jobs, ja.Job)
	}

	return jobs, nil
}

func (s *State) HasAccessToJob(ctx context.Context, userJob string, userGrade int32, targetJob string, level centrum.AccessLevel) (bool, error) {
	// Same job, no need to check access
	if userJob == targetJob {
		return true, nil
	}

	settings, err := s.GetSettings(ctx, userJob)
	if err != nil {
		return false, fmt.Errorf("failed to get settings for job %s. %w", userJob, err)
	}

	if settings == nil || settings.Access == nil {
		return false, nil
	}

	// The source job can have a lower access level than the target job
	for _, ja := range settings.Access.Jobs {
		// Find the target job in the access list and ensure user has access
		if ja.Job != targetJob || ja.MinimumGrade > userGrade {
			continue
		}

		// Retrieve target job settings
		js, err := s.GetSettings(ctx, ja.Job)
		if err != nil {
			return false, fmt.Errorf("failed to get settings from other job %s for job list %s. %w", ja.Job, userJob, err)
		}

		if js == nil || js.Access == nil || js.Access.Jobs == nil {
			continue
		}

		// Check if the target job's share their access on equal or higher level
		if slices.ContainsFunc(js.Access.Jobs, func(j *centrum.JobAccess) bool {
			return j.Job == userJob && j.Access > centrum.AccessLevel_ACCESS_LEVEL_BLOCKED && j.Access >= level
		}) {
			return true, nil
		}
	}

	return false, nil
}
