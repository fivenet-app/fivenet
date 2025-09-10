//nolint:forcetypeassert // This is a generic cache implementation, so there should "never" be a wrong type assertion.
package cache

// A size-bounded concurrent cache that stores up to N key/value pairs.
// It uses github.com/puzpuzpuz/xsync to provide a lock-free concurrent map for
// data storage. The implementation is generic - it works with any comparable key type and
// any value type. Example:
//   cache := lrucache.New[string, int](100) // capacity 100
//   cache.Put("foo", 42)
//   v, ok := cache.Get("foo")
//   cache.Delete("foo")
//
// NOTE: This implementation focuses purely on size-bounded eviction. If you
// need additional features such as per-item TTL, metrics, or callbacks on
// eviction, they can be layered on top. Eviction is random, not true LRU.

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/puzpuzpuz/xsync/v4"
)

// pair stores the user value together with TTL metadata.
type pair[V any] struct {
	value     V
	expiresAt time.Time // zero means no TTL / immortal until size eviction
}

// LRUCache is a concurrency-safe, size-bounded cache with optional TTL per entry.
// Eviction is random when over capacity; no recency is tracked.
type LRUCache[K comparable, V any] struct {
	_        utils.NoCopy            // disallow copying of LRUCache
	capacity int                     // max items retained by cache
	store    *xsync.Map[K, *pair[V]] // key → *pair[V]
}

// NewLRUCache constructs a cache with a fixed positive capacity. A panic is
// raised if capacity ≤ 0.
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		panic("lrucache: capacity must be positive")
	}
	return &LRUCache[K, V]{
		capacity: capacity,
		store:    xsync.NewMap[K, *pair[V]](),
	}
}

// Put inserts or updates a key/value pair with an optional TTL. A ttl of 0
// means the item never expires (except by random eviction if over capacity).
// If inserting causes the cache to exceed its capacity, a random element is evicted.
func (c *LRUCache[K, V]) Put(key K, value V, ttl time.Duration) {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	if elem, ok := c.store.Load(key); ok {
		// Update existing entry
		elem.value = value
		elem.expiresAt = exp
		return
	}

	c.store.Store(key, &pair[V]{value: value, expiresAt: exp})

	if c.store.Size() > c.capacity {
		c.evictLRU()
	}
}

// Get retrieves a value if present and not expired. The boolean result tells
// whether the key was found and valid. No recency is tracked.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	var zero V
	elem, ok := c.store.Load(key)
	if !ok {
		return zero, false
	}

	if c.isExpired(elem) {
		c.Delete(key) // Drop stale entry
		return zero, false
	}

	return elem.value, true
}

// Delete removes a key/value pair regardless of TTL.
func (c *LRUCache[K, V]) Delete(key K) {
	c.store.Delete(key)
}

// Len reports the current number of items, including any yet-uncollected
// expired ones (Cheap O(1)).
func (c *LRUCache[K, V]) Len() int { return c.store.Size() }

// StartJanitor launches a background goroutine that periodically sweeps the
// cache for expired items. It stops when ctx is cancelled. If sweepInterval ≤ 0
// the janitor exits immediately.
func (c *LRUCache[K, V]) StartJanitor(ctx context.Context, sweepInterval time.Duration) {
	if sweepInterval <= 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(sweepInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.cleanupExpired()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// cleanupExpired scans the map and removes all expired entries. Safe to call
// concurrently with other API methods.
func (c *LRUCache[K, V]) cleanupExpired() {
	now := time.Now()
	c.store.Range(func(key K, elem *pair[V]) bool {
		if !elem.expiresAt.IsZero() && elem.expiresAt.Before(now) {
			c.store.Delete(key)
		}
		return true
	})
}

// evictLRU drops a random element (since xsync.Map has no order).
// This is not true LRU eviction.
func (c *LRUCache[K, V]) evictLRU() {
	// xsync.Map does not track recency, so we evict a random element.
	// For true LRU, a different data structure is needed.
	var evicted bool
	c.store.Range(func(key K, _ *pair[V]) bool {
		c.store.Delete(key)
		evicted = true
		return false // stop after one
	})
	_ = evicted // for clarity
}

// isExpired returns true if the pair has a nonzero expiration and is expired.
func (c *LRUCache[K, V]) isExpired(p *pair[V]) bool {
	return !p.expiresAt.IsZero() && time.Now().After(p.expiresAt)
}
