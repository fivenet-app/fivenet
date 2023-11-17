package utils

import (
	"context"
	"sync/atomic"
)

// Tweaked version of https://stackoverflow.com/a/49877632 CC-BY-SA 4.0 [icza](https://stackoverflow.com/users/1705598/icza)

type Broker[T any] struct {
	ctx context.Context

	subs      atomic.Int64
	stopCh    chan struct{}
	publishCh chan T
	subCh     chan chan T
	unsubCh   chan chan T
}

func NewBroker[T any](ctx context.Context) *Broker[T] {
	return &Broker[T]{
		ctx:       ctx,
		stopCh:    make(chan struct{}),
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
	}
}

func (b *Broker[T]) Start() {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				close(msgCh)
			}
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
			b.subs.Add(1)
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
			b.subs.Add(-1)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		case <-b.ctx.Done():
			b.Stop()
			return
		}
	}
}

func (b *Broker[T]) Stop() {
	close(b.stopCh)
}

func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
	close(msgCh)
}

func (b *Broker[T]) Publish(msg T) {
	b.publishCh <- msg
}

func (b *Broker[T]) SubCount() int64 {
	return b.subs.Load()
}
