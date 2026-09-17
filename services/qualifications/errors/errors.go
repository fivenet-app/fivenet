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
	ErrQualificationRequestActive = common.NewI18nErr(
		codes.FailedPrecondition,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationRequestActive.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationRequestActive.title",
		},
	)
	ErrQualificationRequestInvalidTransition = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationRequestInvalidTransition.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationRequestInvalidTransition.title",
		},
	)
	ErrQualificationAlreadySuccessful = common.NewI18nErr(
		codes.FailedPrecondition,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationAlreadySuccessful.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrQualificationAlreadySuccessful.title",
		},
	)
	ErrExamDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamDisabled.content",
		},
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrExamDisabled.title"},
	)
	ErrExamAutoGradingFreeText = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamAutoGradingFreeText.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamAutoGradingFreeText.title",
		},
	)
	ErrExamAutoGradingInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamAutoGradingInvalid.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamAutoGradingInvalid.title",
		},
	)
	ErrExamAutoGradingAnswerLimit = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamAutoGradingAnswerLimit.content",
		},
		&common.I18NItem{
			Key: "errors.qualifications.QualificationsService.ErrExamAutoGradingAnswerLimit.title",
		},
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
