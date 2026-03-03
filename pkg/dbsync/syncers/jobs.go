package syncers

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

type JobsSync struct {
	*Syncer
}

func NewJobsSync(s *Syncer, state *dbsyncconfig.TableSyncState) *JobsSync {
	return &JobsSync{
		Syncer: s,
	}
}

func (s *JobsSync) Sync(ctx context.Context) (int64, error) {
	jobs, err := s.fetchJobs(ctx)
	if err != nil {
		return 0, err
	}

	count := int64(len(jobs))
	s.logger.Debug("jobsSync", zap.Int64("len", count))
	if count == 0 {
		return 0, nil
	}

	hasFilters := len(s.cfg.Tables.Jobs.Filters) > 0
	jobs, err = s.applyFiltersAndRetrieveGrades(ctx, jobs, hasFilters)
	if err != nil {
		return 0, err
	}

	// Log a warning when no jobs are left after filtering
	if hasFilters && count == 0 {
		s.logger.Warn("no jobs left after filtering")
		return 0, nil
	}

	// Sync jobs to FiveNet server
	if err := s.sendData(ctx, &pbsync.SendDataRequest{
		Data: &pbsync.SendDataRequest_Jobs{
			Jobs: &syncdata.DataJobs{
				Jobs: jobs,
			},
		},
	}); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *JobsSync) fetchJobs(ctx context.Context) ([]*jobs.Job, error) {
	limit := s.cfg.Limits.Jobs
	q := s.cfg.Tables.Jobs.GetQuery(limit)
	s.logger.Debug("jobs sync query", zap.String("query", q))

	jobsResults := []*jobs.Job{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &jobsResults); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query jobs. %w", err)
		}
	}

	return jobsResults, nil
}

func (s *JobsSync) applyFiltersAndRetrieveGrades(
	ctx context.Context,
	js []*jobs.Job,
	hasFilters bool,
) ([]*jobs.Job, error) {
	sQuery := s.cfg.Tables.Jobs

	filtered := make([]*jobs.Job, 0, len(js))

	for _, job := range js {
		if hasFilters {
			// Apply filters
			filtered := false
			for _, filter := range sQuery.Filters {
				if filter.CompiledPattern.MatchString(job.Name) {
					switch filter.Action {
					case dbsyncconfig.FilterActionDrop:
						filtered = true

					case dbsyncconfig.FilterActionReplace:
						job.Name = filter.CompiledPattern.ReplaceAllString(
							job.Name,
							filter.Replacement,
						)

					default:
						s.logger.Warn(
							"unknown filter action",
							zap.String("action", string(filter.Action)),
						)
					}

					if filtered {
						break
					}
				}
			}
			if filtered {
				continue
			}
		}

		grades, err := s.getGrades(ctx, job.GetName())
		if err != nil {
			return nil, err
		}
		job.Grades = grades
		filtered = append(filtered, job)
	}

	js = filtered

	return js, nil
}

func (s *JobsSync) getGrades(ctx context.Context, job string) ([]*jobs.JobGrade, error) {
	query := s.cfg.Tables.JobGrades.GetQuery(nil, 200)

	grades := []*jobs.JobGrade{}
	if _, err := qrm.Query(ctx, s.db, query, []any{
		job,
	}, &grades); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query job grades for job %q. %w", job, err)
		}
	}

	return grades, nil
}
