// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/internet/access.proto

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

// Validate checks the field values on PageAccess with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PageAccess) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PageAccess with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PageAccessMultiError, or
// nil if none found.
func (m *PageAccess) ValidateAll() error {
	return m.validate(true)
}

func (m *PageAccess) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetJobs()) > 20 {
		err := PageAccessValidationError{
			field:  "Jobs",
			reason: "value must contain no more than 20 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetJobs() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PageAccessValidationError{
						field:  fmt.Sprintf("Jobs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PageAccessValidationError{
						field:  fmt.Sprintf("Jobs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PageAccessValidationError{
					field:  fmt.Sprintf("Jobs[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(m.GetUsers()) > 20 {
		err := PageAccessValidationError{
			field:  "Users",
			reason: "value must contain no more than 20 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetUsers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PageAccessValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PageAccessValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PageAccessValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PageAccessMultiError(errors)
	}

	return nil
}

// PageAccessMultiError is an error wrapping multiple validation errors
// returned by PageAccess.ValidateAll() if the designated constraints aren't met.
type PageAccessMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PageAccessMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PageAccessMultiError) AllErrors() []error { return m }

// PageAccessValidationError is the validation error returned by
// PageAccess.Validate if the designated constraints aren't met.
type PageAccessValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PageAccessValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PageAccessValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PageAccessValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PageAccessValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PageAccessValidationError) ErrorName() string { return "PageAccessValidationError" }

// Error satisfies the builtin error interface
func (e PageAccessValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPageAccess.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PageAccessValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PageAccessValidationError{}

// Validate checks the field values on PageJobAccess with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PageJobAccess) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PageJobAccess with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PageJobAccessMultiError, or
// nil if none found.
func (m *PageJobAccess) ValidateAll() error {
	return m.validate(true)
}

func (m *PageJobAccess) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for TargetId

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := PageJobAccessValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetMinimumGrade() < 0 {
		err := PageJobAccessValidationError{
			field:  "MinimumGrade",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := AccessLevel_name[int32(m.GetAccess())]; !ok {
		err := PageJobAccessValidationError{
			field:  "Access",
			reason: "value must be one of the defined enum values",
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
					errors = append(errors, PageJobAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PageJobAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PageJobAccessValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.JobLabel != nil {

		if utf8.RuneCountInString(m.GetJobLabel()) > 50 {
			err := PageJobAccessValidationError{
				field:  "JobLabel",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.JobGradeLabel != nil {

		if utf8.RuneCountInString(m.GetJobGradeLabel()) > 50 {
			err := PageJobAccessValidationError{
				field:  "JobGradeLabel",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return PageJobAccessMultiError(errors)
	}

	return nil
}

// PageJobAccessMultiError is an error wrapping multiple validation errors
// returned by PageJobAccess.ValidateAll() if the designated constraints
// aren't met.
type PageJobAccessMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PageJobAccessMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PageJobAccessMultiError) AllErrors() []error { return m }

// PageJobAccessValidationError is the validation error returned by
// PageJobAccess.Validate if the designated constraints aren't met.
type PageJobAccessValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PageJobAccessValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PageJobAccessValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PageJobAccessValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PageJobAccessValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PageJobAccessValidationError) ErrorName() string { return "PageJobAccessValidationError" }

// Error satisfies the builtin error interface
func (e PageJobAccessValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPageJobAccess.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PageJobAccessValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PageJobAccessValidationError{}

// Validate checks the field values on PageUserAccess with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PageUserAccess) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PageUserAccess with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PageUserAccessMultiError,
// or nil if none found.
func (m *PageUserAccess) ValidateAll() error {
	return m.validate(true)
}

func (m *PageUserAccess) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for TargetId

	if m.GetUserId() <= 0 {
		err := PageUserAccessValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := AccessLevel_name[int32(m.GetAccess())]; !ok {
		err := PageUserAccessValidationError{
			field:  "Access",
			reason: "value must be one of the defined enum values",
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
					errors = append(errors, PageUserAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PageUserAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PageUserAccessValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.User != nil {

		if all {
			switch v := interface{}(m.GetUser()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PageUserAccessValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PageUserAccessValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUser()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PageUserAccessValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PageUserAccessMultiError(errors)
	}

	return nil
}

// PageUserAccessMultiError is an error wrapping multiple validation errors
// returned by PageUserAccess.ValidateAll() if the designated constraints
// aren't met.
type PageUserAccessMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PageUserAccessMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PageUserAccessMultiError) AllErrors() []error { return m }

// PageUserAccessValidationError is the validation error returned by
// PageUserAccess.Validate if the designated constraints aren't met.
type PageUserAccessValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PageUserAccessValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PageUserAccessValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PageUserAccessValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PageUserAccessValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PageUserAccessValidationError) ErrorName() string { return "PageUserAccessValidationError" }

// Error satisfies the builtin error interface
func (e PageUserAccessValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPageUserAccess.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PageUserAccessValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PageUserAccessValidationError{}
