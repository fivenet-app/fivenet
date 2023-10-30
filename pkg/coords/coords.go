package coords

import (
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

type Coords[K comparable, V any] struct {
	mutex sync.RWMutex
	tree  *quadtree.Quadtree
}

func New[K comparable, V any]() *Coords[K, V] {
	tree := quadtree.New(orb.Bound{Min: orb.Point{-4000, -4000}, Max: orb.Point{8000, 6000}})
	return &Coords[K, V]{
		mutex: sync.RWMutex{},
		tree:  tree,
	}
}

func (p *Coords[K, V]) Has(point orb.Pointer) (ok bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	found := p.tree.Find(point.Point())
	return found != nil
}

func (p *Coords[K, V]) Add(point orb.Pointer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.tree.Add(point)
}

func (p *Coords[K, V]) Remove(point orb.Pointer, fn quadtree.FilterFunc) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.tree.Remove(point, fn)
}

func (p *Coords[K, V]) Closest(x, y float64) V {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.tree.Find(orb.Point{x, y}).(V)
}

func (p *Coords[K, V]) KNearest(point orb.Point, max int, maxDistance float64) []orb.Pointer {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	points := p.tree.KNearest(nil, point, max, maxDistance)
	return points
}
