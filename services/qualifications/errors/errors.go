package errorsqualifications

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrFailedQuery.title"},
	)
	ErrRequirementsMissing = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrRequirementsMissing.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrRequirementsMissing.title",
		},
	)
	ErrQualificationClosed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationClosed.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationClosed.title",
		},
	)
	ErrExamDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamDisabled.content",
		},
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrExamDisabled.title"},
	)
	ErrRequirementSelfRef = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrRequirementSelfRef.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrRequirementSelfRef.title",
		},
	)
	ErrQualiAccessDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualiAccessDenied.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualiAccessDenied.title",
		},
	)
	ErrQualiUpdateDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrQualiUpdateDenied"},
		nil,
	)
	ErrQualiViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrQualiViewDenied"},
		nil,
	)
)
