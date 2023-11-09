// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/users/jobs.proto

package users

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

// Validate checks the field values on Job with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Job) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Job with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in JobMultiError, or nil if none found.
func (m *Job) ValidateAll() error {
	return m.validate(true)
}

func (m *Job) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) > 50 {
		err := JobValidationError{
			field:  "Name",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetLabel()) > 50 {
		err := JobValidationError{
			field:  "Label",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetGrades() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, JobValidationError{
						field:  fmt.Sprintf("Grades[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, JobValidationError{
						field:  fmt.Sprintf("Grades[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return JobValidationError{
					field:  fmt.Sprintf("Grades[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return JobMultiError(errors)
	}

	return nil
}

// JobMultiError is an error wrapping multiple validation errors returned by
// Job.ValidateAll() if the designated constraints aren't met.
type JobMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m JobMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m JobMultiError) AllErrors() []error { return m }

// JobValidationError is the validation error returned by Job.Validate if the
// designated constraints aren't met.
type JobValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e JobValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e JobValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e JobValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e JobValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e JobValidationError) ErrorName() string { return "JobValidationError" }

// Error satisfies the builtin error interface
func (e JobValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJob.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = JobValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = JobValidationError{}

// Validate checks the field values on JobGrade with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *JobGrade) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on JobGrade with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in JobGradeMultiError, or nil
// if none found.
func (m *JobGrade) ValidateAll() error {
	return m.validate(true)
}

func (m *JobGrade) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetGrade() <= 0 {
		err := JobGradeValidationError{
			field:  "Grade",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetLabel()) > 50 {
		err := JobGradeValidationError{
			field:  "Label",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.JobName != nil {

		if utf8.RuneCountInString(m.GetJobName()) > 50 {
			err := JobGradeValidationError{
				field:  "JobName",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return JobGradeMultiError(errors)
	}

	return nil
}

// JobGradeMultiError is an error wrapping multiple validation errors returned
// by JobGrade.ValidateAll() if the designated constraints aren't met.
type JobGradeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m JobGradeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m JobGradeMultiError) AllErrors() []error { return m }

// JobGradeValidationError is the validation error returned by
// JobGrade.Validate if the designated constraints aren't met.
type JobGradeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e JobGradeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e JobGradeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e JobGradeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e JobGradeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e JobGradeValidationError) ErrorName() string { return "JobGradeValidationError" }

// Error satisfies the builtin error interface
func (e JobGradeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJobGrade.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = JobGradeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = JobGradeValidationError{}

// Validate checks the field values on JobProps with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *JobProps) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on JobProps with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in JobPropsMultiError, or nil
// if none found.
func (m *JobProps) ValidateAll() error {
	return m.validate(true)
}

func (m *JobProps) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := JobPropsValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetTheme()) > 20 {
		err := JobPropsValidationError{
			field:  "Theme",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetLivemapMarkerColor()) != 6 {
		err := JobPropsValidationError{
			field:  "LivemapMarkerColor",
			reason: "value length must be 6 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if !_JobProps_LivemapMarkerColor_Pattern.MatchString(m.GetLivemapMarkerColor()) {
		err := JobPropsValidationError{
			field:  "LivemapMarkerColor",
			reason: "value does not match regex pattern \"^[A-Fa-f0-9]{6}$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetQuickButtons()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, JobPropsValidationError{
					field:  "QuickButtons",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, JobPropsValidationError{
					field:  "QuickButtons",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetQuickButtons()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return JobPropsValidationError{
				field:  "QuickButtons",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetDiscordSyncSettings()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, JobPropsValidationError{
					field:  "DiscordSyncSettings",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, JobPropsValidationError{
					field:  "DiscordSyncSettings",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDiscordSyncSettings()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return JobPropsValidationError{
				field:  "DiscordSyncSettings",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.DiscordGuildId != nil {
		// no validation rules for DiscordGuildId
	}

	if m.DiscordLastSync != nil {

		if all {
			switch v := interface{}(m.GetDiscordLastSync()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, JobPropsValidationError{
						field:  "DiscordLastSync",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, JobPropsValidationError{
						field:  "DiscordLastSync",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDiscordLastSync()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return JobPropsValidationError{
					field:  "DiscordLastSync",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return JobPropsMultiError(errors)
	}

	return nil
}

// JobPropsMultiError is an error wrapping multiple validation errors returned
// by JobProps.ValidateAll() if the designated constraints aren't met.
type JobPropsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m JobPropsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m JobPropsMultiError) AllErrors() []error { return m }

// JobPropsValidationError is the validation error returned by
// JobProps.Validate if the designated constraints aren't met.
type JobPropsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e JobPropsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e JobPropsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e JobPropsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e JobPropsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e JobPropsValidationError) ErrorName() string { return "JobPropsValidationError" }

// Error satisfies the builtin error interface
func (e JobPropsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJobProps.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = JobPropsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = JobPropsValidationError{}

var _JobProps_LivemapMarkerColor_Pattern = regexp.MustCompile("^[A-Fa-f0-9]{6}$")

// Validate checks the field values on QuickButtons with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *QuickButtons) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QuickButtons with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in QuickButtonsMultiError, or
// nil if none found.
func (m *QuickButtons) ValidateAll() error {
	return m.validate(true)
}

func (m *QuickButtons) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PenaltyCalculator

	if len(errors) > 0 {
		return QuickButtonsMultiError(errors)
	}

	return nil
}

// QuickButtonsMultiError is an error wrapping multiple validation errors
// returned by QuickButtons.ValidateAll() if the designated constraints aren't met.
type QuickButtonsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QuickButtonsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QuickButtonsMultiError) AllErrors() []error { return m }

// QuickButtonsValidationError is the validation error returned by
// QuickButtons.Validate if the designated constraints aren't met.
type QuickButtonsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QuickButtonsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QuickButtonsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QuickButtonsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QuickButtonsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QuickButtonsValidationError) ErrorName() string { return "QuickButtonsValidationError" }

// Error satisfies the builtin error interface
func (e QuickButtonsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQuickButtons.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QuickButtonsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QuickButtonsValidationError{}

// Validate checks the field values on DiscordSyncSettings with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DiscordSyncSettings) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DiscordSyncSettings with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DiscordSyncSettingsMultiError, or nil if none found.
func (m *DiscordSyncSettings) ValidateAll() error {
	return m.validate(true)
}

func (m *DiscordSyncSettings) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserInfoSync

	if len(errors) > 0 {
		return DiscordSyncSettingsMultiError(errors)
	}

	return nil
}

// DiscordSyncSettingsMultiError is an error wrapping multiple validation
// errors returned by DiscordSyncSettings.ValidateAll() if the designated
// constraints aren't met.
type DiscordSyncSettingsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DiscordSyncSettingsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DiscordSyncSettingsMultiError) AllErrors() []error { return m }

// DiscordSyncSettingsValidationError is the validation error returned by
// DiscordSyncSettings.Validate if the designated constraints aren't met.
type DiscordSyncSettingsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DiscordSyncSettingsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DiscordSyncSettingsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DiscordSyncSettingsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DiscordSyncSettingsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DiscordSyncSettingsValidationError) ErrorName() string {
	return "DiscordSyncSettingsValidationError"
}

// Error satisfies the builtin error interface
func (e DiscordSyncSettingsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDiscordSyncSettings.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DiscordSyncSettingsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DiscordSyncSettingsValidationError{}
