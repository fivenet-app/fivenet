package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/puzpuzpuz/xsync/v3"
)

func (s *State) GetDispatchesMap(job string) *xsync.MapOf[uint64, *dispatch.Dispatch] {
	store, _ := s.Dispatches.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Dispatch] {
		return xsync.NewMapOf[uint64, *dispatch.Dispatch]()
	})

	return store
}

func (s *State) GetUnitsMap(job string) *xsync.MapOf[uint64, *dispatch.Unit] {
	store, _ := s.Units.LoadOrCompute(job, func() *xsync.MapOf[uint64, *dispatch.Unit] {
		return xsync.NewMapOf[uint64, *dispatch.Unit]()
	})

	return store
}
