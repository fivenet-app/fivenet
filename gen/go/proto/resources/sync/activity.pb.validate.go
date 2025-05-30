// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/sync/activity.proto

package sync

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

// Validate checks the field values on UserOAuth2Conn with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserOAuth2Conn) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserOAuth2Conn with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserOAuth2ConnMultiError,
// or nil if none found.
func (m *UserOAuth2Conn) ValidateAll() error {
	return m.validate(true)
}

func (m *UserOAuth2Conn) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ProviderName

	// no validation rules for Identifier

	// no validation rules for ExternalId

	// no validation rules for Username

	if len(errors) > 0 {
		return UserOAuth2ConnMultiError(errors)
	}

	return nil
}

// UserOAuth2ConnMultiError is an error wrapping multiple validation errors
// returned by UserOAuth2Conn.ValidateAll() if the designated constraints
// aren't met.
type UserOAuth2ConnMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserOAuth2ConnMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserOAuth2ConnMultiError) AllErrors() []error { return m }

// UserOAuth2ConnValidationError is the validation error returned by
// UserOAuth2Conn.Validate if the designated constraints aren't met.
type UserOAuth2ConnValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserOAuth2ConnValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserOAuth2ConnValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserOAuth2ConnValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserOAuth2ConnValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserOAuth2ConnValidationError) ErrorName() string { return "UserOAuth2ConnValidationError" }

// Error satisfies the builtin error interface
func (e UserOAuth2ConnValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserOAuth2Conn.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserOAuth2ConnValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserOAuth2ConnValidationError{}

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

	if m.GetProps() == nil {
		err := UserPropsValidationError{
			field:  "Props",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetProps()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserPropsValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserPropsValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProps()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserPropsValidationError{
				field:  "Props",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Reason != nil {

		if utf8.RuneCountInString(m.GetReason()) > 255 {
			err := UserPropsValidationError{
				field:  "Reason",
				reason: "value length must be at most 255 runes",
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

// Validate checks the field values on ColleagueProps with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ColleagueProps) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ColleagueProps with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ColleaguePropsMultiError,
// or nil if none found.
func (m *ColleagueProps) ValidateAll() error {
	return m.validate(true)
}

func (m *ColleagueProps) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetProps() == nil {
		err := ColleaguePropsValidationError{
			field:  "Props",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetProps()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ColleaguePropsValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ColleaguePropsValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProps()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ColleaguePropsValidationError{
				field:  "Props",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Reason != nil {

		if utf8.RuneCountInString(m.GetReason()) > 255 {
			err := ColleaguePropsValidationError{
				field:  "Reason",
				reason: "value length must be at most 255 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return ColleaguePropsMultiError(errors)
	}

	return nil
}

// ColleaguePropsMultiError is an error wrapping multiple validation errors
// returned by ColleagueProps.ValidateAll() if the designated constraints
// aren't met.
type ColleaguePropsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ColleaguePropsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ColleaguePropsMultiError) AllErrors() []error { return m }

// ColleaguePropsValidationError is the validation error returned by
// ColleagueProps.Validate if the designated constraints aren't met.
type ColleaguePropsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ColleaguePropsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ColleaguePropsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ColleaguePropsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ColleaguePropsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ColleaguePropsValidationError) ErrorName() string { return "ColleaguePropsValidationError" }

// Error satisfies the builtin error interface
func (e ColleaguePropsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sColleagueProps.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ColleaguePropsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ColleaguePropsValidationError{}

// Validate checks the field values on UserUpdate with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserUpdate) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserUpdate with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserUpdateMultiError, or
// nil if none found.
func (m *UserUpdate) ValidateAll() error {
	return m.validate(true)
}

func (m *UserUpdate) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	if m.Group != nil {
		// no validation rules for Group
	}

	if m.Job != nil {

		if utf8.RuneCountInString(m.GetJob()) > 20 {
			err := UserUpdateValidationError{
				field:  "Job",
				reason: "value length must be at most 20 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.JobGrade != nil {
		// no validation rules for JobGrade
	}

	if m.Firstname != nil {
		// no validation rules for Firstname
	}

	if m.Lastname != nil {
		// no validation rules for Lastname
	}

	if len(errors) > 0 {
		return UserUpdateMultiError(errors)
	}

	return nil
}

// UserUpdateMultiError is an error wrapping multiple validation errors
// returned by UserUpdate.ValidateAll() if the designated constraints aren't met.
type UserUpdateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserUpdateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserUpdateMultiError) AllErrors() []error { return m }

// UserUpdateValidationError is the validation error returned by
// UserUpdate.Validate if the designated constraints aren't met.
type UserUpdateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserUpdateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserUpdateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserUpdateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserUpdateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserUpdateValidationError) ErrorName() string { return "UserUpdateValidationError" }

// Error satisfies the builtin error interface
func (e UserUpdateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserUpdate.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserUpdateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserUpdateValidationError{}

// Validate checks the field values on TimeclockUpdate with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *TimeclockUpdate) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TimeclockUpdate with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TimeclockUpdateMultiError, or nil if none found.
func (m *TimeclockUpdate) ValidateAll() error {
	return m.validate(true)
}

func (m *TimeclockUpdate) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Job

	// no validation rules for UserId

	// no validation rules for Start

	if len(errors) > 0 {
		return TimeclockUpdateMultiError(errors)
	}

	return nil
}

// TimeclockUpdateMultiError is an error wrapping multiple validation errors
// returned by TimeclockUpdate.ValidateAll() if the designated constraints
// aren't met.
type TimeclockUpdateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TimeclockUpdateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TimeclockUpdateMultiError) AllErrors() []error { return m }

// TimeclockUpdateValidationError is the validation error returned by
// TimeclockUpdate.Validate if the designated constraints aren't met.
type TimeclockUpdateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TimeclockUpdateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TimeclockUpdateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TimeclockUpdateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TimeclockUpdateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TimeclockUpdateValidationError) ErrorName() string { return "TimeclockUpdateValidationError" }

// Error satisfies the builtin error interface
func (e TimeclockUpdateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTimeclockUpdate.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TimeclockUpdateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TimeclockUpdateValidationError{}
