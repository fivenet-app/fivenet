package cache

import (
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
)

// WithKVPrefix sets a prefix for all keys in the cache.
func WithKVPrefix[T any, U protoutils.ProtoMessageWithMerge[T]](prefix string) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.prefix = prefix
	}
}

// WithTTL sets a time-to-live for cache entries.
func WithTTL[T any, U protoutils.ProtoMessageWithMerge[T]](ttl time.Duration) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.ttl = &ttl
	}
}

// WithJetstreamKV sets a custom NATS JetStream KeyValue store instance for the cache.
func WithJetstreamKV[T any, U protoutils.ProtoMessageWithMerge[T]](
	kv jetstream.KeyValue,
) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.kv = kv
	}
}

// WithIgnoredKeys sets a list of keys to ignore in the cache.
func WithIgnoredKeys[T any, U protoutils.ProtoMessageWithMerge[T]](keys ...string) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.ignoredKeys = keys
	}
}
