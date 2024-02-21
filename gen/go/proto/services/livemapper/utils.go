package livemapper

import "github.com/galexrt/fivenet/gen/go/proto/resources/livemap"

func filterMarkerUpdates(superUser bool, usersJobs map[string]int32, input *livemap.UsersUpdateEvent) *livemap.UsersUpdateEvent {
	return &livemap.UsersUpdateEvent{
		Added:   filterMarkerUpdatesList(superUser, usersJobs, input.Added),
		Updated: filterMarkerUpdatesList(superUser, usersJobs, input.Updated),
		Removed: filterMarkerUpdatesList(superUser, usersJobs, input.Removed),
	}
}

func filterMarkerUpdatesList(superUser bool, usersJobs map[string]int32, list []*livemap.UserMarker) []*livemap.UserMarker {
	out := []*livemap.UserMarker{}

	for _, marker := range list {
		grade, ok := usersJobs[marker.Info.Job]
		if !ok && !superUser {
			continue
		}

		if grade != -1 && marker.User.JobGrade > grade {
			continue
		}

		out = append(out, marker)
	}

	return out
}
