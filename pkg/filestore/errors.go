package filestore

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidUploadMeta = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.Filestore.ErrInvalidUploadMeta"},
		nil,
	)
	// Has param `maxSize`.
	ErrUploadFileTooLarge = common.NewI18nErrFunc(
		codes.ResourceExhausted,
		&common.I18NItem{Key: "errors.Filestore.ErrUploadFileTooLarge"},
		nil,
	)
)
