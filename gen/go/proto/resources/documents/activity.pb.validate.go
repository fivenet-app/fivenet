// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/documents/activity.proto

package documents

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

// Validate checks the field values on DocActivity with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DocActivity) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DocActivity with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DocActivityMultiError, or
// nil if none found.
func (m *DocActivity) ValidateAll() error {
	return m.validate(true)
}

func (m *DocActivity) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DocActivityValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DocActivityValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DocActivityValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for DocumentId

	// no validation rules for ActivityType

	if utf8.RuneCountInString(m.GetCreatorJob()) > 20 {
		err := DocActivityValidationError{
			field:  "CreatorJob",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetReason()) > 255 {
		err := DocActivityValidationError{
			field:  "Reason",
			reason: "value length must be at most 255 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetData()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DocActivityValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DocActivityValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DocActivityValidationError{
				field:  "Data",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.CreatorId != nil {
		// no validation rules for CreatorId
	}

	if m.Creator != nil {

		if all {
			switch v := interface{}(m.GetCreator()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DocActivityValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DocActivityValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreator()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DocActivityValidationError{
					field:  "Creator",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.CreatorJobLabel != nil {

		if utf8.RuneCountInString(m.GetCreatorJobLabel()) > 50 {
			err := DocActivityValidationError{
				field:  "CreatorJobLabel",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return DocActivityMultiError(errors)
	}

	return nil
}

// DocActivityMultiError is an error wrapping multiple validation errors
// returned by DocActivity.ValidateAll() if the designated constraints aren't met.
type DocActivityMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocActivityMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocActivityMultiError) AllErrors() []error { return m }

// DocActivityValidationError is the validation error returned by
// DocActivity.Validate if the designated constraints aren't met.
type DocActivityValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocActivityValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocActivityValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocActivityValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocActivityValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocActivityValidationError) ErrorName() string { return "DocActivityValidationError" }

// Error satisfies the builtin error interface
func (e DocActivityValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocActivity.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocActivityValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocActivityValidationError{}

// Validate checks the field values on DocActivityData with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DocActivityData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DocActivityData with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DocActivityDataMultiError, or nil if none found.
func (m *DocActivityData) ValidateAll() error {
	return m.validate(true)
}

func (m *DocActivityData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	switch v := m.Data.(type) {
	case *DocActivityData_Updated:
		if v == nil {
			err := DocActivityDataValidationError{
				field:  "Data",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetUpdated()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DocActivityDataValidationError{
						field:  "Updated",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DocActivityDataValidationError{
						field:  "Updated",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdated()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DocActivityDataValidationError{
					field:  "Updated",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return DocActivityDataMultiError(errors)
	}

	return nil
}

// DocActivityDataMultiError is an error wrapping multiple validation errors
// returned by DocActivityData.ValidateAll() if the designated constraints
// aren't met.
type DocActivityDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocActivityDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocActivityDataMultiError) AllErrors() []error { return m }

// DocActivityDataValidationError is the validation error returned by
// DocActivityData.Validate if the designated constraints aren't met.
type DocActivityDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocActivityDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocActivityDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocActivityDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocActivityDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocActivityDataValidationError) ErrorName() string { return "DocActivityDataValidationError" }

// Error satisfies the builtin error interface
func (e DocActivityDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocActivityData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocActivityDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocActivityDataValidationError{}

// Validate checks the field values on DocUpdated with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DocUpdated) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DocUpdated with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DocUpdatedMultiError, or
// nil if none found.
func (m *DocUpdated) ValidateAll() error {
	return m.validate(true)
}

func (m *DocUpdated) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Diff

	if len(errors) > 0 {
		return DocUpdatedMultiError(errors)
	}

	return nil
}

// DocUpdatedMultiError is an error wrapping multiple validation errors
// returned by DocUpdated.ValidateAll() if the designated constraints aren't met.
type DocUpdatedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocUpdatedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocUpdatedMultiError) AllErrors() []error { return m }

// DocUpdatedValidationError is the validation error returned by
// DocUpdated.Validate if the designated constraints aren't met.
type DocUpdatedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocUpdatedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocUpdatedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocUpdatedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocUpdatedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocUpdatedValidationError) ErrorName() string { return "DocUpdatedValidationError" }

// Error satisfies the builtin error interface
func (e DocUpdatedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocUpdated.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocUpdatedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocUpdatedValidationError{}
