package errorscitizenstore

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery              = status.Error(codes.Internal, "errors.CitizenStoreService.ErrFailedQuery")
	ErrJobGradeNoPermission     = status.Error(codes.NotFound, "errors.CitizenStoreService.ErrJobGradeNoPermission")
	ErrReasonRequired           = status.Error(codes.InvalidArgument, "errors.CitizenStoreService.ErrReasonRequired")
	ErrPropsWantedDenied        = status.Error(codes.PermissionDenied, "errors.CitizenStoreService.ErrUserPropsWantedDenied")
	ErrPropsJobDenied           = status.Error(codes.PermissionDenied, "errors.CitizenStoreService.ErrPropsJobDenied")
	ErrPropsJobPublic           = status.Error(codes.InvalidArgument, "errors.CitizenStoreService.ErrPropsJobPublic")
	ErrPropsJobInvalid          = status.Error(codes.InvalidArgument, "errors.CitizenStoreService.ErrPropsJobInvalid")
	ErrPropsTrafficPointsDenied = status.Error(codes.PermissionDenied, "errors.CitizenStoreService.ErrPropsTrafficPointsDenied")
	ErrPropsMugShotDenied       = status.Error(codes.PermissionDenied, "errors.CitizenStoreService.ErrPropsMugShotDenied")
)
