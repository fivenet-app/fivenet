package errorswiki

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery  = common.NewI18nErr(codes.Internal, common.NewI18NItem("errors.WikiService.ErrFailedQuery"), nil)
	ErrPageDenied   = common.NewI18nErr(codes.InvalidArgument, common.NewI18NItem("errors.WikiService.ErrPageDenied"), nil)
	ErrPageNotFound = common.NewI18nErr(codes.NotFound, common.NewI18NItem("errors.WikiService.ErrPageNotFound.content"), common.NewI18NItem("errors.WikiService.ErrPageNotFound.title"))

	ErrPageHasChildren = common.NewI18nErr(codes.InvalidArgument, common.NewI18NItem("errors.WikiService.ErrPageHasChildren.content"), common.NewI18NItem("errors.WikiService.ErrPageHasChildren.title"))

	ErrMaxFilesReached = common.NewI18nErr(codes.InvalidArgument, common.NewI18NItem("errors.WikiService.ErrMaxFilesReached.content"), common.NewI18NItem("errors.WikiService.ErrMaxFilesReached.title"))
)
