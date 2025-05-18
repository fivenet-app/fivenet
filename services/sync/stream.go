package sync

import (
	"fmt"
	"strings"

	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Server) Stream(req *pbsync.StreamRequest, srv pbsync.SyncService_StreamServer) error {
	// Setup consumer
	c, err := s.js.CreateConsumer(srv.Context(), strings.ToUpper(string(BaseSubject)), jetstream.ConsumerConfig{
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	cons, err := c.Messages()
	if err != nil {
		return fmt.Errorf("failed to fetch messages consumer. %w", err)
	}
	defer cons.Stop()

	msgCh := make(chan jetstream.Msg, 8)
	go func() {
		for {
			msg, err := cons.Next()
			if err != nil {
				close(msgCh)
				return
			}

			msgCh <- msg
		}
	}()

	for {
		select {
		case <-srv.Context().Done():
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
					return fmt.Errorf("failed to unmarshal dbsync event data: %w", err)
				}

				if dest.UserId == 0 {
					continue
				}

				if err := srv.Send(dest); err != nil {
					return fmt.Errorf("failed to send stream response: %w", err)
				}
			}
		}
	}
}
