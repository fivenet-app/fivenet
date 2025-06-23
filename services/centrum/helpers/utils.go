package helpers

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Helpers struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	tracker  tracker.ITracker
	postals  postals.Postals
	appCfg   appconfig.IConfig

	settings    *settings.SettingsDB
	dispatchers *dispatchers.DispatchersDB
	units       *units.UnitDB

	dispatchLocationsMutex *sync.Mutex
	dispatchLocations      map[string]*coords.Coords[*centrum.Dispatch]

	store *store.Store[centrum.Dispatch, *centrum.Dispatch]
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
	if !settings.Enabled {
		return false
	}

	if settings.Mode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
		return true
	}

	dispatchers, err := s.dispatchers.Get(ctx, job)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return false
	}

	if dispatchers.IsEmpty() {
		if settings.FallbackMode == centrum.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN {
			return true
		}
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
	for i := range dispatchers.Dispatchers {
		if userId == dispatchers.Dispatchers[i].UserId {
			return true
		}
	}

	return false
}

func (s *Helpers) CheckIfUserIsPartOfDispatch(ctx context.Context, userInfo *userinfo.UserInfo, dsp *centrum.Dispatch, dispatcherOkay bool) bool {
	// Check if user is a dispatcher
	if dispatcherOkay && s.CheckIfUserIsDispatcher(ctx, userInfo.Job, userInfo.UserId) {
		return true
	}

	// Iterate over units of dispatch and check if the user is in one of the units
	for i := range dsp.Units {
		unit, err := s.units.Get(ctx, dsp.Units[i].UnitId)
		if unit == nil || err != nil {
			continue
		}

		if s.CheckIfUserPartOfUnit(ctx, userInfo.Job, userInfo.UserId, unit, dispatcherOkay) {
			return true
		}
	}

	return false
}

func (s *Helpers) CheckIfUserPartOfUnit(ctx context.Context, job string, userId int32, unit *centrum.Unit, dispatcherOkay bool) bool {
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
