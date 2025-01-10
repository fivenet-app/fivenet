package dbsync

import (
	"context"
	"io"
	"time"

	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"go.uber.org/zap"
)

func (s *Sync) RunStream(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	for i := 0; i < 3; i++ {
		s.wg.Add(1)
		go s.streamWorker(ctx)
	}

	for {
		if err := s.runStream(ctx); err != nil {
			s.logger.Error("error during sync stream", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(5 * time.Second):
		}
	}
}

func (s *Sync) runStream(ctx context.Context) error {
	stream, err := s.cli.Stream(ctx, &pbsync.StreamRequest{})
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s.streamCh <- msg
	}
}

func (s *Sync) streamWorker(ctx context.Context) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case in := <-s.streamCh:
			if err := s.users.SyncUser(ctx, in.UserId); err != nil {
				s.logger.Error("error during single user sync", zap.Int32("user_id", in.UserId), zap.Error(err))
			}
		}
	}
}
