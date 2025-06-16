package errorsinternet

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery  = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.InternetService.ErrFailedQuery"}, nil)
	ErrFailedSearch = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.InternetService.ErrFailedSearch"}, nil)

	ErrDomainNotTransferable   = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.InternetService.ErrDomainNotTransferable.content"}, &common.I18NItem{Key: "errors.InternetService.ErrDomainNotTransferable.title"})
	ErrDomainWrongTransferCode = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.InternetService.ErrDomainWrongTransferCode.content"}, &common.I18NItem{Key: "errors.InternetService.ErrDomainWrongTransferCode.title"})
)
