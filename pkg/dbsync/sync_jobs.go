package dbsync

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/go-jet/jet/v2/qrm"
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

	limit := 200

	sQuery := s.cfg.Tables.Jobs
	query := prepareStringQuery(sQuery, s.state, 0, limit)

	jobs := []*users.Job{}
	if _, err := qrm.Query(ctx, s.db, query, []interface{}{}, &jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	if len(jobs) == 0 {
		return nil
	}

	// Retrieve grades per job
	var err error
	for k := range jobs {
		jobs[k].Grades, err = s.getGrades(ctx, jobs[k].Name)
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

func (s *jobsSync) getGrades(ctx context.Context, job string) ([]*users.JobGrade, error) {
	sQuery := s.cfg.Tables.JobGrades
	query := prepareStringQuery(sQuery, nil, 0, 200)
	query = strings.ReplaceAll(query, "$jobName", "?")

	grades := []*users.JobGrade{}
	if _, err := qrm.Query(ctx, s.db, query, []interface{}{
		job,
	}, &grades); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return grades, nil
}
