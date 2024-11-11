package errorswiki

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery = status.Error(codes.Internal, "errors.WikiService.ErrFailedQuery")
	ErrPageDenied  = status.Error(codes.InvalidArgument, "errors.WikiService.ErrPageDenied")
)
