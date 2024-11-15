// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/mailer/settings.proto

package mailer

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

// Validate checks the field values on EmailSettings with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *EmailSettings) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EmailSettings with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in EmailSettingsMultiError, or
// nil if none found.
func (m *EmailSettings) ValidateAll() error {
	return m.validate(true)
}

func (m *EmailSettings) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for EmailId

	if len(m.GetBlockedEmails()) > 25 {
		err := EmailSettingsValidationError{
			field:  "BlockedEmails",
			reason: "value must contain no more than 25 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return EmailSettingsMultiError(errors)
	}

	return nil
}

// EmailSettingsMultiError is an error wrapping multiple validation errors
// returned by EmailSettings.ValidateAll() if the designated constraints
// aren't met.
type EmailSettingsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EmailSettingsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EmailSettingsMultiError) AllErrors() []error { return m }

// EmailSettingsValidationError is the validation error returned by
// EmailSettings.Validate if the designated constraints aren't met.
type EmailSettingsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EmailSettingsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EmailSettingsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EmailSettingsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EmailSettingsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EmailSettingsValidationError) ErrorName() string { return "EmailSettingsValidationError" }

// Error satisfies the builtin error interface
func (e EmailSettingsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEmailSettings.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EmailSettingsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EmailSettingsValidationError{}
