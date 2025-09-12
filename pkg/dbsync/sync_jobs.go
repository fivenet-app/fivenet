package dbsync

import (
	"context"
	"errors"
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

	limit := int64(200)

	sQuery := s.cfg.Tables.Jobs
	query := prepareStringQuery(sQuery, s.state, 0, limit)

	jobs := []*jobs.Job{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	s.logger.Debug("jobsSync", zap.Int("len", len(jobs)))

	if len(jobs) == 0 {
		return nil
	}

	// Retrieve grades per job
	var err error
	for k := range jobs {
		jobs[k].Grades, err = s.getGrades(ctx, jobs[k].GetName())
		if err != nil {
			return err
		}
	}

	// Sync jobs to FiveNet server
	if s.cli != nil {
		if _, err := s.cli.SendData(ctx, &pbsync.SendDataRequest{
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
