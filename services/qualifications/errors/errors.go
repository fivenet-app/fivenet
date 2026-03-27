package errorsqualifications

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrFailedQuery"},
		nil,
	)
	ErrRequirementsMissing = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrRequirementsMissing.content"},
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrRequirementsMissing.title"},
	)
	ErrQualificationClosed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrQualificationClosed"},
		nil,
	)
	ErrExamDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrExamDisabled"},
		nil,
	)
	ErrRequirementSelfRef = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.qualifications.QualificationsService.ErrRequirementSelfRef"},
		nil,
	)
	ErrQualiAccessDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.QualificationService.ErrQualiAccessDenied"},
		nil,
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
