package centrum

import (
	"context"
	"sync"

	"github.com/galexrt/fivenet/pkg/events"
)

type Bot struct {
	ctx    context.Context
	mutex  sync.RWMutex
	wg     sync.WaitGroup
	bots   map[string]context.CancelFunc
	events *events.Eventus

	state *state
}

func NewBotManager(ctx context.Context, state *state, eventus *events.Eventus) *Bot {
	return &Bot{
		ctx:    ctx,
		mutex:  sync.RWMutex{},
		wg:     sync.WaitGroup{},
		bots:   map[string]context.CancelFunc{},
		events: eventus,
		state:  state,
	}
}

func (b *Bot) Start(job string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Already a bot active
	if _, ok := b.bots[job]; ok {
		return nil
	}

	ctx, cancel := context.WithCancel(b.ctx)
	b.bots[job] = cancel
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		b.eventLoop(ctx, job)
	}()

	return nil
}

func (b *Bot) Stop(job string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	cancel, ok := b.bots[job]
	if !ok {
		return nil
	}

	cancel()

	return nil
}

func (b *Bot) eventLoop(ctx context.Context, job string) error {
	return nil
	/*
		msgCh := make(chan *nats.Msg, 256)

		sub, err := b.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.>", BaseSubject, job), msgCh, nats.DeliverNew())
		if err != nil {
			return err
		}
		defer sub.Unsubscribe()

		// Watch for events from message queue
		for {
			func() error {
				select {
				case <-b.ctx.Done():
					return nil

				case msg := <-msgCh:
					msg.Ack()

					_, topic, tType := splitSubject(msg.Subject)
					_, _ = topic, tType
				}

				return nil
			}()
		}
	*/
}
