package errorsmailer

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrFailedQuery.title"},
	)
	ErrNoPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrNoPerms.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrNoPerms.title"},
	)

	ErrAddresseAlreadyTaken = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrAddresseAlreadyTaken.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrAddresseAlreadyTaken.title"},
	)
	ErrAddresseInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrAddresseInvalid.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrAddresseInvalid.title"},
	)
	ErrTemplateLimitReached = common.NewI18nErrFunc(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrTemplateLimitReached.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrTemplateLimitReached.title"},
	)
	ErrEmailAccessDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailAccessDenied.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailAccessDenied.title"},
	)
	ErrCantDeleteOwnEmail = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrCantDeleteOwnEmail.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrCantDeleteOwnEmail.title"},
	)
	ErrEmailAccessRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailAccessRequired.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailAccessRequired.title"},
	)
	ErrEmailChangeTooEarly = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailChangeTooEarly.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailChangeTooEarly.title"},
	)
	ErrEmailDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailDisabled.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailDisabled.title"},
	)

	ErrRecipientMinium = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrRecipientMinium.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrRecipientMinium.title"},
	)
	ErrInvalidRecipients = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrInvalidRecipients.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrInvalidRecipients.title"},
	)
	ErrSameAddress = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrSameAddress.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrSameAddress.title"},
	)

	ErrThreadAccessDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrThreadAccessDenied.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrThreadAccessDenied.title"},
	)

	ErrSignatureTooLong = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrSignatureTooLong.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrSignatureTooLong.title"},
	)
)
