package cache

// An LRU (least‑recently‑used) cache that stores up to N key/value pairs.
// It uses github.com/puzpuzpuz/xsync to provide a lock‑free concurrent map for
// data storage and a classic doubly‑linked list (container/list) to track
// recency of use. All list manipulations are protected by a small mutex, while
// the xsync.MapOf shard‑based structure allows high read/write throughput on
// the key/value store itself.
//
// The implementation is generic – it works with any comparable key type and
// any value type. Example:
//   cache := lrucache.New[string, int](100) // capacity 100
//   cache.Put("foo", 42)
//   v, ok := cache.Get("foo")
//   cache.Delete("foo")
//
// The zero allocation behaviour of xsync.MapOf keeps the fast path hot, and
// locking is only taken when we have to update the recency list.
//
// NOTE: This implementation focuses purely on size‑bounded eviction. If you
// need additional features such as per‑item TTL, metrics, or callbacks on
// eviction, they can be layered on top.

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/puzpuzpuz/xsync/v4"
)

// pair stores the user value together with TTL metadata inside the recency
// list element.
type pair[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time // zero means no TTL / immortal until size eviction
}

type LRUCache[K comparable, V any] struct {
	_        utils.NoCopy                 // disallow copying of LRUCache
	capacity int                          // max items retained by LRU policy
	store    *xsync.Map[K, *list.Element] // key → *list.Element(pair)
	list     *list.List                   // recency list; front == MRU
	mu       sync.Mutex                   // protects list manipulations
}

// NewLRUCache constructs an LRU cache with a fixed positive capacity. A panic is
// raised if capacity ≤ 0.
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		panic("lrucache: capacity must be positive")
	}
	return &LRUCache[K, V]{
		capacity: capacity,
		store:    xsync.NewMap[K, *list.Element](),
		list:     list.New(),
	}
}

// Put inserts or updates a key/value pair with an optional TTL. A ttl of 0
// means the item never expires (except by LRU eviction). If inserting causes
// the cache to exceed its capacity the least‑recently‑used element is evicted.
func (c *LRUCache[K, V]) Put(key K, value V, ttl time.Duration) {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.store.Load(key); ok {
		// update existing entry
		p := elem.Value.(*pair[K, V])
		p.value = value
		p.expiresAt = exp
		c.list.MoveToFront(elem)
		return
	}

	elem := c.list.PushFront(&pair[K, V]{key: key, value: value, expiresAt: exp})
	c.store.Store(key, elem)

	if c.list.Len() > c.capacity {
		c.evictLRU()
	}
}

// Get retrieves a value if present and not expired. The boolean result tells
// whether the key was found and valid. A successful Get promotes the item to
// most‑recently‑used.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	var zero V
	elem, ok := c.store.Load(key)
	if !ok {
		return zero, false
	}
	p := elem.Value.(*pair[K, V])

	if c.isExpired(p) {
		c.deleteElement(elem) // drop stale entry
		return zero, false
	}

	c.mu.Lock()
	c.list.MoveToFront(elem)
	c.mu.Unlock()

	return p.value, true
}

// Delete removes a key/value pair regardless of TTL.
func (c *LRUCache[K, V]) Delete(key K) {
	elem, ok := c.store.Load(key)
	if !ok {
		return
	}
	c.deleteElement(elem)
}

// Len reports the current number of items, including any yet‑uncollected
// expired ones (Cheap O(1)).
func (c *LRUCache[K, V]) Len() int { return c.list.Len() }

// StartJanitor launches a background goroutine that periodically sweeps the
// cache for expired items. It stops when ctx is cancelled. If sweepInterval ≤ 0
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

// cleanupExpired scans the list and removes all expired entries. Safe to call
// concurrently with other API methods.
func (c *LRUCache[K, V]) cleanupExpired() {
	now := time.Now()
	var stale []K

	c.mu.Lock()
	for elem := c.list.Back(); elem != nil; {
		prev := elem.Prev()
		p := elem.Value.(*pair[K, V])
		if !p.expiresAt.IsZero() && p.expiresAt.Before(now) {
			c.list.Remove(elem)
			stale = append(stale, p.key)
		}
		elem = prev
	}
	c.mu.Unlock()

	// Do the map deletes outside the list lock.
	for _, k := range stale {
		c.store.Delete(k)
	}
}

// deleteElement removes elem from both list and map. Callers need *not* hold
// c.mu as we lock internally to guard the list.
func (c *LRUCache[K, V]) deleteElement(elem *list.Element) {
	p := elem.Value.(*pair[K, V])
	c.mu.Lock()
	c.list.Remove(elem)
	c.mu.Unlock()
	c.store.Delete(p.key)
}

// evictLRU drops the least‑recently‑used list element. Callers must hold c.mu.
func (c *LRUCache[K, V]) evictLRU() {
	tail := c.list.Back()
	if tail == nil {
		return
	}
	p := tail.Value.(*pair[K, V])
	c.list.Remove(tail)
	c.store.Delete(p.key)
}

func (c *LRUCache[K, V]) isExpired(p *pair[K, V]) bool {
	return !p.expiresAt.IsZero() && time.Now().After(p.expiresAt)
}
