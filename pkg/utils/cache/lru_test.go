package cache

import (
	"context"
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

	// Test LRU eviction
	cache.Put("c", 3, 0) // should evict "a"
	_, ok = cache.Get("a")
	assert.False(t, ok, "expected 'a' to be evicted")
	v, ok = cache.Get("c")
	assert.True(t, ok, "expected 'c' to exist")
	assert.Equal(t, 3, v)
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
