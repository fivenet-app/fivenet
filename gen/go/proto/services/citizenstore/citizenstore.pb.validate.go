// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: services/citizenstore/citizenstore.proto

package citizenstore

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

// Validate checks the field values on ListCitizensRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListCitizensRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListCitizensRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListCitizensRequestMultiError, or nil if none found.
func (m *ListCitizensRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListCitizensRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPagination() == nil {
		err := ListCitizensRequestValidationError{
			field:  "Pagination",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetPagination()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListCitizensRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListCitizensRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListCitizensRequestValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if utf8.RuneCountInString(m.GetSearchName()) > 50 {
		err := ListCitizensRequestValidationError{
			field:  "SearchName",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.Wanted != nil {
		// no validation rules for Wanted
	}

	if m.PhoneNumber != nil {

		if utf8.RuneCountInString(m.GetPhoneNumber()) > 20 {
			err := ListCitizensRequestValidationError{
				field:  "PhoneNumber",
				reason: "value length must be at most 20 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.TrafficPoints != nil {
		// no validation rules for TrafficPoints
	}

	if len(errors) > 0 {
		return ListCitizensRequestMultiError(errors)
	}

	return nil
}

// ListCitizensRequestMultiError is an error wrapping multiple validation
// errors returned by ListCitizensRequest.ValidateAll() if the designated
// constraints aren't met.
type ListCitizensRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListCitizensRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListCitizensRequestMultiError) AllErrors() []error { return m }

// ListCitizensRequestValidationError is the validation error returned by
// ListCitizensRequest.Validate if the designated constraints aren't met.
type ListCitizensRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCitizensRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCitizensRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCitizensRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCitizensRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCitizensRequestValidationError) ErrorName() string {
	return "ListCitizensRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListCitizensRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCitizensRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCitizensRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCitizensRequestValidationError{}

// Validate checks the field values on ListCitizensResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListCitizensResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListCitizensResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListCitizensResponseMultiError, or nil if none found.
func (m *ListCitizensResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListCitizensResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPagination()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListCitizensResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListCitizensResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListCitizensResponseValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetUsers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListCitizensResponseValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListCitizensResponseValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListCitizensResponseValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListCitizensResponseMultiError(errors)
	}

	return nil
}

// ListCitizensResponseMultiError is an error wrapping multiple validation
// errors returned by ListCitizensResponse.ValidateAll() if the designated
// constraints aren't met.
type ListCitizensResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListCitizensResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListCitizensResponseMultiError) AllErrors() []error { return m }

// ListCitizensResponseValidationError is the validation error returned by
// ListCitizensResponse.Validate if the designated constraints aren't met.
type ListCitizensResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCitizensResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCitizensResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCitizensResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCitizensResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCitizensResponseValidationError) ErrorName() string {
	return "ListCitizensResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListCitizensResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCitizensResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCitizensResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCitizensResponseValidationError{}

// Validate checks the field values on GetUserRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetUserRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetUserRequestMultiError,
// or nil if none found.
func (m *GetUserRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUserId() <= 0 {
		err := GetUserRequestValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetUserRequestMultiError(errors)
	}

	return nil
}

// GetUserRequestMultiError is an error wrapping multiple validation errors
// returned by GetUserRequest.ValidateAll() if the designated constraints
// aren't met.
type GetUserRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserRequestMultiError) AllErrors() []error { return m }

// GetUserRequestValidationError is the validation error returned by
// GetUserRequest.Validate if the designated constraints aren't met.
type GetUserRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserRequestValidationError) ErrorName() string { return "GetUserRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetUserRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserRequestValidationError{}

// Validate checks the field values on GetUserResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetUserResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserResponseMultiError, or nil if none found.
func (m *GetUserResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetUser()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetUserResponseValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetUserResponseValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUser()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetUserResponseValidationError{
				field:  "User",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetUserResponseMultiError(errors)
	}

	return nil
}

// GetUserResponseMultiError is an error wrapping multiple validation errors
// returned by GetUserResponse.ValidateAll() if the designated constraints
// aren't met.
type GetUserResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserResponseMultiError) AllErrors() []error { return m }

// GetUserResponseValidationError is the validation error returned by
// GetUserResponse.Validate if the designated constraints aren't met.
type GetUserResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserResponseValidationError) ErrorName() string { return "GetUserResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetUserResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserResponseValidationError{}

// Validate checks the field values on ListUserActivityRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListUserActivityRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListUserActivityRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListUserActivityRequestMultiError, or nil if none found.
func (m *ListUserActivityRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListUserActivityRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPagination() == nil {
		err := ListUserActivityRequestValidationError{
			field:  "Pagination",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetPagination()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListUserActivityRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListUserActivityRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListUserActivityRequestValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.GetUserId() <= 0 {
		err := ListUserActivityRequestValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ListUserActivityRequestMultiError(errors)
	}

	return nil
}

// ListUserActivityRequestMultiError is an error wrapping multiple validation
// errors returned by ListUserActivityRequest.ValidateAll() if the designated
// constraints aren't met.
type ListUserActivityRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListUserActivityRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListUserActivityRequestMultiError) AllErrors() []error { return m }

// ListUserActivityRequestValidationError is the validation error returned by
// ListUserActivityRequest.Validate if the designated constraints aren't met.
type ListUserActivityRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListUserActivityRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListUserActivityRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListUserActivityRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListUserActivityRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListUserActivityRequestValidationError) ErrorName() string {
	return "ListUserActivityRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListUserActivityRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListUserActivityRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListUserActivityRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListUserActivityRequestValidationError{}

// Validate checks the field values on ListUserActivityResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListUserActivityResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListUserActivityResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListUserActivityResponseMultiError, or nil if none found.
func (m *ListUserActivityResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListUserActivityResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPagination()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListUserActivityResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListUserActivityResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListUserActivityResponseValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetActivity() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListUserActivityResponseValidationError{
						field:  fmt.Sprintf("Activity[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListUserActivityResponseValidationError{
						field:  fmt.Sprintf("Activity[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListUserActivityResponseValidationError{
					field:  fmt.Sprintf("Activity[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListUserActivityResponseMultiError(errors)
	}

	return nil
}

// ListUserActivityResponseMultiError is an error wrapping multiple validation
// errors returned by ListUserActivityResponse.ValidateAll() if the designated
// constraints aren't met.
type ListUserActivityResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListUserActivityResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListUserActivityResponseMultiError) AllErrors() []error { return m }

// ListUserActivityResponseValidationError is the validation error returned by
// ListUserActivityResponse.Validate if the designated constraints aren't met.
type ListUserActivityResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListUserActivityResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListUserActivityResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListUserActivityResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListUserActivityResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListUserActivityResponseValidationError) ErrorName() string {
	return "ListUserActivityResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListUserActivityResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListUserActivityResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListUserActivityResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListUserActivityResponseValidationError{}

// Validate checks the field values on SetUserPropsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SetUserPropsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SetUserPropsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SetUserPropsRequestMultiError, or nil if none found.
func (m *SetUserPropsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SetUserPropsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetProps() == nil {
		err := SetUserPropsRequestValidationError{
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
				errors = append(errors, SetUserPropsRequestValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SetUserPropsRequestValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProps()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SetUserPropsRequestValidationError{
				field:  "Props",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.GetReason() != "" {

		if l := utf8.RuneCountInString(m.GetReason()); l < 3 || l > 255 {
			err := SetUserPropsRequestValidationError{
				field:  "Reason",
				reason: "value length must be between 3 and 255 runes, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return SetUserPropsRequestMultiError(errors)
	}

	return nil
}

// SetUserPropsRequestMultiError is an error wrapping multiple validation
// errors returned by SetUserPropsRequest.ValidateAll() if the designated
// constraints aren't met.
type SetUserPropsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SetUserPropsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SetUserPropsRequestMultiError) AllErrors() []error { return m }

// SetUserPropsRequestValidationError is the validation error returned by
// SetUserPropsRequest.Validate if the designated constraints aren't met.
type SetUserPropsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetUserPropsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetUserPropsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetUserPropsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetUserPropsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetUserPropsRequestValidationError) ErrorName() string {
	return "SetUserPropsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e SetUserPropsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetUserPropsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetUserPropsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetUserPropsRequestValidationError{}

// Validate checks the field values on SetUserPropsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SetUserPropsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SetUserPropsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SetUserPropsResponseMultiError, or nil if none found.
func (m *SetUserPropsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SetUserPropsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetProps()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SetUserPropsResponseValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SetUserPropsResponseValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProps()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SetUserPropsResponseValidationError{
				field:  "Props",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SetUserPropsResponseMultiError(errors)
	}

	return nil
}

// SetUserPropsResponseMultiError is an error wrapping multiple validation
// errors returned by SetUserPropsResponse.ValidateAll() if the designated
// constraints aren't met.
type SetUserPropsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SetUserPropsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SetUserPropsResponseMultiError) AllErrors() []error { return m }

// SetUserPropsResponseValidationError is the validation error returned by
// SetUserPropsResponse.Validate if the designated constraints aren't met.
type SetUserPropsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetUserPropsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetUserPropsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetUserPropsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetUserPropsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetUserPropsResponseValidationError) ErrorName() string {
	return "SetUserPropsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e SetUserPropsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetUserPropsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetUserPropsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetUserPropsResponseValidationError{}
