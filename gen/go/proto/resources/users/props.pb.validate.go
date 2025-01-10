// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/users/props.proto

package users

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

// Validate checks the field values on UserProps with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserProps) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserProps with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserPropsMultiError, or nil
// if none found.
func (m *UserProps) ValidateAll() error {
	return m.validate(true)
}

func (m *UserProps) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUserId() <= 0 {
		err := UserPropsValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.UpdatedAt != nil {

		if all {
			switch v := interface{}(m.GetUpdatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserPropsValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Wanted != nil {
		// no validation rules for Wanted
	}

	if m.JobName != nil {
		// no validation rules for JobName
	}

	if m.Job != nil {

		if all {
			switch v := interface{}(m.GetJob()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "Job",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "Job",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetJob()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserPropsValidationError{
					field:  "Job",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.JobGradeNumber != nil {
		// no validation rules for JobGradeNumber
	}

	if m.JobGrade != nil {

		if all {
			switch v := interface{}(m.GetJobGrade()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "JobGrade",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "JobGrade",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetJobGrade()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserPropsValidationError{
					field:  "JobGrade",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.TrafficInfractionPoints != nil {
		// no validation rules for TrafficInfractionPoints
	}

	if m.OpenFines != nil {
		// no validation rules for OpenFines
	}

	if m.BloodType != nil {
		// no validation rules for BloodType
	}

	if m.MugShot != nil {

		if all {
			switch v := interface{}(m.GetMugShot()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "MugShot",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "MugShot",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMugShot()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserPropsValidationError{
					field:  "MugShot",
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
					errors = append(errors, UserPropsValidationError{
						field:  "Attributes",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserPropsValidationError{
						field:  "Attributes",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAttributes()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserPropsValidationError{
					field:  "Attributes",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Email != nil {

		if l := utf8.RuneCountInString(m.GetEmail()); l < 6 || l > 80 {
			err := UserPropsValidationError{
				field:  "Email",
				reason: "value length must be between 6 and 80 runes, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return UserPropsMultiError(errors)
	}

	return nil
}

// UserPropsMultiError is an error wrapping multiple validation errors returned
// by UserProps.ValidateAll() if the designated constraints aren't met.
type UserPropsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserPropsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserPropsMultiError) AllErrors() []error { return m }

// UserPropsValidationError is the validation error returned by
// UserProps.Validate if the designated constraints aren't met.
type UserPropsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserPropsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserPropsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserPropsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserPropsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserPropsValidationError) ErrorName() string { return "UserPropsValidationError" }

// Error satisfies the builtin error interface
func (e UserPropsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserProps.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserPropsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserPropsValidationError{}
