package errorswiki

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery  = common.I18nErr(codes.Internal, common.NewTranslateItem("errors.WikiService.ErrFailedQuery"), nil)
	ErrPageDenied   = common.I18nErr(codes.InvalidArgument, common.NewTranslateItem("errors.WikiService.ErrPageDenied"), nil)
	ErrPageNotFound = common.I18nErr(codes.NotFound, common.NewTranslateItem("errors.WikiService.ErrPageNotFound.content"), common.NewTranslateItem("errors.WikiService.ErrPageNotFound.title"))

	ErrPageHasChildren = common.I18nErr(codes.InvalidArgument, common.NewTranslateItem("errors.WikiService.ErrPageHasChildren.content"), common.NewTranslateItem("errors.WikiService.ErrPageHasChildren.title"))

	ErrMaxFilesReached = common.I18nErr(codes.InvalidArgument, common.NewTranslateItem("errors.WikiService.ErrMaxFilesReached.content"), common.NewTranslateItem("errors.WikiService.ErrMaxFilesReached.title"))
)
