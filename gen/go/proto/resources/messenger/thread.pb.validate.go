// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/messenger/thread.proto

package messenger

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

// Validate checks the field values on Thread with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Thread) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Thread with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ThreadMultiError, or nil if none found.
func (m *Thread) ValidateAll() error {
	return m.validate(true)
}

func (m *Thread) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ThreadValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ThreadValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ThreadValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if l := utf8.RuneCountInString(m.GetTitle()); l < 3 || l > 255 {
		err := ThreadValidationError{
			field:  "Title",
			reason: "value length must be between 3 and 255 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Archived

	if all {
		switch v := interface{}(m.GetUserState()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ThreadValidationError{
					field:  "UserState",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ThreadValidationError{
					field:  "UserState",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUserState()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ThreadValidationError{
				field:  "UserState",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for CreatorJob

	if all {
		switch v := interface{}(m.GetAccess()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ThreadValidationError{
					field:  "Access",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ThreadValidationError{
					field:  "Access",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccess()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ThreadValidationError{
				field:  "Access",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.UpdatedAt != nil {

		if all {
			switch v := interface{}(m.GetUpdatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.DeletedAt != nil {

		if all {
			switch v := interface{}(m.GetDeletedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDeletedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadValidationError{
					field:  "DeletedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.LastMessage != nil {

		if all {
			switch v := interface{}(m.GetLastMessage()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "LastMessage",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "LastMessage",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetLastMessage()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadValidationError{
					field:  "LastMessage",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.CreatorId != nil {
		// no validation rules for CreatorId
	}

	if m.Creator != nil {

		if all {
			switch v := interface{}(m.GetCreator()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreator()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadValidationError{
					field:  "Creator",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ThreadMultiError(errors)
	}

	return nil
}

// ThreadMultiError is an error wrapping multiple validation errors returned by
// Thread.ValidateAll() if the designated constraints aren't met.
type ThreadMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ThreadMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ThreadMultiError) AllErrors() []error { return m }

// ThreadValidationError is the validation error returned by Thread.Validate if
// the designated constraints aren't met.
type ThreadValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ThreadValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ThreadValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ThreadValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ThreadValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ThreadValidationError) ErrorName() string { return "ThreadValidationError" }

// Error satisfies the builtin error interface
func (e ThreadValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sThread.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ThreadValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ThreadValidationError{}

// Validate checks the field values on ThreadUserState with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ThreadUserState) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ThreadUserState with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ThreadUserStateMultiError, or nil if none found.
func (m *ThreadUserState) ValidateAll() error {
	return m.validate(true)
}

func (m *ThreadUserState) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ThreadId

	// no validation rules for UserId

	// no validation rules for Unread

	// no validation rules for Important

	// no validation rules for Favorite

	// no validation rules for Muted

	if m.LastRead != nil {

		if all {
			switch v := interface{}(m.GetLastRead()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ThreadUserStateValidationError{
						field:  "LastRead",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadUserStateValidationError{
						field:  "LastRead",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetLastRead()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadUserStateValidationError{
					field:  "LastRead",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ThreadUserStateMultiError(errors)
	}

	return nil
}

// ThreadUserStateMultiError is an error wrapping multiple validation errors
// returned by ThreadUserState.ValidateAll() if the designated constraints
// aren't met.
type ThreadUserStateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ThreadUserStateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ThreadUserStateMultiError) AllErrors() []error { return m }

// ThreadUserStateValidationError is the validation error returned by
// ThreadUserState.Validate if the designated constraints aren't met.
type ThreadUserStateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ThreadUserStateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ThreadUserStateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ThreadUserStateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ThreadUserStateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ThreadUserStateValidationError) ErrorName() string { return "ThreadUserStateValidationError" }

// Error satisfies the builtin error interface
func (e ThreadUserStateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sThreadUserState.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ThreadUserStateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ThreadUserStateValidationError{}