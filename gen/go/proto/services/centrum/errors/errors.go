package errorscentrum

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery       = status.Error(codes.Internal, "errors.CentrumService.ErrFailedQuery")
	ErrNotPartOfDispatch = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrNotPartOfDispatch")
	ErrNotPartOfUnit     = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrNotPartOfUnit")
	ErrNotOnDuty         = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrNotOnDuty.title;errors.CentrumService.ErrNotOnDuty.content")
	ErrStaticUnit        = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrStaticUnit")

	ErrModeForbidsAction        = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrModeForbidsAction.title;errors.CentrumService.ErrModeForbidsAction.content")
	ErrDispatchAlreadyCompleted = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrDispatchAlreadyCompleted.title;errors.CentrumService.ErrDispatchAlreadyCompleted.content")
)
