package dispatches

import (
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords"
)

func (s *DispatchDB) GetLocations(job string) *coords.Coords[*centrumdispatches.Dispatch] {
	s.dispatchLocationsMutex.Lock()
	defer s.dispatchLocationsMutex.Unlock()

	locations, ok := s.dispatchLocations[job]
	if !ok {
		locations = coords.New[*centrumdispatches.Dispatch]()
		s.dispatchLocations[job] = locations
	}
	return locations
}

func (s *DispatchDB) GetLocationsJob() []string {
	s.dispatchLocationsMutex.Lock()
	defer s.dispatchLocationsMutex.Unlock()

	jobs := make([]string, 0, len(s.dispatchLocations))
	for job := range s.dispatchLocations {
		jobs = append(jobs, job)
	}
	return jobs
}
