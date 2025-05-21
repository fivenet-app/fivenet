package mailer

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
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
				tThreadsRecipients.Email,
			).
			VALUES(
				threadId,
				recipient.EmailId,
				recipient.Email.Email,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	}

	return nil
}

func (s *Server) getThreadRecipients(ctx context.Context, tx qrm.DB, threadId uint64) ([]*mailer.ThreadRecipientEmail, error) {
	tThreadsRecipients := tThreadsRecipients.AS("thread_recipient_email")
	stmt := tThreadsRecipients.
		SELECT(
			tThreadsRecipients.ID,
			tThreadsRecipients.ThreadID,
			tThreadsRecipients.EmailID,
			tEmails.ID,
			tThreadsRecipients.Email.AS("email.email"),
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
	if err := stmt.QueryContext(ctx, tx, &recipients); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return recipients, nil
}

func (s *Server) retrieveRecipientsToEmails(ctx context.Context, senderEmail *mailer.Email, recipients []string) ([]*mailer.ThreadRecipientEmail, error) {
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
			tEmails.Email,
			tEmails.Deactivated,
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
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if len(recipients) != len(dest) {
		return nil, errorsmailer.ErrInvalidRecipients
	}

	// Add email "name" to thread recipient by matching the email via the recipients list
	for _, recipient := range dest {
		if senderEmail.Id == recipient.EmailId {
			return nil, errorsmailer.ErrSameAddress
		}

		idx := slices.IndexFunc(recipients, func(in string) bool {
			return in == recipient.Email.Email
		})
		if idx == -1 {
			return nil, errorsmailer.ErrInvalidRecipients
		}

		recipient.Email = &mailer.Email{
			Id:    recipient.EmailId,
			Email: recipient.Email.Email,
		}
	}

	// The blocklist of receivers is currently checked client-side..

	return dest, nil
}
