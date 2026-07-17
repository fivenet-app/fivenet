package common

import (
	"maps"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrFailedMarshal = []byte("Failed to marshal error")

func NewI18nErr(c codes.Code, content *I18NItem, title *I18NItem) error {
	out, err := protoutils.MarshalToJSON(&Error{
		Title:   title,
		Content: content,
	})
	if err != nil {
		out = ErrFailedMarshal
	}

	return status.Error(c, string(out))
}

// I18nErrFunc is a function that takes params and returns an error.
type I18nErrFunc func(params map[string]any) error

func mergeI18nParams(item *I18NItem, params map[string]any) *I18NItem {
	if item == nil {
		return nil
	}

	merged := maps.Clone(item.GetParameters())
	if merged == nil {
		merged = map[string]string{}
	}

	maps.Copy(merged, utils.ToStringMap(params))

	return NewI18nItemWithParams(item.GetKey(), merged)
}

// NewI18nErrFunc returns a function that creates an i18n error with dynamic params.
// content must be non-nil; title remains optional.
func NewI18nErrFunc(code codes.Code, content *I18NItem, title *I18NItem) I18nErrFunc {
	if content == nil {
		panic("common.NewI18nErrFunc requires non-nil content")
	}

	return func(params map[string]any) error {
		return NewI18nErr(
			code,
			mergeI18nParams(content, params),
			mergeI18nParams(title, params),
		)
	}
}
