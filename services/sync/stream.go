package sync

import (
	"fmt"
	"strings"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
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

	batch, err := consumer.Messages()
	if err != nil {
		return fmt.Errorf("failed to fetch messages consumer. %w", err)
	}
	defer batch.Stop()

	msgCh := make(chan jetstream.Msg, 8)
	go func() {
		for {
			msg, err := batch.Next()
			if err != nil {
				close(msgCh)
				return
			}

			msgCh <- msg
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case msg := <-msgCh:
			// "Forward" dbsync event via this stream
			if msg == nil {
				s.logger.Warn("nil dbsync event received via message queue")
				return nil
			}

			if err := msg.Ack(); err != nil {
				s.logger.Error("failed to ack dbsync event", zap.Error(err))
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
					return fmt.Errorf("failed to send stream response. %w", err)
				}
			}
		}
	}
}
