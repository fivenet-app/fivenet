package helpers

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Helpers struct {
	logger *zap.Logger

	settings    *settings.SettingsDB
	dispatchers *dispatchers.DispatchersDB
	units       *units.UnitDB
}

type Params struct {
	fx.In

	Logger *zap.Logger

	Settings    *settings.SettingsDB
	Dispatchers *dispatchers.DispatchersDB
	Units       *units.UnitDB
}

func New(p Params) *Helpers {
	logger := p.Logger.Named("centrum.state")
	d := &Helpers{
		logger: logger,

		settings:    p.Settings,
		dispatchers: p.Dispatchers,
		units:       p.Units,
	}

	return d
}

func (s *Helpers) CheckIfBotNeeded(ctx context.Context, job string) bool {
	settings, err := s.settings.Get(ctx, job)
	if err != nil {
		s.logger.Error("failed to get centrum settings", zap.String("job", job), zap.Error(err))
		return false
	}

	// If centrum is disabled, why bother with the bot
	if !settings.GetEnabled() {
		return false
	}

	if settings.GetMode() == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
		return true
	}

	dispatchers, err := s.dispatchers.Get(ctx, job)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return false
	}

	if dispatchers.IsEmpty() &&
		settings.GetFallbackMode() == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
		return true
	}

	return false
}

func (s *Helpers) CheckIfUserIsDispatcher(ctx context.Context, job string, userId int32) bool {
	dispatchers, err := s.dispatchers.Get(ctx, job)
	if err != nil {
		return false
	}

	if dispatchers.IsEmpty() {
		return false
	}
	for i := range dispatchers.GetDispatchers() {
		if userId == dispatchers.GetDispatchers()[i].GetUserId() {
			return true
		}
	}

	return false
}

func (s *Helpers) CheckIfUserIsPartOfDispatch(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	dsp *centrum.Dispatch,
	dispatcherOkay bool,
) bool {
	// Check if user is allowed to access the dispatch if the job is not the same as the user's job, need to check the
	// job's dispatch center settings access
	if !dsp.Jobs.ContainsJob(userInfo.GetJob()) {
		ok, err := s.settings.HasAccessToJob(
			ctx,
			userInfo.GetJob(),
			userInfo.GetJobGrade(),
			dsp.GetJobs().GetJobs()[0].GetName(),
			centrum.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_PARTICIPATE,
		)
		if err != nil {
			s.logger.Error(
				"failed to check access to job for dispatch",
				zap.String("job", userInfo.GetJob()),
				zap.Error(err),
			)
			return false
		}
		if !ok {
			return false
		}
	}

	// Check if user is a dispatcher
	if dispatcherOkay && s.CheckIfUserIsDispatcher(ctx, userInfo.GetJob(), userInfo.GetUserId()) {
		return true
	}

	// Iterate over units of dispatch and check if the user is in one of the units
	for i := range dsp.GetUnits() {
		unit, err := s.units.Get(ctx, dsp.GetUnits()[i].GetUnitId())
		if unit == nil || err != nil {
			continue
		}

		if s.CheckIfUserPartOfUnit(
			ctx,
			userInfo.GetJob(),
			userInfo.GetUserId(),
			unit,
			dispatcherOkay,
		) {
			return true
		}
	}

	return false
}

func (s *Helpers) CheckIfUserPartOfUnit(
	ctx context.Context,
	job string,
	userId int32,
	unit *centrum.Unit,
	dispatcherOkay bool,
) bool {
	// Check if user is a dispatcher
	if dispatcherOkay && s.CheckIfUserIsDispatcher(ctx, job, userId) {
		return true
	}

	for i := range unit.GetUsers() {
		if (unit.GetUsers()[i].GetUser() != nil && unit.GetUsers()[i].GetUser().GetUserId() == userId) ||
			unit.GetUsers()[i].GetUserId() == userId {
			return true
		}
	}

	return false
}
