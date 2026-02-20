package dbsync

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Sync) RunStream(ctx context.Context) {
	for range 3 {
		s.wg.Go(func() {
			s.streamWorker(ctx)
		})
	}

	s.logger.Info("starting sync stream")
	for {
		if err := s.runStream(ctx); err != nil {
			s.logger.Warn("error during sync stream, restarting in a second", zap.Error(err))
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
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				s.logger.Error("stream ended with a non-grpc error", zap.Error(err))
				return err
			}

			switch st.Code() {
			case codes.Unavailable:
				s.logger.Debug("stream ended with unavailable code", zap.Error(err))
				return nil

			case codes.Unknown:
				if strings.Contains(
					st.Message(),
					"unexpected HTTP status code received from server: 524",
				) {
					s.logger.Debug(
						"stream ended with gateway timeout (524; Cloudflare?)",
						zap.Error(err),
					)
					return nil
				}

				s.logger.Debug("stream ended with unknown code", zap.Error(err))
				return nil
			}

			return err
		}

		s.streamCh <- msg
	}
}

func (s *Sync) streamWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case in := <-s.streamCh:
			s.logger.Info(
				"received sync stream message",
				zap.Int32("user_id", in.GetUserId()),
			)
			if err := s.users.SyncUser(ctx, in.GetUserId()); err != nil {
				s.logger.Error(
					"error during single user sync",
					zap.Int32("user_id", in.GetUserId()),
					zap.Error(err),
				)
			}
		}
	}
}
