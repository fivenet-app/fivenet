// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/messenger/access.proto

package messenger

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

// Validate checks the field values on ThreadAccess with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ThreadAccess) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ThreadAccess with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ThreadAccessMultiError, or
// nil if none found.
func (m *ThreadAccess) ValidateAll() error {
	return m.validate(true)
}

func (m *ThreadAccess) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetJobs()) > 20 {
		err := ThreadAccessValidationError{
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
					errors = append(errors, ThreadAccessValidationError{
						field:  fmt.Sprintf("Jobs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadAccessValidationError{
						field:  fmt.Sprintf("Jobs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadAccessValidationError{
					field:  fmt.Sprintf("Jobs[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(m.GetUsers()) > 20 {
		err := ThreadAccessValidationError{
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
					errors = append(errors, ThreadAccessValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadAccessValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadAccessValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ThreadAccessMultiError(errors)
	}

	return nil
}

// ThreadAccessMultiError is an error wrapping multiple validation errors
// returned by ThreadAccess.ValidateAll() if the designated constraints aren't met.
type ThreadAccessMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ThreadAccessMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ThreadAccessMultiError) AllErrors() []error { return m }

// ThreadAccessValidationError is the validation error returned by
// ThreadAccess.Validate if the designated constraints aren't met.
type ThreadAccessValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ThreadAccessValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ThreadAccessValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ThreadAccessValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ThreadAccessValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ThreadAccessValidationError) ErrorName() string { return "ThreadAccessValidationError" }

// Error satisfies the builtin error interface
func (e ThreadAccessValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sThreadAccess.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ThreadAccessValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ThreadAccessValidationError{}

// Validate checks the field values on ThreadJobAccess with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ThreadJobAccess) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ThreadJobAccess with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ThreadJobAccessMultiError, or nil if none found.
func (m *ThreadJobAccess) ValidateAll() error {
	return m.validate(true)
}

func (m *ThreadJobAccess) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for ThreadId

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := ThreadJobAccessValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetMinimumGrade() <= 0 {
		err := ThreadJobAccessValidationError{
			field:  "MinimumGrade",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := AccessLevel_name[int32(m.GetAccess())]; !ok {
		err := ThreadJobAccessValidationError{
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
					errors = append(errors, ThreadJobAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadJobAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadJobAccessValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.JobLabel != nil {

		if utf8.RuneCountInString(m.GetJobLabel()) > 50 {
			err := ThreadJobAccessValidationError{
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
			err := ThreadJobAccessValidationError{
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
		return ThreadJobAccessMultiError(errors)
	}

	return nil
}

// ThreadJobAccessMultiError is an error wrapping multiple validation errors
// returned by ThreadJobAccess.ValidateAll() if the designated constraints
// aren't met.
type ThreadJobAccessMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ThreadJobAccessMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ThreadJobAccessMultiError) AllErrors() []error { return m }

// ThreadJobAccessValidationError is the validation error returned by
// ThreadJobAccess.Validate if the designated constraints aren't met.
type ThreadJobAccessValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ThreadJobAccessValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ThreadJobAccessValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ThreadJobAccessValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ThreadJobAccessValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ThreadJobAccessValidationError) ErrorName() string { return "ThreadJobAccessValidationError" }

// Error satisfies the builtin error interface
func (e ThreadJobAccessValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sThreadJobAccess.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ThreadJobAccessValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ThreadJobAccessValidationError{}

// Validate checks the field values on ThreadUserAccess with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ThreadUserAccess) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ThreadUserAccess with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ThreadUserAccessMultiError, or nil if none found.
func (m *ThreadUserAccess) ValidateAll() error {
	return m.validate(true)
}

func (m *ThreadUserAccess) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for ThreadId

	if m.GetUserId() <= 0 {
		err := ThreadUserAccessValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := AccessLevel_name[int32(m.GetAccess())]; !ok {
		err := ThreadUserAccessValidationError{
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
					errors = append(errors, ThreadUserAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadUserAccessValidationError{
						field:  "CreatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadUserAccessValidationError{
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
					errors = append(errors, ThreadUserAccessValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ThreadUserAccessValidationError{
						field:  "User",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUser()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ThreadUserAccessValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ThreadUserAccessMultiError(errors)
	}

	return nil
}

// ThreadUserAccessMultiError is an error wrapping multiple validation errors
// returned by ThreadUserAccess.ValidateAll() if the designated constraints
// aren't met.
type ThreadUserAccessMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ThreadUserAccessMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ThreadUserAccessMultiError) AllErrors() []error { return m }

// ThreadUserAccessValidationError is the validation error returned by
// ThreadUserAccess.Validate if the designated constraints aren't met.
type ThreadUserAccessValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ThreadUserAccessValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ThreadUserAccessValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ThreadUserAccessValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ThreadUserAccessValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ThreadUserAccessValidationError) ErrorName() string { return "ThreadUserAccessValidationError" }

// Error satisfies the builtin error interface
func (e ThreadUserAccessValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sThreadUserAccess.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ThreadUserAccessValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ThreadUserAccessValidationError{}
