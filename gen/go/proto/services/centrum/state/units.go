package state

import (
	"strings"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/puzpuzpuz/xsync/v3"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
)

func (s *State) GetUnit(job string, id uint64) (*dispatch.Unit, bool) {
	units, ok := s.units.Load(job)
	if !ok {
		return nil, false
	}

	unit, ok := units.Load(id)
	if !ok {
		return nil, false
	}

	return proto.Clone(unit).(*dispatch.Unit), true
}

func (s *State) GetUnitIDForUserID(userId int32) (uint64, bool) {
	return s.userIDToUnitID.Load(userId)
}

func (s *State) ListUnits(job string) ([]*dispatch.Unit, bool) {
	us := []*dispatch.Unit{}

	units := s.GetUnitsMap(job)

	units.Range(func(key uint64, unit *dispatch.Unit) bool {
		us = append(us, unit)
		return true
	})

	slices.SortFunc(us, func(a, b *dispatch.Unit) int {
		return strings.Compare(a.Name, b.Name)
	})

	return us, true
}

func (s *State) FilterUnits(job string, statuses []dispatch.StatusUnit, notStatuses []dispatch.StatusUnit) []*dispatch.Unit {
	units, ok := s.ListUnits(job)
	if !ok {
		return nil
	}

	us := []*dispatch.Unit{}
	for i := 0; i < len(units); i++ {
		include := true

		// Include statuses that should be listed
		if len(statuses) > 0 && !utils.InSlice(statuses, units[i].Status.Status) {
			include = false
		} else if len(notStatuses) > 0 {
			// Which statuses to ignore
			for _, status := range notStatuses {
				if units[i].Status != nil && units[i].Status.Status == status {
					include = false
					break
				}
			}
		}

		if include {
			us = append(us, units[i])
		}
	}

	return us
}

func (s *State) DeleteUnit(job string, id uint64) {
	units, ok := s.units.Load(job)
	if ok {
		units.Delete(id)
	}

	s.ClearUnitLock(id)
}

func (s *State) UpdateUnit(job string, unitId uint64, unit *dispatch.Unit) error {
	if u, ok := s.GetUnitsMap(job).LoadOrStore(unitId, unit); ok {
		lock := s.GetUnitLock(unitId)
		lock.Lock()
		defer lock.Unlock()

		u.Merge(unit)
	}

	return nil
}

func (s *State) GetUnitsJobs() []string {
	list := []string{}

	s.units.Range(func(job string, _ *xsync.MapOf[uint64, *dispatch.Unit]) bool {
		list = append(list, job)
		return true
	})

	return list
}
