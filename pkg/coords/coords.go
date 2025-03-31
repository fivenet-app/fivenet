package coords

import (
	"fmt"
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

type ICoords[V orb.Pointer] interface {
	Has(orb.Pointer, quadtree.FilterFunc) bool
	Get(orb.Pointer, quadtree.FilterFunc) V
	Add(orb.Pointer) error
	Remove(orb.Pointer, quadtree.FilterFunc) bool
	Replace(orb.Pointer, quadtree.FilterFunc) error

	Closest(float64, float64) (V, bool)
	KNearest(orb.Pointer, int, quadtree.FilterFunc, float64) []orb.Pointer
}

type Coords[V orb.Pointer] struct {
	ICoords[V]

	mutex sync.Mutex
	tree  *quadtree.Quadtree
}

// CoordsEqualFn must return true if both points are equal
type CoordsEqualFn = func(orb.Pointer, orb.Pointer) bool

func New[V orb.Pointer]() *Coords[V] {
	return NewWithBounds[V](orb.Bound{Min: orb.Point{-9_000, -9_000}, Max: orb.Point{11_000, 11_000}})
}

func NewWithBounds[V orb.Pointer](bounds orb.Bound) *Coords[V] {
	tree := quadtree.New(bounds)
	return &Coords[V]{
		mutex: sync.Mutex{},
		tree:  tree,
	}
}

func (p *Coords[V]) Has(point orb.Pointer, fn quadtree.FilterFunc) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	found := p.tree.Matching(point.Point(), fn)
	return found != nil
}

func (p *Coords[V]) Get(point orb.Pointer, fn quadtree.FilterFunc) V {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	found := p.tree.Matching(point.Point(), fn)
	return found.(V)
}

func (p *Coords[V]) Add(point orb.Pointer) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	err := p.tree.Add(point)
	return err
}

func (p *Coords[V]) Remove(point orb.Pointer, fn quadtree.FilterFunc) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ok := p.tree.Remove(point, fn)
	return ok
}

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

func (p *Coords[V]) Closest(x, y float64) (V, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	point := p.tree.Find(orb.Point{x, y})
	if point == nil {
		return *(new(V)), false
	}
	return point.(V), true
}

func (p *Coords[V]) KNearest(point orb.Pointer, max int, fn quadtree.FilterFunc, maxDistance float64) []orb.Pointer {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	points := p.tree.KNearestMatching(nil, point.Point(), max, fn, maxDistance)
	return points
}
