package state

import (
	"strings"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"golang.org/x/exp/slices"
)

func (s *State) GetUnit(job string, id uint64) (*dispatch.Unit, bool) {
	units, ok := s.Units.Load(job)
	if !ok {
		return nil, false
	}

	return units.Load(id)
}

func (s *State) GetUnitIDForUserID(userId int32) (uint64, bool) {
	return s.UserIDToUnitID.Load(userId)
}

func (s *State) ListUnits(job string) ([]*dispatch.Unit, bool) {
	us := []*dispatch.Unit{}

	units, ok := s.Units.Load(job)
	if !ok {
		return nil, false
	}

	units.Range(func(key uint64, unit *dispatch.Unit) bool {
		us = append(us, unit)
		return true
	})

	slices.SortFunc(us, func(a, b *dispatch.Unit) int {
		return strings.Compare(a.Name, b.Name)
	})

	return us, true
}
