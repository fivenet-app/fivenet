package store

import (
	"github.com/fivenet-app/fivenet/pkg/nats/locks"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
)

func WithKVPrefix[T any, U protoutils.ProtoMessageWithMerge[T]](prefix string) Option[T, U] {
	return func(s *Store[T, U]) {
		s.prefix = prefix
	}
}

func WithLocks[T any, U protoutils.ProtoMessageWithMerge[T]](l *locks.Locks) Option[T, U] {
	return func(s *Store[T, U]) {
		if l == nil {
			s.cl = false
		}

		s.l = l
	}
}

func WithOnUpdateFn[T any, U protoutils.ProtoMessageWithMerge[T]](fn OnUpdateFn[T, U]) Option[T, U] {
	return func(s *Store[T, U]) {
		s.onUpdate = fn
	}
}

func WithOnDeleteFn[T any, U protoutils.ProtoMessageWithMerge[T]](fn OnDeleteFn[T, U]) Option[T, U] {
	return func(s *Store[T, U]) {
		s.onDelete = fn
	}
}

func WithOnNotFoundFn[T any, U protoutils.ProtoMessageWithMerge[T]](fn OnNotFoundFn[T, U]) Option[T, U] {
	return func(s *Store[T, U]) {
		s.onNotFound = fn
	}
}

func WithJetstreamKV[T any, U protoutils.ProtoMessageWithMerge[T]](kv jetstream.KeyValue) Option[T, U] {
	return func(s *Store[T, U]) {
		s.kv = kv
	}
}
