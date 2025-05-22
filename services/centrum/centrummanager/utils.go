package centrummanager

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
)

func (s *Manager) CheckIfUserIsDisponent(ctx context.Context, job string, userId int32) bool {
	disponents, err := s.GetDisponents(ctx, job)
	if err != nil {
		return false
	}

	if disponents == nil || len(disponents.Disponents) == 0 {
		return false
	}
	for i := range disponents.Disponents {
		if userId == disponents.Disponents[i].UserId {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfUserIsPartOfDispatch(ctx context.Context, userInfo *userinfo.UserInfo, dsp *centrum.Dispatch, disponentOkay bool) bool {
	// Check if user is a disponent
	if disponentOkay && s.CheckIfUserIsDisponent(ctx, userInfo.Job, userInfo.UserId) {
		return true
	}

	// Iterate over units of dispatch and check if the user is in one of the units
	for i := range dsp.Units {
		unit, err := s.GetUnit(ctx, dsp.Units[i].Unit.Job, dsp.Units[i].UnitId)
		if unit == nil || err != nil {
			continue
		}

		if s.CheckIfUserPartOfUnit(ctx, userInfo.Job, userInfo.UserId, unit, disponentOkay) {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfUserPartOfUnit(ctx context.Context, userJob string, userId int32, unit *centrum.Unit, disponentOkay bool) bool {
	// Check if user is a disponent
	if disponentOkay && s.CheckIfUserIsDisponent(ctx, userJob, userId) {
		return true
	}

	for i := range unit.Users {
		if (unit.Users[i].User != nil && unit.Users[i].User.UserId == userId) || unit.Users[i].UserId == userId {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfBotNeeded(ctx context.Context, job string) (bool, error) {
	settings, err := s.GetSettings(ctx, job)
	if err != nil {
		return false, err
	}

	// If centrum is disabled, why bother with the bot
	if !settings.Enabled {
		return false, nil
	}

	if settings.Mode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
		return true, nil
	}

	disponents, err := s.GetDisponents(ctx, job)
	if err != nil {
		return false, nil
	}

	if disponents == nil || len(disponents.Disponents) == 0 {
		if settings.FallbackMode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
			return true, nil
		}
	}

	return false, nil
}
