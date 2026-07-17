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

func TestNewErrorReturnsPublicStatusWhenOriginalIsNil(t *testing.T) {
	t.Parallel()

	public := status.Error(codes.InvalidArgument, "bad request")
	err := NewError(nil, public)

	if err != public {
		t.Fatalf("expected nil original to return the public error unchanged")
	}

	if _, ok := err.(IWrappedError); ok {
		t.Fatalf("expected nil original to not return a wrapped error")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected returned error to be a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected status code %v, got %v", codes.InvalidArgument, st.Code())
	}
}
