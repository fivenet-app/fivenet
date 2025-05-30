// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/audit/audit.proto

package audit

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

// Validate checks the field values on AuditEntry with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *AuditEntry) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuditEntry with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AuditEntryMultiError, or
// nil if none found.
func (m *AuditEntry) ValidateAll() error {
	return m.validate(true)
}

func (m *AuditEntry) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, AuditEntryValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, AuditEntryValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AuditEntryValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for UserId

	// no validation rules for UserJob

	// no validation rules for Service

	// no validation rules for Method

	if _, ok := EventType_name[int32(m.GetState())]; !ok {
		err := AuditEntryValidationError{
			field:  "State",
			reason: "value must be one of the defined enum values",
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
					errors = append(errors, AuditEntryValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, AuditEntryValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUser()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AuditEntryValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.TargetUserId != nil {
		// no validation rules for TargetUserId
	}

	if m.TargetUser != nil {

		if all {
			switch v := interface{}(m.GetTargetUser()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, AuditEntryValidationError{
						field:  "TargetUser",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, AuditEntryValidationError{
						field:  "TargetUser",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetTargetUser()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AuditEntryValidationError{
					field:  "TargetUser",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.TargetUserJob != nil {
		// no validation rules for TargetUserJob
	}

	if m.Data != nil {
		// no validation rules for Data
	}

	if len(errors) > 0 {
		return AuditEntryMultiError(errors)
	}

	return nil
}

// AuditEntryMultiError is an error wrapping multiple validation errors
// returned by AuditEntry.ValidateAll() if the designated constraints aren't met.
type AuditEntryMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuditEntryMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuditEntryMultiError) AllErrors() []error { return m }

// AuditEntryValidationError is the validation error returned by
// AuditEntry.Validate if the designated constraints aren't met.
type AuditEntryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuditEntryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuditEntryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuditEntryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuditEntryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuditEntryValidationError) ErrorName() string { return "AuditEntryValidationError" }

// Error satisfies the builtin error interface
func (e AuditEntryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuditEntry.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuditEntryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuditEntryValidationError{}
