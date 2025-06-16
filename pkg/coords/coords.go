package coords

import (
	"fmt"
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

// ICoords defines the interface for coordinate storage and spatial queries.
type ICoords[V orb.Pointer] interface {
	ICoordsRO[V]

	Add(orb.Pointer) error
	Remove(orb.Pointer, quadtree.FilterFunc) bool
	Replace(orb.Pointer, quadtree.FilterFunc) error
}

// Coords provides a thread-safe wrapper around a quadtree for storing and querying coordinates.
type Coords[V orb.Pointer] struct {
	// Embeds the ICoords interface for generic spatial operations
	ICoords[V]

	// mutex protects concurrent access to the quadtree
	mutex sync.RWMutex
	// tree is the underlying quadtree for spatial indexing
	tree *quadtree.Quadtree
}

// CoordsEqualFn must return true if both points are equal.
type CoordsEqualFn = func(orb.Pointer, orb.Pointer) bool

// New creates a new Coords instance with default bounds.
func New[V orb.Pointer]() *Coords[V] {
	return NewWithBounds[V](orb.Bound{Min: orb.Point{-9_000, -9_000}, Max: orb.Point{11_000, 11_000}})
}

// NewWithBounds creates a new Coords instance with the specified bounds.
func NewWithBounds[V orb.Pointer](bounds orb.Bound) *Coords[V] {
	tree := quadtree.New(bounds)
	return &Coords[V]{
		mutex: sync.RWMutex{},
		tree:  tree,
	}
}

// Has returns true if a point matching the filter exists in the quadtree.
func (p *Coords[V]) Has(point orb.Pointer, fn quadtree.FilterFunc) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	found := p.tree.Matching(point, fn)
	return found != nil
}

// Get returns the value of a point matching the filter from the quadtree.
func (p *Coords[V]) Get(point orb.Pointer, fn quadtree.FilterFunc) V {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	found := p.tree.Matching(point, fn)
	return found.(V)
}

// Add inserts a point into the quadtree.
func (p *Coords[V]) Add(point orb.Pointer) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	err := p.tree.Add(point)
	return err
}

// Remove deletes a point matching the filter from the quadtree.
func (p *Coords[V]) Remove(point orb.Pointer, fn quadtree.FilterFunc) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ok := p.tree.Remove(point, fn)
	return ok
}

// Replace removes a point if it exists (and is not equal according to equalFn), then adds the new point.
func (p *Coords[V]) Replace(point orb.Pointer, fn quadtree.FilterFunc, equalFn CoordsEqualFn) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	found := p.tree.Matching(point.Point(), fn)
	if found != nil {
		if equalFn == nil || !equalFn(point, found) {
			p.tree.Remove(point, fn)
		}
	}

	if err := p.tree.Add(point); err != nil {
		return fmt.Errorf("failed to replace point in coords. %w", err)
	}

	return nil
}

// Closest returns the closest point to the given coordinates, if any.
func (p *Coords[V]) Closest(x, y float64) (V, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	point := p.tree.Find(orb.Point{x, y})
	if point == nil {
		return *(new(V)), false
	}
	return point.(V), true
}

// KNearest returns up to max nearest points to the given point, filtered and within maxDistance.
func (p *Coords[V]) KNearest(point orb.Point, max int, fn quadtree.FilterFunc, maxDistance float64) []orb.Pointer {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	points := p.tree.KNearestMatching(nil, point.Point(), max, fn, maxDistance)
	return points
}
