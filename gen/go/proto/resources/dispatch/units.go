package dispatch

import (
	"dario.cat/mergo"
	"github.com/paulmach/orb"
)

func (x *Unit) Update(in *Unit) {
	if x.Id != in.Id {
		return
	}

	err := mergo.Merge(x, in, mergo.WithOverride)
	if err != nil {
		return
	}
}

func (x *UnitStatus) Point() orb.Point {
	if x.X == nil || x.Y == nil {
		return orb.Point{}
	}

	return orb.Point{*x.X, *x.Y}
}
