package common

import (
	"encoding/json"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewI18nErrFuncMergesParamsIntoTitleAndContent(t *testing.T) {
	t.Parallel()

	fn := NewI18nErrFunc(
		codes.InvalidArgument,
		&I18NItem{
			Key: "content.key",
			Parameters: map[string]string{
				"static": "content-static",
			},
		},
		&I18NItem{
			Key: "title.key",
			Parameters: map[string]string{
				"static": "title-static",
			},
		},
	)

	err := fn(map[string]any{
		"dynamic": 42,
	})

	got := &Error{}
	if unmarshalErr := json.Unmarshal(
		[]byte(status.Convert(err).Message()),
		got,
	); unmarshalErr != nil {
		t.Fatalf("unmarshal error payload: %v", unmarshalErr)
	}

	if got.GetContent().GetParameters()["static"] != "content-static" {
		t.Fatalf("content static param lost: %+v", got.GetContent().GetParameters())
	}
	if got.GetContent().GetParameters()["dynamic"] != "42" {
		t.Fatalf("content dynamic param missing: %+v", got.GetContent().GetParameters())
	}
	if got.GetTitle().GetParameters()["static"] != "title-static" {
		t.Fatalf("title static param lost: %+v", got.GetTitle().GetParameters())
	}
	if got.GetTitle().GetParameters()["dynamic"] != "42" {
		t.Fatalf("title dynamic param missing: %+v", got.GetTitle().GetParameters())
	}
}

func TestNewI18nErrFuncPanicsOnNilContent(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected NewI18nErrFunc to panic when content is nil")
		}
	}()

	_ = NewI18nErrFunc(codes.InvalidArgument, nil, nil)
}
