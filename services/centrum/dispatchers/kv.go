package dispatchers

import (
	"context"
	"errors"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *DispatchersDB) updateDispatchersInKV(
	ctx context.Context,
	job string,
	dispatchers []*jobscolleagues.Colleague,
) error {
	dspers := &centrumdispatchers.Dispatchers{
		Job:         job,
		Dispatchers: dispatchers,
	}
	s.enricher.EnrichJobName(dspers)
	if err := s.store.Put(ctx, job, dspers); err != nil {
		return err
	}

	return nil
}

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
