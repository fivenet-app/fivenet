package manager

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
)

func (s *Manager) CheckIfUserIsDisponent(job string, userId int32) bool {
	disponents, err := s.GetDisponents(job)
	if err != nil {
		return false
	}

	if len(disponents) == 0 {
		return false
	}

	for i := 0; i < len(disponents); i++ {
		if userId == disponents[i].UserId {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfUserIsPartOfDispatch(userInfo *userinfo.UserInfo, dsp *centrum.Dispatch, disponentOkay bool) bool {
	// Check if user is a disponent
	if disponentOkay && s.CheckIfUserIsDisponent(userInfo.Job, userInfo.UserId) {
		return true
	}

	// Iterate over units of dispatch and check if the user is in one of the units
	for i := 0; i < len(dsp.Units); i++ {
		unit, err := s.GetUnit(dsp.Units[i].Unit.Job, dsp.Units[i].UnitId)
		if unit == nil || err != nil {
			continue
		}

		if s.CheckIfUserPartOfUnit(userInfo.Job, userInfo.UserId, unit, disponentOkay) {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfUserPartOfUnit(job string, userId int32, unit *centrum.Unit, disponentOkay bool) bool {
	// Check if user is a disponent
	if disponentOkay && s.CheckIfUserIsDisponent(job, userId) {
		return true
	}

	for i := 0; i < len(unit.Users); i++ {
		if (unit.Users[i].User != nil && unit.Users[i].User.UserId == userId) || unit.Users[i].UserId == userId {
			return true
		}
	}

	return false
}
func (s *Manager) CheckIfBotNeeded(job string) bool {
	settings := s.GetSettings(job)

	if settings.Mode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
		return true
	}

	disponents, err := s.GetDisponents(job)
	if err != nil {
		return false
	}

	if len(disponents) == 0 {
		if settings.FallbackMode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
			return true
		}
	}

	return false
}
