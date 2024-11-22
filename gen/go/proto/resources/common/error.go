package common

import (
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ErrorKeyFormat = "errors.%s.%s"

var errFailedMarshal = []byte("Failed to marshal error")

func I18nErr(c codes.Code, content *TranslateItem, title *TranslateItem) error {
	out, err := protoutils.Marshal(&Error{
		Title:   title,
		Content: content,
	})
	if err != nil {
		out = errFailedMarshal
	}

	return status.Error(c, string(out))
}
