package access

import (
	"context"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
)

func resetGroupedAccesses(t *testing.T) {
	t.Helper()

	groupedAccessesMu.Lock()
	original := groupedAccesses
	groupedAccesses = map[string]GroupedAccess{}
	groupedAccessesMu.Unlock()

	t.Cleanup(func() {
		groupedAccessesMu.Lock()
		groupedAccesses = original
		groupedAccessesMu.Unlock()
	})
}

func TestGetAccessMissing(t *testing.T) {
	t.Parallel()

	resetGroupedAccesses(t)

	if _, ok := GetAccess("missing"); ok {
		t.Fatal("expected missing access to return ok=false")
	}
}

func TestRegisterAndGetAccess(t *testing.T) {
	t.Parallel()

	resetGroupedAccesses(t)

	RegisterAccess(
		"demo",
		&GroupedAccessAdapter{
			CanUserAccessTargetFn: func(context.Context, int64, *userinfo.UserInfo, int32) (bool, error) {
				return true, nil
			},
		},
	)

	got, ok := GetAccess("demo")
	if !ok {
		t.Fatal("expected registered access to be found")
	}
	if got == nil {
		t.Fatal("expected registered access to be non-nil")
	}
}
