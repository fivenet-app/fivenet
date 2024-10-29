// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/centrum/general.proto

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

// Validate checks the field values on Attributes with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Attributes) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Attributes with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AttributesMultiError, or
// nil if none found.
func (m *Attributes) ValidateAll() error {
	return m.validate(true)
}

func (m *Attributes) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return AttributesMultiError(errors)
	}

	return nil
}

// AttributesMultiError is an error wrapping multiple validation errors
// returned by Attributes.ValidateAll() if the designated constraints aren't met.
type AttributesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AttributesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AttributesMultiError) AllErrors() []error { return m }

// AttributesValidationError is the validation error returned by
// Attributes.Validate if the designated constraints aren't met.
type AttributesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AttributesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AttributesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AttributesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AttributesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AttributesValidationError) ErrorName() string { return "AttributesValidationError" }

// Error satisfies the builtin error interface
func (e AttributesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAttributes.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AttributesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AttributesValidationError{}

// Validate checks the field values on Disponents with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Disponents) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Disponents with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DisponentsMultiError, or
// nil if none found.
func (m *Disponents) ValidateAll() error {
	return m.validate(true)
}

func (m *Disponents) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := DisponentsValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetDisponents() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DisponentsValidationError{
						field:  fmt.Sprintf("Disponents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DisponentsValidationError{
						field:  fmt.Sprintf("Disponents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DisponentsValidationError{
					field:  fmt.Sprintf("Disponents[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return DisponentsMultiError(errors)
	}

	return nil
}

// DisponentsMultiError is an error wrapping multiple validation errors
// returned by Disponents.ValidateAll() if the designated constraints aren't met.
type DisponentsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DisponentsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DisponentsMultiError) AllErrors() []error { return m }

// DisponentsValidationError is the validation error returned by
// Disponents.Validate if the designated constraints aren't met.
type DisponentsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DisponentsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DisponentsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DisponentsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DisponentsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DisponentsValidationError) ErrorName() string { return "DisponentsValidationError" }

// Error satisfies the builtin error interface
func (e DisponentsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDisponents.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DisponentsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DisponentsValidationError{}

// Validate checks the field values on UserUnitMapping with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserUnitMapping) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserUnitMapping with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserUnitMappingMultiError, or nil if none found.
func (m *UserUnitMapping) ValidateAll() error {
	return m.validate(true)
}

func (m *UserUnitMapping) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UnitId

	// no validation rules for Job

	if m.GetUserId() < 0 {
		err := UserUnitMappingValidationError{
			field:  "UserId",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserUnitMappingValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserUnitMappingValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserUnitMappingValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UserUnitMappingMultiError(errors)
	}

	return nil
}

// UserUnitMappingMultiError is an error wrapping multiple validation errors
// returned by UserUnitMapping.ValidateAll() if the designated constraints
// aren't met.
type UserUnitMappingMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserUnitMappingMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserUnitMappingMultiError) AllErrors() []error { return m }

// UserUnitMappingValidationError is the validation error returned by
// UserUnitMapping.Validate if the designated constraints aren't met.
type UserUnitMappingValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserUnitMappingValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserUnitMappingValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserUnitMappingValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserUnitMappingValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserUnitMappingValidationError) ErrorName() string { return "UserUnitMappingValidationError" }

// Error satisfies the builtin error interface
func (e UserUnitMappingValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserUnitMapping.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserUnitMappingValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserUnitMappingValidationError{}
