package centrumstate

import (
	"strconv"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/pkg/coords"
)

func JobIdKey(job string, id uint64) string {
	return job + "." + strconv.Itoa(int(id))
}

func userIdKey(id int32) string {
	return strconv.Itoa(int(id))
}

func (s *State) GetDispatchLocations(job string) (*coords.Coords[*centrum.Dispatch], bool) {
	s.dispatchLocationsMutex.RLock()
	defer s.dispatchLocationsMutex.RUnlock()

	locations, ok := s.dispatchLocations[job]
	return locations, ok
}
