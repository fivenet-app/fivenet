package filestore

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidUploadMeta = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.Filestore.ErrInvalidUploadMeta"}, nil)
	// Has param `maxSize`
	ErrUploadFileTooLarge = common.I18nErr(codes.ResourceExhausted, &common.TranslateItem{Key: "errors.Filestore.ErrUploadFileTooLarge"}, nil)
)
