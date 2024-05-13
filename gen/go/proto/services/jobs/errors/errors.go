package errorsjobs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery        = status.Error(codes.Internal, "errors.JobsService.ErrFailedQuery")
	ErrPropsAbsenceDenied = status.Error(codes.PermissionDenied, "errors.JobsService.ErrPropsAbsenceDenied")
	ErrPropsNoteDenied    = status.Error(codes.PermissionDenied, "errors.JobsService.ErrPropsNoteDenied")
	ErrReasonRequired     = status.Error(codes.InvalidArgument, "errors.JobsService.ErrReasonRequired")
)
