package mailer

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/utils"
)

func (s *Server) sendUpdate(ctx context.Context, event *mailer.MailerEvent, emailIds []uint64) error {
	emailIds = utils.RemoveSliceDuplicates(emailIds)

	for _, emailId := range emailIds {
		// TODO how to handle sending out the updates to users
		if _, err := s.js.PublishAsyncProto(ctx, fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, emailId), &notifications.UserEvent{
			Data: &notifications.UserEvent_Mailer{
				Mailer: event,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
