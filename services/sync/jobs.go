package sync

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	jobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) SendJobs(
	ctx context.Context,
	req *pbsync.SendJobsRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())

	rowsAffected, err := s.handleJobsData(ctx, req.GetJobs())
	if err != nil {
		return nil, fmt.Errorf("failed to handle jobs data. %w", err)
	}

	return &pbsync.SendDataResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (s *Server) handleJobsData(
	ctx context.Context,
	jobs []*jobs.Job,
) (int64, error) {
	if len(jobs) == 0 {
		return 0, nil
	}

	tJobs := table.FivenetJobs

	stmt := tJobs.
		INSERT(
			tJobs.Name,
			tJobs.Label,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobs.Name.SET(mysql.RawString("VALUES(`name`)")),
			tJobs.Label.SET(mysql.RawString("VALUES(`label`)")),
		)

	for _, job := range jobs {
		stmt = stmt.VALUES(
			job.GetName(),
			job.GetLabel(),
		)
	}

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, fmt.Errorf("failed to execute job insert statement. %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for job insert. %w", err)
	}

	for _, job := range jobs {
		rowCounts, err := s.handleJobGrades(ctx, job)
		if err != nil {
			return 0, fmt.Errorf("failed to handle job grades for job %s. %w", job.GetName(), err)
		}

		rowsAffected += rowCounts
	}

	return rowsAffected, nil
}

func (s *Server) handleJobGrades(ctx context.Context, job *jobs.Job) (int64, error) {
	if len(job.GetGrades()) == 0 {
		return 0, nil
	}

	rowsAffectedCount := int64(0)

	tJobsGrades := table.FivenetJobsGrades.AS("job_grade")

	selectStmt := tJobsGrades.
		SELECT(
			tJobsGrades.JobName.AS("job_grade.job_name"),
			tJobsGrades.Grade,
			tJobsGrades.Label,
		).
		FROM(tJobsGrades).
		ORDER_BY(
			tJobsGrades.Grade.ASC(),
		).
		WHERE(tJobsGrades.JobName.EQ(mysql.String(job.GetName())))

	currentGrades := []*jobs.JobGrade{}
	if err := selectStmt.QueryContext(ctx, s.db, &currentGrades); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, fmt.Errorf(
				"failed to query current job grades for job %s. %w",
				job.GetName(),
				err,
			)
		}
	}

	toCreate, toUpdate, toDelete := []*jobs.JobGrade{}, []*jobs.JobGrade{}, []*jobs.JobGrade{}
	if len(currentGrades) == 0 {
		toCreate = job.GetGrades()
	} else {
		// Update cache
		foundTracker := []int{}
		for _, cg := range currentGrades {
			var found *jobs.JobGrade
			var foundIdx int

			for i, ug := range job.GetGrades() {
				if cg.GetGrade() != ug.GetGrade() {
					continue
				}

				found = ug
				foundIdx = i
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete = append(toDelete, cg)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cg.GetLabel() != found.GetLabel() {
				cg.Label = found.GetLabel()
				changed = true
			}

			if changed {
				toUpdate = append(toUpdate, cg)
			}
		}

		for i, uj := range job.GetGrades() {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate = append(toCreate, uj)
			}
		}
	}

	tJobsGrades = table.FivenetJobsGrades

	if len(toCreate) > 0 {
		stmt := tJobsGrades.
			INSERT(
				tJobsGrades.JobName,
				tJobsGrades.Grade,
				tJobsGrades.Label,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tJobsGrades.JobName.SET(mysql.RawString("VALUES(`job_name`)")),
				tJobsGrades.Grade.SET(mysql.RawInt("VALUES(`grade`)")),
				tJobsGrades.Label.SET(mysql.RawString("VALUES(`label`)")),
			)

		for _, grade := range toCreate {
			stmt = stmt.VALUES(
				grade.GetJobName(),
				grade.GetGrade(),
				grade.GetLabel(),
			)
		}

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return 0, fmt.Errorf("failed to execute job grades insert statement. %w", err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve rows affected for job grades insert. %w", err)
		}

		rowsAffectedCount += rowsAffected
	}

	if len(toUpdate) > 0 {
		for _, grade := range toUpdate {
			stmt := tJobsGrades.
				UPDATE(
					tJobsGrades.JobName,
					tJobsGrades.Grade,
					tJobsGrades.Label,
				).
				SET(
					grade.GetJobName(),
					grade.GetGrade(),
					grade.GetLabel(),
				).
				WHERE(mysql.AND(
					tJobsGrades.JobName.EQ(mysql.String(job.GetName())),
					tJobsGrades.Grade.EQ(mysql.Int32(grade.GetGrade())),
				))

			res, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, fmt.Errorf(
					"failed to execute job grades update statement for grade %d. %w",
					grade.GetGrade(),
					err,
				)
			}
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf(
					"failed to retrieve rows affected for job grades update. %w",
					err,
				)
			}

			rowsAffectedCount += rowsAffected
		}
	}

	if len(toDelete) > 0 {
		for _, grade := range toDelete {
			stmt := tJobsGrades.
				DELETE().
				WHERE(mysql.AND(
					tJobsGrades.JobName.EQ(mysql.String(job.GetName())),
					tJobsGrades.Grade.EQ(mysql.Int32(grade.GetGrade())),
				)).
				LIMIT(1)

			res, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, fmt.Errorf(
					"failed to execute job grades delete statement for grade %d. %w",
					grade.GetGrade(),
					err,
				)
			}
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf(
					"failed to retrieve rows affected for job grades delete. %w",
					err,
				)
			}

			rowsAffectedCount += rowsAffected
		}
	}

	return rowsAffectedCount, nil
}
