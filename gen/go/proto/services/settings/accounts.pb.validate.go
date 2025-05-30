// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: services/settings/accounts.proto

package settings

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

// Validate checks the field values on ListAccountsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListAccountsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListAccountsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListAccountsRequestMultiError, or nil if none found.
func (m *ListAccountsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListAccountsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPagination() == nil {
		err := ListAccountsRequestValidationError{
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
				errors = append(errors, ListAccountsRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListAccountsRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListAccountsRequestValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Sort != nil {

		if all {
			switch v := interface{}(m.GetSort()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListAccountsRequestValidationError{
						field:  "Sort",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListAccountsRequestValidationError{
						field:  "Sort",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetSort()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListAccountsRequestValidationError{
					field:  "Sort",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.License != nil {

		if utf8.RuneCountInString(m.GetLicense()) > 64 {
			err := ListAccountsRequestValidationError{
				field:  "License",
				reason: "value length must be at most 64 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Enabled != nil {
		// no validation rules for Enabled
	}

	if len(errors) > 0 {
		return ListAccountsRequestMultiError(errors)
	}

	return nil
}

// ListAccountsRequestMultiError is an error wrapping multiple validation
// errors returned by ListAccountsRequest.ValidateAll() if the designated
// constraints aren't met.
type ListAccountsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListAccountsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListAccountsRequestMultiError) AllErrors() []error { return m }

// ListAccountsRequestValidationError is the validation error returned by
// ListAccountsRequest.Validate if the designated constraints aren't met.
type ListAccountsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListAccountsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListAccountsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListAccountsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListAccountsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListAccountsRequestValidationError) ErrorName() string {
	return "ListAccountsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListAccountsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListAccountsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListAccountsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListAccountsRequestValidationError{}

// Validate checks the field values on ListAccountsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListAccountsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListAccountsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListAccountsResponseMultiError, or nil if none found.
func (m *ListAccountsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListAccountsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPagination()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListAccountsResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListAccountsResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListAccountsResponseValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetAccounts() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListAccountsResponseValidationError{
						field:  fmt.Sprintf("Accounts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListAccountsResponseValidationError{
						field:  fmt.Sprintf("Accounts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListAccountsResponseValidationError{
					field:  fmt.Sprintf("Accounts[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListAccountsResponseMultiError(errors)
	}

	return nil
}

// ListAccountsResponseMultiError is an error wrapping multiple validation
// errors returned by ListAccountsResponse.ValidateAll() if the designated
// constraints aren't met.
type ListAccountsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListAccountsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListAccountsResponseMultiError) AllErrors() []error { return m }

// ListAccountsResponseValidationError is the validation error returned by
// ListAccountsResponse.Validate if the designated constraints aren't met.
type ListAccountsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListAccountsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListAccountsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListAccountsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListAccountsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListAccountsResponseValidationError) ErrorName() string {
	return "ListAccountsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListAccountsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListAccountsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListAccountsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListAccountsResponseValidationError{}

// Validate checks the field values on UpdateAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateAccountRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateAccountRequestMultiError, or nil if none found.
func (m *UpdateAccountRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateAccountRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() <= 0 {
		err := UpdateAccountRequestValidationError{
			field:  "Id",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.Enabled != nil {
		// no validation rules for Enabled
	}

	if m.LastChar != nil {
		// no validation rules for LastChar
	}

	if len(errors) > 0 {
		return UpdateAccountRequestMultiError(errors)
	}

	return nil
}

// UpdateAccountRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateAccountRequest.ValidateAll() if the designated
// constraints aren't met.
type UpdateAccountRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateAccountRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateAccountRequestMultiError) AllErrors() []error { return m }

// UpdateAccountRequestValidationError is the validation error returned by
// UpdateAccountRequest.Validate if the designated constraints aren't met.
type UpdateAccountRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateAccountRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateAccountRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateAccountRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateAccountRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateAccountRequestValidationError) ErrorName() string {
	return "UpdateAccountRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateAccountRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateAccountRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateAccountRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateAccountRequestValidationError{}

// Validate checks the field values on UpdateAccountResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateAccountResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateAccountResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateAccountResponseMultiError, or nil if none found.
func (m *UpdateAccountResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateAccountResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAccount()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateAccountResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateAccountResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccount()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateAccountResponseValidationError{
				field:  "Account",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdateAccountResponseMultiError(errors)
	}

	return nil
}

// UpdateAccountResponseMultiError is an error wrapping multiple validation
// errors returned by UpdateAccountResponse.ValidateAll() if the designated
// constraints aren't met.
type UpdateAccountResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateAccountResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateAccountResponseMultiError) AllErrors() []error { return m }

// UpdateAccountResponseValidationError is the validation error returned by
// UpdateAccountResponse.Validate if the designated constraints aren't met.
type UpdateAccountResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateAccountResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateAccountResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateAccountResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateAccountResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateAccountResponseValidationError) ErrorName() string {
	return "UpdateAccountResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateAccountResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateAccountResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateAccountResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateAccountResponseValidationError{}

// Validate checks the field values on DisconnectOAuth2ConnectionRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *DisconnectOAuth2ConnectionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DisconnectOAuth2ConnectionRequest
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// DisconnectOAuth2ConnectionRequestMultiError, or nil if none found.
func (m *DisconnectOAuth2ConnectionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DisconnectOAuth2ConnectionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() <= 0 {
		err := DisconnectOAuth2ConnectionRequestValidationError{
			field:  "Id",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetProviderName()) > 255 {
		err := DisconnectOAuth2ConnectionRequestValidationError{
			field:  "ProviderName",
			reason: "value length must be at most 255 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DisconnectOAuth2ConnectionRequestMultiError(errors)
	}

	return nil
}

// DisconnectOAuth2ConnectionRequestMultiError is an error wrapping multiple
// validation errors returned by
// DisconnectOAuth2ConnectionRequest.ValidateAll() if the designated
// constraints aren't met.
type DisconnectOAuth2ConnectionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DisconnectOAuth2ConnectionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DisconnectOAuth2ConnectionRequestMultiError) AllErrors() []error { return m }

// DisconnectOAuth2ConnectionRequestValidationError is the validation error
// returned by DisconnectOAuth2ConnectionRequest.Validate if the designated
// constraints aren't met.
type DisconnectOAuth2ConnectionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DisconnectOAuth2ConnectionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DisconnectOAuth2ConnectionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DisconnectOAuth2ConnectionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DisconnectOAuth2ConnectionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DisconnectOAuth2ConnectionRequestValidationError) ErrorName() string {
	return "DisconnectOAuth2ConnectionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DisconnectOAuth2ConnectionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDisconnectOAuth2ConnectionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DisconnectOAuth2ConnectionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DisconnectOAuth2ConnectionRequestValidationError{}

// Validate checks the field values on DisconnectOAuth2ConnectionResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *DisconnectOAuth2ConnectionResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DisconnectOAuth2ConnectionResponse
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// DisconnectOAuth2ConnectionResponseMultiError, or nil if none found.
func (m *DisconnectOAuth2ConnectionResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DisconnectOAuth2ConnectionResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DisconnectOAuth2ConnectionResponseMultiError(errors)
	}

	return nil
}

// DisconnectOAuth2ConnectionResponseMultiError is an error wrapping multiple
// validation errors returned by
// DisconnectOAuth2ConnectionResponse.ValidateAll() if the designated
// constraints aren't met.
type DisconnectOAuth2ConnectionResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DisconnectOAuth2ConnectionResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DisconnectOAuth2ConnectionResponseMultiError) AllErrors() []error { return m }

// DisconnectOAuth2ConnectionResponseValidationError is the validation error
// returned by DisconnectOAuth2ConnectionResponse.Validate if the designated
// constraints aren't met.
type DisconnectOAuth2ConnectionResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DisconnectOAuth2ConnectionResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DisconnectOAuth2ConnectionResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DisconnectOAuth2ConnectionResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DisconnectOAuth2ConnectionResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DisconnectOAuth2ConnectionResponseValidationError) ErrorName() string {
	return "DisconnectOAuth2ConnectionResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DisconnectOAuth2ConnectionResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDisconnectOAuth2ConnectionResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DisconnectOAuth2ConnectionResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DisconnectOAuth2ConnectionResponseValidationError{}

// Validate checks the field values on DeleteAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteAccountRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteAccountRequestMultiError, or nil if none found.
func (m *DeleteAccountRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteAccountRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() <= 0 {
		err := DeleteAccountRequestValidationError{
			field:  "Id",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteAccountRequestMultiError(errors)
	}

	return nil
}

// DeleteAccountRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteAccountRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteAccountRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteAccountRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteAccountRequestMultiError) AllErrors() []error { return m }

// DeleteAccountRequestValidationError is the validation error returned by
// DeleteAccountRequest.Validate if the designated constraints aren't met.
type DeleteAccountRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteAccountRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteAccountRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteAccountRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteAccountRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteAccountRequestValidationError) ErrorName() string {
	return "DeleteAccountRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteAccountRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteAccountRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteAccountRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteAccountRequestValidationError{}

// Validate checks the field values on DeleteAccountResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteAccountResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteAccountResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteAccountResponseMultiError, or nil if none found.
func (m *DeleteAccountResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteAccountResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DeleteAccountResponseMultiError(errors)
	}

	return nil
}

// DeleteAccountResponseMultiError is an error wrapping multiple validation
// errors returned by DeleteAccountResponse.ValidateAll() if the designated
// constraints aren't met.
type DeleteAccountResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteAccountResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteAccountResponseMultiError) AllErrors() []error { return m }

// DeleteAccountResponseValidationError is the validation error returned by
// DeleteAccountResponse.Validate if the designated constraints aren't met.
type DeleteAccountResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteAccountResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteAccountResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteAccountResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteAccountResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteAccountResponseValidationError) ErrorName() string {
	return "DeleteAccountResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteAccountResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteAccountResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteAccountResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteAccountResponseValidationError{}
