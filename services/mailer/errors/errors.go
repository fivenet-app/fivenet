package errorsmailer

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrFailedQuery"},
		nil,
	)
	ErrNoPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrNoPerms"},
		nil,
	)

	ErrAddresseAlreadyTaken = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrAddresseAlreadyTaken.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrAddresseAlreadyTaken.title"},
	)
	ErrAddresseInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrAddresseInvalid"},
		nil,
	)
	ErrTemplateLimitReached = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrTemplateLimitReached"},
		nil,
	)
	ErrEmailAccessDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailAccessDenied"},
		nil,
	)
	ErrCantDeleteOwnEmail = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrCantDeleteOwnEmail"},
		nil,
	)
	ErrEmailAccessRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrEmailAccessRequired"},
		nil,
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
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrRecipientMinium"},
		nil,
	)
	ErrInvalidRecipients = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrInvalidRecipients"},
		nil,
	)
	ErrSameAddress = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrSameAddress"},
		nil,
	)

	ErrThreadAccessDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrThreadAccessDenied.content"},
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrThreadAccessDenied.title"},
	)

	ErrSignatureTooLong = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.mailer.MailerService.ErrSignatureTooLong.content"},
		nil,
	)
)
