package centrumstate

import (
	"strconv"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords"
)

func JobIdKey(job string, id uint64) string {
	return job + "." + strconv.FormatUint(id, 10)
}

func (s *State) GetDispatchLocations(job string) (*coords.Coords[*centrum.Dispatch], bool) {
	s.dispatchLocationsMutex.RLock()
	defer s.dispatchLocationsMutex.RUnlock()

	locations, ok := s.dispatchLocations[job]
	if !ok {
		locations = coords.New[*centrum.Dispatch]()
		s.dispatchLocations[job] = locations
	}
	return locations, ok
}
