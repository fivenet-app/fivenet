// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/documents/pins.proto

package documents

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

// Validate checks the field values on DocumentPin with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DocumentPin) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DocumentPin with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DocumentPinMultiError, or
// nil if none found.
func (m *DocumentPin) ValidateAll() error {
	return m.validate(true)
}

func (m *DocumentPin) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetDocumentId() <= 0 {
		err := DocumentPinValidationError{
			field:  "DocumentId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for State

	// no validation rules for CreatorId

	if m.Job != nil {

		if utf8.RuneCountInString(m.GetJob()) > 20 {
			err := DocumentPinValidationError{
				field:  "Job",
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
			err := DocumentPinValidationError{
				field:  "UserId",
				reason: "value must be greater than 0",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.CreatedAt != nil {

		if all {
			switch v := interface{}(m.GetCreatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DocumentPinValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DocumentPinValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DocumentPinValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return DocumentPinMultiError(errors)
	}

	return nil
}

// DocumentPinMultiError is an error wrapping multiple validation errors
// returned by DocumentPin.ValidateAll() if the designated constraints aren't met.
type DocumentPinMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocumentPinMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocumentPinMultiError) AllErrors() []error { return m }

// DocumentPinValidationError is the validation error returned by
// DocumentPin.Validate if the designated constraints aren't met.
type DocumentPinValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocumentPinValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocumentPinValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocumentPinValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocumentPinValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocumentPinValidationError) ErrorName() string { return "DocumentPinValidationError" }

// Error satisfies the builtin error interface
func (e DocumentPinValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocumentPin.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocumentPinValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocumentPinValidationError{}
