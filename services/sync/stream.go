package sync

import (
	"errors"
	"fmt"
	"strings"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Server) Stream(req *pbsync.StreamRequest, srv pbsync.SyncService_StreamServer) error {
	ctx := srv.Context()

	// Update last (seen) dbsync version when set
	if req.Version != nil && req.GetVersion() != "" {
		ver := req.GetVersion()
		s.lastDBSyncVersion.Store(&ver)
	}

	// Setup consumer
	consumer, err := s.js.CreateOrUpdateConsumer(
		ctx,
		strings.ToUpper(string(BaseSubject)),
		jetstream.ConsumerConfig{
			Durable:           instance.ID() + "_sync",
			FilterSubject:     fmt.Sprintf("%s.>", BaseSubject),
			DeliverPolicy:     jetstream.DeliverNewPolicy,
			InactiveThreshold: 1 * time.Minute, // Close consumer if inactive for 1 minute
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer. %w", err)
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		msgs, err := consumer.Messages(
			jetstream.PullMaxMessages(1),
			jetstream.WithMessagesErrOnMissingHeartbeat(false),
		)
		if err != nil {
			return err
		}
		defer msgs.Stop()

		for {
			msg, err := msgs.Next(jetstream.NextContext(gctx))
			if err != nil {
				if protoutils.IsContextCanceled(err) ||
					errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					return nil
				}
				return err
			}

			// "Forward" dbsync event via this stream
			if msg == nil {
				s.logger.Warn("nil dbsync event received via message queue")
				return nil
			}
			if err := msg.Ack(); err != nil {
				s.logger.Error("failed to ack dbsync event", zap.Error(err))
				continue
			}

			_, topic := splitSubject(msg.Subject())
			switch topic {
			case TopicUser:
				dest := &pbsync.StreamResponse{}
				if err := protojson.Unmarshal(msg.Data(), dest); err != nil {
					return fmt.Errorf("failed to unmarshal dbsync event data. %w", err)
				}

				if dest.GetUserId() == 0 {
					continue
				}

				if err := srv.Send(dest); err != nil {
					if protoutils.IsContextCanceled(err) {
						return nil
					}
					return fmt.Errorf("failed to send stream response. %w", err)
				}

			default:
				s.logger.Warn(
					"received dbsync event with unknown topic",
					zap.String("topic", string(topic)),
				)
			}
		}
	})

	return g.Wait()
}
