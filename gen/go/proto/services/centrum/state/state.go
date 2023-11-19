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
	js nats.JetStreamContext

	settings   *xsync.MapOf[string, *centrum.Settings]
	disponents *store.Store[centrum.Disponents, *centrum.Disponents]
	units      *store.Store[centrum.Unit, *centrum.Unit]
	dispatches *store.Store[centrum.Dispatch, *centrum.Dispatch]

	dispatchLocations map[string]*coords.Coords[*centrum.Dispatch]

	userIDToUnitID *store.Store[centrum.UserUnitMapping, *centrum.UserUnitMapping]
}

func New(logger *zap.Logger, js nats.JetStreamContext, cfg *config.Config) (*State, error) {
	locs := map[string]*coords.Coords[*centrum.Dispatch]{}
	for _, job := range cfg.Game.Livemap.Jobs {
		locs[job] = coords.New[*centrum.Dispatch]()
	}

	disponents, err := store.New[centrum.Disponents, *centrum.Disponents](logger, js, "centrum_disponents")
	if err != nil {
		return nil, err
	}

	units, err := store.New[centrum.Unit, *centrum.Unit](logger, js, "centrum_units")
	if err != nil {
		return nil, err
	}

	dispatches, err := store.New[centrum.Dispatch, *centrum.Dispatch](logger, js, "centrum_dispatches")
	if err != nil {
		return nil, err
	}

	userIDToUnitID, err := store.New[centrum.UserUnitMapping, *centrum.UserUnitMapping](logger, js, "centrum_usersunits")
	if err != nil {
		return nil, err
	}

	go userIDToUnitID.Start(context.TODO())
	go disponents.Start(context.TODO())
	go units.Start(context.TODO())
	go dispatches.Start(context.TODO())

	s := &State{
		js: js,

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
