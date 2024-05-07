package errorsauth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAccountCreateFailed = status.Error(codes.InvalidArgument, "errors.AuthService.ErrAccountCreateFailed")
	ErrAccountExistsFailed = status.Error(codes.InvalidArgument, "errors.AuthService.ErrAccountExistsFailed")
	ErrInvalidLogin        = status.Error(codes.InvalidArgument, "errors.AuthService.ErrInvalidLogin")
	ErrNoAccount           = status.Error(codes.InvalidArgument, "errors.AuthService.ErrNoAccount")
	ErrNoCharFound         = status.Error(codes.NotFound, "errors.AuthService.ErrNoCharFound")
	ErrGenericLogin        = status.Error(codes.Internal, "errors.AuthService.ErrGenericLogin")
	ErrUnableToChooseChar  = status.Error(codes.PermissionDenied, "errors.AuthService.ErrUnableToChooseChar")
	ErrUpdateAccount       = status.Error(codes.InvalidArgument, "errors.AuthService.ErrUpdateAccount")
	ErrChangePassword      = status.Error(codes.InvalidArgument, "errors.AuthService.ErrChangePassword")
	ErrForgotPassword      = status.Error(codes.InvalidArgument, "errors.AuthService.ErrForgotPassword")
	ErrSignupDisabled      = status.Error(codes.InvalidArgument, "errors.AuthService.ErrSignupDisabled")
	ErrAccountDuplicate    = status.Error(codes.InvalidArgument, "errors.AuthService.ErrAccountDuplicate")
	ErrChangeUsername      = status.Error(codes.InvalidArgument, "errors.AuthService.ErrChangeUsername")
	ErrBadUsername         = status.Error(codes.InvalidArgument, "errors.AuthService.ErrBadUsername")
	ErrCharLock            = status.Error(codes.InvalidArgument, "errors.AuthService.ErrCharLock.title;errors.AuthService.ErrCharLock.content")
)
