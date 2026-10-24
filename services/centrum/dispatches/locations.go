package dispatches

import (
	"context"
	"errors"
	"fmt"
	"strings"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
)

// RebuildDispatchLocations resets and repopulates the in-memory spatial index
// from the current dispatch projections. Callers should establish component
// watchers before invoking this so updates that race the rebuild remain queued.
func (s *DispatchDB) RebuildDispatchLocations(ctx context.Context) error {
	s.dispatchLocationsMutex.Lock()
	s.dispatchLocations = make(map[string]*coords.Coords[*centrumdispatches.Dispatch])
	s.dispatchLocationEntries = make(map[int64]map[string]*centrumdispatches.Dispatch)
	s.dispatchLocationsMutex.Unlock()

	var errs error
	s.specStore.Range(func(_ string, spec *centrumdispatches.DispatchSpec) bool {
		if spec == nil {
			return true
		}
		dispatch, err := s.Get(ctx, spec.GetId())
		if err != nil {
			errs = errors.Join(
				errs,
				fmt.Errorf(
					"failed to load dispatch %d for spatial index rebuild: %w",
					spec.GetId(),
					err,
				),
			)
			return true
		}
		if err := s.UpdateDispatchLocation(dispatch); err != nil {
			errs = errors.Join(
				errs,
				fmt.Errorf(
					"failed to index dispatch %d during spatial index rebuild: %w",
					spec.GetId(),
					err,
				),
			)
		}
		return true
	})

	return errs
}

// UpdateDispatchLocation reconciles one dispatch in the deduplication index.
// The ID map preserves the previous point so moved dispatches can be removed.
func (s *DispatchDB) UpdateDispatchLocation(dispatch *centrumdispatches.Dispatch) error {
	if dispatch == nil {
		return nil
	}

	s.dispatchLocationsMutex.Lock()
	defer s.dispatchLocationsMutex.Unlock()

	s.removeDispatchLocationLocked(dispatch.GetId())
	if dispatch.GetStatus() != nil &&
		centrumutils.IsStatusDispatchComplete(dispatch.GetStatus().GetStatus()) {
		return nil
	}

	jobs := make(map[string]struct{})
	for _, job := range dispatch.GetJobs().GetJobStrings() {
		job = strings.TrimSpace(job)
		if job != "" {
			jobs[job] = struct{}{}
		}
	}
	if len(jobs) == 0 {
		return nil
	}

	entries := make(map[string]*centrumdispatches.Dispatch, len(jobs))
	for job := range jobs {
		locations, ok := s.dispatchLocations[job]
		if !ok {
			locations = coords.New[*centrumdispatches.Dispatch]()
			s.dispatchLocations[job] = locations
		}
		if err := locations.Add(dispatch); err != nil {
			return fmt.Errorf("failed to add dispatch to spatial index: %w", err)
		}
		entries[job] = dispatch
	}
	s.dispatchLocationEntries[dispatch.GetId()] = entries
	return nil
}

// RemoveDispatchLocation removes all indexed entries for a dispatch.
func (s *DispatchDB) RemoveDispatchLocation(id int64) {
	s.dispatchLocationsMutex.Lock()
	defer s.dispatchLocationsMutex.Unlock()
	s.removeDispatchLocationLocked(id)
}

func (s *DispatchDB) removeDispatchLocationLocked(id int64) {
	entries, ok := s.dispatchLocationEntries[id]
	if !ok {
		return
	}
	for job, dispatch := range entries {
		if locations, exists := s.dispatchLocations[job]; exists {
			locations.Remove(dispatch, centrumdispatches.DispatchPointMatchFn(id))
		}
	}
	delete(s.dispatchLocationEntries, id)
}

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
