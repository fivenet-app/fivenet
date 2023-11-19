package manager

import (
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) registerSubscriptions() error {
	if _, err := s.js.Subscribe(fmt.Sprintf("%s.*.%s.>", eventscentrum.BaseSubject, eventscentrum.TopicGeneral), s.watchTopicGeneral, nats.DeliverLastPerSubject()); err != nil {
		s.logger.Error("failed to subscribe to centrum general topic", zap.Error(err))
		return err
	}

	return nil
}

func (s *Manager) watchTopicGeneral(msg *nats.Msg) {
	job, _, tType := eventscentrum.SplitSubject(msg.Subject)

	meta, _ := msg.Metadata()
	s.logger.Debug("received general message", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
		zap.String("job", job), zap.String("type", string(tType)))

	switch tType {
	case eventscentrum.TypeGeneralSettings:
		var dest centrum.Settings
		if err := proto.Unmarshal(msg.Data, &dest); err != nil {
			s.logger.Error("failed to unmarshal settings message", zap.Error(err))
			return
		}

		if err := s.State.UpdateSettings(job, &dest); err != nil {
			s.logger.Error("failed to update settings", zap.Error(err))
			return
		}
	}

	s.logger.Debug("handled general message", zap.String("job", job), zap.String("type", string(tType)))
}
