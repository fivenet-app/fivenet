package errorsinternet

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery  = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.InternetService.ErrFailedQuery"}, nil)
	ErrFailedSearch = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.InternetService.ErrFailedSearch"}, nil)

	ErrDomainNotTransferable   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.InternetService.ErrDomainNotTransferable.content"}, &common.TranslateItem{Key: "errors.InternetService.ErrDomainNotTransferable.title"})
	ErrDomainWrongTransferCode = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.InternetService.ErrDomainWrongTransferCode.content"}, &common.TranslateItem{Key: "errors.InternetService.ErrDomainWrongTransferCode.title"})
)
