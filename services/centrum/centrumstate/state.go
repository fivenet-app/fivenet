package centrumstate

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/paulmach/orb"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var StateModule = fx.Module("centrum_state",
	fx.Provide(
		New,
	))

type State struct {
	js *events.JSWrapper

	logger *zap.Logger

	settings   *xsync.Map[string, *centrum.Settings]
	disponents *store.Store[centrum.Disponents, *centrum.Disponents]
	units      *store.Store[centrum.Unit, *centrum.Unit]
	dispatches *store.Store[centrum.Dispatch, *centrum.Dispatch]

	dispatchLocationsMutex *sync.RWMutex
	dispatchLocations      map[string]*coords.Coords[*centrum.Dispatch]

	userIDToUnitID *store.Store[centrum.UserUnitMapping, *centrum.UserUnitMapping]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
}

func New(p Params) (*State, error) {
	logger := p.Logger.Named("centrum_state")

	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &State{
		js: p.JS,

		logger: logger,

		settings: xsync.NewMap[string, *centrum.Settings](),

		dispatchLocationsMutex: &sync.RWMutex{},
		dispatchLocations:      map[string]*coords.Coords[*centrum.Dispatch]{},
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		disponents, err := store.New[centrum.Disponents, *centrum.Disponents](ctxStartup, logger, p.JS, "centrum_disponents")
		if err != nil {
			return err
		}

		units, err := store.New[centrum.Unit, *centrum.Unit](ctxStartup, logger, p.JS, "centrum_units")
		if err != nil {
			return err
		}

		userIDToUnitID, err := store.New[centrum.UserUnitMapping, *centrum.UserUnitMapping](ctxStartup, logger, p.JS, "centrum_usersunits")
		if err != nil {
			return err
		}

		dispatches, err := store.New[centrum.Dispatch, *centrum.Dispatch](ctxCancel, logger, p.JS, "centrum_dispatches",
			store.WithOnUpdateFn(func(_ *store.Store[centrum.Dispatch, *centrum.Dispatch], dsp *centrum.Dispatch) (*centrum.Dispatch, error) {
				if dsp == nil {
					logger.Warn("unable to update dispatch location, got nil dispatch")
					return dsp, nil
				}

				if locs, ok := s.GetDispatchLocations(dsp.Job); ok && locs != nil {
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
			}),

			store.WithOnDeleteFn(func(_ *store.Store[centrum.Dispatch, *centrum.Dispatch], entry jetstream.KeyValueEntry, dsp *centrum.Dispatch) error {
				key := entry.Key()
				if dsp == nil && key == "" {
					logger.Warn("unable to delete dispatch location, got nil dispatch item", zap.String("store_dispatch_key", entry.Key()))
					return nil
				}

				var job string
				var dspId uint64
				if dsp != nil {
					job = dsp.Job
					dspId = dsp.Id
				} else {
					split := strings.Split(key, ".")
					if len(split) < 2 {
						logger.Warn("unable to delete dispatch location, invalid key", zap.String("store_dispatch_key", entry.Key()))
						return nil
					}

					job = split[0]
					dspId, err = strconv.ParseUint(split[1], 10, 64)
					if err != nil {
						logger.Warn("unable to delete dispatch location, fallback to key failed", zap.String("store_dispatch_key", entry.Key()))
						return nil
					}
				}

				if locs, ok := s.GetDispatchLocations(job); ok && locs != nil {
					if locs.Has(dsp, centrum.DispatchPointMatchFn(dspId)) {
						locs.Remove(dsp, centrum.DispatchPointMatchFn(dspId))
					}
				}

				return nil
			}),
		)
		if err != nil {
			return err
		}

		if err := userIDToUnitID.Start(ctxCancel, false); err != nil {
			return err
		}
		s.userIDToUnitID = userIDToUnitID

		if err := disponents.Start(ctxCancel, false); err != nil {
			return err
		}
		s.disponents = disponents

		if err := units.Start(ctxCancel, false); err != nil {
			return err
		}
		s.units = units

		if err := dispatches.Start(ctxCancel, false); err != nil {
			return err
		}
		s.dispatches = dispatches

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s, nil
}

// Expose the stores for deeper interaction with updates

func (s *State) UnitsStore() *store.Store[centrum.Unit, *centrum.Unit] {
	return s.units
}

func (s *State) DispatchesStore() *store.Store[centrum.Dispatch, *centrum.Dispatch] {
	return s.dispatches
}
