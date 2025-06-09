package eventscentrum

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	BaseSubject events.Subject = "centrum"

	TopicGeneral           events.Topic = "general"
	TypeGeneralSettings    events.Type  = "settings"
	TypeGeneralDispatchers events.Type  = "dispatchers"

	TopicDispatch       events.Topic = "dispatch"
	TypeDispatchDeleted events.Type  = "deleted"
	TypeDispatchCreated events.Type  = "created"
	TypeDispatchUpdated events.Type  = "updated"
	TypeDispatchStatus  events.Type  = "status"

	TopicUnit       events.Topic = "unit"
	TypeUnitCreated events.Type  = "created"
	TypeUnitDeleted events.Type  = "deleted"
	TypeUnitUpdated events.Type  = "updated"
	TypeUnitStatus  events.Type  = "status"
)

func SplitSubject(subject string) (string, events.Topic, events.Type) {
	split := strings.Split(subject, ".")
	if len(split) < 3 {
		return "", "", ""
	}

	return split[1], events.Topic(split[2]), events.Type(split[3])
}

// Structure: "BASE_SUJBECT.JOB.TOPIC.TYPE.ID"
func BuildSubject(topic events.Topic, tType events.Type, job string) string {
	return fmt.Sprintf("%s.%s.%s.%s", BaseSubject, job, topic, tType)
}

func RegisterStream(ctx context.Context, js *events.JSWrapper) (jetstream.StreamConfig, error) {
	cfg := jetstream.StreamConfig{
		Name:        "CENTRUM",
		Description: "Centrum Events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      60 * time.Second,
		Storage:     jetstream.MemoryStorage,
		Duplicates:  10 * time.Second,
	}
	if _, err := js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
