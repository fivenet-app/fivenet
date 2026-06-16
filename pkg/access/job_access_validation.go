package access

import (
	"errors"
	"slices"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/nats-io/nats.go/jetstream"
)

type jobGetter interface {
	Get(job string) (*jobs.Job, error)
}

func ValidateJobAccessEntries(
	js jobGetter,
	in *[]*resourcesaccess.JobAccess,
	fixEntries bool,
) (bool, error) {
	if js == nil {
		return true, nil
	}

	jobMap := make(map[string]*jobs.Job)
	valid := true
	*in = slices.DeleteFunc(*in, func(ja *resourcesaccess.JobAccess) bool {
		if !valid {
			return false
		}

		if ja.GetJob() == "" || ja.GetMinimumGrade() < 0 {
			return true
		}

		j, ok := jobMap[ja.GetJob()]
		if !ok {
			var err error
			j, err = js.Get(ja.GetJob())
			if err != nil {
				if fixEntries && errors.Is(err, jetstream.ErrKeyNotFound) {
					return true
				}

				valid = false
				return true
			}

			jobMap[ja.GetJob()] = j
		}

		if int(ja.GetMinimumGrade()) >= len(j.Grades) {
			if !slices.ContainsFunc(j.Grades, func(g *jobs.JobGrade) bool {
				return g.GetGrade() == ja.GetMinimumGrade()
			}) {
				if fixEntries {
					if len(j.Grades) == 0 {
						return true
					}

					jg := j.Grades[len(j.Grades)-1]
					ja.SetMinimumGrade(jg.Grade)
					return false
				}

				valid = false
				return true
			}
		}

		return false
	})

	return valid, nil
}
