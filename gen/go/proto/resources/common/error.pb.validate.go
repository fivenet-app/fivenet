// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/common/error.proto

package common

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Error with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Error) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Error with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ErrorMultiError, or nil if none found.
func (m *Error) ValidateAll() error {
	return m.validate(true)
}

func (m *Error) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetContent()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ErrorValidationError{
					field:  "Content",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ErrorValidationError{
					field:  "Content",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetContent()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ErrorValidationError{
				field:  "Content",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Title != nil {

		if all {
			switch v := interface{}(m.GetTitle()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ErrorValidationError{
						field:  "Title",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ErrorValidationError{
						field:  "Title",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetTitle()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ErrorValidationError{
					field:  "Title",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ErrorMultiError(errors)
	}

	return nil
}

// ErrorMultiError is an error wrapping multiple validation errors returned by
// Error.ValidateAll() if the designated constraints aren't met.
type ErrorMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ErrorMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ErrorMultiError) AllErrors() []error { return m }

// ErrorValidationError is the validation error returned by Error.Validate if
// the designated constraints aren't met.
type ErrorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ErrorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ErrorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ErrorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ErrorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ErrorValidationError) ErrorName() string { return "ErrorValidationError" }

// Error satisfies the builtin error interface
func (e ErrorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sError.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ErrorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ErrorValidationError{}
