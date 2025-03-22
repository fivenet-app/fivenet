package centrumbrokers

import (
	"context"
	"sync"

	pbcentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"go.uber.org/fx"
)

type Brokers struct {
	ctx context.Context

	brokersMutex *sync.RWMutex
	brokers      map[string]*broker.Broker[*pbcentrum.StreamResponse]
	brokerCancel map[string]context.CancelFunc
}

type Params struct {
	fx.In

	LC fx.Lifecycle
}

func New(p Params) *Brokers {
	ctxCancel, cancel := context.WithCancel(context.Background())

	b := &Brokers{
		ctx:          ctxCancel,
		brokersMutex: &sync.RWMutex{},
		brokers:      map[string]*broker.Broker[*pbcentrum.StreamResponse]{},
		brokerCancel: map[string]context.CancelFunc{},
	}

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		return nil
	}))

	return b
}

func (s *Brokers) GetJobBroker(job string) (*broker.Broker[*pbcentrum.StreamResponse], bool) {
	s.brokersMutex.RLock()
	defer s.brokersMutex.RUnlock()

	broker, ok := s.brokers[job]
	return broker, ok
}

func (s *Brokers) CreateJobBroker(job string) *broker.Broker[*pbcentrum.StreamResponse] {
	s.brokersMutex.Lock()
	defer s.brokersMutex.Unlock()

	ctx, cancel := context.WithCancel(s.ctx)
	s.brokerCancel[job] = cancel

	broker := broker.New[*pbcentrum.StreamResponse]()
	go broker.Start(ctx)
	s.brokers[job] = broker

	return broker
}

func (s *Brokers) GetOrCreateJobBroker(job string) *broker.Broker[*pbcentrum.StreamResponse] {
	if broker, ok := s.GetJobBroker(job); ok {
		return broker
	}

	broker := s.CreateJobBroker(job)
	return broker
}

func (s *Brokers) RemoveJobBroker(job string) {
	_, ok := s.GetJobBroker(job)
	if !ok {
		return
	}

	s.brokersMutex.Lock()
	defer s.brokersMutex.Unlock()

	s.brokerCancel[job]()
	delete(s.brokers, job)
}
