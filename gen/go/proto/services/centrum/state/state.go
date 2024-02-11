package state

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/coords"
	"github.com/galexrt/fivenet/pkg/nats/store"
	"github.com/nats-io/nats.go"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var StateModule = fx.Module("centrum_state", fx.Provide(
	New,
))

type State struct {
	ctx context.Context
	js  nats.JetStreamContext

	logger *zap.Logger

	settings   *xsync.MapOf[string, *centrum.Settings]
	disponents *store.Store[centrum.Disponents, *centrum.Disponents]
	units      *store.Store[centrum.Unit, *centrum.Unit]
	dispatches *store.Store[centrum.Dispatch, *centrum.Dispatch]

	dispatchLocations map[string]*coords.Coords[*centrum.Dispatch]

	userIDToUnitID *store.Store[centrum.UserUnitMapping, *centrum.UserUnitMapping]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     nats.JetStreamContext
	Config *config.Config
}

func New(p Params) (*State, error) {
	locs := map[string]*coords.Coords[*centrum.Dispatch]{}
	for _, job := range p.Config.Game.Livemap.Jobs {
		locs[job] = coords.New[*centrum.Dispatch]()
	}

	disponents, err := store.New[centrum.Disponents, *centrum.Disponents](p.Logger, p.JS, "centrum_disponents")
	if err != nil {
		return nil, err
	}

	units, err := store.New[centrum.Unit, *centrum.Unit](p.Logger, p.JS, "centrum_units")
	if err != nil {
		return nil, err
	}

	dispatches, err := store.New[centrum.Dispatch, *centrum.Dispatch](p.Logger, p.JS, "centrum_dispatches")
	if err != nil {
		return nil, err
	}

	userIDToUnitID, err := store.New[centrum.UserUnitMapping, *centrum.UserUnitMapping](p.Logger, p.JS, "centrum_usersunits")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		if err := userIDToUnitID.Start(ctx); err != nil {
			return err
		}

		if err := disponents.Start(ctx); err != nil {
			return err
		}

		if err := units.Start(ctx); err != nil {
			return err
		}

		if err := dispatches.Start(ctx); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	s := &State{
		ctx: ctx,
		js:  p.JS,

		logger: p.Logger.Named("centrum_state"),

		settings:   xsync.NewMapOf[string, *centrum.Settings](),
		disponents: disponents,

		units: units,

		dispatches:        dispatches,
		dispatchLocations: locs,

		userIDToUnitID: userIDToUnitID,
	}

	return s, nil
}

// Expose the stores for deeper interaction with updates

func (s *State) UnitsStore() *store.Store[centrum.Unit, *centrum.Unit] {
	return s.units
}

func (s *State) DispatchesStore() *store.Store[centrum.Dispatch, *centrum.Dispatch] {
	return s.dispatches
}
