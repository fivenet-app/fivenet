// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/centrum/attributes.proto

package centrum

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

// Validate checks the field values on UnitAttributes with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UnitAttributes) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UnitAttributes with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UnitAttributesMultiError,
// or nil if none found.
func (m *UnitAttributes) ValidateAll() error {
	return m.validate(true)
}

func (m *UnitAttributes) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetList() {
		_, _ = idx, item

		if _, ok := UnitAttribute_name[int32(item)]; !ok {
			err := UnitAttributesValidationError{
				field:  fmt.Sprintf("List[%v]", idx),
				reason: "value must be one of the defined enum values",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return UnitAttributesMultiError(errors)
	}

	return nil
}

// UnitAttributesMultiError is an error wrapping multiple validation errors
// returned by UnitAttributes.ValidateAll() if the designated constraints
// aren't met.
type UnitAttributesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UnitAttributesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UnitAttributesMultiError) AllErrors() []error { return m }

// UnitAttributesValidationError is the validation error returned by
// UnitAttributes.Validate if the designated constraints aren't met.
type UnitAttributesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UnitAttributesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UnitAttributesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UnitAttributesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UnitAttributesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UnitAttributesValidationError) ErrorName() string { return "UnitAttributesValidationError" }

// Error satisfies the builtin error interface
func (e UnitAttributesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUnitAttributes.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UnitAttributesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UnitAttributesValidationError{}

// Validate checks the field values on DispatchAttributes with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DispatchAttributes) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DispatchAttributes with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DispatchAttributesMultiError, or nil if none found.
func (m *DispatchAttributes) ValidateAll() error {
	return m.validate(true)
}

func (m *DispatchAttributes) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetList() {
		_, _ = idx, item

		if _, ok := DispatchAttribute_name[int32(item)]; !ok {
			err := DispatchAttributesValidationError{
				field:  fmt.Sprintf("List[%v]", idx),
				reason: "value must be one of the defined enum values",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return DispatchAttributesMultiError(errors)
	}

	return nil
}

// DispatchAttributesMultiError is an error wrapping multiple validation errors
// returned by DispatchAttributes.ValidateAll() if the designated constraints
// aren't met.
type DispatchAttributesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DispatchAttributesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DispatchAttributesMultiError) AllErrors() []error { return m }

// DispatchAttributesValidationError is the validation error returned by
// DispatchAttributes.Validate if the designated constraints aren't met.
type DispatchAttributesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DispatchAttributesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DispatchAttributesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DispatchAttributesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DispatchAttributesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DispatchAttributesValidationError) ErrorName() string {
	return "DispatchAttributesValidationError"
}

// Error satisfies the builtin error interface
func (e DispatchAttributesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDispatchAttributes.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DispatchAttributesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DispatchAttributesValidationError{}
