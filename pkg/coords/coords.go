package coords

import (
	"fmt"
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"github.com/puzpuzpuz/xsync/v3"
)

type CoordsKeyFn[K comparable] func(orb.Point) K

type Coords[K comparable, V any] struct {
	mutex       sync.Mutex
	tree        *quadtree.Quadtree
	targets     *xsync.MapOf[K, V]
	coordsKeyFn CoordsKeyFn[K]
}

func New[K comparable, V any](fn CoordsKeyFn[K]) *Coords[K, V] {
	tree := quadtree.New(orb.Bound{Min: orb.Point{-4000, -4000}, Max: orb.Point{8000, 6000}})
	return &Coords[K, V]{
		mutex:       sync.Mutex{},
		tree:        tree,
		targets:     xsync.NewMapOf[K, V](),
		coordsKeyFn: fn,
	}
}

func (p *Coords[K, V]) Has(x float64, y float64) (ok bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok = p.targets.Load(p.coordsKeyFn(orb.Point{x, y}))
	return
}

func (p *Coords[K, V]) Add(x float64, y float64, data V) {
	point := orb.Point{x, y}
	p.targets.Store(p.coordsKeyFn(point), data)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.tree.Add(point)
}

func (p *Coords[K, V]) Remove(x float64, y float64, fn quadtree.FilterFunc) bool {
	point := orb.Point{x, y}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.targets.Delete(p.coordsKeyFn(point))
	return p.tree.Remove(point, fn)
}

func (p *Coords[K, V]) Closest(x, y float64) (V, bool) {
	nearest := p.tree.Find(orb.Point{x, y})
	key := p.coordsKeyFn(nearest.Point())

	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.targets.Load(key)
}

func (p *Coords[K, V]) KNearest(point orb.Point, max int, maxDistance float64) []V {
	points := p.tree.KNearest(nil, point, max, maxDistance)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	results := make([]V, len(points))
	for i := 0; i < len(points); i++ {
		key := p.coordsKeyFn(points[i].Point())

		val, ok := p.targets.Load(key)
		if ok {
			results[i] = val
		}
	}

	return results
}

func GetCoordsKey(point orb.Point) string {
	return fmt.Sprintf("%f-%f", point.X(), point.Y())
}
