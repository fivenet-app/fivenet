// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/internet/internet.proto

package internet

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

// Validate checks the field values on Domain with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Domain) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Domain with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in DomainMultiError, or nil if none found.
func (m *Domain) ValidateAll() error {
	return m.validate(true)
}

func (m *Domain) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Domain

	if len(errors) > 0 {
		return DomainMultiError(errors)
	}

	return nil
}

// DomainMultiError is an error wrapping multiple validation errors returned by
// Domain.ValidateAll() if the designated constraints aren't met.
type DomainMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DomainMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DomainMultiError) AllErrors() []error { return m }

// DomainValidationError is the validation error returned by Domain.Validate if
// the designated constraints aren't met.
type DomainValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DomainValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DomainValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DomainValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DomainValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DomainValidationError) ErrorName() string { return "DomainValidationError" }

// Error satisfies the builtin error interface
func (e DomainValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDomain.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DomainValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DomainValidationError{}

// Validate checks the field values on WebPage with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *WebPage) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WebPage with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in WebPageMultiError, or nil if none found.
func (m *WebPage) ValidateAll() error {
	return m.validate(true)
}

func (m *WebPage) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for DomainId

	// no validation rules for Url

	if len(errors) > 0 {
		return WebPageMultiError(errors)
	}

	return nil
}

// WebPageMultiError is an error wrapping multiple validation errors returned
// by WebPage.ValidateAll() if the designated constraints aren't met.
type WebPageMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WebPageMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WebPageMultiError) AllErrors() []error { return m }

// WebPageValidationError is the validation error returned by WebPage.Validate
// if the designated constraints aren't met.
type WebPageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WebPageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WebPageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WebPageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WebPageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WebPageValidationError) ErrorName() string { return "WebPageValidationError" }

// Error satisfies the builtin error interface
func (e WebPageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWebPage.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WebPageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WebPageValidationError{}