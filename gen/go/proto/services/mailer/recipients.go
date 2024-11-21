package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	jet "github.com/go-jet/jet/v2/mysql"
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

func (s *Server) getThreadRecipients(ctx context.Context, threadId uint64) ([]*mailer.ThreadRecipientEmail, error) {
	tThreadsRecipients := tThreadsRecipients.AS("thread_recipient_email")
	stmt := tThreadsRecipients.
		SELECT(
			tThreadsRecipients.ID,
			tThreadsRecipients.ThreadID,
			tThreadsRecipients.EmailID,
			tEmails.ID,
			tEmails.Email,
			tEmails.Internal,
		).
		FROM(
			tThreadsRecipients.
				INNER_JOIN(tEmails,
					tEmails.ID.EQ(tThreadsRecipients.EmailID).
						AND(tEmails.Deactivated.IS_FALSE()),
				),
		).
		WHERE(jet.AND(
			tThreadsRecipients.ThreadID.EQ(jet.Uint64(threadId)),
			tEmails.DeletedAt.IS_NULL(),
		))

	recipients := []*mailer.ThreadRecipientEmail{}
	if err := stmt.QueryContext(ctx, s.db, &recipients); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return recipients, nil
}

func (s *Server) resolveRecipientsToEmails(ctx context.Context, senderEmail *mailer.Email, recipients []string) ([]*mailer.ThreadRecipientEmail, error) {
	if len(recipients) == 0 {
		return nil, errorsmailer.ErrRecipientMinium
	}

	emails := make([]jet.Expression, len(recipients))
	for idx := range recipients {
		emails[idx] = jet.String(recipients[idx])
	}

	stmt := tEmails.
		SELECT(
			tEmails.ID.AS("thread_recipient_email.email_id"),
			tEmails.Deactivated,
			tEmails.Internal,
		).
		FROM(tEmails).
		WHERE(jet.AND(
			tEmails.Email.IN(emails...),
			tEmails.DeletedAt.IS_NULL(),
		)).
		LIMIT(int64(len(recipients)))

	dest := []*mailer.ThreadRecipientEmail{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if len(recipients) != len(dest) {
		return nil, errorsmailer.ErrInvalidRecipients
	}

	for _, recipient := range dest {
		// How should we handle internal email addresses of recipients?
		if senderEmail.Id == recipient.EmailId {
			return nil, errorsmailer.ErrSameAddress
		}
	}

	// TODO check blocklist of receivers

	return dest, nil
}