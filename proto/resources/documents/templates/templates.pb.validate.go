// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/documents/templates/templates.proto

package templates

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

// Validate checks the field values on TemplateData with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *TemplateData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TemplateData with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TemplateDataMultiError, or
// nil if none found.
func (m *TemplateData) ValidateAll() error {
	return m.validate(true)
}

func (m *TemplateData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetActiveChar()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TemplateDataValidationError{
					field:  "ActiveChar",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TemplateDataValidationError{
					field:  "ActiveChar",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetActiveChar()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TemplateDataValidationError{
				field:  "ActiveChar",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(m.GetUsers()) > 3 {
		err := TemplateDataValidationError{
			field:  "Users",
			reason: "value must contain no more than 3 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetUsers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, TemplateDataValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, TemplateDataValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TemplateDataValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return TemplateDataMultiError(errors)
	}

	return nil
}

// TemplateDataMultiError is an error wrapping multiple validation errors
// returned by TemplateData.ValidateAll() if the designated constraints aren't met.
type TemplateDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TemplateDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TemplateDataMultiError) AllErrors() []error { return m }

// TemplateDataValidationError is the validation error returned by
// TemplateData.Validate if the designated constraints aren't met.
type TemplateDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TemplateDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TemplateDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TemplateDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TemplateDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TemplateDataValidationError) ErrorName() string { return "TemplateDataValidationError" }

// Error satisfies the builtin error interface
func (e TemplateDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTemplateData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TemplateDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TemplateDataValidationError{}
