package errorsdmv

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrFailedQuery = status.Error(codes.Internal, "errors.DMVService.ErrFailedQuery")
