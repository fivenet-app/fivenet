package errorswiki

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		common.NewI18nItem("errors.WikiService.ErrFailedQuery"),
		nil,
	)
	ErrPageDenied = common.NewI18nErr(
		codes.InvalidArgument,
		common.NewI18nItem("errors.WikiService.ErrPageDenied"),
		nil,
	)
	ErrPageNotFound = common.NewI18nErr(
		codes.NotFound,
		common.NewI18nItem("errors.WikiService.ErrPageNotFound.content"),
		common.NewI18nItem("errors.WikiService.ErrPageNotFound.title"),
	)

	ErrPageHasChildren = common.NewI18nErr(
		codes.InvalidArgument,
		common.NewI18nItem("errors.WikiService.ErrPageHasChildren.content"),
		common.NewI18nItem("errors.WikiService.ErrPageHasChildren.title"),
	)

	ErrMaxFilesReached = common.NewI18nErr(
		codes.InvalidArgument,
		common.NewI18nItem("errors.WikiService.ErrMaxFilesReached.content"),
		common.NewI18nItem("errors.WikiService.ErrMaxFilesReached.title"),
	)
)
