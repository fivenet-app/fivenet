package errorsqualifications

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery         = status.Error(codes.Internal, "errors.QualificationsService.ErrFailedQuery")
	ErrRequirementsMissing = status.Error(codes.InvalidArgument, "errors.QualificationsService.ErrRequirementsMissing.title;errors.QualificationsService.ErrRequirementsMissing.content")
	ErrQualificationClosed = status.Error(codes.InvalidArgument, "errors.QualificationsService.ErrQualificationClosed")

	ErrExamDisabled = status.Error(codes.InvalidArgument, "errors.QualificationsService.ErrExamDisabled")
)
