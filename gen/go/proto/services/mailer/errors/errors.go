package errorsmailer

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery = status.Error(codes.Internal, "errors.MailerService.ErrFailedQuery")
	ErrNoPerms     = status.Error(codes.InvalidArgument, "errors.MailerService.ErrNoPerms")

	ErrAddresseAlreadyTaken = status.Error(codes.InvalidArgument, "errors.MailerService.ErrAddresseAlreadyTaken.title;errors.MailerService.ErrAddresseAlreadyTaken.content")
	ErrAddresseInvalid      = status.Error(codes.InvalidArgument, "errors.MailerService.ErrAddresseInvalid")
	ErrTemplateLimitReached = status.Error(codes.InvalidArgument, "errors.MailerService.ErrTemplateLimitReached")
	ErrEmailAccessDenied    = status.Error(codes.InvalidArgument, "errors.MailerService.ErrEmailAccessDenied")
	ErrCantDeleteOwnEmail   = status.Error(codes.InvalidArgument, "errors.MailerService.ErrCantDeleteOwnEmail")
	ErrEmailAccessRequired  = status.Error(codes.InvalidArgument, "errors.MailerService.ErrEmailAccessRequired")

	ErrRecipientMinium   = status.Error(codes.InvalidArgument, "errors.MailerService.ErrRecipientMinium")
	ErrInvalidRecipients = status.Error(codes.InvalidArgument, "errors.MailerService.ErrInvalidRecipients")
	ErrSameAddress       = status.Error(codes.InvalidArgument, "errors.MailerService.ErrSameAddress")

	ErrThreadAccessDenied = status.Error(codes.InvalidArgument, "errors.MailerService.ErrThreadAccessDenied.title;errors.MailerService.ErrThreadAccessDenied.content")
)
