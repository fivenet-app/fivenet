package state

import (
	"sync"

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
	Settings   *xsync.MapOf[string, *dispatch.Settings]
	Disponents *xsync.MapOf[string, []*users.UserShort]

	Units      *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]]
	UnitsLocks *xsync.MapOf[uint64, *sync.Mutex]

	Dispatches        *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]]
	DispatchLocations map[string]*coords.Coords[*dispatch.Dispatch]

	userIDToUnitID *xsync.MapOf[int32, uint64]
}

func New(cfg *config.Config) *State {
	locs := map[string]*coords.Coords[*dispatch.Dispatch]{}
	for _, job := range cfg.Game.Livemap.Jobs {
		locs[job] = coords.New[*dispatch.Dispatch]()
	}

	return &State{
		Settings:   xsync.NewMapOf[string, *dispatch.Settings](),
		Disponents: xsync.NewMapOf[string, []*users.UserShort](),

		Units:      xsync.NewMapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]](),
		UnitsLocks: xsync.NewMapOf[uint64, *sync.Mutex](),

		Dispatches:        xsync.NewMapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]](),
		DispatchLocations: locs,

		userIDToUnitID: xsync.NewMapOf[int32, uint64](),
	}
}
