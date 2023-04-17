// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/notifications/notifications.proto

package notifications

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

// Validate checks the field values on Notification with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Notification) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Notification with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in NotificationMultiError, or
// nil if none found.
func (m *Notification) ValidateAll() error {
	return m.validate(true)
}

func (m *Notification) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, NotificationValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, NotificationValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return NotificationValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetReadAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, NotificationValidationError{
					field:  "ReadAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, NotificationValidationError{
					field:  "ReadAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetReadAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return NotificationValidationError{
				field:  "ReadAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for UserId

	if utf8.RuneCountInString(m.GetTitle()) > 255 {
		err := NotificationValidationError{
			field:  "Title",
			reason: "value length must be at most 255 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetType()) > 128 {
		err := NotificationValidationError{
			field:  "Type",
			reason: "value length must be at most 128 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetContent()) > 512 {
		err := NotificationValidationError{
			field:  "Content",
			reason: "value length must be at most 512 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.Data != nil {

		if len(m.GetData()) > 512 {
			err := NotificationValidationError{
				field:  "Data",
				reason: "value length must be at most 512 bytes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return NotificationMultiError(errors)
	}

	return nil
}

// NotificationMultiError is an error wrapping multiple validation errors
// returned by Notification.ValidateAll() if the designated constraints aren't met.
type NotificationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m NotificationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m NotificationMultiError) AllErrors() []error { return m }

// NotificationValidationError is the validation error returned by
// Notification.Validate if the designated constraints aren't met.
type NotificationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NotificationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NotificationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NotificationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NotificationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NotificationValidationError) ErrorName() string { return "NotificationValidationError" }

// Error satisfies the builtin error interface
func (e NotificationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNotification.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NotificationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NotificationValidationError{}
