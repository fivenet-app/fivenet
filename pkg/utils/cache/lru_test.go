package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLRUCacheBasicOperations(t *testing.T) {
	cache := NewLRUCache[string, int](2)

	// Test Put and Get
	cache.Put("a", 1, 0)
	cache.Put("b", 2, 0)
	v, ok := cache.Get("a")
	assert.True(t, ok, "expected 'a' to exist")
	assert.Equal(t, 1, v)
	v, ok = cache.Get("b")
	assert.True(t, ok, "expected 'b' to exist")
	assert.Equal(t, 2, v)

	// Test random eviction
	cache.Put("c", 3, 0) // should evict one of "a" or "b"
	count := 0
	if _, ok := cache.Get("a"); ok {
		count++
	}
	if _, ok := cache.Get("b"); ok {
		count++
	}
	if _, ok := cache.Get("c"); ok {
		count++
	}
	assert.Equal(t, 2, count, "expected exactly 2 items to remain in the cache")
}

func TestLRUCacheDelete(t *testing.T) {
	cache := NewLRUCache[string, int](2)
	cache.Put("a", 1, 0)
	cache.Delete("a")
	_, ok := cache.Get("a")
	assert.False(t, ok, "expected 'a' to be deleted")
}

func TestLRUCacheTTL(t *testing.T) {
	cache := NewLRUCache[string, int](2)
	cache.Put("a", 1, 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, ok := cache.Get("a")
	assert.False(t, ok, "expected 'a' to be expired")
}

func TestLRUCacheJanitor(t *testing.T) {
	cache := NewLRUCache[string, int](2)
	cache.Put("a", 1, 10*time.Millisecond)
	ctx, cancel := context.WithCancel(t.Context())
	cache.StartJanitor(ctx, 5*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	cancel()
	_, ok := cache.Get("a")
	assert.False(t, ok, "expected 'a' to be expired and cleaned by janitor")
}

func TestLRUCacheLen(t *testing.T) {
	cache := NewLRUCache[string, int](2)
	cache.Put("a", 1, 0)
	cache.Put("b", 2, 0)
	l := cache.Len()
	assert.Equal(t, 2, l, "expected len 2")
	cache.Put("c", 3, 0)
	l = cache.Len()
	assert.Equal(t, 2, l, "expected len 2 after eviction")
}

func TestLRUCacheEvictionRandomness(t *testing.T) {
	cache := NewLRUCache[int, int](1)
	cache.Put(1, 100, 0)
	cache.Put(2, 200, 0)
	v1, ok1 := cache.Get(1)
	v2, ok2 := cache.Get(2)
	// Only one should exist, but which is not defined
	assert.True(t, (ok1 && !ok2) || (!ok1 && ok2), "only one key should remain")
	if ok1 {
		assert.Equal(t, 100, v1)
	}
	if ok2 {
		assert.Equal(t, 200, v2)
	}
}

func TestLRUCacheNoTTL(t *testing.T) {
	cache := NewLRUCache[string, int](2)
	cache.Put("a", 1, 0)
	cache.Put("b", 2, 0)
	v, ok := cache.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	v, ok = cache.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, v)
}

func TestLRUCacheOverwrite(t *testing.T) {
	cache := NewLRUCache[string, int](2)
	cache.Put("a", 1, 0)
	cache.Put("a", 2, 0)
	v, ok := cache.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 2, v)
}

func TestLRUCacheConcurrentAccess(t *testing.T) {
	cache := NewLRUCache[int, int](100)
	wg := sync.WaitGroup{}
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Put(i, i*10, 0)
		}(i)
	}
	wg.Wait()
	for i := range 100 {
		v, ok := cache.Get(i)
		if ok {
			assert.Equal(t, i*10, v)
		}
	}
}

func TestLRUCacheDeleteNonexistent(t *testing.T) {
	cache := NewLRUCache[string, int](2)
	cache.Delete("notfound")
	// Should not panic or error
}
