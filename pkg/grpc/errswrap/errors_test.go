package errswrap

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewErrorUnwrapsOriginalError(t *testing.T) {
	t.Parallel()

	orig := errors.New("base error")
	wrapped := NewError(orig, status.Error(codes.InvalidArgument, "bad request"))

	if !errors.Is(wrapped, orig) {
		t.Fatalf("expected wrapped error to unwrap to original error")
	}
}
