package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/coords"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
)

var StateModule = fx.Module("centrum_state", fx.Provide(
	New,
))

type State struct {
	settings   *xsync.MapOf[string, *dispatch.Settings]
	disponents *xsync.MapOf[string, []*users.UserShort]

	units      *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]]
	unitsLocks *xsync.MapOf[uint64, *sync.Mutex]

	dispatches        *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]]
	dispatchLocations map[string]*coords.Coords[*dispatch.Dispatch]

	userIDToUnitID *xsync.MapOf[int32, uint64]
}

func New(cfg *config.Config) *State {
	locs := map[string]*coords.Coords[*dispatch.Dispatch]{}
	for _, job := range cfg.Game.Livemap.Jobs {
		locs[job] = coords.New[*dispatch.Dispatch]()
	}

	s := &State{
		settings:   xsync.NewMapOf[string, *dispatch.Settings](),
		disponents: xsync.NewMapOf[string, []*users.UserShort](),

		units:      xsync.NewMapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]](),
		unitsLocks: xsync.NewMapOf[uint64, *sync.Mutex](),

		dispatches:        xsync.NewMapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]](),
		dispatchLocations: locs,

		userIDToUnitID: xsync.NewMapOf[int32, uint64](),
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)

			units, ok := s.units.Load("ambulance")
			if !ok {
				continue
			}

			unit, ok := units.Load(5)
			if !ok {
				continue
			}

			fmt.Printf("UNIT %d, %s - %p\n", len(unit.Users), time.Now().String(), unit)
			for _, user := range unit.Users {
				fmt.Printf("UNIT User - %d %d - %p\n", user.UnitId, user.UserId, user)
			}
		}
	}()

	return s
}
