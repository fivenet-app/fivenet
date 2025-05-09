package cache

import (
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
)

func WithKVPrefix[T any, U protoutils.ProtoMessageWithMerge[T]](prefix string) Option[T, U] {
	return func(s *Cache[T, U]) {
		s.prefix = prefix
	}
}

func WithTTL[T any, U protoutils.ProtoMessageWithMerge[T]](ttl time.Duration) Option[T, U] {
	return func(s *Cache[T, U]) {
		s.ttl = &ttl
	}
}
