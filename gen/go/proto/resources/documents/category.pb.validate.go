// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/documents/category.proto

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

// Validate checks the field values on Category with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Category) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Category with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CategoryMultiError, or nil
// if none found.
func (m *Category) ValidateAll() error {
	return m.validate(true)
}

func (m *Category) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CategoryValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CategoryValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CategoryValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if l := utf8.RuneCountInString(m.GetName()); l < 3 || l > 128 {
		err := CategoryValidationError{
			field:  "Name",
			reason: "value length must be between 3 and 128 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.DeletedAt != nil {

		if all {
			switch v := interface{}(m.GetDeletedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CategoryValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CategoryValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDeletedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CategoryValidationError{
					field:  "DeletedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Description != nil {

		if utf8.RuneCountInString(m.GetDescription()) > 255 {
			err := CategoryValidationError{
				field:  "Description",
				reason: "value length must be at most 255 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Job != nil {

		if utf8.RuneCountInString(m.GetJob()) > 20 {
			err := CategoryValidationError{
				field:  "Job",
				reason: "value length must be at most 20 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Color != nil {

		if l := utf8.RuneCountInString(m.GetColor()); l < 3 || l > 7 {
			err := CategoryValidationError{
				field:  "Color",
				reason: "value length must be between 3 and 7 runes, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Icon != nil {

		if utf8.RuneCountInString(m.GetIcon()) > 128 {
			err := CategoryValidationError{
				field:  "Icon",
				reason: "value length must be at most 128 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return CategoryMultiError(errors)
	}

	return nil
}

// CategoryMultiError is an error wrapping multiple validation errors returned
// by Category.ValidateAll() if the designated constraints aren't met.
type CategoryMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CategoryMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CategoryMultiError) AllErrors() []error { return m }

// CategoryValidationError is the validation error returned by
// Category.Validate if the designated constraints aren't met.
type CategoryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CategoryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CategoryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CategoryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CategoryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CategoryValidationError) ErrorName() string { return "CategoryValidationError" }

// Error satisfies the builtin error interface
func (e CategoryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCategory.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CategoryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CategoryValidationError{}
