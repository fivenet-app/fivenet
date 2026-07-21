//nolint:forbidigo // This is part of a CLI tool that uses `fmt.Println` for output.
package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	resourcescommon "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type ToolsNotificationsCmd struct {
	Send NotificationSendCmd `cmd:"" help:"Send a test notification over JetStream without storing it in the database."`
}

type NotificationSendCmd struct {
	UserID int32 `help:"Target user ID." required:""`

	Type     string `help:"Notification type."     default:"info"    enum:"error,warning,info,success"`
	Category string `help:"Notification category." default:"general" enum:"general,document,calendar"`

	TitleKey   string `help:"Translation key for the title."                 default:"notifications.system.test_notification.title"`
	ContentKey string `help:"Translation key for the content."               default:"notifications.system.test_notification.content"`
	Index      int    `help:"Index parameter used by the default title key." default:"1"`
}

func (c *NotificationSendCmd) Run(_ *kong.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nc, js, err := connectNotificationJetStream(cfg.Config)
	if err != nil {
		return err
	}
	defer nc.Close()

	if err := ensureNotificationStream(ctx, js, cfg.Config.NATS.Replicas); err != nil {
		return err
	}

	not, err := c.buildNotification()
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, c.UserID)
	data, err := protoutils.MarshalToJSON(&notificationsevents.UserEvent{
		Data: &notificationsevents.UserEvent_Notification{
			Notification: not,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification event: %w", err)
	}

	ack, err := js.Publish(ctx, subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish notification event: %w", err)
	}

	fmt.Printf(
		"Published notification to %s (stream: %s, seq: %d)\n",
		subject,
		ack.Stream,
		ack.Sequence,
	)

	return nil
}

func (c *NotificationSendCmd) buildNotification() (*resourcesnotifications.Notification, error) {
	typ, err := parseNotificationType(c.Type)
	if err != nil {
		return nil, err
	}

	category, err := parseNotificationCategory(c.Category)
	if err != nil {
		return nil, err
	}

	return &resourcesnotifications.Notification{
		UserId: c.UserID,
		Title: resourcescommon.NewI18nItemWithParams(
			c.TitleKey,
			map[string]string{"index": strconv.Itoa(c.Index)},
		),
		Type: typ,
		Content: resourcescommon.NewI18nItemWithParams(
			c.ContentKey,
			map[string]string{"type": typ.String()},
		),
		Category: category,
	}, nil
}

func parseNotificationType(v string) (resourcesnotifications.NotificationType, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "error":
		return resourcesnotifications.NotificationType_NOTIFICATION_TYPE_ERROR, nil
	case "warning":
		return resourcesnotifications.NotificationType_NOTIFICATION_TYPE_WARNING, nil
	case "info":
		return resourcesnotifications.NotificationType_NOTIFICATION_TYPE_INFO, nil
	case "success":
		return resourcesnotifications.NotificationType_NOTIFICATION_TYPE_SUCCESS, nil
	default:
		return resourcesnotifications.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			fmt.Errorf("unsupported notification type %q", v)
	}
}

func parseNotificationCategory(v string) (resourcesnotifications.NotificationCategory, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "general":
		return resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL, nil
	case "document":
		return resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT, nil
	case "calendar":
		return resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR, nil
	default:
		return resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_UNSPECIFIED,
			fmt.Errorf("unsupported notification category %q", v)
	}
}

func connectNotificationJetStream(cfg *config.Config) (*nats.Conn, jetstream.JetStream, error) {
	opts := []nats.Option{
		nats.Name(version.ProjectName),
	}

	if cfg.NATS.NKey != nil {
		nKeyOpt, err := nats.NkeyOptionFromSeed(*cfg.NATS.NKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read nats nkey: %w", err)
		}
		opts = append(opts, nKeyOpt)
	} else if cfg.NATS.Creds != nil {
		opts = append(opts, nats.UserCredentials(*cfg.NATS.Creds))
	}

	nc, err := nats.Connect(cfg.NATS.URL, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := jetstream.New(
		nc,
		jetstream.WithPublishAsyncMaxPending(cfg.NATS.PublishAsyncMaxPending),
	)
	if err != nil {
		nc.Close()
		return nil, nil, fmt.Errorf("failed to initialize jetstream: %w", err)
	}

	return nc, js, nil
}

func ensureNotificationStream(
	ctx context.Context,
	js jetstream.JetStream,
	replicas int,
) error {
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:              notifi.StreamName,
		Description:       "User, Job, Object and System notification events",
		Subjects:          []string{fmt.Sprintf("%s.>", notifi.BaseSubject)},
		Retention:         jetstream.InterestPolicy,
		Discard:           jetstream.DiscardOld,
		MaxAge:            15 * time.Minute,
		MaxMsgsPerSubject: 3,
		Duplicates:        time.Minute,
		Replicas:          replicas,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update notification stream: %w", err)
	}

	return nil
}
