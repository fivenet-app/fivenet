package livemap

import (
	"slices"
	"sync"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
)

type markerPublicCache struct {
	mu      sync.RWMutex
	active  []*livemapmarkers.MarkerMarker
	deleted []int64
}

func newMarkerPublicCache() *markerPublicCache {
	return &markerPublicCache{
		active:  []*livemapmarkers.MarkerMarker{},
		deleted: []int64{},
	}
}

func (c *markerPublicCache) Snapshot() ([]*livemapmarkers.MarkerMarker, []int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.active) == 0 {
		return []*livemapmarkers.MarkerMarker{}, slices.Clone(c.deleted)
	}

	return slices.Clone(c.active), slices.Clone(c.deleted)
}

func (c *markerPublicCache) Replace(active []*livemapmarkers.MarkerMarker, deleted []int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.active = slices.Clone(active)
	c.deleted = slices.Clone(deleted)
}

func (c *markerPublicCache) Apply(marker *livemapmarkers.MarkerMarker) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := marker.GetId()
	c.active = slices.DeleteFunc(c.active, func(existing *livemapmarkers.MarkerMarker) bool {
		return existing.GetId() == id
	})
	c.deleted = slices.DeleteFunc(c.deleted, func(existing int64) bool {
		return existing == id
	})

	if !marker.GetPublic() {
		return
	}

	if marker.GetDeletedAt() == nil {
		c.active = append(c.active, marker)
		return
	}

	c.deleted = append(c.deleted, id)
}
