package centrummanager

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"go.uber.org/zap"
)

func (s *Manager) CheckIfUserIsDispatcher(ctx context.Context, job string, userId int32) bool {
	dispatchers, err := s.GetDispatchers(ctx, job)
	if err != nil {
		return false
	}

	if dispatchers.IsEmpty() {
		return false
	}
	for i := range dispatchers.Dispatchers {
		if userId == dispatchers.Dispatchers[i].UserId {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfUserIsPartOfDispatch(ctx context.Context, userInfo *userinfo.UserInfo, dsp *centrum.Dispatch, dispatcherOkay bool) bool {
	// Check if user is a dispatcher
	if dispatcherOkay && s.CheckIfUserIsDispatcher(ctx, userInfo.Job, userInfo.UserId) {
		return true
	}

	// Iterate over units of dispatch and check if the user is in one of the units
	for i := range dsp.Units {
		unit, err := s.GetUnit(ctx, dsp.Units[i].Unit.Job, dsp.Units[i].UnitId)
		if unit == nil || err != nil {
			continue
		}

		if s.CheckIfUserPartOfUnit(ctx, userInfo.Job, userInfo.UserId, unit, dispatcherOkay) {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfUserPartOfUnit(ctx context.Context, job string, userId int32, unit *centrum.Unit, dispatcherOkay bool) bool {
	// Check if user is a dispatcher
	if dispatcherOkay && s.CheckIfUserIsDispatcher(ctx, job, userId) {
		return true
	}

	for i := range unit.Users {
		if (unit.Users[i].User != nil && unit.Users[i].User.UserId == userId) || unit.Users[i].UserId == userId {
			return true
		}
	}

	return false
}

func (s *Manager) CheckIfBotNeeded(ctx context.Context, job string) bool {
	settings, err := s.GetSettings(ctx, job)
	if err != nil {
		s.logger.Error("failed to get centrum settings", zap.String("job", job), zap.Error(err))
		return false
	}

	// If centrum is disabled, why bother with the bot
	if !settings.Enabled {
		return false
	}

	if settings.Mode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
		return true
	}

	dispatchers, err := s.GetDispatchers(ctx, job)
	if err != nil {
		return false
	}

	if dispatchers.IsEmpty() {
		if settings.FallbackMode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
			return true
		}
	}

	return false
}
