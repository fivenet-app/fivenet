package sync

import (
	"fmt"
	"strings"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Server) Stream(req *StreamRequest, srv grpc.ServerStreamingServer[StreamResponse]) error {
	// Setup consumer
	c, err := s.js.CreateConsumer(srv.Context(), strings.ToUpper(string(BaseSubject)), jetstream.ConsumerConfig{
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		return err
	}

	cons, err := c.Messages()
	if err != nil {
		return err
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
				dest := &StreamResponse{}
				if err := protojson.Unmarshal(msg.Data(), dest); err != nil {
					return err
				}

				if dest.UserId == 0 {
					continue
				}

				if err := srv.Send(dest); err != nil {
					return err
				}
			}
		}
	}
}
