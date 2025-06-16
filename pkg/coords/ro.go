package coords

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

type ICoordsRO[V orb.Pointer] interface {
	Has(orb.Pointer, quadtree.FilterFunc) bool
	Get(orb.Pointer, quadtree.FilterFunc) V

	Closest(float64, float64) (V, bool)
	KNearest(orb.Pointer, int, quadtree.FilterFunc, float64) []orb.Pointer
}

// CoordsRO provides a read-only wrapper around a quadtree for spatial queries.
type CoordsRO[V orb.Pointer] struct {
	// Embeds the ICoords interface for generic spatial operations
	ICoords[V]

	// tree is the underlying quadtree for spatial indexing
	tree *quadtree.Quadtree
}

// NewReadOnly creates a new read-only coordinate store from a slice of points.
// Returns an error if any point cannot be added to the quadtree.
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

// Has returns true if a point matching the filter exists in the quadtree.
func (p *CoordsRO[V]) Has(point orb.Pointer, fn quadtree.FilterFunc) bool {
	found := p.tree.Matching(point.Point(), fn)
	return found != nil
}

// Get returns the value of a point matching the filter from the quadtree.
func (p *CoordsRO[V]) Get(point orb.Pointer, fn quadtree.FilterFunc) V {
	found := p.tree.Matching(point.Point(), fn)
	return found.(V)
}

// Closest returns the closest point to the given coordinates, if any.
func (p *CoordsRO[V]) Closest(x, y float64) (V, bool) {
	point := p.tree.Find(orb.Point{x, y})
	if point == nil {
		return *(new(V)), false
	}
	return point.(V), true
}

// KNearest returns up to max nearest points to the given point, filtered and within maxDistance.
func (p *CoordsRO[V]) KNearest(point orb.Pointer, max int, fn quadtree.FilterFunc, maxDistance float64) []orb.Pointer {
	points := p.tree.KNearestMatching(nil, point.Point(), max, fn, maxDistance)
	return points
}
