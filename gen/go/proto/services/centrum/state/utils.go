package state

import (
	"sync"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/coords"
	"github.com/puzpuzpuz/xsync/v3"
)

func (s *State) GetDispatchesMap(job string) *xsync.MapOf[uint64, *dispatch.Dispatch] {
	store, _ := s.dispatches.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Dispatch] {
		return xsync.NewMapOf[uint64, *dispatch.Dispatch]()
	})

	return store
}

func (s *State) GetUnitsMap(job string) *xsync.MapOf[uint64, *dispatch.Unit] {
	store, _ := s.units.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Unit] {
		return xsync.NewMapOf[uint64, *dispatch.Unit]()
	})

	return store
}

func (s *State) GetUnitLock(unitId uint64) *sync.Mutex {
	lock, _ := s.unitsLocks.LoadOrCompute(unitId, func() *sync.Mutex {
		return &sync.Mutex{}
	})

	return lock
}

func (s *State) ClearUnitLock(unitId uint64) {
	s.unitsLocks.Delete(unitId)
}

func (s *State) GetDispatchLocations(job string) *coords.Coords[*dispatch.Dispatch] {
	return s.dispatchLocations[job]
}
