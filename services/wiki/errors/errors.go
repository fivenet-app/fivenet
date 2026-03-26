package errorswiki

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		common.NewI18nItem("errors.wiki.WikiService.ErrFailedQuery"),
		nil,
	)
	ErrPageDenied = common.NewI18nErr(
		codes.InvalidArgument,
		common.NewI18nItem("errors.wiki.WikiService.ErrPageDenied"),
		nil,
	)
	ErrPageNotFound = common.NewI18nErr(
		codes.NotFound,
		common.NewI18nItem("errors.wiki.WikiService.ErrPageNotFound.content"),
		common.NewI18nItem("errors.wiki.WikiService.ErrPageNotFound.title"),
	)

	ErrPageHasChildren = common.NewI18nErr(
		codes.InvalidArgument,
		common.NewI18nItem("errors.wiki.WikiService.ErrPageHasChildren.content"),
		common.NewI18nItem("errors.wiki.WikiService.ErrPageHasChildren.title"),
	)

	ErrMaxFilesReached = common.NewI18nErr(
		codes.InvalidArgument,
		common.NewI18nItem("errors.wiki.WikiService.ErrMaxFilesReached.content"),
		common.NewI18nItem("errors.wiki.WikiService.ErrMaxFilesReached.title"),
	)
)
