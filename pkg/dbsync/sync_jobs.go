package dbsync

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/sync"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

type jobsSync struct {
	*syncer

	state *TableSyncState
}

func newJobsSync(s *syncer, state *TableSyncState) *jobsSync {
	return &jobsSync{
		syncer: s,
		state:  state,
	}
}

func (s *jobsSync) Sync(ctx context.Context) error {
	if !s.cfg.Tables.Jobs.Enabled {
		return nil
	}

	jobs, err := s.fetchJobs(ctx)
	if err != nil {
		return err
	}

	s.logger.Debug("jobsSync", zap.Int("len", len(jobs)))

	if len(jobs) == 0 {
		return nil
	}

	hasFilters := len(s.cfg.Tables.Jobs.Filters) > 0
	jobs, err = s.applyFiltersAndRetrieveGrades(ctx, jobs, hasFilters)
	if err != nil {
		return err
	}

	// Log a warning when no jobs are left after filtering
	if hasFilters && len(jobs) == 0 {
		s.logger.Warn("no jobs left after filtering")
		return nil
	}

	// Sync jobs to FiveNet server
	if s.cli != nil {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Jobs{
				Jobs: &sync.DataJobs{
					Jobs: jobs,
				},
			},
		}); err != nil {
			return err
		}
	}

	s.state.Set(0, nil)

	return nil
}

func (s *jobsSync) fetchJobs(ctx context.Context) ([]*jobs.Job, error) {
	limit := int64(200)
	sQuery := s.cfg.Tables.Jobs
	query := prepareStringQuery(sQuery.DBSyncTable, s.state, 0, limit)

	jobs := []*jobs.Job{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return jobs, nil
}

func (s *jobsSync) applyFiltersAndRetrieveGrades(
	ctx context.Context,
	jobs []*jobs.Job,
	hasFilters bool,
) ([]*jobs.Job, error) {
	sQuery := s.cfg.Tables.Jobs

outer:
	for k := range jobs {
		if hasFilters {
			// Apply filters
			for _, filter := range sQuery.Filters {
				if filter.compiledPattern.MatchString(jobs[k].Name) {
					switch filter.Action {
					case FilterActionDrop:
						jobs = slices.Delete(jobs, k, 1)
						continue outer

					case FilterActionReplace:
						jobs[k].Name = filter.compiledPattern.ReplaceAllString(
							jobs[k].Name,
							filter.Replacement,
						)

					default:
						s.logger.Warn("unknown filter action", zap.String("action", string(filter.Action)))
					}
					continue
				}
			}
		}

		grades, err := s.getGrades(ctx, jobs[k].GetName())
		if err != nil {
			return nil, err
		}
		jobs[k].Grades = grades
	}

	return jobs, nil
}

func (s *jobsSync) getGrades(ctx context.Context, job string) ([]*jobs.JobGrade, error) {
	sQuery := s.cfg.Tables.JobGrades
	query := prepareStringQuery(sQuery, nil, 0, 200)
	query = strings.ReplaceAll(query, "$jobName", "?")

	grades := []*jobs.JobGrade{}
	if _, err := qrm.Query(ctx, s.db, query, []any{
		job,
	}, &grades); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return grades, nil
}
