package postals

import "github.com/paulmach/orb"

type Postal struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Code *string `json:"code"`
}

func (p *Postal) Point() orb.Point {
	return orb.Point{p.X, p.Y}
}
