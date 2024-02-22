package livemapper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	natsutils "github.com/galexrt/fivenet/pkg/nats"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	BaseSubject events.Subject = "livemap"

	MarkerUpdate events.Type = "marker_update"
)

func (s *Server) registerEvents(ctx context.Context) error {
	cfg := &nats.StreamConfig{
		Name:        "LIVEMAP",
		Description: natsutils.Description,
		Retention:   nats.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     nats.DiscardOld,
		MaxAge:      2 * time.Minute,
	}
	if _, err := natsutils.CreateOrUpdateStream(ctx, s.js, cfg); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendUpdateEvent(tType events.Type, event proto.Message) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(fmt.Sprintf("%s.%s", BaseSubject, tType), data); err != nil {
		return err
	}

	return nil
}

func (s *Server) watchForEvents(msg *nats.Msg) {
	split := strings.Split(msg.Subject, ".")
	if len(split) < 2 {
		return
	}

	tType := events.Type(split[1])
	if tType == MarkerUpdate {
		if err := s.refreshData(); err != nil {
			s.logger.Error("failed to refresh livemap markers cache", zap.Error(err))
			return
		}

		// Send marker update when data has been refreshed and we have at least one subscriber
		if s.broker.SubCount() <= 0 {
			return
		}
		s.broker.Publish(&brokerEvent{Send: MarkerUpdate})
	}
}
