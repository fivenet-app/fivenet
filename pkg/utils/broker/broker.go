package broker

import (
	"context"
	"sync/atomic"
)

// Broker provides a simple publish/subscribe message broker for generic types.
type Broker[T any] struct {
	// number of active subscribers
	subs atomic.Int64
	// channel for publishing messages
	publishCh chan T
	// channel for new subscriptions
	subCh chan chan T
	// channel for unsubscriptions
	unsubCh chan chan T
}

// New creates a new Broker instance.
func New[T any]() *Broker[T] {
	return &Broker[T]{
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
	}
}

// Start runs the broker event loop, handling subscriptions, unsubscriptions, and publishing.
func (b *Broker[T]) Start(ctx context.Context) {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-ctx.Done():
			// Close all subscriber channels on shutdown
			for msgCh := range subs {
				close(msgCh)
			}
			return

		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
			b.subs.Add(1)

		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
			close(msgCh)
			b.subs.Add(-1)

		case msg := <-b.publishCh:
			for msgCh := range subs {
				// Non-blocking send to avoid blocking the broker if a subscriber is slow
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

// Subscribe registers a new subscriber and returns its message channel.
func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, 7)
	b.subCh <- msgCh
	return msgCh
}

// Unsubscribe removes a subscriber and closes its channel.
func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

// Publish sends a message to all subscribers.
func (b *Broker[T]) Publish(msg T) {
	b.publishCh <- msg
}

// SubCount returns the current number of subscribers.
func (b *Broker[T]) SubCount() int64 {
	return b.subs.Load()
}
