package errorsinternet

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery  = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.InternetService.ErrFailedQuery"}, nil)
	ErrFailedSearch = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.InternetService.ErrFailedSearch"}, nil)
)
