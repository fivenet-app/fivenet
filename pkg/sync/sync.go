package sync

import (
	"database/sql"

	"github.com/galexrt/rphub/model"
	"go.uber.org/zap"
)

type Sync struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Sync {
	return &Sync{
		logger: logger,
	}
}

func (s *Sync) SyncJobs(db *sql.DB) {
	res, err := db.Query("SELECT `name`, `label` FROM jobs")
	if err != nil {
		s.logger.Error("failed to query fivem jobs", zap.Error(err))
		return
	}
	defer res.Close()

	for res.Next() {

		job := &model.Job{}
		err := res.Scan(&job.Name, &job.Label)
		if err != nil {
			s.logger.Error("failed to scan fivem data into job object", zap.Error(err))
			continue
		}

		result := model.DB.Create(job)
		if result.Error != nil {
			s.logger.Error("failed to create job in our database", zap.Error(result.Error))
			continue
		}
	}
}

func (s *Sync) SyncJobGrades(db *sql.DB) {
	res, err := db.Query("SELECT `job_name`, `grade`, `label` FROM job_grades")
	if err != nil {
		s.logger.Error("failed to query fivem job_grades", zap.Error(err))
		return
	}
	defer res.Close()

	for res.Next() {

		jobGrade := &model.JobGrade{}
		jobName := ""
		err := res.Scan(&jobName, &jobGrade.Grade, &jobGrade.Label)
		if err != nil {
			s.logger.Error("failed to scan fivem data into job grade object", zap.Error(err))
			continue
		}

		job := &model.Job{}
		result := model.DB.Where("name = ?", jobName).First(job)
		if result.Error != nil {
			s.logger.Error("failed to find job in our database", zap.Error(result.Error))
			continue
		}
		jobGrade.JobID = job.ID

		result = model.DB.Create(jobGrade)
		if result.Error != nil {
			s.logger.Error("failed to create job grade in our database", zap.Error(result.Error))
			continue
		}
	}
}

func (s *Sync) SyncUsers(db *sql.DB) {
	// TODO
}
