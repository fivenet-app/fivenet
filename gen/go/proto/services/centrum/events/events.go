package eventscentrum

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	"github.com/nats-io/nats.go"
)

const (
	BaseSubject events.Subject = "centrum"

	TopicGeneral          events.Topic = "general"
	TypeGeneralSettings   events.Type  = "settings"
	TypeGeneralDisponents events.Type  = "disponents"

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

func GetEventTypeFromSubject(subject string) (events.Topic, events.Type) {
	_, topic, eType := SplitSubject(subject)
	return topic, eType
}

func SplitSubject(subject string) (string, events.Topic, events.Type) {
	split := strings.Split(subject, ".")
	if len(split) < 3 {
		return "", "", ""
	}

	return split[1], events.Topic(split[2]), events.Type(split[3])
}

// Structure: "centrum.JOB.TOPIC.TYPE.ID"
func BuildSubject(topic events.Topic, tType events.Type, job string, id uint64) string {
	format := "%s.%s." + string(topic) + "." + string(tType)
	return fmt.Sprintf(format+".%d", BaseSubject, job, id)
}

func RegisterStreams(ctx context.Context, ev *events.Eventus) error {
	cfg := &nats.StreamConfig{
		Name:        "CENTRUM",
		Description: "FiveNet Centrum events.",
		Retention:   nats.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     nats.DiscardOld,
		MaxAge:      120 * time.Second,
		Storage:     nats.MemoryStorage,
	}

	if _, err := ev.JS.UpdateStream(cfg); err != nil {
		if !errors.Is(nats.ErrStreamNotFound, err) {
			return err
		}

		if _, err := ev.JS.AddStream(cfg); err != nil {
			return err
		}
	}

	return nil
}
