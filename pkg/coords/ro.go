package coords

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

type CoordsRO[V orb.Pointer] struct {
	ICoords[V]

	tree *quadtree.Quadtree
}

func NewReadOnly[V orb.Pointer](points []V) (*CoordsRO[V], error) {
	tree := quadtree.New(orb.Bound{Min: orb.Point{-9_000, -9_000}, Max: orb.Point{11_000, 11_000}})
	cs := &CoordsRO[V]{
		tree: tree,
	}

	for _, point := range points {
		cs.Add(point)
	}

	return cs, nil
}

func (p *CoordsRO[V]) Has(point orb.Pointer, fn quadtree.FilterFunc) bool {
	found := p.tree.Matching(point.Point(), fn)
	return found != nil
}

func (p *CoordsRO[V]) Closest(x, y float64) V {
	point := p.tree.Find(orb.Point{x, y})
	return point.(V)
}

func (p *CoordsRO[V]) KNearest(point orb.Pointer, max int, fn quadtree.FilterFunc, maxDistance float64) []orb.Pointer {
	points := p.tree.KNearestMatching(nil, point.Point(), max, fn, maxDistance)
	return points
}
