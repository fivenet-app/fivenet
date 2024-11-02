// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/wiki/page.proto

package wiki

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

// Validate checks the field values on Page with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Page) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Page with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PageMultiError, or nil if none found.
func (m *Page) ValidateAll() error {
	return m.validate(true)
}

func (m *Page) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Job

	// no validation rules for Path

	// no validation rules for ContentType

	if all {
		switch v := interface{}(m.GetMeta()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PageValidationError{
					field:  "Meta",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PageValidationError{
					field:  "Meta",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMeta()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PageValidationError{
				field:  "Meta",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Content

	if all {
		switch v := interface{}(m.GetAccess()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PageValidationError{
					field:  "Access",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PageValidationError{
					field:  "Access",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccess()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PageValidationError{
				field:  "Access",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PageMultiError(errors)
	}

	return nil
}

// PageMultiError is an error wrapping multiple validation errors returned by
// Page.ValidateAll() if the designated constraints aren't met.
type PageMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PageMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PageMultiError) AllErrors() []error { return m }

// PageValidationError is the validation error returned by Page.Validate if the
// designated constraints aren't met.
type PageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PageValidationError) ErrorName() string { return "PageValidationError" }

// Error satisfies the builtin error interface
func (e PageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPage.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PageValidationError{}

// Validate checks the field values on PageMeta with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PageMeta) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PageMeta with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PageMetaMultiError, or nil
// if none found.
func (m *PageMeta) ValidateAll() error {
	return m.validate(true)
}

func (m *PageMeta) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Title

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PageMetaValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PageMetaValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PageMetaValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PageMetaValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PageMetaValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PageMetaValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetAuthor()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PageMetaValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PageMetaValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAuthor()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PageMetaValidationError{
				field:  "Author",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Description

	if len(errors) > 0 {
		return PageMetaMultiError(errors)
	}

	return nil
}

// PageMetaMultiError is an error wrapping multiple validation errors returned
// by PageMeta.ValidateAll() if the designated constraints aren't met.
type PageMetaMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PageMetaMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PageMetaMultiError) AllErrors() []error { return m }

// PageMetaValidationError is the validation error returned by
// PageMeta.Validate if the designated constraints aren't met.
type PageMetaValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PageMetaValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PageMetaValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PageMetaValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PageMetaValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PageMetaValidationError) ErrorName() string { return "PageMetaValidationError" }

// Error satisfies the builtin error interface
func (e PageMetaValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPageMeta.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PageMetaValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PageMetaValidationError{}
