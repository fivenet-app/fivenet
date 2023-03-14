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

// Validate checks the field values on FindUsersRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *FindUsersRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FindUsersRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FindUsersRequestMultiError, or nil if none found.
func (m *FindUsersRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *FindUsersRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetOffset() < 0 {
		err := FindUsersRequestValidationError{
			field:  "Offset",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetOrderBy()) > 3 {
		err := FindUsersRequestValidationError{
			field:  "OrderBy",
			reason: "value must contain no more than 3 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetOrderBy() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, FindUsersRequestValidationError{
						field:  fmt.Sprintf("OrderBy[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, FindUsersRequestValidationError{
						field:  fmt.Sprintf("OrderBy[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FindUsersRequestValidationError{
					field:  fmt.Sprintf("OrderBy[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if utf8.RuneCountInString(m.GetSearchName()) > 50 {
		err := FindUsersRequestValidationError{
			field:  "SearchName",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Wanted

	if len(errors) > 0 {
		return FindUsersRequestMultiError(errors)
	}

	return nil
}

// FindUsersRequestMultiError is an error wrapping multiple validation errors
// returned by FindUsersRequest.ValidateAll() if the designated constraints
// aren't met.
type FindUsersRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FindUsersRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FindUsersRequestMultiError) AllErrors() []error { return m }

// FindUsersRequestValidationError is the validation error returned by
// FindUsersRequest.Validate if the designated constraints aren't met.
type FindUsersRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FindUsersRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FindUsersRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FindUsersRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FindUsersRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FindUsersRequestValidationError) ErrorName() string { return "FindUsersRequestValidationError" }

// Error satisfies the builtin error interface
func (e FindUsersRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFindUsersRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FindUsersRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FindUsersRequestValidationError{}

// Validate checks the field values on FindUsersResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *FindUsersResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FindUsersResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FindUsersResponseMultiError, or nil if none found.
func (m *FindUsersResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *FindUsersResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TotalCount

	// no validation rules for Offset

	// no validation rules for End

	for idx, item := range m.GetUsers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, FindUsersResponseValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, FindUsersResponseValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FindUsersResponseValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return FindUsersResponseMultiError(errors)
	}

	return nil
}

// FindUsersResponseMultiError is an error wrapping multiple validation errors
// returned by FindUsersResponse.ValidateAll() if the designated constraints
// aren't met.
type FindUsersResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FindUsersResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FindUsersResponseMultiError) AllErrors() []error { return m }

// FindUsersResponseValidationError is the validation error returned by
// FindUsersResponse.Validate if the designated constraints aren't met.
type FindUsersResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FindUsersResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FindUsersResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FindUsersResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FindUsersResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FindUsersResponseValidationError) ErrorName() string {
	return "FindUsersResponseValidationError"
}

// Error satisfies the builtin error interface
func (e FindUsersResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFindUsersResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FindUsersResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FindUsersResponseValidationError{}

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

// Validate checks the field values on GetUserActivityRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserActivityRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserActivityRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserActivityRequestMultiError, or nil if none found.
func (m *GetUserActivityRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserActivityRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUserId() <= 0 {
		err := GetUserActivityRequestValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetUserActivityRequestMultiError(errors)
	}

	return nil
}

// GetUserActivityRequestMultiError is an error wrapping multiple validation
// errors returned by GetUserActivityRequest.ValidateAll() if the designated
// constraints aren't met.
type GetUserActivityRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserActivityRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserActivityRequestMultiError) AllErrors() []error { return m }

// GetUserActivityRequestValidationError is the validation error returned by
// GetUserActivityRequest.Validate if the designated constraints aren't met.
type GetUserActivityRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserActivityRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserActivityRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserActivityRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserActivityRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserActivityRequestValidationError) ErrorName() string {
	return "GetUserActivityRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserActivityRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserActivityRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserActivityRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserActivityRequestValidationError{}

// Validate checks the field values on GetUserActivityResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserActivityResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserActivityResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserActivityResponseMultiError, or nil if none found.
func (m *GetUserActivityResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserActivityResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetActivity() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetUserActivityResponseValidationError{
						field:  fmt.Sprintf("Activity[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetUserActivityResponseValidationError{
						field:  fmt.Sprintf("Activity[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetUserActivityResponseValidationError{
					field:  fmt.Sprintf("Activity[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetUserActivityResponseMultiError(errors)
	}

	return nil
}

// GetUserActivityResponseMultiError is an error wrapping multiple validation
// errors returned by GetUserActivityResponse.ValidateAll() if the designated
// constraints aren't met.
type GetUserActivityResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserActivityResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserActivityResponseMultiError) AllErrors() []error { return m }

// GetUserActivityResponseValidationError is the validation error returned by
// GetUserActivityResponse.Validate if the designated constraints aren't met.
type GetUserActivityResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserActivityResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserActivityResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserActivityResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserActivityResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserActivityResponseValidationError) ErrorName() string {
	return "GetUserActivityResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserActivityResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserActivityResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserActivityResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserActivityResponseValidationError{}

// Validate checks the field values on GetUserDocumentsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserDocumentsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserDocumentsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserDocumentsRequestMultiError, or nil if none found.
func (m *GetUserDocumentsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserDocumentsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetOffset() < 0 {
		err := GetUserDocumentsRequestValidationError{
			field:  "Offset",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetUserId() <= 0 {
		err := GetUserDocumentsRequestValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetUserDocumentsRequestMultiError(errors)
	}

	return nil
}

// GetUserDocumentsRequestMultiError is an error wrapping multiple validation
// errors returned by GetUserDocumentsRequest.ValidateAll() if the designated
// constraints aren't met.
type GetUserDocumentsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserDocumentsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserDocumentsRequestMultiError) AllErrors() []error { return m }

// GetUserDocumentsRequestValidationError is the validation error returned by
// GetUserDocumentsRequest.Validate if the designated constraints aren't met.
type GetUserDocumentsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserDocumentsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserDocumentsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserDocumentsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserDocumentsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserDocumentsRequestValidationError) ErrorName() string {
	return "GetUserDocumentsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserDocumentsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserDocumentsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserDocumentsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserDocumentsRequestValidationError{}

// Validate checks the field values on GetUserDocumentsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserDocumentsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserDocumentsResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserDocumentsResponseMultiError, or nil if none found.
func (m *GetUserDocumentsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserDocumentsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TotalCount

	// no validation rules for Offset

	// no validation rules for End

	for idx, item := range m.GetDocuments() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetUserDocumentsResponseValidationError{
						field:  fmt.Sprintf("Documents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetUserDocumentsResponseValidationError{
						field:  fmt.Sprintf("Documents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetUserDocumentsResponseValidationError{
					field:  fmt.Sprintf("Documents[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetUserDocumentsResponseMultiError(errors)
	}

	return nil
}

// GetUserDocumentsResponseMultiError is an error wrapping multiple validation
// errors returned by GetUserDocumentsResponse.ValidateAll() if the designated
// constraints aren't met.
type GetUserDocumentsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserDocumentsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserDocumentsResponseMultiError) AllErrors() []error { return m }

// GetUserDocumentsResponseValidationError is the validation error returned by
// GetUserDocumentsResponse.Validate if the designated constraints aren't met.
type GetUserDocumentsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserDocumentsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserDocumentsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserDocumentsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserDocumentsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserDocumentsResponseValidationError) ErrorName() string {
	return "GetUserDocumentsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserDocumentsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserDocumentsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserDocumentsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserDocumentsResponseValidationError{}

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
