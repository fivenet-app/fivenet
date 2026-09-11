package broker

import (
	"context"
	"sync/atomic"
)

// Broker provides a simple publish/subscribe message broker for generic types.
type Broker[T any] struct {
	// number of active subscribers
	subs atomic.Int64
	// queue size assigned to each subscriber
	subscriberBuffer int
	// close slow subscribers instead of silently dropping a message
	closeSlowSubscribers bool
	// channel for publishing messages
	publishCh chan T
	// channel for new subscriptions
	subCh chan subscription[T]
	// channel for unsubscriptions
	unsubCh chan chan T
}

type subscription[T any] struct {
	channel chan T
	ready   chan struct{}
}

// New creates a new Broker instance.
func New[T any]() *Broker[T] {
	return NewWithBuffer[T](7)
}

// NewWithBuffer creates a broker with the given per-subscriber queue size.
// Slow subscribers still drop messages rather than blocking every publisher.
func NewWithBuffer[T any](subscriberBuffer int) *Broker[T] {
	return newBroker[T](subscriberBuffer, false)
}

// NewWithResyncOnSlowSubscriber creates a broker that closes a subscriber's
// channel when its queue overflows. Consumers can use the close as a resync
// signal instead of continuing with incomplete state.
func NewWithResyncOnSlowSubscriber[T any](subscriberBuffer int) *Broker[T] {
	return newBroker[T](subscriberBuffer, true)
}

func newBroker[T any](subscriberBuffer int, closeSlowSubscribers bool) *Broker[T] {
	if subscriberBuffer < 1 {
		subscriberBuffer = 1
	}

	return &Broker[T]{
		subscriberBuffer:     subscriberBuffer,
		closeSlowSubscribers: closeSlowSubscribers,
		publishCh:            make(chan T, 1),
		subCh:                make(chan subscription[T], 1),
		unsubCh:              make(chan chan T, 1),
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

		case sub := <-b.subCh:
			subs[sub.channel] = struct{}{}
			b.subs.Add(1)
			close(sub.ready)

		case msgCh := <-b.unsubCh:
			if _, ok := subs[msgCh]; ok {
				delete(subs, msgCh)
				close(msgCh)
				b.subs.Add(-1)
			}

		case msg := <-b.publishCh:
			for msgCh := range subs {
				// Non-blocking send to avoid blocking the broker if a subscriber is slow
				select {
				case msgCh <- msg:
				default:
					if b.closeSlowSubscribers {
						delete(subs, msgCh)
						close(msgCh)
						b.subs.Add(-1)
					}
				}
			}
		}
	}
}

// Subscribe registers a new subscriber and returns its message channel.
func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, b.subscriberBuffer)
	ready := make(chan struct{})
	b.subCh <- subscription[T]{channel: msgCh, ready: ready}
	<-ready
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
