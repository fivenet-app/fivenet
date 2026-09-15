package centrum

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"go.uber.org/zap"
)

func (s *Server) startFeedHub(ctx context.Context) {
	s.wg.Go(func() { s.runCentrumEventFeed(ctx) })
	for _, feed := range feeds {
		s.wg.Go(func() { s.runKVFeed(ctx, feed) })
	}
}

func (s *Server) runCentrumEventFeed(ctx context.Context) {
	for {
		err := s.consumeCentrumEvents(ctx)
		if ctx.Err() == nil && !protoutils.IsContextCanceled(err) {
			if err != nil {
				s.logger.Error("centrum event feed stopped", zap.Error(err))
			}
			s.publishFeedResync("events")
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Server) runKVFeed(ctx context.Context, feed feedCfg) {
	for {
		err := s.consumeKVFeed(ctx, feed)
		if ctx.Err() == nil && !protoutils.IsContextCanceled(err) {
			if err != nil {
				s.logger.Error(
					"centrum KV feed stopped",
					zap.Error(err),
					zap.String("bucket", feed.Bucket),
				)
			}
			s.publishFeedResync(feed.StreamName)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}
