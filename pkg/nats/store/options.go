package store

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/locks"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
)

// WithKVPrefix sets a prefix for all keys in the store.
func WithKVPrefix[T any, U protoutils.ProtoMessageWithMerge[T]](prefix string) Option[T, U] {
	return func(s *Store[T, U]) {
		s.prefix = prefix
	}
}

// WithLocks sets the Locks instance to use for distributed locking.
func WithLocks[T any, U protoutils.ProtoMessageWithMerge[T]](l *locks.Locks) Option[T, U] {
	return func(s *Store[T, U]) {
		if l == nil {
			s.cl = false
		}

		s.l = l
	}
}

// WithOnUpdateFn sets a callback function to be called on update events (only local).
func WithOnUpdateFn[T any, U protoutils.ProtoMessageWithMerge[T]](fn OnUpdateFn[T, U]) Option[T, U] {
	return func(s *Store[T, U]) {
		s.onUpdate = fn
	}
}

// WithOnDeleteFn sets a callback function to be called on delete events (only local).
func WithOnDeleteFn[T any, U protoutils.ProtoMessageWithMerge[T]](fn OnDeleteFn[T, U]) Option[T, U] {
	return func(s *Store[T, U]) {
		s.onDelete = fn
	}
}

// WithOnCreatedFn sets a callback function to be called on create events (only local).
func WithOnCreatedFn[T any, U protoutils.ProtoMessageWithMerge[T]](fn OnCreatedFn[T, U]) Option[T, U] {
	return func(s *Store[T, U]) {
		s.onCreated = fn
	}
}

// WithJetstreamKV sets a custom NATS JetStream KeyValue store instance.
func WithJetstreamKV[T any, U protoutils.ProtoMessageWithMerge[T]](kv jetstream.KeyValue) Option[T, U] {
	return func(s *Store[T, U]) {
		s.kv = kv
	}
}

// WithIgnoredKeys sets a list of keys to ignore in the store.
func WithIgnoredKeys[T any, U protoutils.ProtoMessageWithMerge[T]](keys ...string) Option[T, U] {
	return func(s *Store[T, U]) {
		s.ignoredKeys = keys
	}
}
