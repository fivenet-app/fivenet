package cache

import (
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
)

func WithKVPrefix[T any, U protoutils.ProtoMessageWithMerge[T]](prefix string) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.prefix = prefix
	}
}

func WithTTL[T any, U protoutils.ProtoMessageWithMerge[T]](ttl time.Duration) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.ttl = &ttl
	}
}

func WithJetstreamKV[T any, U protoutils.ProtoMessageWithMerge[T]](kv jetstream.KeyValue) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.kv = kv
	}
}

func WithIgnoredKeys[T any, U protoutils.ProtoMessageWithMerge[T]](keys ...string) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.ignoredKeys = keys
	}
}
