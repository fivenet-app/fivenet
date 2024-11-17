package mailer

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) handleRecipientsChanges(ctx context.Context, tx qrm.DB, threadId uint64, recipients []*mailer.ThreadRecipientEmail) error {
	if len(recipients) == 0 {
		return nil
	}

	for _, recipient := range recipients {
		stmt := tThreadsRecipients.
			INSERT(
				tThreadsRecipients.ThreadID,
				tThreadsRecipients.EmailID,
			).
			VALUES(
				threadId,
				recipient.EmailId,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	}

	return nil
}
