package errorsdmv

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var ErrFailedQuery = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.DMVService.ErrFailedQuery"}, nil)
