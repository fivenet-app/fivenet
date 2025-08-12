package dispatchers

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *DispatchersDB) updateDispatchersInKV(
	ctx context.Context,
	job string,
	dispatchers []*jobs.Colleague,
) error {
	if err := s.store.Put(ctx, job, &centrum.Dispatchers{
		Job:         job,
		Dispatchers: dispatchers,
	}); err != nil {
		return err
	}

	return nil
}

func (s *DispatchersDB) Get(ctx context.Context, job string) (*centrum.Dispatchers, error) {
	dispatchers, err := s.store.GetOrLoad(ctx, job)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil, err
	}

	if dispatchers == nil {
		dispatchers = &centrum.Dispatchers{
			Job: job,
		}
		s.enricher.EnrichJobName(dispatchers)
	}

	return dispatchers, nil
}
