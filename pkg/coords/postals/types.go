package postals

import "github.com/paulmach/orb"

// Postal represents a postal coordinate with an optional code.
type Postal struct {
	// X is the X coordinate (longitude)
	X float64 `json:"x"`
	// Y is the Y coordinate (latitude)
	Y float64 `json:"y"`
	// Code is the optional postal code
	Code *string `json:"code"`
}

// Point returns the orb.Point representation of the postal coordinates.
func (p *Postal) Point() orb.Point {
	return orb.Point{p.X, p.Y}
}
