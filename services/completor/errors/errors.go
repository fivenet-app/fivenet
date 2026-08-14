package errorscompletor

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var ErrFailedSearch = common.NewI18nErr(
	codes.Internal,
	&common.I18NItem{Key: "errors.completor.CompletorService.ErrFailedSearch.content"},
	&common.I18NItem{Key: "errors.completor.CompletorService.ErrFailedSearch.title"},
)
