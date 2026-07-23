package utils

import (
	"slices"
	"strings"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
)

// NormalizeUserJobs keeps the scalar job fields and job list aligned.
//
// When jobs are present, the scalar job/job_grade fields are derived from the
// primary job entry if they are empty. When no usable jobs remain after cleanup,
// a single fallback job is created from the scalar fields or the provided
// fallback values.
func NormalizeUserJobs(user *syncdata.DataUser, fallbackJob string, fallbackGrade int32) {
	if user == nil {
		return
	}

	jobs := user.GetJobs()
	normalizedJobs := make([]*users.UserJob, 0, len(jobs))
	for _, job := range jobs {
		if strings.TrimSpace(job.GetJob()) == "" {
			continue
		}
		normalizedJobs = append(normalizedJobs, job)
	}

	if len(normalizedJobs) == 0 {
		if user.GetJob() == "" {
			user.SetJob(fallbackJob)
			user.SetJobGrade(fallbackGrade)
		}
		user.SetJobs([]*users.UserJob{
			{
				UserId:    user.GetUserId(),
				Job:       user.GetJob(),
				Grade:     user.GetJobGrade(),
				IsPrimary: true,
			},
		})
		return
	}

	slices.SortFunc(normalizedJobs, func(a, b *users.UserJob) int {
		if a.GetIsPrimary() && !b.GetIsPrimary() {
			return -1
		}
		if !a.GetIsPrimary() && b.GetIsPrimary() {
			return 1
		}
		return strings.Compare(a.GetJob(), b.GetJob())
	})

	if len(normalizedJobs) == 1 && user.GetJob() != "" {
		normalizedJobs[0].UserId = user.GetUserId()
		normalizedJobs[0].Job = user.GetJob()
		normalizedJobs[0].Grade = user.GetJobGrade()
		normalizedJobs[0].IsPrimary = true
		user.SetJobs(normalizedJobs)
		return
	}

	if user.GetJob() == "" {
		user.SetJob(normalizedJobs[0].GetJob())
		user.SetJobGrade(normalizedJobs[0].GetGrade())
	}

	user.SetJobs(normalizedJobs)

	foundPrimary := false
	primaryJob := user.GetJob()
	for _, job := range user.GetJobs() {
		if job.GetJob() == primaryJob {
			foundPrimary = true
			job.IsPrimary = true
		} else {
			job.IsPrimary = false
		}
		job.UpdatedAt = nil
	}

	if !foundPrimary {
		user.Jobs[0].IsPrimary = true
	}
}
