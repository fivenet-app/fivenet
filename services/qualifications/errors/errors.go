package errorsqualifications

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery         = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.QualificationsService.ErrFailedQuery"}, nil)
	ErrRequirementsMissing = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.QualificationsService.ErrRequirementsMissing.content"}, &common.TranslateItem{Key: "errors.QualificationsService.ErrRequirementsMissing.title"})
	ErrQualificationClosed = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.QualificationsService.ErrQualificationClosed"}, nil)
	ErrExamDisabled        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.QualificationsService.ErrExamDisabled"}, nil)
	ErrRequirementSelfRef  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.QualificationsService.ErrRequirementSelfRef"}, nil)
	ErrQualiAccessDenied   = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.QualificationService.ErrQualiAccessDenied"}, nil)
)
