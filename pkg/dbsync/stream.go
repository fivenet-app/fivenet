package dbsync

import (
	"context"
	"io"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Sync) RunStream(ctx context.Context) {
	defer s.wg.Done()

	for range 3 {
		s.wg.Add(1)
		go s.streamWorker(ctx)
	}

	for {
		if err := s.runStream(ctx); err != nil {
			s.logger.Error("error during sync stream, restarting in a second", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(1 * time.Second):
		}
	}
}

func (s *Sync) runStream(ctx context.Context) error {
	stream, err := s.syncCli.Stream(ctx, &pbsync.StreamRequest{
		Version: &version.Version,
	})
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			st := status.Convert(err)
			if st.Code() == codes.Unavailable {
				s.logger.Debug("stream ended with unavailable code", zap.Error(err))
				return nil
			}

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
