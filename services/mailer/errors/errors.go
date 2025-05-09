package errorsmailer

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.MailerService.ErrFailedQuery"}, nil)
	ErrNoPerms     = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrNoPerms"}, nil)

	ErrAddresseAlreadyTaken = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrAddresseAlreadyTaken.content"}, &common.TranslateItem{Key: "errors.MailerService.ErrAddresseAlreadyTaken.title"})
	ErrAddresseInvalid      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrAddresseInvalid"}, nil)
	ErrTemplateLimitReached = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrTemplateLimitReached"}, nil)
	ErrEmailAccessDenied    = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrEmailAccessDenied"}, nil)
	ErrCantDeleteOwnEmail   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrCantDeleteOwnEmail"}, nil)
	ErrEmailAccessRequired  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrEmailAccessRequired"}, nil)
	ErrEmailChangeTooEarly  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrEmailChangeTooEarly.content"}, &common.TranslateItem{Key: "errors.MailerService.ErrEmailChangeTooEarly.title"})
	ErrEmailDisabled        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrEmailDisabled.content"}, &common.TranslateItem{Key: "errors.MailerService.ErrEmailDisabled.title"})

	ErrRecipientMinium   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrRecipientMinium"}, nil)
	ErrInvalidRecipients = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrInvalidRecipients"}, nil)
	ErrSameAddress       = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrSameAddress"}, nil)

	ErrThreadAccessDenied = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.MailerService.ErrThreadAccessDenied.content"}, &common.TranslateItem{Key: "errors.MailerService.ErrThreadAccessDenied.title"})
)
