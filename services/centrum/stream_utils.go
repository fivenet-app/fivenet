package centrum

import (
	"fmt"

	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
)

func centrumSubjects(jobs []string) []string {
	subs := make([]string, 0, len(jobs))
	for _, j := range jobs {
		subs = append(subs, fmt.Sprintf("%s.%s.>", eventscentrum.BaseSubject, j))
	}
	return subs
}

func kvSubjects(bucket string, jobs []string, noWildcard bool) []string {
	subs := make([]string, 0, len(jobs))
	for _, j := range jobs {
		// JetStream turns each KV key into: $KV.<BUCKET>.<key>
		subject := fmt.Sprintf("$KV.%s.%s", bucket, j)
		if !noWildcard {
			subject += ".*"
		}

		subs = append(subs, subject)
	}
	return subs
}
