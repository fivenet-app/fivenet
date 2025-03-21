package broker

import (
	"context"
	"sync/atomic"
)

// Tweaked version of https://stackoverflow.com/a/49877632 CC-BY-SA 4.0 [icza](https://stackoverflow.com/users/1705598/icza)

type Broker[T any] struct {
	subs      atomic.Int64
	publishCh chan T
	subCh     chan chan T
	unsubCh   chan chan T
}

func New[T any]() *Broker[T] {
	return &Broker[T]{
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
	}
}

func (b *Broker[T]) Start(ctx context.Context) {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-ctx.Done():
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
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, 7)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

func (b *Broker[T]) Publish(msg T) {
	b.publishCh <- msg
}

func (b *Broker[T]) SubCount() int64 {
	return b.subs.Load()
}
