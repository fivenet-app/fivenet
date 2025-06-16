package errorsvehicles

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var ErrFailedQuery = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.VehiclesService.ErrFailedQuery"}, nil)
