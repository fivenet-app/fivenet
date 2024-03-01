// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/vehicles/vehicles.proto

package vehicles

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

// Validate checks the field values on Vehicle with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Vehicle) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Vehicle with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in VehicleMultiError, or nil if none found.
func (m *Vehicle) ValidateAll() error {
	return m.validate(true)
}

func (m *Vehicle) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Plate

	// no validation rules for Type

	if all {
		switch v := interface{}(m.GetOwner()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, VehicleValidationError{
					field:  "Owner",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, VehicleValidationError{
					field:  "Owner",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOwner()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return VehicleValidationError{
				field:  "Owner",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Model != nil {
		// no validation rules for Model
	}

	if len(errors) > 0 {
		return VehicleMultiError(errors)
	}

	return nil
}

// VehicleMultiError is an error wrapping multiple validation errors returned
// by Vehicle.ValidateAll() if the designated constraints aren't met.
type VehicleMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m VehicleMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m VehicleMultiError) AllErrors() []error { return m }

// VehicleValidationError is the validation error returned by Vehicle.Validate
// if the designated constraints aren't met.
type VehicleValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e VehicleValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e VehicleValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e VehicleValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e VehicleValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e VehicleValidationError) ErrorName() string { return "VehicleValidationError" }

// Error satisfies the builtin error interface
func (e VehicleValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sVehicle.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = VehicleValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = VehicleValidationError{}
