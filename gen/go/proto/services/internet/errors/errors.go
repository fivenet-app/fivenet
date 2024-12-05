package errorsinternet

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedSearch   = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.InternetService.ErrFailedSearch"}, nil)
	ErrDomainNotFound = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.InternetService.ErrDomainNotFound"}, nil)
	ErrPageNotFound   = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.InternetService.ErrPageNotFound"}, nil)
)
