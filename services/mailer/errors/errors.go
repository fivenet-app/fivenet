package errorsmailer

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.MailerService.ErrFailedQuery"}, nil)
	ErrNoPerms     = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrNoPerms"}, nil)

	ErrAddresseAlreadyTaken = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrAddresseAlreadyTaken.content"}, &common.I18NItem{Key: "errors.MailerService.ErrAddresseAlreadyTaken.title"})
	ErrAddresseInvalid      = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrAddresseInvalid"}, nil)
	ErrTemplateLimitReached = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrTemplateLimitReached"}, nil)
	ErrEmailAccessDenied    = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrEmailAccessDenied"}, nil)
	ErrCantDeleteOwnEmail   = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrCantDeleteOwnEmail"}, nil)
	ErrEmailAccessRequired  = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrEmailAccessRequired"}, nil)
	ErrEmailChangeTooEarly  = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrEmailChangeTooEarly.content"}, &common.I18NItem{Key: "errors.MailerService.ErrEmailChangeTooEarly.title"})
	ErrEmailDisabled        = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrEmailDisabled.content"}, &common.I18NItem{Key: "errors.MailerService.ErrEmailDisabled.title"})

	ErrRecipientMinium   = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrRecipientMinium"}, nil)
	ErrInvalidRecipients = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrInvalidRecipients"}, nil)
	ErrSameAddress       = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrSameAddress"}, nil)

	ErrThreadAccessDenied = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.MailerService.ErrThreadAccessDenied.content"}, &common.I18NItem{Key: "errors.MailerService.ErrThreadAccessDenied.title"})
)
