package errorsjobs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery       = status.Error(codes.Internal, "errors.JobsService.ErrFailedQuery")
	ErrPropsWantedDenied = status.Error(codes.PermissionDenied, "errors.JobsService.ErrPropsWantedDenied")
	ErrReasonRequired    = status.Error(codes.InvalidArgument, "errors.JobsService.ErrReasonRequired")
)
