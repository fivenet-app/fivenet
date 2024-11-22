package errorscompletor

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var ErrFailedSearch = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.CompletorService.ErrFailedSearch"}, nil)
