package errorswiki

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.I18nErr(codes.Internal, common.NewTranslateItem("errors.WikiService.ErrFailedQuery"), nil)
	ErrPageDenied  = common.I18nErr(codes.InvalidArgument, common.NewTranslateItem("errors.WikiService.ErrPageDenied"), nil)
)
