package errorsmailer

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery          = status.Error(codes.Internal, "errors.MailerService.ErrFailedQuery")
	ErrNoPerms              = status.Error(codes.InvalidArgument, "errors.MailerService.ErrNoPerms")
	ErrTemplateLimitReached = status.Error(codes.InvalidArgument, "errors.MailerService.ErrTemplateLimitReached")
)
