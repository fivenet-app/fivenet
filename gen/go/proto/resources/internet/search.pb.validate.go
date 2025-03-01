// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/internet/search.proto

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

// Validate checks the field values on SearchResult with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SearchResult) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SearchResult with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SearchResultMultiError, or
// nil if none found.
func (m *SearchResult) ValidateAll() error {
	return m.validate(true)
}

func (m *SearchResult) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Title

	// no validation rules for Description

	// no validation rules for DomainId

	// no validation rules for Path

	if m.Domain != nil {

		if all {
			switch v := interface{}(m.GetDomain()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, SearchResultValidationError{
						field:  "Domain",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, SearchResultValidationError{
						field:  "Domain",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDomain()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SearchResultValidationError{
					field:  "Domain",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return SearchResultMultiError(errors)
	}

	return nil
}

// SearchResultMultiError is an error wrapping multiple validation errors
// returned by SearchResult.ValidateAll() if the designated constraints aren't met.
type SearchResultMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SearchResultMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SearchResultMultiError) AllErrors() []error { return m }

// SearchResultValidationError is the validation error returned by
// SearchResult.Validate if the designated constraints aren't met.
type SearchResultValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SearchResultValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SearchResultValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SearchResultValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SearchResultValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SearchResultValidationError) ErrorName() string { return "SearchResultValidationError" }

// Error satisfies the builtin error interface
func (e SearchResultValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSearchResult.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SearchResultValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SearchResultValidationError{}
