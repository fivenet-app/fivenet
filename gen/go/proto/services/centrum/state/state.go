package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
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
	Dispatches *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]]

	UserIDToUnitID *xsync.MapOf[int32, uint64]
}

func New() *State {
	return &State{
		Settings:   xsync.NewMapOf[string, *dispatch.Settings](),
		Disponents: xsync.NewMapOf[string, []*users.UserShort](),
		Units:      xsync.NewMapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]](),
		Dispatches: xsync.NewMapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]](),

		UserIDToUnitID: xsync.NewMapOf[int32, uint64](),
	}
}
