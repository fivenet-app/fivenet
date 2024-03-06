package state

import (
	"context"
	"sync"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/config/appconfig"
	"github.com/galexrt/fivenet/pkg/coords"
	"github.com/galexrt/fivenet/pkg/nats/store"
	"github.com/nats-io/nats.go"
	"github.com/paulmach/orb"
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

	dispatchLocationsMutex sync.RWMutex
	dispatchLocations      map[string]*coords.Coords[*centrum.Dispatch]

	userIDToUnitID *store.Store[centrum.UserUnitMapping, *centrum.UserUnitMapping]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	JS        nats.JetStreamContext
	AppConfig *appconfig.Config
}

func New(p Params) (*State, error) {
	logger := p.Logger.Named("centrum_state")

	dspLocs := map[string]*coords.Coords[*centrum.Dispatch]{}

	disponents, err := store.New[centrum.Disponents, *centrum.Disponents](logger, p.JS, "centrum_disponents")
	if err != nil {
		return nil, err
	}

	units, err := store.New[centrum.Unit, *centrum.Unit](logger, p.JS, "centrum_units")
	if err != nil {
		return nil, err
	}

	userIDToUnitID, err := store.New[centrum.UserUnitMapping, *centrum.UserUnitMapping](logger, p.JS, "centrum_usersunits")
	if err != nil {
		return nil, err
	}

	s := &State{
		js: p.JS,

		logger: logger,

		settings:   xsync.NewMapOf[string, *centrum.Settings](),
		disponents: disponents,

		units: units,

		dispatchLocationsMutex: sync.RWMutex{},
		dispatchLocations:      dspLocs,

		userIDToUnitID: userIDToUnitID,
	}

	dispatches, err := store.New[centrum.Dispatch, *centrum.Dispatch](logger, p.JS, "centrum_dispatches",
		func(st *store.Store[centrum.Dispatch, *centrum.Dispatch]) error {
			st.OnUpdate = func(dsp *centrum.Dispatch) (*centrum.Dispatch, error) {
				if dsp == nil {
					logger.Warn("unable to update dispatch location, got nil dispatch")
					return dsp, nil
				}

				s.dispatchLocationsMutex.Lock()
				defer s.dispatchLocationsMutex.Unlock()

				if locs := dspLocs[dsp.Job]; locs != nil {
					if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
						if locs.Has(dsp, centrum.DispatchPointMatchFn(dsp.Id)) {
							locs.Remove(dsp, centrum.DispatchPointMatchFn(dsp.Id))
						}
					} else {
						if err := locs.Replace(dsp, func(p orb.Pointer) bool {
							return p.(*centrum.Dispatch).Id == dsp.Id
						}, func(p1, p2 orb.Pointer) bool {
							return p1.Point().Equal(p2.Point())
						}); err != nil {
							logger.Error("failed to add non-existant dispatch to locations", zap.Uint64("dispatch_id", dsp.Id))
						}
					}
				}

				return dsp, nil
			}
			return nil
		},
		func(st *store.Store[centrum.Dispatch, *centrum.Dispatch]) error {
			st.OnDelete = func(entry nats.KeyValueEntry, dsp *centrum.Dispatch) error {
				if dsp == nil {
					logger.Warn("unable to delete dispatch location, got nil dispatch item", zap.String("store_dispatch_key", entry.Key()))
					return nil
				}

				s.dispatchLocationsMutex.Lock()
				defer s.dispatchLocationsMutex.Unlock()

				if loc := dspLocs[dsp.Job]; loc != nil {
					if loc.Has(dsp, centrum.DispatchPointMatchFn(dsp.Id)) {
						loc.Remove(dsp, centrum.DispatchPointMatchFn(dsp.Id))
					}
				}

				return nil
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	s.dispatches = dispatches

	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		s.handleAppConfigUpdate(p.AppConfig.Get())

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

		// Handle app config updates
		go func() {
			configUpdateCh := p.AppConfig.Subscribe()
			for {
				select {
				case <-s.ctx.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					s.handleAppConfigUpdate(cfg)
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s, nil
}

func (s *State) handleAppConfigUpdate(appCfg *appconfig.Cfg) {
	s.dispatchLocationsMutex.Lock()
	defer s.dispatchLocationsMutex.Unlock()

	for _, job := range appCfg.UserTracker.LivemapJobs {
		if _, ok := s.dispatchLocations[job]; !ok {
			s.dispatchLocations[job] = coords.New[*centrum.Dispatch]()
		}
	}
}

// Expose the stores for deeper interaction with updates

func (s *State) UnitsStore() *store.Store[centrum.Unit, *centrum.Unit] {
	return s.units
}

func (s *State) DispatchesStore() *store.Store[centrum.Dispatch, *centrum.Dispatch] {
	return s.dispatches
}
