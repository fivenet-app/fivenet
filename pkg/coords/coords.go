package coords

import (
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

type ICoords[V orb.Pointer] interface {
	Has(orb.Pointer, quadtree.FilterFunc) bool
	Add(orb.Pointer) error
	Remove(orb.Pointer, quadtree.FilterFunc) bool
	Closest(float64, float64) V
	KNearest(orb.Pointer, int, quadtree.FilterFunc, float64) []orb.Pointer
}

type Coords[V orb.Pointer] struct {
	mutex sync.RWMutex
	tree  *quadtree.Quadtree

	ICoords[V]
}

func New[V orb.Pointer]() *Coords[V] {
	return NewWithBounds[V](orb.Bound{Min: orb.Point{-9_000, -9_000}, Max: orb.Point{11_000, 11_000}})
}

func NewWithBounds[V orb.Pointer](bounds orb.Bound) *Coords[V] {
	tree := quadtree.New(bounds)
	return &Coords[V]{
		mutex: sync.RWMutex{},
		tree:  tree,
	}
}

func (p *Coords[V]) Has(point orb.Pointer, fn quadtree.FilterFunc) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	found := p.tree.Matching(point.Point(), fn)
	return found != nil
}

func (p *Coords[V]) Add(point orb.Pointer) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.tree.Add(point)
}

func (p *Coords[V]) Remove(point orb.Pointer, fn quadtree.FilterFunc) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ok := p.tree.Remove(point, fn)
	return ok
}

func (p *Coords[V]) Closest(x, y float64) V {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	point := p.tree.Find(orb.Point{x, y})
	return point.(V)
}

func (p *Coords[V]) KNearest(point orb.Pointer, max int, fn quadtree.FilterFunc, maxDistance float64) []orb.Pointer {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	points := p.tree.KNearestMatching(nil, point.Point(), max, fn, maxDistance)
	return points
}
