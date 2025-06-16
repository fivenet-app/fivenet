package centrumstate

import (
	"context"
	"fmt"
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

	settings    *store.Store[centrum.Settings, *centrum.Settings]
	dispatchers *store.Store[centrum.Dispatchers, *centrum.Dispatchers]
	units       *store.Store[centrum.Unit, *centrum.Unit]
	dispatches  *store.Store[centrum.Dispatch, *centrum.Dispatch]

	dispatchLocationsMutex *sync.RWMutex
	dispatchLocations      map[string]*coords.Coords[*centrum.Dispatch]
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

		dispatchLocationsMutex: &sync.RWMutex{},
		dispatchLocations:      map[string]*coords.Coords[*centrum.Dispatch]{},
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		settings, err := store.New[centrum.Settings, *centrum.Settings](ctxStartup, logger, p.JS, "centrum_settings")
		if err != nil {
			return err
		}

		dispatchers, err := store.New[centrum.Dispatchers, *centrum.Dispatchers](ctxStartup, logger, p.JS, "centrum_dispatchers")
		if err != nil {
			return err
		}

		units, err := store.New[centrum.Unit, *centrum.Unit](ctxStartup, logger, p.JS, "centrum_units")
		if err != nil {
			return err
		}

		dispatches, err := store.New[centrum.Dispatch, *centrum.Dispatch](ctxCancel, logger, p.JS, "centrum_dispatches")
		if err != nil {
			return err
		}

		if err := settings.Start(ctxCancel, false); err != nil {
			return err
		}
		s.settings = settings

		if err := dispatchers.Start(ctxCancel, false); err != nil {
			return err
		}
		s.dispatchers = dispatchers

		if err := units.Start(ctxCancel, false); err != nil {
			return err
		}
		s.units = units

		if err := dispatches.Start(ctxCancel, false); err != nil {
			return err
		}
		s.dispatches = dispatches

		if err := s.runWatchKV(ctxCancel); err != nil {
			return fmt.Errorf("failed to run watch KV. %w", err)
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s, nil
}

func (s *State) runWatchKV(ctx context.Context) error {
	watchCh, err := s.dispatches.WatchAll(ctx)
	if err != nil {
		return err
	}

	go s.watchKV(ctx, watchCh)
	return nil
}

func (s *State) watchKV(ctx context.Context, watchCh <-chan *store.KeyValueEntry[centrum.Dispatch, *centrum.Dispatch]) {
	for {
		select {
		case <-ctx.Done():
			return

		case event := <-watchCh:
			if event == nil {
				s.logger.Error("received nil user changes event, skipping")
				continue
			}

			if event.Operation() == jetstream.KeyValuePut {
				dsp, err := event.Value()
				if err != nil {
					s.logger.Error("failed to get dispatch from store event", zap.Error(err))
					continue
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
							s.logger.Error("failed to add non-existant dispatch to locations", zap.Uint64("dispatch_id", dsp.Id))
						}
					}
				}
			} else if event.Operation() == jetstream.KeyValueDelete || event.Operation() == jetstream.KeyValuePurge {
				key := event.Key()
				if key == "" {
					s.logger.Warn("unable to delete dispatch location, got nil dispatch item", zap.String("store_dispatch_key", key))
					continue
				}

				split := strings.Split(key, ".")
				if len(split) < 2 {
					s.logger.Warn("unable to delete dispatch location, invalid key", zap.String("store_dispatch_key", key))
					continue
				}

				job := split[0]
				dspId, err := strconv.ParseUint(split[1], 10, 64)
				if err != nil {
					s.logger.Warn("unable to delete dispatch location, fallback to key failed", zap.String("store_dispatch_key", key))
					continue
				}

				if locs, ok := s.GetDispatchLocations(job); ok && locs != nil {
					locs.Remove(nil, centrum.DispatchPointMatchFn(dspId))
				}
			}
		}
	}
}

// Expose the stores for deeper interaction for updates

func (s *State) UnitsStore() *store.Store[centrum.Unit, *centrum.Unit] {
	return s.units
}

func (s *State) DispatchesStore() *store.Store[centrum.Dispatch, *centrum.Dispatch] {
	return s.dispatches
}
