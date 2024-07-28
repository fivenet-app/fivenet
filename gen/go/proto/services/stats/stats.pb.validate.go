// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: services/stats/stats.proto

package stats

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

// Validate checks the field values on GetStatsRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetStatsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetStatsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetStatsRequestMultiError, or nil if none found.
func (m *GetStatsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetStatsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetStatsRequestMultiError(errors)
	}

	return nil
}

// GetStatsRequestMultiError is an error wrapping multiple validation errors
// returned by GetStatsRequest.ValidateAll() if the designated constraints
// aren't met.
type GetStatsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetStatsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetStatsRequestMultiError) AllErrors() []error { return m }

// GetStatsRequestValidationError is the validation error returned by
// GetStatsRequest.Validate if the designated constraints aren't met.
type GetStatsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetStatsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetStatsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetStatsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetStatsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetStatsRequestValidationError) ErrorName() string { return "GetStatsRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetStatsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetStatsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetStatsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetStatsRequestValidationError{}

// Validate checks the field values on GetStatsResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetStatsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetStatsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetStatsResponseMultiError, or nil if none found.
func (m *GetStatsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetStatsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	{
		sorted_keys := make([]string, len(m.GetStats()))
		i := 0
		for key := range m.GetStats() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetStats()[key]
			_ = val

			// no validation rules for Stats[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, GetStatsResponseValidationError{
							field:  fmt.Sprintf("Stats[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, GetStatsResponseValidationError{
							field:  fmt.Sprintf("Stats[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return GetStatsResponseValidationError{
						field:  fmt.Sprintf("Stats[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if len(errors) > 0 {
		return GetStatsResponseMultiError(errors)
	}

	return nil
}

// GetStatsResponseMultiError is an error wrapping multiple validation errors
// returned by GetStatsResponse.ValidateAll() if the designated constraints
// aren't met.
type GetStatsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetStatsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetStatsResponseMultiError) AllErrors() []error { return m }

// GetStatsResponseValidationError is the validation error returned by
// GetStatsResponse.Validate if the designated constraints aren't met.
type GetStatsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetStatsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetStatsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetStatsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetStatsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetStatsResponseValidationError) ErrorName() string { return "GetStatsResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetStatsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetStatsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetStatsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetStatsResponseValidationError{}
