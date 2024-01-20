package errorsrector

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery       = status.Error(codes.Internal, "errors.RectorService.ErrFailedQuery")
	ErrInvalidRequest    = status.Error(codes.InvalidArgument, "errors.RectorService.ErrInvalidRequest")
	ErrNoPermission      = status.Error(codes.PermissionDenied, "errors.RectorService.ErrNoPermission")
	ErrRoleAlreadyExists = status.Error(codes.InvalidArgument, "errors.RectorService.ErrRoleAlreadyExists")
	ErrOwnRoleDeletion   = status.Error(codes.InvalidArgument, "errors.RectorService.ErrOwnRoleDeletion")
	ErrInvalidAttrs      = status.Error(codes.InvalidArgument, "errors.RectorService.ErrInvalidAttrs")
	ErrInvalidPerms      = status.Error(codes.InvalidArgument, "errors.RectorService.ErrInvalidPerms")
)
