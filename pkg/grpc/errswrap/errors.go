// Package errswrap provides utilities for wrapping errors with gRPC status and original error information.
package errswrap

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCOrigErrorTag is the key used to tag the original error in gRPC metadata.
const GRPCOrigErrorTag = "grpc.error_orig"

var ErrInternalServer = common.NewI18nErr(
	codes.Internal,
	&common.I18NItem{Key: "errors.general.internal_error.content"},
	&common.I18NItem{Key: "errors.general.internal_error.title"},
)

var errInternalServerStatus = status.Convert(ErrInternalServer)

// IWrappedError defines an interface for errors that wrap a gRPC status and an original error.
type IWrappedError interface {
	// GRPCStatus returns the gRPC status associated with the error.
	GRPCStatus() *status.Status
	// OrigErr returns the original error wrapped by this error.
	OrigErr() error
}

// Error wraps a pointer to a gRPC status and an original error.
// It implements error and Status interfaces.
// A nil *Error should never be returned by this package.
type Error struct {
	// s is the gRPC status associated with the error.
	s *status.Status
	// e is the original error wrapped by this Error.
	e error
}

// Unwrap returns the original wrapped error.
func (e *Error) Unwrap() error {
	return e.e
}

// OrigErr returns the original error wrapped by this Error.
func (e *Error) OrigErr() error {
	return e.e
}

// Error returns the error string of the original error.
func (e *Error) Error() string {
	return e.e.Error()
}

// GRPCStatus returns the gRPC status represented by s.
func (e *Error) GRPCStatus() *status.Status {
	return e.s
}

// NewError returns the public status unchanged when original is nil; otherwise it wraps
// the original cause together with the sanitized gRPC status.
func NewError(original error, public error) error {
	if original == nil {
		return public
	}

	publicStatus := status.Convert(public)
	if public == nil || publicStatus.Code() == codes.OK {
		publicStatus = errInternalServerStatus
	}

	return &Error{
		s: publicStatus,
		e: original,
	}
}
