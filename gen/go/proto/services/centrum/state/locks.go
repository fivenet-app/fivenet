package state

import "sync"

func (s *State) GetUnitLock(unitId uint64) *sync.Mutex {
	lock, _ := s.UnitsLocks.LoadOrCompute(unitId, func() *sync.Mutex {
		return &sync.Mutex{}
	})

	return lock
}
