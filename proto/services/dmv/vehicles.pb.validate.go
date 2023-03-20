// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: services/dmv/vehicles.proto

package dmv

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

// Validate checks the field values on FindVehiclesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *FindVehiclesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FindVehiclesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FindVehiclesRequestMultiError, or nil if none found.
func (m *FindVehiclesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *FindVehiclesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return FindVehiclesRequestMultiError(errors)
	}

	return nil
}

// FindVehiclesRequestMultiError is an error wrapping multiple validation
// errors returned by FindVehiclesRequest.ValidateAll() if the designated
// constraints aren't met.
type FindVehiclesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FindVehiclesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FindVehiclesRequestMultiError) AllErrors() []error { return m }

// FindVehiclesRequestValidationError is the validation error returned by
// FindVehiclesRequest.Validate if the designated constraints aren't met.
type FindVehiclesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FindVehiclesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FindVehiclesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FindVehiclesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FindVehiclesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FindVehiclesRequestValidationError) ErrorName() string {
	return "FindVehiclesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e FindVehiclesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFindVehiclesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FindVehiclesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FindVehiclesRequestValidationError{}

// Validate checks the field values on FindVehiclesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *FindVehiclesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FindVehiclesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FindVehiclesResponseMultiError, or nil if none found.
func (m *FindVehiclesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *FindVehiclesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return FindVehiclesResponseMultiError(errors)
	}

	return nil
}

// FindVehiclesResponseMultiError is an error wrapping multiple validation
// errors returned by FindVehiclesResponse.ValidateAll() if the designated
// constraints aren't met.
type FindVehiclesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FindVehiclesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FindVehiclesResponseMultiError) AllErrors() []error { return m }

// FindVehiclesResponseValidationError is the validation error returned by
// FindVehiclesResponse.Validate if the designated constraints aren't met.
type FindVehiclesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FindVehiclesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FindVehiclesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FindVehiclesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FindVehiclesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FindVehiclesResponseValidationError) ErrorName() string {
	return "FindVehiclesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e FindVehiclesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFindVehiclesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FindVehiclesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FindVehiclesResponseValidationError{}
