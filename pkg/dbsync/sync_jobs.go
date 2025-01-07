package dbsync

import (
	"context"
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

func NewJobsSync(s *syncer, state *TableSyncState) (ISyncer, error) {
	return &jobsSync{
		syncer: s,
		state:  state,
	}, nil
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
		return err
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
		if _, err := s.cli.SyncData(ctx, &pbsync.SyncDataRequest{
			Data: &pbsync.SyncDataRequest_Jobs{
				Jobs: &sync.DataJobs{
					Jobs: jobs,
				},
			},
		}); err != nil {
			return err
		}
	}

	s.state.Set(s.cfg.Tables.Jobs.IDField, 0, nil)

	return nil
}

func (s *jobsSync) getGrades(ctx context.Context, job string) ([]*users.JobGrade, error) {
	grades := []*users.JobGrade{}

	sQuery := s.cfg.Tables.JobGrades
	query := prepareStringQuery(sQuery, nil, 0, 200)
	query = strings.ReplaceAll(query, "$jobName", "?")

	if _, err := qrm.Query(ctx, s.db, query, []interface{}{
		job,
	}, &grades); err != nil {
		return nil, err
	}

	return grades, nil
}
