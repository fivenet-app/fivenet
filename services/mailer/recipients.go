package mailer

import (
	"context"
	"slices"

	maileremails "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/emails"
	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
)

func (s *Server) retrieveRecipientsToEmails(
	ctx context.Context,
	senderEmail *maileremails.Email,
	recipients []string,
) ([]*mailerthreads.ThreadRecipientEmail, error) {
	if len(recipients) == 0 {
		return nil, errorsmailer.ErrRecipientMinium
	}

	dest, err := s.store.ListRecipientsByEmails(ctx, s.db, recipients)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if len(recipients) != len(dest) {
		return nil, errorsmailer.ErrInvalidRecipients
	}

	// Add email "name" to thread recipient by matching the email via the recipients list
	for _, recipient := range dest {
		if senderEmail.GetId() == recipient.GetEmailId() {
			return nil, errorsmailer.ErrSameAddress
		}

		idx := slices.IndexFunc(recipients, func(in string) bool {
			return in == recipient.GetEmail().GetEmail()
		})
		if idx == -1 {
			return nil, errorsmailer.ErrInvalidRecipients
		}

		recipient.Email = &maileremails.Email{
			Id:    recipient.GetEmailId(),
			Email: recipient.GetEmail().GetEmail(),
		}
	}

	// The blocklist of receivers is currently checked client-side..

	return dest, nil
}
