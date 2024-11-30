// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: services/notificator/notificator.proto

package notificator

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

	notifications "github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
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

	_ = notifications.NotificationCategory(0)
)

// Validate checks the field values on GetNotificationsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetNotificationsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetNotificationsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetNotificationsRequestMultiError, or nil if none found.
func (m *GetNotificationsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetNotificationsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPagination() == nil {
		err := GetNotificationsRequestValidationError{
			field:  "Pagination",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetPagination()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetNotificationsRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetNotificationsRequestValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetNotificationsRequestValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(m.GetCategories()) > 4 {
		err := GetNotificationsRequestValidationError{
			field:  "Categories",
			reason: "value must contain no more than 4 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetCategories() {
		_, _ = idx, item

		if _, ok := notifications.NotificationCategory_name[int32(item)]; !ok {
			err := GetNotificationsRequestValidationError{
				field:  fmt.Sprintf("Categories[%v]", idx),
				reason: "value must be one of the defined enum values",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.IncludeRead != nil {
		// no validation rules for IncludeRead
	}

	if len(errors) > 0 {
		return GetNotificationsRequestMultiError(errors)
	}

	return nil
}

// GetNotificationsRequestMultiError is an error wrapping multiple validation
// errors returned by GetNotificationsRequest.ValidateAll() if the designated
// constraints aren't met.
type GetNotificationsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetNotificationsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetNotificationsRequestMultiError) AllErrors() []error { return m }

// GetNotificationsRequestValidationError is the validation error returned by
// GetNotificationsRequest.Validate if the designated constraints aren't met.
type GetNotificationsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetNotificationsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetNotificationsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetNotificationsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetNotificationsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetNotificationsRequestValidationError) ErrorName() string {
	return "GetNotificationsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetNotificationsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetNotificationsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetNotificationsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetNotificationsRequestValidationError{}

// Validate checks the field values on GetNotificationsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetNotificationsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetNotificationsResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetNotificationsResponseMultiError, or nil if none found.
func (m *GetNotificationsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetNotificationsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPagination()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetNotificationsResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetNotificationsResponseValidationError{
					field:  "Pagination",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPagination()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetNotificationsResponseValidationError{
				field:  "Pagination",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetNotifications() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetNotificationsResponseValidationError{
						field:  fmt.Sprintf("Notifications[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetNotificationsResponseValidationError{
						field:  fmt.Sprintf("Notifications[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetNotificationsResponseValidationError{
					field:  fmt.Sprintf("Notifications[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetNotificationsResponseMultiError(errors)
	}

	return nil
}

// GetNotificationsResponseMultiError is an error wrapping multiple validation
// errors returned by GetNotificationsResponse.ValidateAll() if the designated
// constraints aren't met.
type GetNotificationsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetNotificationsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetNotificationsResponseMultiError) AllErrors() []error { return m }

// GetNotificationsResponseValidationError is the validation error returned by
// GetNotificationsResponse.Validate if the designated constraints aren't met.
type GetNotificationsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetNotificationsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetNotificationsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetNotificationsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetNotificationsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetNotificationsResponseValidationError) ErrorName() string {
	return "GetNotificationsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetNotificationsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetNotificationsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetNotificationsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetNotificationsResponseValidationError{}

// Validate checks the field values on MarkNotificationsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *MarkNotificationsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MarkNotificationsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// MarkNotificationsRequestMultiError, or nil if none found.
func (m *MarkNotificationsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *MarkNotificationsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetIds()) > 0 {

		if l := len(m.GetIds()); l < 1 || l > 20 {
			err := MarkNotificationsRequestValidationError{
				field:  "Ids",
				reason: "value must contain between 1 and 20 items, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.All != nil {
		// no validation rules for All
	}

	if len(errors) > 0 {
		return MarkNotificationsRequestMultiError(errors)
	}

	return nil
}

// MarkNotificationsRequestMultiError is an error wrapping multiple validation
// errors returned by MarkNotificationsRequest.ValidateAll() if the designated
// constraints aren't met.
type MarkNotificationsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MarkNotificationsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MarkNotificationsRequestMultiError) AllErrors() []error { return m }

// MarkNotificationsRequestValidationError is the validation error returned by
// MarkNotificationsRequest.Validate if the designated constraints aren't met.
type MarkNotificationsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MarkNotificationsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MarkNotificationsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MarkNotificationsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MarkNotificationsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MarkNotificationsRequestValidationError) ErrorName() string {
	return "MarkNotificationsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e MarkNotificationsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMarkNotificationsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MarkNotificationsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MarkNotificationsRequestValidationError{}

// Validate checks the field values on MarkNotificationsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *MarkNotificationsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MarkNotificationsResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// MarkNotificationsResponseMultiError, or nil if none found.
func (m *MarkNotificationsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *MarkNotificationsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Updated

	if len(errors) > 0 {
		return MarkNotificationsResponseMultiError(errors)
	}

	return nil
}

// MarkNotificationsResponseMultiError is an error wrapping multiple validation
// errors returned by MarkNotificationsResponse.ValidateAll() if the
// designated constraints aren't met.
type MarkNotificationsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MarkNotificationsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MarkNotificationsResponseMultiError) AllErrors() []error { return m }

// MarkNotificationsResponseValidationError is the validation error returned by
// MarkNotificationsResponse.Validate if the designated constraints aren't met.
type MarkNotificationsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MarkNotificationsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MarkNotificationsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MarkNotificationsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MarkNotificationsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MarkNotificationsResponseValidationError) ErrorName() string {
	return "MarkNotificationsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e MarkNotificationsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMarkNotificationsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MarkNotificationsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MarkNotificationsResponseValidationError{}

// Validate checks the field values on StreamRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *StreamRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StreamRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in StreamRequestMultiError, or
// nil if none found.
func (m *StreamRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StreamRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return StreamRequestMultiError(errors)
	}

	return nil
}

// StreamRequestMultiError is an error wrapping multiple validation errors
// returned by StreamRequest.ValidateAll() if the designated constraints
// aren't met.
type StreamRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StreamRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StreamRequestMultiError) AllErrors() []error { return m }

// StreamRequestValidationError is the validation error returned by
// StreamRequest.Validate if the designated constraints aren't met.
type StreamRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StreamRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StreamRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StreamRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StreamRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StreamRequestValidationError) ErrorName() string { return "StreamRequestValidationError" }

// Error satisfies the builtin error interface
func (e StreamRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStreamRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StreamRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StreamRequestValidationError{}

// Validate checks the field values on StreamResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *StreamResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StreamResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in StreamResponseMultiError,
// or nil if none found.
func (m *StreamResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StreamResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for NotificationCount

	switch v := m.Data.(type) {
	case *StreamResponse_UserEvent:
		if v == nil {
			err := StreamResponseValidationError{
				field:  "Data",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetUserEvent()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "UserEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "UserEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUserEvent()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return StreamResponseValidationError{
					field:  "UserEvent",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *StreamResponse_JobEvent:
		if v == nil {
			err := StreamResponseValidationError{
				field:  "Data",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetJobEvent()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "JobEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "JobEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetJobEvent()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return StreamResponseValidationError{
					field:  "JobEvent",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *StreamResponse_JobGradeEvent:
		if v == nil {
			err := StreamResponseValidationError{
				field:  "Data",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetJobGradeEvent()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "JobGradeEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "JobGradeEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetJobGradeEvent()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return StreamResponseValidationError{
					field:  "JobGradeEvent",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *StreamResponse_SystemEvent:
		if v == nil {
			err := StreamResponseValidationError{
				field:  "Data",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetSystemEvent()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "SystemEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "SystemEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetSystemEvent()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return StreamResponseValidationError{
					field:  "SystemEvent",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *StreamResponse_MailerEvent:
		if v == nil {
			err := StreamResponseValidationError{
				field:  "Data",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetMailerEvent()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "MailerEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, StreamResponseValidationError{
						field:  "MailerEvent",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMailerEvent()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return StreamResponseValidationError{
					field:  "MailerEvent",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		_ = v // ensures v is used
	}

	if m.Restart != nil {
		// no validation rules for Restart
	}

	if len(errors) > 0 {
		return StreamResponseMultiError(errors)
	}

	return nil
}

// StreamResponseMultiError is an error wrapping multiple validation errors
// returned by StreamResponse.ValidateAll() if the designated constraints
// aren't met.
type StreamResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StreamResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StreamResponseMultiError) AllErrors() []error { return m }

// StreamResponseValidationError is the validation error returned by
// StreamResponse.Validate if the designated constraints aren't met.
type StreamResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StreamResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StreamResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StreamResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StreamResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StreamResponseValidationError) ErrorName() string { return "StreamResponseValidationError" }

// Error satisfies the builtin error interface
func (e StreamResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStreamResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StreamResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StreamResponseValidationError{}
