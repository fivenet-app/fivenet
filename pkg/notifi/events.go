package notifi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	StreamName = "NOTIFI"

	BaseSubject events.Subject = "notifi"

	// User notifications topic
	UserTopic events.Topic = "user"
	// Job event topic
	JobTopic events.Topic = "job"
	// Job Grade event topic
	JobGradeTopic events.Topic = "job_grade"
	// System event topic
	SystemTopic events.Topic = "sys"
	// Mailer event topic
	MailerTopic events.Topic = "mailer"
)

func (n *Notifi) registerEvents(ctx context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        StreamName,
		Description: "User and System Notification events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      30 * time.Minute,
	}
	if _, err := n.js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return err
	}

	return nil
}

func SplitSubject(in string) (events.Subject, events.Topic, []string) {
	parts := strings.Split(in, ".")

	if len(parts) >= 3 {
		return events.Subject(parts[0]), events.Topic(parts[1]), parts[2:]
	} else if len(parts) == 2 {
		return events.Subject(parts[0]), events.Topic(parts[1]), nil
	}

	return events.Subject(""), events.Topic(""), nil
}
