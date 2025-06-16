package notifi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

// StreamName is the name of the JetStream stream for notifications.
const (
	StreamName = "NOTIFI"

	// BaseSubject is the root subject for all notification events.
	BaseSubject events.Subject = "notifi"

	// UserTopic is the topic for user notifications.
	UserTopic events.Topic = "user"
	// JobTopic is the topic for job event notifications.
	JobTopic events.Topic = "job"
	// JobGradeTopic is the topic for job grade event notifications.
	JobGradeTopic events.Topic = "job_grade"
	// SystemTopic is the topic for system event notifications.
	SystemTopic events.Topic = "sys"
	// MailerTopic is the topic for mailer event notifications.
	MailerTopic events.Topic = "mailer"
)

// registerEvents creates or updates the JetStream stream for notification events.
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
		return fmt.Errorf("failed to create/update stream for notification events. %w", err)
	}

	return nil
}

// SplitSubject splits a subject string into its base subject, topic, and any additional parts.
// Returns the subject, topic, and a slice of remaining parts (if any).
func SplitSubject(in string) (events.Subject, events.Topic, []string) {
	parts := strings.Split(in, ".")

	if len(parts) >= 3 {
		return events.Subject(parts[0]), events.Topic(parts[1]), parts[2:]
	} else if len(parts) == 2 {
		return events.Subject(parts[0]), events.Topic(parts[1]), nil
	}

	return events.Subject(""), events.Topic(""), nil
}
