package state

import (
	"strings"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/utils"
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
