package dbsync

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
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

	cfg := s.cfg.Load()
	api := cfg.Destination.API
	s.logger.Info("starting sync stream", zap.String("host", api.URL))
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
		if errors.Is(err, io.EOF) || protoutils.IsContextCanceled(err) {
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

		select {
		case <-ctx.Done():
			return nil

		case s.streamCh <- msg:
		}
	}
}

func (s *Sync) streamWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-s.streamCh:
			if !ok {
				return
			}
			if msg == nil {
				s.logger.Warn("received nil dbsync stream response")
				continue
			}

			switch data := msg.GetPayload().(type) {
			case nil:
				s.logger.Warn("received dbsync stream response without payload (nil)")
				continue

			case *pbsync.StreamResponse_UserId:
				s.logger.Info(
					"received single user sync request",
					zap.Int32("user_id", data.UserId),
				)
				if err := s.users.SyncUser(ctx, data.UserId); err != nil {
					s.logger.Error(
						"error during single user sync",
						zap.Int32("user_id", data.UserId),
						zap.Error(err),
					)
				}

			default:
				s.logger.Warn(
					"received unknown dbsync stream response payload",
					zap.Any("payload", data),
				)
			}
		}
	}
}
