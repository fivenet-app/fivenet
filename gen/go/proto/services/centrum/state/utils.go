package state

import (
	"strconv"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/pkg/coords"
)

func JobIdKey(job string, id uint64) string {
	return job + "." + strconv.Itoa(int(id))
}

func userIdKey(id int32) string {
	return strconv.Itoa(int(id))
}

func (s *State) GetDispatchLocations(job string) *coords.Coords[*centrum.Dispatch] {
	s.dispatchLocationsMutex.RLock()
	defer s.dispatchLocationsMutex.RUnlock()

	return s.dispatchLocations[job]
}
