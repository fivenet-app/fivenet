package mailer

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
)

func (s *Server) sendUpdate(
	ctx context.Context,
	event *mailer.MailerEvent,
	emailIds ...uint64,
) error {
	emailIds = utils.RemoveSliceDuplicates(emailIds)

	for _, emailId := range emailIds {
		if _, err := s.js.PublishAsyncProto(ctx, fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.MailerTopic, emailId), event); err != nil {
			return err
		}
	}

	return nil
}
