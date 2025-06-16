package errswrap

import (
	"google.golang.org/grpc/status"
)

// GRPCOrigErrorTag is the key used to tag the original error in gRPC metadata.
const GRPCOrigErrorTag = "grpc.error_orig"

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

// NewError creates a new Error that wraps an original error and a gRPC status error.
// If the original error is nil, it returns the status error as is.
func NewError(e error, s error) error {
	if e == nil {
		return s
	}

	return &Error{
		s: status.Convert(s),
		e: e,
	}
}
