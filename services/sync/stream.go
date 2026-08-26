package sync

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (s *Server) Stream(srv pbsync.SyncService_StreamServer) error {
	ctx := srv.Context()
	s.markDBSyncStreamConnected()
	defer s.markDBSyncStreamDisconnected()

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
		for {
			req, err := srv.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || protoutils.IsContextCanceled(err) {
					return nil
				}
				return err
			}
			if req == nil {
				continue
			}

			s.handleStreamRequest(req)
		}
	})

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
				continue
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
					s.logger.Error("failed to unmarshal dbsync event data", zap.Error(err))
					continue
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

func (s *Server) handleStreamRequest(req *pbsync.StreamRequest) {
	if req == nil {
		return
	}

	if ver := req.GetVersion(); ver != "" {
		s.lastDBSyncVersion.Store(&ver)
	}

	if req.GetSyncState() != nil {
		s.lastDBSyncState.Store(proto.Clone(req.GetSyncState()).(*pbsync.ClientSyncState))

		s.logger.Debug(
			"received dbsync sync state",
			zap.Int("tables", len(req.GetSyncState().GetTables())),
		)
	}
}
