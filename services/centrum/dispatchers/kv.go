package dispatchers

import (
	"context"
	"errors"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *DispatchersDB) Get(
	ctx context.Context,
	job string,
) (*centrumdispatchers.Dispatchers, error) {
	dispatchers, err := s.store.GetOrLoad(ctx, job)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil, err
	}

	if dispatchers == nil {
		dispatchers = &centrumdispatchers.Dispatchers{
			Job: job,
		}
	}

	s.enricher.EnrichJobName(dispatchers)

	return dispatchers, nil
}

func (s *DispatchersDB) Range(fn func(key string, value *centrumdispatchers.Dispatchers) bool) {
	s.store.Range(fn)
}
