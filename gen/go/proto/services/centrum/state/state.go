package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/puzpuzpuz/xsync/v2"
	"go.uber.org/fx"
)

var StateModule = fx.Module("centrum_state", fx.Provide(
	NewState,
))

type State struct {
	Settings   *xsync.MapOf[string, *dispatch.Settings]
	Disponents *xsync.MapOf[string, []*users.UserShort]
	Units      *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]]
	Dispatches *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]]

	UserIDToUnitID *xsync.MapOf[int32, uint64]
}

func NewState() *State {
	return &State{
		Settings:   xsync.NewMapOf[*dispatch.Settings](),
		Disponents: xsync.NewMapOf[[]*users.UserShort](),
		Units:      xsync.NewMapOf[*xsync.MapOf[uint64, *dispatch.Unit]](),
		Dispatches: xsync.NewMapOf[*xsync.MapOf[uint64, *dispatch.Dispatch]](),

		UserIDToUnitID: xsync.NewIntegerMapOf[int32, uint64](),
	}
}
