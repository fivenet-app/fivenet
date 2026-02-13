package mailer

import (
	"context"
	"fmt"

	mailerevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
)

func (s *Server) sendUpdate(
	ctx context.Context,
	event *mailerevents.MailerEvent,
	emailIds ...int64,
) error {
	emailIds = utils.RemoveSliceDuplicates(emailIds)

	for _, emailId := range emailIds {
		if _, err := s.js.PublishAsyncProto(ctx, fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.MailerTopic, emailId), event); err != nil {
			return err
		}
	}

	return nil
}
