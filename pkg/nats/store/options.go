package store

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/locks"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
)

// WithKVPrefix sets a prefix for all keys in the store.
func WithKVPrefix[T any, U protoutils.ProtoMessageWithMerge[T]](prefix string) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.prefix = prefix
	}
}

// WithLocks sets the Locks instance to use for distributed locking.
func WithLocks[T any, U protoutils.ProtoMessageWithMerge[T]](l *locks.Locks) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		if l == nil {
			s.cl = false
		}

		s.l = l
	}
}

// WithIgnoredKeys sets a list of keys to ignore in the store.
func WithIgnoredKeys[T any, U protoutils.ProtoMessageWithMerge[T]](keys ...string) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.ignoredKeys = keys
	}
}

// WithJetstreamKV sets a custom NATS JetStream KeyValue store instance.
func WithJetstreamKV[T any, U protoutils.ProtoMessageWithMerge[T]](
	kv jetstream.KeyValue,
) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.kv = kv
	}
}

// WithOnUpdateFn sets a callback function to be called on update events (only local).
func WithOnUpdateFn[T any, U protoutils.ProtoMessageWithMerge[T]](
	fn OnUpdateFn[T, U],
) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.onUpdate = fn
	}
}

// WithOnDeleteFn sets a callback function to be called on delete events (only local).
func WithOnDeleteFn[T any, U protoutils.ProtoMessageWithMerge[T]](
	fn OnDeleteFn[T, U],
) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.onDelete = fn
	}
}

// WithOnCreatedFn sets a callback function to be called on create events (only local).
func WithOnCreatedFn[T any, U protoutils.ProtoMessageWithMerge[T]](
	fn OnCreatedFn[T, U],
) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.onCreated = fn
	}
}

// WithOnRemoteUpdatedFn sets a callback function to be called on remove update events (only local).
func WithOnRemoteUpdatedFn[T any, U protoutils.ProtoMessageWithMerge[T]](
	fn OnUpdateFn[T, U],
) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.onRemoteUpdate = fn
	}
}

// WithOnRemoteDeletedFn sets a callback function to be called on create events (only local).
func WithOnRemoteDeletedFn[T any, U protoutils.ProtoMessageWithMerge[T]](
	fn OnDeleteFn[T, U],
) Option[T, U] {
	return func(s *Store[T, U], _ *jetstream.KeyValueConfig) {
		s.onRemoteDeletion = fn
	}
}

// WithKVConfig currently only used to set the TTL of the KeyValue store.
func WithKVConfig[T any, U protoutils.ProtoMessageWithMerge[T]](
	config jetstream.KeyValueConfig,
) Option[T, U] {
	return func(s *Store[T, U], kvConfig *jetstream.KeyValueConfig) {
		kvConfig.TTL = config.TTL
	}
}
