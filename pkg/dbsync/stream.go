package dbsync

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbsync/statusmapper"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
	s.logger.Info("starting sync stream", zap.String("host", cfg.Destination.API.URL))
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
	g, gctx := errgroup.WithContext(ctx)

	stream, err := s.syncCli.Stream(gctx)
	if err != nil {
		return err
	}

	if err := stream.Send(s.buildStreamRequest()); err != nil {
		return err
	}

	ticker := time.NewTicker(s.cfg.Load().GetSyncStateInterval())
	defer ticker.Stop()

	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil

			case <-ticker.C:
				if err := stream.Send(s.buildStreamRequest()); err != nil {
					// Keep the stream from reconnecting on transient state-send failures; the Recv loop reconnects on errors.
					s.logger.Warn("failed to send dbsync sync state", zap.Error(err))
				}
			}
		}
	})

	g.Go(func() error {
		for {
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || protoutils.IsContextCanceled(err) {
					return nil
				}

				st, ok := status.FromError(err)
				if !ok {
					s.logger.Error("stream ended with a non-grpc error", zap.Error(err))
					return err
				}

				switch st.Code() {
				case codes.Unavailable:
					return err

				case codes.Unknown:
					if strings.Contains(
						st.Message(),
						"unexpected HTTP status code received from server: 524",
					) {
						return nil
					}

					return nil
				}

				return err
			}

			select {
			case <-gctx.Done():
				return nil

			case s.streamCh <- msg:
			}
		}
	})

	return g.Wait()
}

func (s *Sync) buildStreamRequest() *pbsync.StreamRequest {
	req := &pbsync.StreamRequest{}
	req.SetVersion(version.Version)

	if state := s.buildStreamState(); state != nil {
		req.SetSyncState(state)
	}

	return req
}

func (s *Sync) buildStreamState() *pbsync.ClientSyncState {
	snapshot := statusmapper.FromRuntimeState(s.state, s.cfg.Load())
	return snapshot.ToClientSyncState()
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
