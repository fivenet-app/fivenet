package errswrap

import (
	"google.golang.org/grpc/status"
)

const GRPCOrigErrorTag = "grpc.error_orig"

type IWrappedError interface {
	GRPCStatus() *status.Status
	GetOrigErr() error
}

// Error wraps a pointer of a status proto and an original error. It implements error and Status,
// and a nil *Error should never be returned by this package.
type Error struct {
	s  *status.Status
	se error

	e error
}

func (e *Error) OrigError() error {
	return e.e
}

func (e *Error) Error() string {
	return e.s.String()
}

// GRPCStatus returns the Status represented by s.
func (e *Error) GRPCStatus() *status.Status {
	return e.s
}

func NewError(s error, e error) error {
	return &Error{
		s:  status.Convert(s),
		se: s,
		e:  e,
	}
}
