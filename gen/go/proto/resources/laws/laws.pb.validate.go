// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/laws/laws.proto

package laws

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

// Validate checks the field values on LawBook with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LawBook) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LawBook with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in LawBookMultiError, or nil if none found.
func (m *LawBook) ValidateAll() error {
	return m.validate(true)
}

func (m *LawBook) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if l := utf8.RuneCountInString(m.GetName()); l < 3 || l > 128 {
		err := LawBookValidationError{
			field:  "Name",
			reason: "value length must be between 3 and 128 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetLaws() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, LawBookValidationError{
						field:  fmt.Sprintf("Laws[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, LawBookValidationError{
						field:  fmt.Sprintf("Laws[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LawBookValidationError{
					field:  fmt.Sprintf("Laws[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.CreatedAt != nil {

		if all {
			switch v := interface{}(m.GetCreatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, LawBookValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, LawBookValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LawBookValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.UpdatedAt != nil {

		if all {
			switch v := interface{}(m.GetUpdatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, LawBookValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, LawBookValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LawBookValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Description != nil {

		if utf8.RuneCountInString(m.GetDescription()) > 255 {
			err := LawBookValidationError{
				field:  "Description",
				reason: "value length must be at most 255 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return LawBookMultiError(errors)
	}

	return nil
}

// LawBookMultiError is an error wrapping multiple validation errors returned
// by LawBook.ValidateAll() if the designated constraints aren't met.
type LawBookMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LawBookMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LawBookMultiError) AllErrors() []error { return m }

// LawBookValidationError is the validation error returned by LawBook.Validate
// if the designated constraints aren't met.
type LawBookValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LawBookValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LawBookValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LawBookValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LawBookValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LawBookValidationError) ErrorName() string { return "LawBookValidationError" }

// Error satisfies the builtin error interface
func (e LawBookValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLawBook.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LawBookValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LawBookValidationError{}

// Validate checks the field values on Law with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Law) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Law with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in LawMultiError, or nil if none found.
func (m *Law) ValidateAll() error {
	return m.validate(true)
}

func (m *Law) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for LawbookId

	if l := utf8.RuneCountInString(m.GetName()); l < 3 || l > 128 {
		err := LawValidationError{
			field:  "Name",
			reason: "value length must be between 3 and 128 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.CreatedAt != nil {

		if all {
			switch v := interface{}(m.GetCreatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, LawValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, LawValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LawValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.UpdatedAt != nil {

		if all {
			switch v := interface{}(m.GetUpdatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, LawValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, LawValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LawValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Description != nil {

		if utf8.RuneCountInString(m.GetDescription()) > 1024 {
			err := LawValidationError{
				field:  "Description",
				reason: "value length must be at most 1024 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Hint != nil {

		if utf8.RuneCountInString(m.GetHint()) > 512 {
			err := LawValidationError{
				field:  "Hint",
				reason: "value length must be at most 512 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Fine != nil {
		// no validation rules for Fine
	}

	if m.DetentionTime != nil {
		// no validation rules for DetentionTime
	}

	if m.StvoPoints != nil {
		// no validation rules for StvoPoints
	}

	if len(errors) > 0 {
		return LawMultiError(errors)
	}

	return nil
}

// LawMultiError is an error wrapping multiple validation errors returned by
// Law.ValidateAll() if the designated constraints aren't met.
type LawMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LawMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LawMultiError) AllErrors() []error { return m }

// LawValidationError is the validation error returned by Law.Validate if the
// designated constraints aren't met.
type LawValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LawValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LawValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LawValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LawValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LawValidationError) ErrorName() string { return "LawValidationError" }

// Error satisfies the builtin error interface
func (e LawValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLaw.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LawValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LawValidationError{}
