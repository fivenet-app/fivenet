package coords

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

type CoordsRO[V orb.Pointer] struct {
	ICoords[V]

	tree *quadtree.Quadtree
}

func NewReadOnly[V orb.Pointer](points []V) (*CoordsRO[V], error) {
	tree := quadtree.New(orb.Bound{Min: orb.Point{-9_000, -9_000}, Max: orb.Point{11_000, 11_000}})

	for _, point := range points {
		if err := tree.Add(point); err != nil {
			return nil, fmt.Errorf("failed to add point (%+v) to tree. %w", point, err)
		}
	}

	return &CoordsRO[V]{
		tree: tree,
	}, nil
}

func (p *CoordsRO[V]) Has(point orb.Pointer, fn quadtree.FilterFunc) bool {
	found := p.tree.Matching(point.Point(), fn)
	return found != nil
}

func (p *CoordsRO[V]) Get(point orb.Pointer, fn quadtree.FilterFunc) V {
	found := p.tree.Matching(point.Point(), fn)
	return found.(V)
}

func (p *CoordsRO[V]) Closest(x, y float64) V {
	point := p.tree.Find(orb.Point{x, y})
	return point.(V)
}

func (p *CoordsRO[V]) KNearest(point orb.Pointer, max int, fn quadtree.FilterFunc, maxDistance float64) []orb.Pointer {
	points := p.tree.KNearestMatching(nil, point.Point(), max, fn, maxDistance)
	return points
}
