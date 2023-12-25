package errswrap

import (
	"google.golang.org/grpc/status"
)

const GRPCOrigErrorTag = "grpc.error_orig"

type IWrappedError interface {
	GRPCStatus() *status.Status
	OrigErr() error
}

// Error wraps a pointer of a status proto and an original error. It implements error and Status,
// and a nil *Error should never be returned by this package.
type Error struct {
	s *status.Status
	e error
}

func (e *Error) OrigErr() error {
	return e.e
}

func (e *Error) Error() string {
	return e.e.Error()
}

// GRPCStatus returns the Status represented by s.
func (e *Error) GRPCStatus() *status.Status {
	return e.s
}

func NewError(s error, e error) error {
	if e == nil {
		return s
	}

	return &Error{
		s: status.Convert(s),
		e: e,
	}
}
