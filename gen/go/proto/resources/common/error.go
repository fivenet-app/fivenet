package common

import (
	"maps"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errFailedMarshal = []byte("Failed to marshal error")

func NewI18nErr(c codes.Code, content *I18NItem, title *I18NItem) error {
	out, err := protoutils.MarshalToPJSON(&Error{
		Title:   title,
		Content: content,
	})
	if err != nil {
		out = errFailedMarshal
	}

	return status.Error(c, string(out))
}

// I18nErrFunc is a function that takes params and returns an error.
type I18nErrFunc func(params map[string]any) error

// NewI18nErrFunc returns a function that creates an i18n error with dynamic params.
func NewI18nErrFunc(c codes.Code, content *I18NItem, title *I18NItem) I18nErrFunc {
	return func(params map[string]any) error {
		merged := make(map[string]string)
		if len(content.Parameters) > 0 {
			maps.Copy(merged, content.Parameters)
		}
		maps.Copy(merged, utils.ToStringMap(params))

		return NewI18nErr(
			c,
			NewI18NItemWithParams(content.Key, merged),
			title,
		)
	}
}
