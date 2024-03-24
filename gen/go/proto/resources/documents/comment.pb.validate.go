// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/documents/comment.proto

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

// Validate checks the field values on Comment with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Comment) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Comment with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in CommentMultiError, or nil if none found.
func (m *Comment) ValidateAll() error {
	return m.validate(true)
}

func (m *Comment) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for DocumentId

	if utf8.RuneCountInString(m.GetComment()) < 3 {
		err := CommentValidationError{
			field:  "Comment",
			reason: "value length must be at least 3 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetComment()) > 2048 {
		err := CommentValidationError{
			field:  "Comment",
			reason: "value length must be at most 2048 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetCreatorJob()) > 20 {
		err := CommentValidationError{
			field:  "CreatorJob",
			reason: "value length must be at most 20 runes",
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
					errors = append(errors, CommentValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CommentValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CommentValidationError{
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
					errors = append(errors, CommentValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CommentValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CommentValidationError{
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
					errors = append(errors, CommentValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CommentValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDeletedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CommentValidationError{
					field:  "DeletedAt",
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
					errors = append(errors, CommentValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CommentValidationError{
						field:  "Creator",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreator()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CommentValidationError{
					field:  "Creator",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CommentMultiError(errors)
	}

	return nil
}

// CommentMultiError is an error wrapping multiple validation errors returned
// by Comment.ValidateAll() if the designated constraints aren't met.
type CommentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CommentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CommentMultiError) AllErrors() []error { return m }

// CommentValidationError is the validation error returned by Comment.Validate
// if the designated constraints aren't met.
type CommentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CommentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CommentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CommentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CommentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CommentValidationError) ErrorName() string { return "CommentValidationError" }

// Error satisfies the builtin error interface
func (e CommentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sComment.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CommentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CommentValidationError{}
