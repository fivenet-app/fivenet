// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/centrum/units.proto

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

// Validate checks the field values on Unit with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Unit) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Unit with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in UnitMultiError, or nil if none found.
func (m *Unit) ValidateAll() error {
	return m.validate(true)
}

func (m *Unit) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := UnitValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetName()); l < 3 || l > 24 {
		err := UnitValidationError{
			field:  "Name",
			reason: "value length must be between 3 and 24 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetInitials()); l < 2 || l > 4 {
		err := UnitValidationError{
			field:  "Initials",
			reason: "value length must be between 2 and 4 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetColor()) != 7 {
		err := UnitValidationError{
			field:  "Color",
			reason: "value length must be 7 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if !_Unit_Color_Pattern.MatchString(m.GetColor()) {
		err := UnitValidationError{
			field:  "Color",
			reason: "value does not match regex pattern \"^#[A-Fa-f0-9]{6}$\"",
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
					errors = append(errors, UnitValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetAccess()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UnitValidationError{
					field:  "Access",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UnitValidationError{
					field:  "Access",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccess()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UnitValidationError{
				field:  "Access",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.CreatedAt != nil {

		if all {
			switch v := interface{}(m.GetCreatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.UpdatedAt != nil {

		if all {
			switch v := interface{}(m.GetUpdatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Description != nil {

		if utf8.RuneCountInString(m.GetDescription()) > 255 {
			err := UnitValidationError{
				field:  "Description",
				reason: "value length must be at most 255 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Status != nil {

		if all {
			switch v := interface{}(m.GetStatus()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "Status",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "Status",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetStatus()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitValidationError{
					field:  "Status",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Attributes != nil {

		if all {
			switch v := interface{}(m.GetAttributes()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "Attributes",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitValidationError{
						field:  "Attributes",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAttributes()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitValidationError{
					field:  "Attributes",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.HomePostal != nil {

		if utf8.RuneCountInString(m.GetHomePostal()) > 48 {
			err := UnitValidationError{
				field:  "HomePostal",
				reason: "value length must be at most 48 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return UnitMultiError(errors)
	}

	return nil
}

// UnitMultiError is an error wrapping multiple validation errors returned by
// Unit.ValidateAll() if the designated constraints aren't met.
type UnitMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UnitMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UnitMultiError) AllErrors() []error { return m }

// UnitValidationError is the validation error returned by Unit.Validate if the
// designated constraints aren't met.
type UnitValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UnitValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UnitValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UnitValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UnitValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UnitValidationError) ErrorName() string { return "UnitValidationError" }

// Error satisfies the builtin error interface
func (e UnitValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUnit.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UnitValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UnitValidationError{}

var _Unit_Color_Pattern = regexp.MustCompile("^#[A-Fa-f0-9]{6}$")

// Validate checks the field values on UnitAssignments with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UnitAssignments) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UnitAssignments with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UnitAssignmentsMultiError, or nil if none found.
func (m *UnitAssignments) ValidateAll() error {
	return m.validate(true)
}

func (m *UnitAssignments) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UnitId

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := UnitAssignmentsValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
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
					errors = append(errors, UnitAssignmentsValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitAssignmentsValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitAssignmentsValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UnitAssignmentsMultiError(errors)
	}

	return nil
}

// UnitAssignmentsMultiError is an error wrapping multiple validation errors
// returned by UnitAssignments.ValidateAll() if the designated constraints
// aren't met.
type UnitAssignmentsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UnitAssignmentsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UnitAssignmentsMultiError) AllErrors() []error { return m }

// UnitAssignmentsValidationError is the validation error returned by
// UnitAssignments.Validate if the designated constraints aren't met.
type UnitAssignmentsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UnitAssignmentsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UnitAssignmentsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UnitAssignmentsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UnitAssignmentsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UnitAssignmentsValidationError) ErrorName() string { return "UnitAssignmentsValidationError" }

// Error satisfies the builtin error interface
func (e UnitAssignmentsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUnitAssignments.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UnitAssignmentsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UnitAssignmentsValidationError{}

// Validate checks the field values on UnitAssignment with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UnitAssignment) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UnitAssignment with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UnitAssignmentMultiError,
// or nil if none found.
func (m *UnitAssignment) ValidateAll() error {
	return m.validate(true)
}

func (m *UnitAssignment) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UnitId

	if m.GetUserId() < 0 {
		err := UnitAssignmentValidationError{
			field:  "UserId",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.User != nil {

		if all {
			switch v := interface{}(m.GetUser()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitAssignmentValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitAssignmentValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUser()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitAssignmentValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UnitAssignmentMultiError(errors)
	}

	return nil
}

// UnitAssignmentMultiError is an error wrapping multiple validation errors
// returned by UnitAssignment.ValidateAll() if the designated constraints
// aren't met.
type UnitAssignmentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UnitAssignmentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UnitAssignmentMultiError) AllErrors() []error { return m }

// UnitAssignmentValidationError is the validation error returned by
// UnitAssignment.Validate if the designated constraints aren't met.
type UnitAssignmentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UnitAssignmentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UnitAssignmentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UnitAssignmentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UnitAssignmentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UnitAssignmentValidationError) ErrorName() string { return "UnitAssignmentValidationError" }

// Error satisfies the builtin error interface
func (e UnitAssignmentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUnitAssignment.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UnitAssignmentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UnitAssignmentValidationError{}

// Validate checks the field values on UnitStatus with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UnitStatus) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UnitStatus with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UnitStatusMultiError, or
// nil if none found.
func (m *UnitStatus) ValidateAll() error {
	return m.validate(true)
}

func (m *UnitStatus) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for UnitId

	if _, ok := StatusUnit_name[int32(m.GetStatus())]; !ok {
		err := UnitStatusValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.CreatedAt != nil {

		if all {
			switch v := interface{}(m.GetCreatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitStatusValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Unit != nil {

		if all {
			switch v := interface{}(m.GetUnit()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "Unit",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "Unit",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUnit()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitStatusValidationError{
					field:  "Unit",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Reason != nil {

		if utf8.RuneCountInString(m.GetReason()) > 255 {
			err := UnitStatusValidationError{
				field:  "Reason",
				reason: "value length must be at most 255 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Code != nil {

		if utf8.RuneCountInString(m.GetCode()) > 20 {
			err := UnitStatusValidationError{
				field:  "Code",
				reason: "value length must be at most 20 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.UserId != nil {

		if m.GetUserId() <= 0 {
			err := UnitStatusValidationError{
				field:  "UserId",
				reason: "value must be greater than 0",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.User != nil {

		if all {
			switch v := interface{}(m.GetUser()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUser()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitStatusValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.X != nil {
		// no validation rules for X
	}

	if m.Y != nil {
		// no validation rules for Y
	}

	if m.Postal != nil {

		if utf8.RuneCountInString(m.GetPostal()) > 48 {
			err := UnitStatusValidationError{
				field:  "Postal",
				reason: "value length must be at most 48 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.CreatorId != nil {

		if m.GetCreatorId() <= 0 {
			err := UnitStatusValidationError{
				field:  "CreatorId",
				reason: "value must be greater than 0",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Creator != nil {

		if all {
			switch v := interface{}(m.GetCreator()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UnitStatusValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreator()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UnitStatusValidationError{
					field:  "Creator",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UnitStatusMultiError(errors)
	}

	return nil
}

// UnitStatusMultiError is an error wrapping multiple validation errors
// returned by UnitStatus.ValidateAll() if the designated constraints aren't met.
type UnitStatusMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UnitStatusMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UnitStatusMultiError) AllErrors() []error { return m }

// UnitStatusValidationError is the validation error returned by
// UnitStatus.Validate if the designated constraints aren't met.
type UnitStatusValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UnitStatusValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UnitStatusValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UnitStatusValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UnitStatusValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UnitStatusValidationError) ErrorName() string { return "UnitStatusValidationError" }

// Error satisfies the builtin error interface
func (e UnitStatusValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUnitStatus.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UnitStatusValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UnitStatusValidationError{}
