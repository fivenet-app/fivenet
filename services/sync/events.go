package sync

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	BaseSubject events.Subject = "dbsync"

	TopicUser events.Topic = "user"
)

func splitSubject(subject string) (string, events.Topic) {
	split := strings.Split(subject, ".")
	if len(split) < 2 {
		return "", ""
	}

	return split[1], events.Topic(split[2])
}

// Structure: "BASE_SUJBECT.TOPIC"
func BuildSubject(topic events.Topic) string {
	return fmt.Sprintf("%s.%s", BaseSubject, topic)
}

func (s *Server) registerStream(ctx context.Context, js *events.JSWrapper) (jetstream.StreamConfig, error) {
	cfg := jetstream.StreamConfig{
		Name:        strings.ToUpper(string(BaseSubject)),
		Description: "DBSync Events",
		Retention:   jetstream.LimitsPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		Storage:     jetstream.MemoryStorage,
		MaxAge:      300 * time.Second, // 5 minutes
		Duplicates:  15 * time.Second,
	}
	if _, err := js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
